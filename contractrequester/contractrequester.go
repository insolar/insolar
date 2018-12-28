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
	"fmt"
	"sync"

	"github.com/insolar/insolar/ledger/storage/record"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/core/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/pkg/errors"
)

// ContractRequester helps to call contracts
type ContractRequester struct {
	MessageBus                 core.MessageBus                 `inject:""`
	PulseStorage               core.PulseStorage               `inject:""`
	PlatformCryptographyScheme core.PlatformCryptographyScheme `inject:""`
	ResultMutex                sync.Mutex
	ResultMap                  map[core.RecordRef]chan *message.ReturnResults
}

// New creates new ContractRequester
func New() (*ContractRequester, error) {
	return &ContractRequester{
		ResultMap: make(map[core.RecordRef]chan *message.ReturnResults),
	}, nil
}

func (cr *ContractRequester) Start(ctx context.Context) error {
	cr.MessageBus.MustRegister(core.TypeReturnResults, cr.HandleReceiveResultMessage)
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
func (cr *ContractRequester) SendRequest(ctx context.Context, ref *core.RecordRef, method string, argsIn []interface{}) (core.Reply, error) {
	args, err := core.MarshalArgs(argsIn...)
	if err != nil {
		return nil, errors.Wrap(err, "[ ContractRequester::SendRequest ] Can't marshal")
	}

	bm := &message.BaseLogicMessage{
		Nonce: randomUint64(),
	}
	routResult, err := cr.CallMethod(ctx, bm, false, ref, method, args, nil)
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
	}
	if mustPrototype != nil {
		msg.ProxyPrototype = *mustPrototype
	}

	var ch chan *message.ReturnResults
	var reqRef core.RecordRef

	if !async {
		pu, _ := cr.PulseStorage.Current(ctx)
		reqRef = *record.NewRecordRefFromMessage(cr.PlatformCryptographyScheme, pu.PulseNumber, msg)

		cr.ResultMutex.Lock()
		ch = make(chan *message.ReturnResults, 1)
		cr.ResultMap[reqRef] = ch
		cr.ResultMutex.Unlock()
	}

	var res core.Reply
	var err error

	utils.MeasureExecutionTime(ctx, "ContractRequester.CallMethod mb.Send, msg.method="+msg.Method,
		func() {
			res, err = mb.Send(ctx, msg, nil)
		})
	if err != nil {
		return nil, errors.Wrap(err, "couldn't dispatch event")
	}

	r, ok := res.(*reply.RegisterRequest)
	if !ok {
		return nil, errors.New("Got not reply.RegisterRequest in reply for CallMethod")
	} else if !async && !r.Request.Equal(reqRef) {
		return nil, errors.Errorf("RecordId was predicted incorrectly we=%s their=%s", reqRef.String(), r.Request.String())
	}

	if async {
		return res, nil
	}

	inslogger.FromContext(ctx).Debug("Waiting for Method results ref=", r.Request)

	var result *reply.CallMethod
	err = nil

	utils.MeasureExecutionTime(ctx, "ContractRequester.CallMethod select",
		func() {
			select {
			case ret := <-ch:
				inslogger.FromContext(ctx).Debug("GOT Method results")
				if ret.Error != "" {
					err = errors.New(ret.Error)
					return
				}
				retReply, ok := ret.Reply.(*reply.CallMethod)
				if !ok {
					err = errors.New("Reply is not CallMethod")
					return

				}
				result = &reply.CallMethod{
					Request: r.Request,
					Result:  retReply.Result,
				}
			case <-ctx.Done():
				cr.ResultMutex.Lock()
				delete(cr.ResultMap, reqRef)
				cr.ResultMutex.Unlock()
				err = errors.New("canceled")
			}
		})

	utils.MeasureInfo(ctx, fmt.Sprintf("ContractRequester.CallMethod select returned result = %v, error = %v", result, err))
	return result, err
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

	var reqRef core.RecordRef
	var ch chan *message.ReturnResults
	if !async {
		pu, _ := cr.PulseStorage.Current(ctx)
		reqRef = *record.NewRecordRefFromMessage(cr.PlatformCryptographyScheme, pu.PulseNumber, msg)

		cr.ResultMutex.Lock()
		ch = make(chan *message.ReturnResults, 1)
		cr.ResultMap[reqRef] = ch
		cr.ResultMutex.Unlock()
	}

	res, err := mb.Send(ctx, msg, nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't save new object as delegate")
	}

	r, ok := res.(*reply.RegisterRequest)
	if !ok {
		return nil, errors.New("Got not reply.CallConstructor in reply for CallConstructor")
	} else if !async && !r.Request.Equal(reqRef) {
		return nil, errors.Errorf("RecordId was predicted incorrectly we=%s their=%s", reqRef.String(), r.Request.String())
	}

	if async {
		return &r.Request, nil
	}

	inslogger.FromContext(ctx).Debug("Waiting for constructor results req=", r.Request)

	select {
	case ret := <-ch:
		inslogger.FromContext(ctx).Debug("GOT Constructor results")
		if ret.Error != "" {
			return nil, errors.New(ret.Error)
		}
		return &r.Request, nil
	case <-ctx.Done():

		cr.ResultMutex.Lock()
		delete(cr.ResultMap, reqRef)
		cr.ResultMutex.Unlock()

		return nil, errors.New("canceled")
	}
}

func (cr *ContractRequester) HandleReceiveResultMessage(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	msg, ok := parcel.Message().(*message.ReturnResults)
	if !ok {
		return nil, errors.New("ReceiveResult() accepts only message.ReturnResults")
	}

	utils.MeasureExecutionTime(ctx, "ContractRequester.ReceiveResult",
		func() {
			cr.ResultMutex.Lock()
			defer cr.ResultMutex.Unlock()

			log := inslogger.FromContext(ctx)
			c, ok := cr.ResultMap[msg.Request]
			if !ok {
				log.Info("oops unwaited results req=", msg.Request)
				return
			}
			inslogger.FromContext(ctx).Debug("Got wanted results seq=", msg.Request)

			c <- msg
			delete(cr.ResultMap, msg.Request)
		})
	return &reply.OK{}, nil
}
