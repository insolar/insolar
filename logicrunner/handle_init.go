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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
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
}

type Init struct {
	dep *Dependencies

	Message bus.Message
}

func (s *Init) Present(ctx context.Context, f flow.Flow) error {
	switch s.Message.Parcel.Message().Type() {
	case insolar.TypeCallMethod, insolar.TypeCallConstructor:
		h := &HandleCall{
			dep:     s.dep,
			Message: s.Message,
		}
		return f.Handle(ctx, h.Present)
	case insolar.TypePendingFinished:
		h := &HandlePendingFinished{
			dep:     s.dep,
			Message: s.Message,
		}
		return f.Handle(ctx, h.Present)
	case insolar.TypeStillExecuting:
		h := &HandleStillExecuting{
			dep:     s.dep,
			Message: s.Message,
		}
		return f.Handle(ctx, h.Present)
	default:
		return fmt.Errorf("[ Init.Present ] no handler for message type %s", s.Message.Parcel.Message().Type().String())
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
