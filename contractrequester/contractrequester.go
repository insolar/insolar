//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package contractrequester

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/api"
	"github.com/insolar/insolar/insolar/bus"
	busMeta "github.com/insolar/insolar/insolar/bus/meta"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
)

// ContractRequester helps to call contracts
type ContractRequester struct {
	Sender                     bus.Sender                         `inject:""`
	PulseAccessor              pulse.Accessor                     `inject:""`
	JetCoordinator             jet.Coordinator                    `inject:""`
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`

	UnwantedResponseCallback func(ctx context.Context, msg payload.Payload) error

	ResultMutex sync.Mutex
	ResultMap   map[[insolar.RecordHashSize]byte]chan *payload.ReturnResults

	inRequestResultsRouter *message.Router

	// callTimeout is mainly needed for unit tests which
	// sometimes may unpredictably fail on CI with a default timeout
	callTimeout time.Duration
}

// ensure that ContractRequester implements insolar.ContractRequester
var _ insolar.ContractRequester = &ContractRequester{}

// New creates new ContractRequester
func New(ctx context.Context, subscriber message.Subscriber, b *bus.Bus) (*ContractRequester, error) {
	wmLogger := log.NewWatermillLogAdapter(inslogger.FromContext(ctx))
	inRouter, err := message.NewRouter(message.RouterConfig{}, wmLogger)
	if err != nil {
		return nil, err
	}

	cr := &ContractRequester{
		ResultMap:              make(map[[insolar.RecordHashSize]byte]chan *payload.ReturnResults),
		callTimeout:            25 * time.Second,
		inRequestResultsRouter: inRouter,
	}

	inRouter.AddMiddleware(
		middleware.InstantAck,
		b.IncomingMessageRouter,
	)

	inRouter.AddNoPublisherHandler(
		"IncomingRequestResults",
		bus.TopicIncomingRequestResults,
		subscriber,
		cr.ReceiveResult,
	)

	startRouter(ctx, inRouter)

	return cr, nil
}

func randomUint64() uint64 {
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return binary.LittleEndian.Uint64(buf)
}

// SendRequest makes synchronously call to method of contract by its ref without additional information
func (cr *ContractRequester) SendRequest(ctx context.Context, ref *insolar.Reference, method string, argsIn []interface{}) (insolar.Reply, error) {
	pulse, err := cr.PulseAccessor.Latest(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "[ ContractRequester::SendRequest ] Couldn't fetch current pulse")
	}

	r, _, err := cr.SendRequestWithPulse(ctx, ref, method, argsIn, pulse.PulseNumber)
	return r, err
}

func (cr *ContractRequester) SendRequestWithPulse(ctx context.Context, ref *insolar.Reference, method string, argsIn []interface{}, pulse insolar.PulseNumber) (insolar.Reply, *insolar.Reference, error) {
	ctx, span := instracer.StartSpan(ctx, "SendRequest "+method)
	defer span.End()

	args, err := insolar.Serialize(argsIn)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ ContractRequester::SendRequest ] Can't marshal")
	}

	msg := &payload.CallMethod{
		Request: &record.IncomingRequest{
			Object:       ref,
			Method:       method,
			Arguments:    args,
			APIRequestID: utils.TraceID(ctx),
			Reason:       api.MakeReason(pulse, args),
			APINode:      cr.JetCoordinator.Me(),
		},
	}

	routResult, ref, err := cr.Call(ctx, msg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ ContractRequester::SendRequest ] Can't route call")
	}

	return routResult, ref, nil
}

func (cr *ContractRequester) calcRequestHash(request record.IncomingRequest) ([insolar.RecordHashSize]byte, error) {
	var hash [insolar.RecordHashSize]byte

	virtRec := record.Wrap(&request)
	buf, err := virtRec.Marshal()
	if err != nil {
		return hash, errors.Wrap(err, "[ ContractRequester::calcRequestHash ] Failed to marshal record")
	}

	hasher := cr.PlatformCryptographyScheme.ReferenceHasher()
	copy(hash[:], hasher.Hash(buf)[0:insolar.RecordHashSize])
	return hash, nil
}

func (cr *ContractRequester) checkCall(_ context.Context, msg payload.CallMethod) error {
	switch {
	case msg.Request.Caller.IsEmpty() && msg.Request.APINode.IsEmpty():
		return errors.New("either Caller or APINode should be set, both empty")
	case !msg.Request.Caller.IsEmpty() && !msg.Request.APINode.IsEmpty():
		return errors.New("either Caller or APINode should be set, both set")
	}

	return nil
}

func (cr *ContractRequester) createResultWaiter(
	request record.IncomingRequest,
) (
	chan *payload.ReturnResults, [insolar.RecordHashSize]byte, error,
) {
	cr.ResultMutex.Lock()
	defer cr.ResultMutex.Unlock()

	reqHash, err := cr.calcRequestHash(request)
	if err != nil {
		return nil, reqHash, errors.Wrap(err, "failed to calculate hash")
	}

	ch := make(chan *payload.ReturnResults, 1)
	cr.ResultMap[reqHash] = ch

	return ch, reqHash, nil
}

func (cr *ContractRequester) Call(ctx context.Context, inMsg insolar.Payload) (insolar.Reply, *insolar.Reference, error) {
	msg := inMsg.(*payload.CallMethod)

	ctx, span := instracer.StartSpan(ctx, "ContractRequester.Call")
	defer span.End()

	async := msg.Request.ReturnMode == record.ReturnNoWait

	if msg.Request.Nonce == 0 {
		msg.Request.Nonce = randomUint64()
	}

	logger := inslogger.FromContext(ctx)
	logger.Debug("about to send call to method ", msg.Request.Method)

	err := cr.checkCall(ctx, *msg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "incorrect request")
	}

	var ch chan *payload.ReturnResults
	var reqHash [insolar.RecordHashSize]byte

	if !async {
		var err error
		ch, reqHash, err = cr.createResultWaiter(*msg.Request)
		if err != nil {
			return nil, nil, errors.Wrap(err, "can't create waiter record")
		}
	}

	sender := bus.NewRetrySender(cr.Sender, cr.PulseAccessor, 5, 1)

	message, err := payload.NewMessage(msg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to marshal payload")
	}

	target := record.CalculateRequestAffinityRef(msg.Request, msg.PulseNumber, cr.PlatformCryptographyScheme)

	resp, done := sender.SendRole(ctx, message, insolar.DynamicRoleVirtualExecutor, *target)
	defer done()

	rawResponse := <-resp

	replyData, err := deserializePayload(rawResponse)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to deserialize payload")
	}

	var (
		res insolar.Reply
		r   reply.RegisterRequest
	)

	switch res := replyData.(type) {
	// early result
	case *reply.CallMethod:
		logger.Debug("early result for request, not registered")
		if !async {
			cr.ResultMutex.Lock()
			defer cr.ResultMutex.Unlock()

			delete(cr.ResultMap, reqHash)
			close(ch)
		}

		return res, nil, nil
	default:
		return nil, nil, errors.New("Got not reply.RegisterRequest in reply for CallMethod")
		// request register
	case *reply.RegisterRequest:
		r = *res
	}

	if async {
		return res, &r.Request, nil
	}

	if !bytes.Equal(r.Request.Record().Hash(), reqHash[:]) {
		return nil, &r.Request, errors.New("Registered request has different hash")
	}

	ctx, cancel := context.WithTimeout(ctx, cr.callTimeout)
	defer cancel()

	logger = logger.WithFields(
		map[string]interface{}{
			"called_request": r.Request.String(),
			"called_method":  msg.Request.Method,
		},
	)
	logger.Debug("waiting results of request")

	select {
	case ret := <-ch:
		logger.Debug("got results of request")
		if ret.Error != "" {
			return nil, &r.Request, errors.Wrap(errors.New(ret.Error), "CallMethod returns error")
		}
		res, err := reply.Deserialize(bytes.NewReader(ret.Reply))
		if err != nil {
			return nil, &r.Request, errors.Wrap(err, "failed to deserialize reply")
		}
		return res, &r.Request, nil
	case <-ctx.Done():
		cr.ResultMutex.Lock()
		defer cr.ResultMutex.Unlock()

		delete(cr.ResultMap, reqHash)
		logger.Error("request timeout")
		return nil, &r.Request, errors.Errorf("request to contract was canceled: timeout of %s was exceeded", cr.callTimeout)
	}
}

func (cr *ContractRequester) result(ctx context.Context, msg *payload.ReturnResults) error {
	cr.ResultMutex.Lock()
	defer cr.ResultMutex.Unlock()

	var reqHash [insolar.RecordHashSize]byte
	copy(reqHash[:], msg.RequestRef.Record().Hash())
	c, ok := cr.ResultMap[reqHash]
	if !ok {
		inslogger.FromContext(ctx).Info("unwaited results of request ", msg.RequestRef.String())
		if cr.UnwantedResponseCallback != nil {
			return cr.UnwantedResponseCallback(ctx, msg)
		}
		inslogger.FromContext(ctx).Warn("drop unwanted ", msg.RequestRef.String())
		return nil
	}

	c <- msg
	delete(cr.ResultMap, reqHash)
	return nil
}

func (cr *ContractRequester) ReceiveResult(msg *message.Message) error {
	if msg == nil {
		return errors.New("can't deserialize payload of nil message")
	}

	ctx := inslogger.ContextWithTrace(context.Background(), msg.Metadata.Get(busMeta.TraceID))

	parentSpan, err := instracer.Deserialize([]byte(msg.Metadata.Get(busMeta.SpanData)))
	if err == nil {
		ctx = instracer.WithParentSpan(ctx, parentSpan)
	} else {
		inslogger.FromContext(ctx).Error(err)
	}

	for k, v := range msg.Metadata {
		if k == busMeta.SpanData || k == busMeta.TraceID {
			continue
		}
		ctx, _ = inslogger.WithField(ctx, k, v)
	}

	meta := payload.Meta{}
	err = meta.Unmarshal(msg.Payload)
	if err != nil {
		return err
	}

	payloadType, err := payload.UnmarshalType(meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload type")
	}

	ctx, logger := inslogger.WithField(ctx, "msg_type", payloadType.String())

	logger.Debug("Start to handle new message")

	if payloadType != payload.TypeReturnResults {
		return errors.Errorf("unexpected payload type %s", payloadType)
	}

	res := payload.ReturnResults{}
	if err = res.Unmarshal(meta.Payload); err != nil {
		return errors.Wrap(err, "failed to unmarshal payload.ReturnResults")
	}

	ctx, span := instracer.StartSpan(ctx, "ContractRequester.ReceiveResult")
	defer span.End()

	if err = cr.result(ctx, &res); err != nil {
		return errors.Wrap(err, "[ ReceiveResult ]")
	}

	return nil
}

func (cr *ContractRequester) Stop() error {
	errIn := cr.inRequestResultsRouter.Close()

	return errors.Wrap(errIn, "Error while closing router")
}

func startRouter(ctx context.Context, router *message.Router) {
	go func() {
		if err := router.Run(ctx); err != nil {
			inslogger.FromContext(ctx).Error("Error while running router", err)
		}
	}()
	<-router.Running()
}

func deserializePayload(msg *message.Message) (insolar.Reply, error) {
	if msg == nil {
		return nil, errors.New("can't deserialize payload of nil message")
	}
	meta := payload.Meta{}
	err := meta.Unmarshal(msg.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "can't deserialize message payload")
	}

	if msg.Metadata.Get(busMeta.Type) == busMeta.TypeReply {
		rep, err := reply.Deserialize(bytes.NewBuffer(meta.Payload))
		if err != nil {
			return nil, errors.Wrap(err, "can't deserialize payload to reply")
		}
		return rep, nil
	}

	payloadType, err := payload.UnmarshalType(meta.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal payload type")
	}
	if payloadType != payload.TypeError {
		return nil, errors.Errorf("message bus receive unexpected payload type: %s", payloadType)
	}

	pl, err := payload.UnmarshalFromMeta(msg.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal error")
	}
	p, ok := pl.(*payload.Error)
	if !ok {
		return nil, errors.Errorf("unexpected error type %T", pl)
	}
	return nil, errors.New(p.Text)
}
