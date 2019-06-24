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

	wmessage "github.com/ThreeDotsLabs/watermill/message"

	"github.com/insolar/insolar/insolar"
	wbus "github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/pkg/errors"
)

type Init struct {
	dep     *proc.Dependencies
	message *wmessage.Message
	meta    payload.Meta
	sender  wbus.Sender
}

func NewInit(dep *proc.Dependencies, sender wbus.Sender, msg *wmessage.Message) *Init {
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
	if err != nil {
		logger.Error(errors.Wrap(err, "handling error"))
	}
	return err
}

func (s *Init) handle(ctx context.Context, f flow.Flow) error {
	var err error

	meta := payload.Meta{}
	err = meta.Unmarshal(s.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal meta")
	}
	payloadType, err := payload.UnmarshalType(meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload type")
	}

	ctx, _ = inslogger.WithField(ctx, "msg_type", payloadType.String())

	s.meta = meta

	if err == nil {
		switch payloadType {
		case payload.TypeGetObject:
			h := NewGetObject(s.dep, meta, false)
			err = f.Handle(ctx, h.Present)
		case payload.TypePassState:
			h := NewPassState(s.dep, meta)
			err = f.Handle(ctx, h.Present)
		case payload.TypeGetCode:
			h := NewGetCode(s.dep, meta, false)
			err = f.Handle(ctx, h.Present)
		case payload.TypeSetCode:
			h := NewSetCode(s.dep, meta, false)
			err = f.Handle(ctx, h.Present)
		case payload.TypePass:
			err = s.handlePass(ctx, f, meta)
		case payload.TypeError:
			err = f.Handle(ctx, NewError(s.message).Present)
		default:
			err = fmt.Errorf("no handler for message type %s", payloadType.String())
		}
		if err != nil {
			s.replyError(ctx, meta, err)
		}
		return err
	}

	parcel, err := message.DeserializeParcel(bytes.NewBuffer(meta.Payload))
	if err != nil {
		return errors.Wrap(err, "can't deserialize payload")
	}

	ctx, span := instracer.StartSpan(ctx, fmt.Sprintf("Present %v", parcel.Message().Type().String()))
	defer span.End()

	msgType := s.message.Metadata.Get(wbus.MetaType)
	switch msgType {
	case insolar.TypeSetRecord.String():
		msg := parcel.Message().(*message.SetRecord)
		h := NewSetRecord(s.dep, meta, msg)
		return f.Handle(ctx, h.Present)
	case insolar.TypeSetBlob.String():
		msg := parcel.Message().(*message.SetBlob)
		h := NewSetBlob(s.dep, meta, msg)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetRequest.String():
		msg := parcel.Message().(*message.GetRequest)
		h := NewGetRequest(s.dep, meta, msg.Request)
		return f.Handle(ctx, h.Present)
	case insolar.TypeUpdateObject.String():
		msg := parcel.Message().(*message.UpdateObject)
		h := NewUpdateObject(s.dep, meta, msg)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetChildren.String():
		h := NewGetChildren(s.dep, meta, parcel)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetDelegate.String():
		h := NewGetDelegate(s.dep, meta, parcel)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetPendingRequests.String():
		h := NewGetPendingRequests(s.dep, meta, parcel)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetPendingRequestID.String():
		h := NewGetPendingRequestID(s.dep, meta, parcel)
		return f.Handle(ctx, h.Present)
	case insolar.TypeRegisterChild.String():
		msg := parcel.Message().(*message.RegisterChild)
		h := NewRegisterChild(s.dep, meta, msg, parcel.Pulse())
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetJet.String():
		msg := parcel.Message().(*message.GetJet)
		h := NewGetJet(s.dep, meta, msg)
		return f.Handle(ctx, h.Present)
	case insolar.TypeHotRecords.String():
		msg := parcel.Message().(*message.HotData)
		h := NewHotData(s.dep, meta, msg)
		return f.Handle(ctx, h.Present)
	case insolar.TypeGetPendingFilament.String():
		msg := parcel.Message().(*message.GetPendingFilament)
		h := NewGetPendingFilament(s.dep, meta, msg)
		return f.Handle(ctx, h.Present)
	default:
		return fmt.Errorf("no handler for message type %s", msgType)
	}
}

func (s *Init) handlePass(ctx context.Context, f flow.Flow, meta payload.Meta) error {
	var err error
	pl, err := payload.Unmarshal(meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal pass payload")
	}
	pass, ok := pl.(*payload.Pass)
	if !ok {
		return errors.New("wrong pass payload")
	}

	originMeta := payload.Meta{}
	err = originMeta.Unmarshal(pass.Origin)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload type")
	}
	payloadType, err := payload.UnmarshalType(originMeta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload type")
	}

	ctx, _ = inslogger.WithField(ctx, "msg_type_original", payloadType.String())

	if originMeta.Pulse != meta.Pulse {
		s.replyError(ctx, originMeta, flow.ErrCancelled)
		return flow.ErrCancelled
	}

	switch payloadType {
	case payload.TypeGetObject:
		h := NewGetObject(s.dep, originMeta, true)
		err = f.Handle(ctx, h.Present)
	case payload.TypeGetCode:
		h := NewGetCode(s.dep, originMeta, true)
		err = f.Handle(ctx, h.Present)
	case payload.TypeSetCode:
		h := NewSetCode(s.dep, originMeta, true)
		err = f.Handle(ctx, h.Present)
	default:
		err = fmt.Errorf("no handler for message type %s", payloadType.String())
	}
	if err != nil {
		s.replyError(ctx, originMeta, err)
	}

	return err
}

func (s *Init) Past(ctx context.Context, f flow.Flow) error {
	return f.Procedure(ctx, &proc.ReturnReply{
		Message: s.meta,
		Err:     errors.New("no past handler"),
		Sender:  s.dep.Sender,
	}, false)
}

func (s *Init) replyError(ctx context.Context, replyTo payload.Meta, err error) {
	errCode := payload.CodeUnknown
	if err == flow.ErrCancelled {
		errCode = payload.CodeFlowCanceled
	}
	errMsg, err := payload.NewMessage(&payload.Error{Text: err.Error(), Code: uint32(errCode)})
	if err != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to reply error"))
	}
	go s.sender.Reply(ctx, replyTo, errMsg)
}
