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
	"bytes"
	"context"
	"fmt"

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"

	"github.com/insolar/insolar/insolar"
	wmBus "github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/pkg/errors"
)

type Init struct {
	Dep *proc.Dependencies

	Message watermillMsg.Message
}

func (s *Init) Future(ctx context.Context, f flow.Flow) error {
	return f.Migrate(ctx, s.Present)
}

func (s *Init) Present(ctx context.Context, f flow.Flow) error {
	meta := payload.Meta{}
	err := meta.Unmarshal(s.Message.Payload)
	if err != nil {
		return errors.Wrap(err, "can't deserialize meta payload")
	}
	parcel, err := message.DeserializeParcel(bytes.NewBuffer(meta.Payload))
	if err != nil {
		return errors.Wrap(err, "can't deserialize payload")
	}

	msgType := s.Message.Metadata.Get(wmBus.MetaType)
	switch msgType {
	case insolar.TypeGetObject.String():
		fmt.Println("TypeGetObject gets inited")
		h := &GetObject{
			dep:     s.Dep,
			Message: &s.Message,
			Parcel:  parcel,
		}
		err = f.Handle(ctx, h.Present)
		fmt.Println("TypeGetObject gets error - ", err)
		return err
	case insolar.TypeSetRecord.String():
		msg := parcel.Message().(*message.SetRecord)
		h := NewSetRecord(s.Dep, &s.Message, msg)
		return f.Handle(ctx, h.Present)
	case insolar.TypeSetBlob.String():
		msg := parcel.Message().(*message.SetBlob)
		h := NewSetBlob(s.Dep, &s.Message, msg)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetCode.String():
		msg := parcel.Message().(*message.GetCode)
		h := NewGetCode(s.Dep, &s.Message, msg.Code)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetRequest.String():
		msg := parcel.Message().(*message.GetRequest)
		h := NewGetRequest(s.Dep, &s.Message, msg.Request)
		return f.Handle(ctx, h.Present)
	case insolar.TypeUpdateObject.String():
		msg := parcel.Message().(*message.UpdateObject)
		h := NewUpdateObject(s.Dep, &s.Message, msg)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetChildren.String():
		h := NewGetChildren(s.Dep, &s.Message, parcel)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetDelegate.String():
		h := NewGetDelegate(s.Dep, &s.Message, parcel)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetPendingRequests.String():
		h := NewGetPendingRequests(s.Dep, &s.Message, parcel)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetPendingRequestID.String():
		h := NewGetPendingRequestID(s.Dep, &s.Message, parcel)
		return f.Handle(ctx, h.Present)
	case insolar.TypeRegisterChild.String():
		msg := parcel.Message().(*message.RegisterChild)
		h := NewRegisterChild(s.Dep, &s.Message, msg, parcel.Pulse())
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetJet.String():
		msg := parcel.Message().(*message.GetJet)
		h := NewGetJet(s.Dep, &s.Message, msg)
		return f.Handle(ctx, h.Present)
	case insolar.TypeHotRecords.String():
		msg := parcel.Message().(*message.HotData)
		h := NewHotData(s.Dep, &s.Message, msg)
		return f.Handle(ctx, h.Present)
	default:
		return fmt.Errorf("no handler for message type %s", msgType)
	}
}

func (s *Init) Past(ctx context.Context, f flow.Flow) error {
	return f.Procedure(ctx, &proc.ReturnReply{
		Message: &s.Message,
		Err:     errors.New("no past handler"),
		Sender:  s.Dep.Sender,
	}, false)
}
