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
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/api"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/bus/meta"
	busMeta "github.com/insolar/insolar/insolar/bus/meta"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/dispatcher"
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
	Sender                     bus.Sender
	PulseAccessor              pulse.Accessor
	JetCoordinator             jet.Coordinator
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme

	// TODO: remove this hack in INS-3341
	// we need ResultMatcher, not Logicrunner
	LR insolar.LogicRunner

	FlowDispatcher dispatcher.Dispatcher

	ResultMutex sync.Mutex
	ResultMap   map[[insolar.RecordHashSize]byte]chan *payload.ReturnResults

	// callTimeout is mainly needed for unit tests which
	// sometimes may unpredictably fail on CI with a default timeout
	callTimeout time.Duration
}

// New creates new ContractRequester
func New(
	sender bus.Sender,
	pulses pulse.Accessor,
	jetCoordinator jet.Coordinator,
	pcs insolar.PlatformCryptographyScheme,
) (*ContractRequester, error) {
	cr := &ContractRequester{
		ResultMap:   make(map[[insolar.RecordHashSize]byte]chan *payload.ReturnResults),
		callTimeout: 25 * time.Second,

		Sender:                     sender,
		PulseAccessor:              pulses,
		JetCoordinator:             jetCoordinator,
		PlatformCryptographyScheme: pcs,
	}

	handle := func(msg *message.Message) *handleResults {
		return &handleResults{
			cr:      cr,
			Message: msg,
		}
	}

	cr.FlowDispatcher = dispatcher.NewDispatcher(cr.PulseAccessor,
		func(msg *message.Message) flow.Handle {
			return handle(msg).Present
		}, func(msg *message.Message) flow.Handle {
			return handle(msg).Future
		}, func(msg *message.Message) flow.Handle {
			return handle(msg).Past
		})
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
	if msg.PulseNumber == 0 {
		pulseObject, err := cr.PulseAccessor.Latest(ctx)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to get latest pulse")
		}
		msg.PulseNumber = pulseObject.PulseNumber
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

	sender := bus.NewRetrySender(cr.Sender, cr.PulseAccessor, 1, 1)

	messagePayload, err := payload.NewMessage(msg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to marshal payload")
	}

	target := record.CalculateRequestAffinityRef(msg.Request, msg.PulseNumber, cr.PlatformCryptographyScheme)

	resp, done := sender.SendRole(ctx, messagePayload, insolar.DynamicRoleVirtualExecutor, *target)
	defer done()
	rawResponse, ok := <-resp
	if !ok {
		return nil, nil, errors.New("no reply")
	}

	if rawResponse.Metadata.Get(meta.Type) != meta.TypeReply {
		data, err := payload.UnmarshalFromMeta(rawResponse.Payload)
		if err != nil {
			return nil, nil, errors.Wrap(err, "bad reply")
		}

		responseErr, isError := data.(*payload.Error)
		if !isError {
			return nil, nil, errors.Errorf("not a reply in reply, message data is %T", data)
		}

		return nil, nil, errors.Wrap(errors.New(responseErr.Text), "got reply with error")
	}

	replyData, err := reply.UnmarshalFromMeta(rawResponse.Payload)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to unmarshal reply")
	}

	switch replyTyped := replyData.(type) {
	// early result
	case *reply.CallMethod:
		return cr.handleCallMethod(ctx, replyTyped, reqHash, ch, async)
	case *reply.RegisterRequest:
		ctx, _ = inslogger.WithFields(ctx,
			map[string]interface{}{
				"called_request": replyTyped.Request.String(),
				"called_method":  msg.Request.Method,
			},
		)
		return cr.handleRegisterResult(ctx, replyTyped, reqHash, ch, async)
	default:
		return nil, nil, errors.Errorf("Got not reply.RegisterRequest in reply for CallMethod %T", replyData)
	}
}

func (cr *ContractRequester) handleCallMethod(ctx context.Context, r *reply.CallMethod,
	reqHash [insolar.RecordHashSize]byte, ch chan *payload.ReturnResults, async bool) (insolar.Reply, *insolar.Reference, error) {
	inslogger.FromContext(ctx).Debug("early result for request, not registered")
	if !async {
		cr.ResultMutex.Lock()
		defer cr.ResultMutex.Unlock()

		delete(cr.ResultMap, reqHash)
		close(ch)
	}

	return r, nil, nil
}

func (cr *ContractRequester) handleRegisterResult(ctx context.Context, r *reply.RegisterRequest,
	reqHash [insolar.RecordHashSize]byte, ch chan *payload.ReturnResults, async bool) (insolar.Reply, *insolar.Reference, error) {
	logger := inslogger.FromContext(ctx)

	if async {
		return r, &r.Request, nil
	}

	if !bytes.Equal(r.Request.GetLocal().Hash(), reqHash[:]) {
		return nil, &r.Request, errors.New("Registered request has different hash")
	}

	ctx, cancel := context.WithTimeout(ctx, cr.callTimeout)
	defer cancel()

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
	copy(reqHash[:], msg.RequestRef.GetLocal().Hash())
	c, ok := cr.ResultMap[reqHash]
	if !ok {
		inslogger.FromContext(ctx).Info("unwanted results of request ", msg.RequestRef.String())
		if cr.LR != nil {
			return cr.LR.AddUnwantedResponse(ctx, msg)
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
		log.Error("can't deserialize payload of nil message")
		return nil
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

	payloadMeta := &payload.Meta{}
	err = payloadMeta.Unmarshal(msg.Payload)
	if err != nil {
		return err
	}

	err = cr.handleMessage(ctx, payloadMeta)
	if err != nil {
		bus.ReplyError(ctx, cr.Sender, *payloadMeta, err)
		return nil
	}
	cr.Sender.Reply(ctx, *payloadMeta, bus.ReplyAsMessage(ctx, &reply.OK{}))
	return nil
}

func (cr *ContractRequester) handleMessage(ctx context.Context, payloadMeta *payload.Meta) error {
	payloadType, err := payload.UnmarshalType(payloadMeta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload type")
	}

	ctx, logger := inslogger.WithField(ctx, "msg_type", payloadType.String())

	logger.Debug("Start to handle new message")

	if payloadType != payload.TypeReturnResults {
		return errors.Errorf("unexpected payload type %s", payloadType)
	}

	res := payload.ReturnResults{}
	err = res.Unmarshal(payloadMeta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload.ReturnResults")
	}

	ctx, span := instracer.StartSpan(ctx, "ContractRequester.ReceiveResult")
	defer span.End()

	err = cr.result(ctx, &res)
	if err != nil {
		return errors.Wrap(err, "[ ReceiveResult ]")
	}

	return nil
}
