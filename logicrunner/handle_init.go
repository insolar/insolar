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

package logicrunner

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/writecontroller"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

const InnerMsgTopic = "InnerMsg"

type Dependencies struct {
	Publisher        message.Publisher
	StateStorage     StateStorage
	ResultsMatcher   ResultMatcher
	lr               *LogicRunner
	Sender           bus.Sender
	JetStorage       jet.Storage
	WriteAccessor    writecontroller.Accessor
	OutgoingSender   OutgoingRequestSender
	RequestsExecutor RequestsExecutor
}

type Init struct {
	dep *Dependencies

	Message *message.Message
}

func (s *Init) Future(ctx context.Context, f flow.Flow) error {
	return f.Migrate(ctx, s.Present)
}

func sendErrorMessage(ctx context.Context, sender bus.Sender, meta payload.Meta, err error) error {
	repMsg, err := payload.NewMessage(&payload.Error{Text: err.Error()})
	if err != nil {
		return err
	}

	go sender.Reply(ctx, meta, repMsg)
	return nil
}

func (s *Init) Present(ctx context.Context, f flow.Flow) error {
	msgType := s.Message.Metadata.Get(bus.MetaType)
	if msgType != "" {
		return fmt.Errorf("[ Init.handleParcel ] no handler for message type %s", s.Message.Metadata.Get(bus.MetaType))
	}

	var err error

	meta := payload.Meta{}
	err = meta.Unmarshal(s.Message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal meta")
	}
	payloadType, err := payload.UnmarshalType(meta.Payload)
	if err != nil {
		inslogger.FromContext(ctx).WithField("metaPayload", meta.Payload).Info("payload")
		return errors.Wrap(err, "failed to unmarshal payload type")
	}

	ctx, _ = inslogger.WithField(ctx, "msg_type", payloadType.String())

	ctx, span := instracer.StartSpan(ctx, "HandleCall.Present")
	span.AddAttributes(
		trace.StringAttribute("msg.Type", payloadType.String()),
	)
	defer span.End()

	switch payloadType {
	case payload.TypeSagaCallAcceptNotification:
		h := &HandleSagaCallAcceptNotification{
			dep:  s.dep,
			meta: meta,
		}
		return f.Handle(ctx, h.Present)
	case payload.TypeAbandonedRequestsNotification:
		h := &HandleAbandonedRequestsNotification{
			dep:  s.dep,
			meta: meta,
		}
		return f.Handle(ctx, h.Present)
	case payload.TypeUpdateJet:
		h := &HandleUpdateJet{
			dep:  s.dep,
			meta: meta,
		}
		return f.Handle(ctx, h.Present)
	case payload.TypePendingFinished:
		h := &HandlePendingFinished{
			dep:     s.dep,
			Message: meta,
		}
		return f.Handle(ctx, h.Present)
	case payload.TypeExecutorResults:
		h := &HandleExecutorResults{
			dep:     s.dep,
			Message: meta,
		}
		return f.Handle(ctx, h.Present)
	case payload.TypeStillExecuting:
		h := &HandleStillExecuting{
			dep:     s.dep,
			Message: meta,
		}
		return f.Handle(ctx, h.Present)
	case payload.TypeCallMethod:
		h := &HandleCall{
			dep:     s.dep,
			Message: meta,
		}
		return f.Handle(ctx, h.Present)
	case payload.TypeAdditionalCallFromPreviousExecutor:
		h := &HandleAdditionalCallFromPreviousExecutor{
			dep:     s.dep,
			Message: meta,
		}
		return f.Handle(ctx, h.Present)
	default:
		return fmt.Errorf("[ Init.Present ] no handler for message type %s", msgType)
	}
}

func (s *Init) Past(ctx context.Context, f flow.Flow) error {
	msgType := s.Message.Metadata.Get(bus.MetaType)
	if msgType != "" {
		return fmt.Errorf("[ Init.handleParcel ] no handler for message type %s", s.Message.Metadata.Get(bus.MetaType))
	}

	var err error

	meta := payload.Meta{}
	err = meta.Unmarshal(s.Message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal meta")
	}
	payloadType, err := payload.UnmarshalType(meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload type")
	}

	ctx, _ = inslogger.WithField(ctx, "msg_type", payloadType.String())

	if payloadType == payload.TypeCallMethod {
		meta := payload.Meta{}
		err := meta.Unmarshal(s.Message.Payload)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal meta")
		}

		errMsg, err := payload.NewMessage(&payload.Error{Text: "flow cancelled: Incorrect message pulse, get message from past on virtual node", Code: uint32(payload.CodeFlowCanceled)})
		if err != nil {
			inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to reply error"))
		}

		go s.dep.Sender.Reply(ctx, meta, errMsg)

		return nil
	}

	return s.Present(ctx, f)
}
