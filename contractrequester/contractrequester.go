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
	"github.com/insolar/insolar/platformpolicy"
)

// ContractRequester helps to call contracts
type ContractRequester struct {
	MessageBus     insolar.MessageBus `inject:""`
	ResultMutex    sync.Mutex
	ResultMap      map[[insolar.RecordHashSize]byte]chan *message.ReturnResults
	PulseAccessor  pulse.Accessor  `inject:""`
	JetCoordinator jet.Coordinator `inject:""`
	// callTimeout is mainly needed for unit tests which
	// sometimes may unpredictably fail on CI with a default timeout
	callTimeout time.Duration
}

// New creates new ContractRequester
func New() (*ContractRequester, error) {
	return &ContractRequester{
		ResultMap:   make(map[[insolar.RecordHashSize]byte]chan *message.ReturnResults),
		callTimeout: 25 * time.Second,
	}, nil
}

func (cr *ContractRequester) SetCallTimeout(timeout time.Duration) {
	cr.callTimeout = timeout
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
	return cr.SendRequestWithPulse(ctx, ref, method, argsIn, pulse.PulseNumber)
}

func (cr *ContractRequester) SendRequestWithPulse(ctx context.Context, ref *insolar.Reference, method string, argsIn []interface{}, pulse insolar.PulseNumber) (insolar.Reply, error) {
	ctx, span := instracer.StartSpan(ctx, "SendRequest "+method)
	defer span.End()

	args, err := insolar.MarshalArgs(argsIn...)
	if err != nil {
		return nil, errors.Wrap(err, "[ ContractRequester::SendRequest ] Can't marshal")
	}

	msg := &message.CallMethod{
		IncomingRequest: record.IncomingRequest{
			Object:       ref,
			Method:       method,
			Arguments:    args,
			APIRequestID: utils.TraceID(ctx),
			Reason:       api.MakeReason(pulse, args),
		},
	}

	routResult, err := cr.CallMethod(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(err, "[ ContractRequester::SendRequest ] Can't route call")
	}

	return routResult, nil
}

func (cr *ContractRequester) calcRequestHash(request record.IncomingRequest, hash *[insolar.RecordHashSize]byte) error {
	virtRec := record.Wrap(request)
	buf, err := virtRec.Marshal()
	if err != nil {
		return errors.Wrap(err, "[ ContractRequester::calcRequestHash ] Failed to marshal record")
	}

	hasher := platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher()
	copy(hash[:], hasher.Hash(buf)[0:insolar.RecordHashSize])
	return nil
}

func (cr *ContractRequester) Call(ctx context.Context, inMsg insolar.Message) (insolar.Reply, error) {
	ctx, span := instracer.StartSpan(ctx, "ContractRequester.Call")
	defer span.End()

	msg := inMsg.(*message.CallMethod)

	async := msg.ReturnMode == record.ReturnNoWait

	if msg.Nonce == 0 {
		msg.Nonce = randomUint64()
	}
	msg.Sender = cr.JetCoordinator.Me()

	var ch chan *message.ReturnResults
	var reqHash [insolar.RecordHashSize]byte

	if !async {
		cr.ResultMutex.Lock()
		err := cr.calcRequestHash(msg.IncomingRequest, &reqHash)
		if err != nil {
			return nil, errors.Wrap(err, "[ ContractRequester::Call ] Failed to calculate hash")
		}
		ch = make(chan *message.ReturnResults, 1)
		cr.ResultMap[reqHash] = ch

		cr.ResultMutex.Unlock()
	}

	sender := messagebus.BuildSender(
		cr.MessageBus.Send,
		messagebus.RetryIncorrectPulse(cr.PulseAccessor),
		messagebus.RetryFlowCancelled(cr.PulseAccessor),
	)
	res, err := sender(ctx, msg, nil)

	if err != nil {
		return nil, errors.Wrap(err, "couldn't dispatch event")
	}

	r, ok := res.(*reply.RegisterRequest)
	if !ok {
		return nil, errors.New("Got not reply.RegisterRequest in reply for CallMethod")
	}

	if async {
		return res, nil
	}

	ctx, cancel := context.WithTimeout(ctx, cr.callTimeout)
	defer cancel()

	ctx, _ = inslogger.WithField(ctx, "request", r.Request.String())
	ctx, logger := inslogger.WithField(ctx, "method", msg.Method)

	logger.Debug("Waiting results of request")

	select {
	case ret := <-ch:
		logger.Debug("Got results of request")
		if ret.Error != "" {
			return nil, errors.Wrap(errors.New(ret.Error), "CallMethod returns error")
		}
		return ret.Reply, nil
	case <-ctx.Done():
		cr.ResultMutex.Lock()
		delete(cr.ResultMap, reqHash)
		cr.ResultMutex.Unlock()
		return nil, errors.Errorf("request to contract was canceled: timeout of %s was exceeded", cr.callTimeout)
	}
}

func (cr *ContractRequester) CallMethod(ctx context.Context, inMsg insolar.Message) (insolar.Reply, error) {
	return cr.Call(ctx, inMsg)
}

func (cr *ContractRequester) CallConstructor(ctx context.Context, inMsg insolar.Message) (*insolar.Reference, error) {
	res, err := cr.Call(ctx, inMsg)
	if err != nil {
		return nil, err
	}

	rep, ok := res.(*reply.CallConstructor)
	if !ok {
		return nil, errors.New("Reply is not CallConstructor")
	}
	return rep.Object, nil
}

func (cr *ContractRequester) ReceiveResult(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg, ok := parcel.Message().(*message.ReturnResults)
	if !ok {
		return nil, errors.New("ReceiveResult() accepts only message.ReturnResults")
	}

	ctx, span := instracer.StartSpan(ctx, "ContractRequester.ReceiveResult")
	defer span.End()

	cr.ResultMutex.Lock()
	defer cr.ResultMutex.Unlock()

	var reqHash [insolar.RecordHashSize]byte
	copy(reqHash[:], msg.RequestRef.Record().Hash())
	c, ok := cr.ResultMap[reqHash]
	if !ok {
		inslogger.FromContext(ctx).Warn("unwaited results of request ", msg.RequestRef.String())
		return &reply.OK{}, nil
	}

	c <- msg
	delete(cr.ResultMap, reqHash)

	return &reply.OK{}, nil
}
