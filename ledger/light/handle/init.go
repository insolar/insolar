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

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/bus/meta"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/ledger/light/proc"
)

type Init struct {
	dep     *proc.Dependencies
	message *message.Message
	sender  bus.Sender
}

func NewInit(dep *proc.Dependencies, sender bus.Sender, msg *message.Message) *Init {
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
		if err == flow.ErrCancelled {
			logger.Info("Handling flow cancelled")
		} else {
			logger.Error("Handling error: ", err.Error())
		}
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

	logger := inslogger.FromContext(ctx)
	logger.Debug("Start to handle new message")

	switch payloadType {
	case payload.TypeGetObject:
		h := NewGetObject(s.dep, meta, false)
		err = f.Handle(ctx, h.Present)
	case payload.TypeGetRequest:
		h := NewGetRequest(s.dep, meta, false)
		err = f.Handle(ctx, h.Present)
	case payload.TypeGetRequestInfo:
		h := NewGetRequestInfo(s.dep, meta)
		err = f.Handle(ctx, h.Present)
	case payload.TypeGetFilament:
		h := NewGetFilament(s.dep, meta)
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
	case payload.TypeSetIncomingRequest:
		h := NewSetIncomingRequest(s.dep, meta, false)
		err = f.Handle(ctx, h.Present)
	case payload.TypeSetOutgoingRequest:
		h := NewSetOutgoingRequest(s.dep, meta, false)
		err = f.Handle(ctx, h.Present)
	case payload.TypeSetResult:
		h := NewSetResult(s.dep, meta, false)
		err = f.Handle(ctx, h.Present)
	case payload.TypeActivate:
		h := NewActivateObject(s.dep, meta, false)
		err = f.Handle(ctx, h.Present)
	case payload.TypeDeactivate:
		h := NewDeactivateObject(s.dep, meta, false)
		err = f.Handle(ctx, h.Present)
	case payload.TypeUpdate:
		h := NewUpdateObject(s.dep, meta, false)
		err = f.Handle(ctx, h.Present)
	case payload.TypeGetPendings:
		h := NewGetPendings(s.dep, meta, false)
		err = f.Handle(ctx, h.Present)
	case payload.TypeHasPendings:
		h := NewHasPendings(s.dep, meta, false)
		err = f.Handle(ctx, h.Present)
	case payload.TypeGetJet:
		h := NewGetJet(s.dep, meta, false)
		err = f.Handle(ctx, h.Present)
	case payload.TypePass:
		err = s.handlePass(ctx, f, meta)
	case payload.TypeError:
		err = f.Handle(ctx, NewError(s.message).Present)
	case payload.TypeHotObjects:
		err = f.Handle(ctx, NewHotObjects(s.dep, meta).Present)
	case payload.TypeGetPulse:
		err = f.Handle(ctx, NewGetPulse(s.dep, meta).Present)
	default:
		err = fmt.Errorf("no handler for message type %s", payloadType.String())
	}
	if err != nil {
		bus.ReplyError(ctx, s.sender, meta, err)
		s.errorMetrics(ctx, payloadType.String(), err)
	}
	return err
}

func (s *Init) errorMetrics(ctx context.Context, msgType string, err error) {
	if err == nil {
		return
	}
	errCode := payload.CodeUnknown
	if err == flow.ErrCancelled {
		errCode = payload.CodeFlowCanceled
	}
	cause := errors.Cause(err)
	insError, ok := cause.(*payload.CodedError)
	if ok {
		errCode = insError.GetCode()
	}

	ctx = insmetrics.InsertTag(ctx, KeyErrorCode, errCode.String())
	ctx = insmetrics.InsertTag(ctx, KeyMsgType, msgType)
	stats.Record(ctx, statHandlerError.M(1))
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
		bus.ReplyError(ctx, s.sender, originMeta, flow.ErrCancelled)
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
	case payload.TypeSetIncomingRequest:
		h := NewSetIncomingRequest(s.dep, originMeta, true)
		err = f.Handle(ctx, h.Present)
	case payload.TypeSetOutgoingRequest:
		h := NewSetOutgoingRequest(s.dep, originMeta, true)
		err = f.Handle(ctx, h.Present)
	case payload.TypeSetResult:
		h := NewSetResult(s.dep, originMeta, true)
		err = f.Handle(ctx, h.Present)
	case payload.TypeActivate:
		h := NewActivateObject(s.dep, originMeta, true)
		err = f.Handle(ctx, h.Present)
	case payload.TypeDeactivate:
		h := NewDeactivateObject(s.dep, originMeta, true)
		err = f.Handle(ctx, h.Present)
	case payload.TypeUpdate:
		h := NewUpdateObject(s.dep, originMeta, true)
		err = f.Handle(ctx, h.Present)
	case payload.TypeGetPendings:
		h := NewGetPendings(s.dep, originMeta, true)
		err = f.Handle(ctx, h.Present)
	case payload.TypeHasPendings:
		h := NewHasPendings(s.dep, originMeta, true)
		err = f.Handle(ctx, h.Present)
	case payload.TypeGetJet:
		h := NewGetJet(s.dep, originMeta, true)
		err = f.Handle(ctx, h.Present)
	case payload.TypeGetRequest:
		h := NewGetRequest(s.dep, originMeta, true)
		err = f.Handle(ctx, h.Present)
	default:
		err = fmt.Errorf("no handler for message type %s", payloadType.String())
	}
	if err != nil {
		bus.ReplyError(ctx, s.sender, originMeta, err)
		s.errorMetrics(ctx, payloadType.String(), err)
	}

	return err
}

func (s *Init) Past(ctx context.Context, f flow.Flow) error {
	msgType := s.message.Metadata.Get(meta.Type)
	if msgType != "" {
		return s.Present(ctx, f)
	}

	meta := payload.Meta{}
	err := meta.Unmarshal(s.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal meta")
	}

	payloadType, err := payload.UnmarshalType(meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload type")
	}

	if payloadType == payload.TypePass {
		pl, err := payload.Unmarshal(meta.Payload)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal pass payload")
		}
		pass, ok := pl.(*payload.Pass)
		if !ok {
			return fmt.Errorf("unexpected pass type %T", pl)
		}
		originMeta := payload.Meta{}
		err = originMeta.Unmarshal(pass.Origin)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal payload type")
		}

		pt, err := payload.UnmarshalType(originMeta.Payload)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal payload type")
		}
		payloadType = pt
		meta = originMeta
	}

	// Only allow read operations in the past.
	switch payloadType {
	case
		payload.TypeGetObject,
		payload.TypeGetCode,
		payload.TypeGetPendings,
		payload.TypeHasPendings,
		payload.TypeGetJet,
		payload.TypeGetRequest,
		payload.TypePassState,
		payload.TypeGetRequestInfo,
		payload.TypeGetFilament,
		payload.TypeGetPulse:
		return s.Present(ctx, f)
	}

	bus.ReplyError(ctx, s.sender, meta, flow.ErrCancelled)
	return nil
}
