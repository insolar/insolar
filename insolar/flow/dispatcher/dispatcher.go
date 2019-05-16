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

package dispatcher

import (
	"context"
	"strconv"
	"sync/atomic"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	wmBus "github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/flow/internal/pulse"
	"github.com/insolar/insolar/insolar/flow/internal/thread"
)

type Dispatcher struct {
	handles struct {
		present flow.MakeHandle
		future  flow.MakeHandle
	}
	controller         *thread.Controller
	currentPulseNumber uint32
}

func NewDispatcher(present flow.MakeHandle, future flow.MakeHandle) *Dispatcher {
	d := &Dispatcher{
		controller: thread.NewController(),
	}
	d.handles.present = present
	d.handles.future = future
	d.currentPulseNumber = insolar.FirstPulseNumber
	return d
}

// ChangePulse is a handle for pulse change vent.
func (d *Dispatcher) ChangePulse(ctx context.Context, pulse insolar.Pulse) {
	d.controller.Pulse()
	atomic.StoreUint32(&d.currentPulseNumber, uint32(pulse.PulseNumber))
}

func (d *Dispatcher) getHandleByPulse(msgPulseNumber insolar.PulseNumber) flow.MakeHandle {
	if uint32(msgPulseNumber) > atomic.LoadUint32(&d.currentPulseNumber) {
		return d.handles.future
	}
	return d.handles.present
}

func (d *Dispatcher) WrapBusHandle(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg := bus.Message{
		ReplyTo: make(chan bus.Reply, 1),
		Parcel:  parcel,
	}

	ctx = pulse.ContextWith(ctx, parcel.Pulse())

	f := thread.NewThread(msg, d.controller)
	handle := d.getHandleByPulse(parcel.Pulse())

	err := f.Run(ctx, handle(msg))
	var rep bus.Reply
	select {
	case rep = <-msg.ReplyTo:
		return rep.Reply, rep.Err
	default:
	}

	if err != nil {
		return nil, err
	}

	return nil, errors.New("no reply from handler")
}

func (d *Dispatcher) InnerSubscriber(watermillMsg *message.Message) ([]*message.Message, error) {
	msg := bus.Message{
		WatermillMsg: watermillMsg,
	}

	ctx := context.Background()
	ctx = inslogger.ContextWithTrace(ctx, watermillMsg.Metadata.Get(wmBus.MetaTraceID))
	logger := inslogger.FromContext(ctx)
	go func() {
		f := thread.NewThread(msg, d.controller)
		err := f.Run(ctx, d.handles.present(msg))
		if err != nil {
			logger.Error("Handling failed", err)
		}
	}()
	return nil, nil
}

// Process handles incoming message.
func (d *Dispatcher) Process(msg *message.Message) ([]*message.Message, error) {
	ctx := msg.Context()
	msgBus := bus.Message{
		WatermillMsg: msg,
		ReplyTo:      make(chan bus.Reply),
	}
	p, err := pulseFromString(msg.Metadata.Get(wmBus.MetaPulse))
	if err != nil {
		return nil, errors.Wrap(err, "can't get pulse from string")
	}
	ctx, logger := inslogger.WithField(ctx, "pulse", msg.Metadata.Get(wmBus.MetaPulse))
	ctx = pulse.ContextWith(ctx, p)
	go func() {
		f := thread.NewThread(msgBus, d.controller)
		handle := d.getHandleByPulse(p)
		err := f.Run(ctx, handle(msgBus))
		if err != nil {
			logger.Error("Handling failed", err)
		}
	}()

	// TODO: move this logic to specific function and use it instead writing to ReplyTo channel
	// now its here for simplicity of moving only one message type (GetObject)
	rep := <-msgBus.ReplyTo
	var resInBytes []byte
	var replyType string
	if rep.Err != nil {
		resInBytes, err = wmBus.ErrorToBytes(rep.Err)
		if err != nil {
			return nil, errors.Wrap(err, "can't convert error to bytes")
		}
		replyType = wmBus.TypeError
	} else {
		resInBytes = reply.ToBytes(rep.Reply)
		replyType = string(rep.Reply.Type())
	}
	resAsMsg := message.NewMessage(watermill.NewUUID(), resInBytes)
	resAsMsg.Metadata.Set(wmBus.MetaType, replyType)
	receiver := msgBus.WatermillMsg.Metadata.Get(wmBus.MetaSender)
	resAsMsg.Metadata.Set(wmBus.MetaReceiver, receiver)
	return []*message.Message{resAsMsg}, nil

}

func pulseFromString(p string) (insolar.PulseNumber, error) {
	u64, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		return insolar.PulseNumber(0), errors.Wrap(err, "can't convert string value to pulse")
	}
	pInt := uint32(u64)
	return insolar.PulseNumber(pInt), nil
}
