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

	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/messagebus"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/api"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

// ContractRequester helps to call contracts
type ContractRequester struct {
	MessageBus                 insolar.MessageBus                 `inject:""`
	PulseAccessor              pulse.Accessor                     `inject:""`
	JetCoordinator             jet.Coordinator                    `inject:""`
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`
	lr                         insolar.LogicRunner

	ResultMutex sync.Mutex
	ResultMap   map[[insolar.RecordHashSize]byte]chan *message.ReturnResults

	// callTimeout is mainly needed for unit tests which
	// sometimes may unpredictably fail on CI with a default timeout
	callTimeout time.Duration
}

// New creates new ContractRequester
func New(lr insolar.LogicRunner) (*ContractRequester, error) {
	return &ContractRequester{
		ResultMap:   make(map[[insolar.RecordHashSize]byte]chan *message.ReturnResults),
		callTimeout: 25 * time.Second,
		lr:          lr,
	}, nil
}

func (cr *ContractRequester) Start(ctx context.Context) error {
	cr.MessageBus.MustRegister(insolar.TypeReturnResults, cr.ReceiveResult)
	return nil
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

	msg := &message.CallMethod{
		IncomingRequest: record.IncomingRequest{
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

func (cr *ContractRequester) checkCall(_ context.Context, msg *message.CallMethod) error {
	switch {
	case msg.Caller.IsEmpty() && msg.APINode.IsEmpty():
		return errors.New("either Caller or APINode should be set, both empty")
	case !msg.Caller.IsEmpty() && !msg.APINode.IsEmpty():
		return errors.New("either Caller or APINode should be set, both set")
	}

	return nil
}

func (cr *ContractRequester) createResultWaiter(
	request record.IncomingRequest,
) (
	chan *message.ReturnResults, [insolar.RecordHashSize]byte, error,
) {
	cr.ResultMutex.Lock()
	defer cr.ResultMutex.Unlock()

	reqHash, err := cr.calcRequestHash(request)
	if err != nil {
		return nil, reqHash, errors.Wrap(err, "failed to calculate hash")
	}

	ch := make(chan *message.ReturnResults, 1)
	cr.ResultMap[reqHash] = ch

	return ch, reqHash, nil
}

func (cr *ContractRequester) Call(ctx context.Context, inMsg insolar.Message) (insolar.Reply, *insolar.Reference, error) {
	ctx, span := instracer.StartSpan(ctx, "ContractRequester.Call")
	defer span.End()

	msg := inMsg.(*message.CallMethod)

	async := msg.ReturnMode == record.ReturnNoWait

	if msg.Nonce == 0 {
		msg.Nonce = randomUint64()
	}

	logger := inslogger.FromContext(ctx)
	logger.Debug("about to send call to method ", msg.Method)

	err := cr.checkCall(ctx, msg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "incorrect request")
	}

	var ch chan *message.ReturnResults
	var reqHash [insolar.RecordHashSize]byte

	if !async {
		var err error
		ch, reqHash, err = cr.createResultWaiter(msg.IncomingRequest)
		if err != nil {
			return nil, nil, errors.Wrap(err, "can't create waiter record")
		}
	}

	sender := messagebus.BuildSender(
		cr.MessageBus.Send,
		messagebus.RetryIncorrectPulse(cr.PulseAccessor),
		messagebus.RetryFlowCancelled(cr.PulseAccessor),
	)

	res, err := sender(ctx, msg, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "couldn't dispatch event")
	}

	if _, earlyResult := res.(*reply.CallMethod); earlyResult {
		logger.Debug("early result for request, not registered")
		if !async {
			cr.ResultMutex.Lock()
			defer cr.ResultMutex.Unlock()

			delete(cr.ResultMap, reqHash)
			close(ch)
		}

		return res, nil, nil
	}

	r, ok := res.(*reply.RegisterRequest)
	if !ok {
		return nil, nil, errors.New("Got not reply.RegisterRequest in reply for CallMethod")
	}

	if async {
		return res, &r.Request, nil
	}

	if !bytes.Equal(r.Request.Record().Hash(), reqHash[:]) {
		return nil, &r.Request, errors.New("Registered request has different hash")
	}

	ctx, cancel := context.WithTimeout(ctx, cr.callTimeout)
	defer cancel()

	ctx, logger = inslogger.WithFields(
		ctx,
		map[string]interface{}{
			"called_request": r.Request.String(),
			"called_method":  msg.Method,
		},
	)
	logger.Debug("waiting results of request")

	select {
	case ret := <-ch:
		logger.Debug("got results of request")
		if ret.Error != "" {
			return nil, &r.Request, errors.Wrap(errors.New(ret.Error), "CallMethod returns error")
		}
		return ret.Reply, &r.Request, nil
	case <-ctx.Done():
		cr.ResultMutex.Lock()
		defer cr.ResultMutex.Unlock()

		delete(cr.ResultMap, reqHash)
		logger.Error("request timeout")
		return nil, &r.Request, errors.Errorf("request to contract was canceled: timeout of %s was exceeded", cr.callTimeout)
	}
}

func (cr *ContractRequester) result(ctx context.Context, msg *message.ReturnResults) error {
	cr.ResultMutex.Lock()
	defer cr.ResultMutex.Unlock()

	var reqHash [insolar.RecordHashSize]byte
	copy(reqHash[:], msg.RequestRef.Record().Hash())
	c, ok := cr.ResultMap[reqHash]
	if !ok {
		inslogger.FromContext(ctx).Info("unwaited results of request ", msg.RequestRef.String())
		if cr.lr != nil {
			return cr.lr.AddUnwantedResponse(ctx, msg)
		}
		inslogger.FromContext(ctx).Warn("drop unwanted ", msg.RequestRef.String())
		return nil
	}

	c <- msg
	delete(cr.ResultMap, reqHash)
	return nil
}

func (cr *ContractRequester) ReceiveResult(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg, ok := parcel.Message().(*message.ReturnResults)
	if !ok {
		return nil, errors.New("ReceiveResult() accepts only message.ReturnResults")
	}

	ctx, span := instracer.StartSpan(ctx, "ContractRequester.ReceiveResult")
	defer span.End()

	err := cr.result(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(err, "[ ReceiveResult ]")
	}

	return &reply.OK{}, nil
}
