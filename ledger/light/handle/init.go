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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/pkg/errors"
)

type Init struct {
	Dep *proc.Dependencies

	Message bus.Message
}

func (s *Init) Future(ctx context.Context, f flow.Flow) error {
	return f.Migrate(ctx, s.Present)
}

func (s *Init) Present(ctx context.Context, f flow.Flow) error {
	switch s.Message.Parcel.Message().Type() {
	case insolar.TypeGetObject:
		h := &GetObject{
			dep:     s.Dep,
			Message: s.Message,
		}
		return f.Handle(ctx, h.Present)
	case insolar.TypeSetRecord:
		msg := s.Message.Parcel.Message().(*message.SetRecord)
		h := NewSetRecord(s.Dep, s.Message.ReplyTo, msg)
		return f.Handle(ctx, h.Present)
	case insolar.TypeSetBlob:
		msg := s.Message.Parcel.Message().(*message.SetBlob)
		h := NewSetBlob(s.Dep, s.Message.ReplyTo, msg)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetCode:
		msg := s.Message.Parcel.Message().(*message.GetCode)
		h := NewGetCode(s.Dep, s.Message.ReplyTo, msg.Code)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetRequest:
		msg := s.Message.Parcel.Message().(*message.GetRequest)
		h := NewGetRequest(s.Dep, s.Message.ReplyTo, msg.Request)
		return f.Handle(ctx, h.Present)
	case insolar.TypeUpdateObject:
		msg := s.Message.Parcel.Message().(*message.UpdateObject)
		h := NewUpdateObject(s.Dep, s.Message.ReplyTo, msg)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetPendingRequests:
		h := NewGetPendingRequests(s.Dep, s.Message.ReplyTo, s.Message.Parcel)
		return f.Handle(ctx, h.Present)
	case insolar.TypeRegisterChild:
		msg := s.Message.Parcel.Message().(*message.RegisterChild)
		h := NewRegisterChild(s.Dep, s.Message.ReplyTo, msg, s.Message.Parcel.Pulse())
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetJet:
		msg := s.Message.Parcel.Message().(*message.GetJet)
		h := NewGetJet(s.Dep, s.Message.ReplyTo, msg)
		return f.Handle(ctx, h.Present)
	case insolar.TypeHotRecords:
		msg := s.Message.Parcel.Message().(*message.HotData)
		h := NewHotData(s.Dep, s.Message.ReplyTo, msg)
		return f.Handle(ctx, h.Present)
	default:
		return fmt.Errorf("no handler for message type %s", s.Message.Parcel.Message().Type().String())
	}
}

func (s *Init) Past(ctx context.Context, f flow.Flow) error {
	return f.Procedure(ctx, &proc.ReturnReply{
		ReplyTo: s.Message.ReplyTo,
		Err:     errors.New("no past handler"),
	}, false)
}
