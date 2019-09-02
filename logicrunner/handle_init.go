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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/bus/meta"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/writecontroller"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

const InnerMsgTopic = "InnerMsg"

type Dependencies struct {
	ArtifactManager  artifacts.Client
	Publisher        message.Publisher
	StateStorage     StateStorage
	ResultsMatcher   ResultMatcher
	Sender           bus.Sender
	JetStorage       jet.Storage
	JetCoordinator   jet.Coordinator
	WriteAccessor    writecontroller.Accessor
	OutgoingSender   OutgoingRequestSender
	RequestsExecutor RequestsExecutor
	PulseAccessor    pulse.Accessor
}

type Init struct {
	dep *Dependencies

	Message *message.Message
}

func (s *Init) Future(ctx context.Context, f flow.Flow) error {
	return f.Migrate(ctx, s.Present)
}

func (s *Init) replyError(ctx context.Context, meta payload.Meta, err error) {
	errCode := uint32(payload.CodeUnknown)

	// Throwing custom error code
	cause := errors.Cause(err)
	insError, ok := cause.(*payload.CodedError)
	if ok {
		errCode = insError.GetCode()
	}

	// todo refactor this #INS-3191
	if cause == flow.ErrCancelled {
		errCode = uint32(payload.CodeFlowCanceled)
	}
	errMsg, newErr := payload.NewMessage(&payload.Error{Text: err.Error(), Code: errCode})
	if newErr != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to reply error"))
	}
	s.dep.Sender.Reply(ctx, meta, errMsg)
}

func (s *Init) Present(ctx context.Context, f flow.Flow) error {
	msgType := s.Message.Metadata.Get(meta.Type)
	if msgType != "" {
		return fmt.Errorf("[ Init.handleParcel ] no handler for message type %s", s.Message.Metadata.Get(meta.Type))
	}

	var err error

	originMeta := payload.Meta{}
	err = originMeta.Unmarshal(s.Message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal meta")
	}
	payloadType, err := payload.UnmarshalType(originMeta.Payload)
	if err != nil {
		inslogger.FromContext(ctx).WithField("metaPayload", originMeta.Payload).Info("payload")
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
			meta: originMeta,
		}
		err = f.Handle(ctx, h.Present)
	case payload.TypeAbandonedRequestsNotification:
		h := &HandleAbandonedRequestsNotification{
			dep:  s.dep,
			meta: originMeta,
		}
		err = f.Handle(ctx, h.Present)
	case payload.TypeUpdateJet:
		h := &HandleUpdateJet{
			dep:  s.dep,
			meta: originMeta,
		}
		err = f.Handle(ctx, h.Present)
	case payload.TypePendingFinished:
		h := &HandlePendingFinished{
			dep:     s.dep,
			Message: originMeta,
		}
		err = f.Handle(ctx, h.Present)
	case payload.TypeExecutorResults:
		h := &HandleExecutorResults{
			dep:  s.dep,
			meta: originMeta,
		}
		err = f.Handle(ctx, h.Present)
	case payload.TypeStillExecuting:
		h := &HandleStillExecuting{
			dep:     s.dep,
			Message: originMeta,
		}
		err = f.Handle(ctx, h.Present)
	case payload.TypeCallMethod:
		h := &HandleCall{
			dep:     s.dep,
			Message: originMeta,
		}
		err = f.Handle(ctx, h.Present)
	case payload.TypeAdditionalCallFromPreviousExecutor:
		h := &HandleAdditionalCallFromPreviousExecutor{
			dep:     s.dep,
			Message: originMeta,
		}
		err = f.Handle(ctx, h.Present)
	default:
		err = fmt.Errorf("[ Init.Present ] no handler for message type %s", msgType)
	}
	if err != nil {
		bus.ReplyError(ctx, s.dep.Sender, originMeta, err)
	}
	return err
}

func (s *Init) Past(ctx context.Context, f flow.Flow) error {
	msgType := s.Message.Metadata.Get(meta.Type)
	if msgType != "" {
		return fmt.Errorf("[ Init.handleParcel ] no handler for message type %s", s.Message.Metadata.Get(meta.Type))
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
