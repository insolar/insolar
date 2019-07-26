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
	"bytes"
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	insolarMsg "github.com/insolar/insolar/insolar/message"
)

const InnerMsgTopic = "InnerMsg"

const (
	getLedgerPendingRequestMsg = "GetLedgerPendingRequest"
)

type Dependencies struct {
	Publisher      message.Publisher
	StateStorage   StateStorage
	ResultsMatcher ResultMatcher
	lr             *LogicRunner
	Sender         bus.Sender
}

type Init struct {
	dep *Dependencies

	Message *message.Message
}

func (s *Init) Future(ctx context.Context, f flow.Flow) error {
	return f.Migrate(ctx, s.Present)
}

func (s *Init) Present(ctx context.Context, f flow.Flow) error {
	msgType := s.Message.Metadata.Get(bus.MetaType)
	if msgType != "" {
		return s.handleParcel(ctx, f)
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
	default:
		return fmt.Errorf("[ Init.Present ] no handler for message type %s", msgType)
	}
}

func (s *Init) handleParcel(ctx context.Context, f flow.Flow) error {
	var err error

	meta := payload.Meta{}
	err = meta.Unmarshal(s.Message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal meta")
	}

	parcel, err := insolarMsg.DeserializeParcel(bytes.NewBuffer(meta.Payload))
	if err != nil {
		return errors.Wrap(err, "can't deserialize payload")
	}

	msgType := s.Message.Metadata.Get(bus.MetaType)

	switch msgType {
	case insolar.TypeCallMethod.String():
		h := &HandleCall{
			dep:     s.dep,
			Message: meta,
			Parcel:  parcel,
		}
		return f.Handle(ctx, h.Present)
	case insolar.TypeAdditionalCallFromPreviousExecutor.String():
		h := &HandleAdditionalCallFromPreviousExecutor{
			dep:     s.dep,
			Message: meta,
			Parcel:  parcel,
		}
		return f.Handle(ctx, h.Present)
	case insolar.TypePendingFinished.String():
		h := &HandlePendingFinished{
			dep:     s.dep,
			Message: meta,
			Parcel:  parcel,
		}
		return f.Handle(ctx, h.Present)
	case insolar.TypeStillExecuting.String():
		h := &HandleStillExecuting{
			dep:     s.dep,
			Message: meta,
			Parcel:  parcel,
		}
		return f.Handle(ctx, h.Present)
	case insolar.TypeExecutorResults.String():
		h := &HandleExecutorResults{
			dep:     s.dep,
			Message: meta,
			Parcel:  parcel,
		}
		return f.Handle(ctx, h.Present)
	default:
		return fmt.Errorf("[ Init.handleParcel ] no handler for message type %s", msgType)
	}
}

func (s *Init) Past(ctx context.Context, f flow.Flow) error {
	msgType := s.Message.Metadata.Get(bus.MetaType)

	if msgType == insolar.TypeCallMethod.String() {
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

type InnerInit struct {
	dep *Dependencies

	Message *message.Message
}

func (s *InnerInit) Present(ctx context.Context, f flow.Flow) error {
	switch s.Message.Metadata.Get(bus.MetaType) {
	default:
		return fmt.Errorf("[ InnerInit.Present ] no handler for message type %s", s.Message.Metadata.Get("Type"))
	}
}
