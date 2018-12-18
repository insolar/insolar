/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package contractrequester

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/pkg/errors"
)

// ContractRequester helps to call contracts
type ContractRequester struct {
	MessageBus   core.MessageBus   `inject:""`
	PulseManager core.PulseManager `inject:""`
	ResultMutex  sync.Mutex
	ResultMap    map[core.RecordRef]chan *message.ReturnResults
}

// New creates new ContractRequester
func New() (*ContractRequester, error) {
	return &ContractRequester{
		ResultMap: make(map[core.RecordRef]chan *message.ReturnResults),
	}, nil
}

func (cr *ContractRequester) Start(ctx context.Context) error {
	cr.MessageBus.MustRegister(core.TypeReturnResults, cr.ReceiveResult)
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

func (cr *ContractRequester) routeCall(ctx context.Context, ref core.RecordRef, method string, args core.Arguments) (core.Reply, error) {
	if cr.MessageBus == nil {
		return nil, errors.New("[ ContractRequester::routeCall ] message bus was not set during initialization")
	}

	e := &message.CallMethod{
		BaseLogicMessage: message.BaseLogicMessage{Nonce: randomUint64()},
		ObjectRef:        ref,
		Method:           method,
		Arguments:        args,
	}

	currentPulse, err := cr.PulseManager.Current(ctx)
	if err != nil {
		return nil, err
	}

	res, err := cr.MessageBus.Send(ctx, e, *currentPulse, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "[ ContractRequester::routeCall ] couldn't send message: %s", ref.String())
	}

	return res, nil
}

// SendRequest makes call to method of contract by its ref
func (cr *ContractRequester) SendRequest(ctx context.Context, ref *core.RecordRef, method string, argsIn []interface{}) (core.Reply, error) {
	args, err := core.MarshalArgs(argsIn...)
	if err != nil {
		return nil, errors.Wrap(err, "[ ContractRequester::SendRequest ] Can't marshal")
	}

	routResult, err := cr.routeCall(ctx, *ref, method, args)
	if err != nil {
		return nil, errors.Wrap(err, "[ ContractRequester::SendRequest ] Can't route call")
	}

	return routResult, nil
}

func (cr *ContractRequester) CallMethod(ctx context.Context, base core.Message, async bool, ref *core.RecordRef, method string, argsIn core.Arguments, mustPrototype *core.RecordRef) (core.Reply, error) {
	baseMessage, ok := base.(*message.BaseLogicMessage)
	if !ok {
		return nil, errors.New("Wrong type for BaseMessage")
	}
	log := inslogger.FromContext(ctx)

	mb := core.MessageBusFromContext(ctx, cr.MessageBus)
	if mb == nil {
		log.Debug("Context doesn't provide MessageBus")
		mb = cr.MessageBus
	}

	var mode message.MethodReturnMode
	if async {
		mode = message.ReturnNoWait
	} else {
		mode = message.ReturnResult
	}
	msg := &message.CallMethod{
		BaseLogicMessage: *baseMessage,
		ReturnMode:       mode,
		ObjectRef:        *ref,
		Method:           method,
		Arguments:        argsIn,
		ProxyPrototype:   *mustPrototype,
	}

	currentSlotPulse, err := cr.PulseManager.Current(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get pulse")
	}

	res, err := mb.Send(ctx, msg, *currentSlotPulse, nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't dispatch event")
	}

	r, ok := res.(*reply.CallMethod)
	if !ok {
		return nil, errors.New("Got not reply.CallMethod in reply for CallMethod")
	}

	if async {
		return res, nil

	}
	cr.ResultMutex.Lock()
	ch := make(chan *message.ReturnResults)
	cr.ResultMap[r.Request] = ch
	cr.ResultMutex.Unlock()
	inslogger.FromContext(ctx).Debug("Waiting for Method results ref=", r.Request)

	ret := <-ch

	return &reply.CallMethod{
		Request: r.Request,
		Result:  ret.Result,
	}, nil
}

func (cr *ContractRequester) CallConstructor(ctx context.Context, base core.Message, async bool,
	prototype *core.RecordRef, to *core.RecordRef, method string,
	argsIn core.Arguments, saveAs int) (*core.RecordRef, error) {
	baseMessage, ok := base.(*message.BaseLogicMessage)
	if !ok {
		return nil, errors.New("Wrong type for BaseMessage")
	}

	mb := core.MessageBusFromContext(ctx, cr.MessageBus)
	if mb == nil {
		return nil, errors.New("No access to message bus")
	}

	msg := &message.CallConstructor{
		BaseLogicMessage: *baseMessage,
		PrototypeRef:     *prototype,
		ParentRef:        *to,
		Name:             method,
		Arguments:        argsIn,
		SaveAs:           message.SaveAs(saveAs),
	}

	currentSlotPulse, err := cr.PulseManager.Current(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get pulse")
	}
	res, err := mb.Send(ctx, msg, *currentSlotPulse, nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't save new object as delegate")
	}

	r, ok := res.(*reply.CallConstructor)
	if !ok {
		return nil, errors.New("Got not reply.CallConstructor in reply for CallConstructor")
	}

	if async {
		return r.Object, nil
	}

	cr.ResultMutex.Lock()
	ch := make(chan *message.ReturnResults)
	cr.ResultMap[*r.Object] = ch
	cr.ResultMutex.Unlock()
	inslogger.FromContext(ctx).Debug("Waiting for constructor results ref=", *r.Object)

	<-ch

	return r.Object, nil
}

func (cr *ContractRequester) ReceiveResult(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	msg, ok := parcel.Message().(*message.ReturnResults)
	if !ok {
		return nil, errors.New("ReceiveResult() accepts only message.ReturnResults")
	}

	cr.ResultMutex.Lock()
	defer cr.ResultMutex.Unlock()

	c, ok := cr.ResultMap[msg.Request]
	if !ok {
		inslogger.FromContext(ctx).Info("oops unwaited results ref=", msg.Request)
		return &reply.OK{}, nil
	}
	inslogger.FromContext(ctx).Debug("Got wanted results")

	c <- msg
	delete(cr.ResultMap, msg.Request)
	return &reply.OK{}, nil
}
