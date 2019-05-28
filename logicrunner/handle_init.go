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
	"github.com/insolar/insolar/insolar/payload"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	insolarMsg "github.com/insolar/insolar/insolar/message"
)

const InnerMsgTopic = "InnerMsg"
const MessageTypeField = "Type"

const (
	processExecutionQueueMsg   = "ProcessExecutionQueue"
	getLedgerPendingRequestMsg = "GetLedgerPendingRequest"
)

type Dependencies struct {
	Publisher message.Publisher
	lr        *LogicRunner
	Bus       bus.Sender
}

type Init struct {
	dep *Dependencies

	Message *message.Message
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
	parcel, err := insolarMsg.DeserializeParcel(bytes.NewBuffer(meta.Payload))
	if err != nil {
		return errors.Wrap(err, "can't deserialize payload")
	}

	msgType := s.Message.Metadata.Get(bus.MetaType)

	switch msgType {
	case insolar.TypeCallMethod.String():
		h := &HandleCall{
			dep:     s.dep,
			Message: s.Message,
			Parcel:  parcel,
		}
		return f.Handle(ctx, h.Present)
	case insolar.TypePendingFinished.String():
		h := &HandlePendingFinished{
			dep:     s.dep,
			Message: s.Message,
			Parcel:  parcel,
		}
		return f.Handle(ctx, h.Present)
	case insolar.TypeStillExecuting.String():
		h := &HandleStillExecuting{
			dep:     s.dep,
			Message: s.Message,
			Parcel:  parcel,
		}
		return f.Handle(ctx, h.Present)
	case insolar.TypeAbandonedRequestsNotification.String():
		h := &HandleAbandonedRequestsNotification{
			dep:     s.dep,
			Message: s.Message,
			Parcel:  parcel,
		}
		return f.Handle(ctx, h.Present)
	case insolar.TypeExecutorResults.String():
		h := &HandleExecutorResults{
			dep:     s.dep,
			Message: s.Message,
			Parcel:  parcel,
		}
		return f.Handle(ctx, h.Present)
	default:
		return fmt.Errorf("[ Init.Present ] no handler for message type %s", msgType)
	}
}

type InnerInit struct {
	dep *Dependencies

	Message *message.Message
}

func (s *InnerInit) Present(ctx context.Context, f flow.Flow) error {
	switch s.Message.Metadata.Get(MessageTypeField) {
	case processExecutionQueueMsg:
		h := ProcessExecutionQueue{
			dep:     s.dep,
			Message: s.Message,
		}
		return f.Handle(ctx, h.Present)
	case getLedgerPendingRequestMsg:
		h := GetLedgerPendingRequest{
			dep:     s.dep,
			Message: s.Message,
		}
		return f.Handle(ctx, h.Present)
	default:
		return fmt.Errorf("[ InnerInit.Present ] no handler for message type %s", s.Message.Metadata.Get("Type"))
	}
}
