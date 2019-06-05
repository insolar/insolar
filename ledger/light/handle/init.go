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

package handle

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	wmessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"

	"github.com/insolar/insolar/insolar"
	wbus "github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/pkg/errors"
)

type Init struct {
	dep     *proc.Dependencies
	message bus.Message
	sender  wbus.Sender
}

func NewInit(dep *proc.Dependencies, sender wbus.Sender, msg bus.Message) *Init {
	return &Init{
		dep:     dep,
		sender:  sender,
		message: msg,
	}
}

func (s *Init) Future(ctx context.Context, f flow.Flow) error {
	return f.Migrate(ctx, s.Present)
}

func (s *Init) Present(ctx context.Context, f flow.Flow) error {
	logger := inslogger.FromContext(ctx)
	err := s.handle(ctx, f)
	if err != nil && s.message.WatermillMsg != nil {
		errMsg, err := payload.NewMessage(&payload.Error{Text: err.Error()})
		if err != nil {
			logger.Error(errors.Wrap(err, "failed to reply error"))
			return err
		}
		go s.sender.Reply(ctx, s.message.WatermillMsg, errMsg)
	}
	return err
}

func (s *Init) handle(ctx context.Context, f flow.Flow) error {
	if s.message.WatermillMsg != nil {
		payloadType, err := payload.UnmarshalTypeFromMeta(s.message.WatermillMsg.Payload)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal payload type")
		}
		switch payloadType {
		case payload.TypeGetObject:
			h := NewGetObject(s.dep, s.message.WatermillMsg, false)
			return f.Handle(ctx, h.Present)
		case payload.TypePassState:
			h := NewPassState(s.dep, s.message.WatermillMsg)
			return f.Handle(ctx, h.Present)
		case payload.TypeGetCode:
			h := NewGetCode(s.dep, s.message.WatermillMsg, false)
			return f.Handle(ctx, h.Present)
		case payload.TypeSetCode:
			h := NewSetCode(s.dep, s.message.WatermillMsg, false)
			return f.Handle(ctx, h.Present)
		case payload.TypePass:
			return s.handlePass(ctx, f)
		case payload.TypeError:
			return f.Handle(ctx, NewError(s.message.WatermillMsg).Present)
		default:
			return fmt.Errorf("no handler for message type %s", payloadType.String())
		}
	}

	switch s.message.Parcel.Message().Type() {
	case insolar.TypeSetRecord:
		msg := s.message.Parcel.Message().(*message.SetRecord)
		h := NewSetRecord(s.dep, s.message.ReplyTo, msg)
		return f.Handle(ctx, h.Present)
	case insolar.TypeSetBlob:
		msg := s.message.Parcel.Message().(*message.SetBlob)
		h := NewSetBlob(s.dep, s.message.ReplyTo, msg)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetRequest:
		msg := s.message.Parcel.Message().(*message.GetRequest)
		h := NewGetRequest(s.dep, s.message.ReplyTo, msg.Request)
		return f.Handle(ctx, h.Present)
	case insolar.TypeUpdateObject:
		msg := s.message.Parcel.Message().(*message.UpdateObject)
		h := NewUpdateObject(s.dep, s.message.ReplyTo, msg)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetChildren:
		h := NewGetChildren(s.dep, s.message.ReplyTo, s.message)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetDelegate:
		h := NewGetDelegate(s.dep, s.message.ReplyTo, s.message.Parcel)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetPendingRequests:
		h := NewGetPendingRequests(s.dep, s.message.ReplyTo, s.message.Parcel)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetPendingRequestID:
		h := NewGetPendingRequestID(s.dep, s.message.ReplyTo, s.message.Parcel)
		return f.Handle(ctx, h.Present)
	case insolar.TypeRegisterChild:
		msg := s.message.Parcel.Message().(*message.RegisterChild)
		h := NewRegisterChild(s.dep, s.message.ReplyTo, msg, s.message.Parcel.Pulse())
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetJet:
		msg := s.message.Parcel.Message().(*message.GetJet)
		h := NewGetJet(s.dep, s.message.ReplyTo, msg)
		return f.Handle(ctx, h.Present)
	case insolar.TypeHotRecords:
		msg := s.message.Parcel.Message().(*message.HotData)
		h := NewHotData(s.dep, s.message.ReplyTo, msg)
		return f.Handle(ctx, h.Present)
	default:
		return fmt.Errorf("no handler for message type %s", s.message.Parcel.Message().Type().String())
	}
}

func (s *Init) handlePass(ctx context.Context, f flow.Flow) error {
	pl, err := payload.UnmarshalFromMeta(s.message.WatermillMsg.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal pass payload")
	}
	pass, ok := pl.(*payload.Pass)
	if !ok {
		return errors.New("wrong pass payload")
	}

	payloadType, err := payload.UnmarshalTypeFromMeta(pass.Origin)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload type")
	}
	origin := wmessage.NewMessage(watermill.NewUUID(), pass.Origin)
	middleware.SetCorrelationID(string(pass.CorrelationID), origin)

	switch payloadType {
	case payload.TypeGetObject:
		h := NewGetObject(s.dep, origin, true)
		return f.Handle(ctx, h.Present)
	case payload.TypeGetCode:
		h := NewGetCode(s.dep, origin, true)
		return f.Handle(ctx, h.Present)
	case payload.TypeSetCode:
		h := NewSetCode(s.dep, origin, true)
		return f.Handle(ctx, h.Present)
	}
	return nil
}

func (s *Init) Past(ctx context.Context, f flow.Flow) error {
	return f.Procedure(ctx, &proc.ReturnReply{
		ReplyTo: s.message.ReplyTo,
		Err:     errors.New("no past handler"),
	}, false)
}
