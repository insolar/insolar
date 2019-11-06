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

package common

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/conveyor"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/dispatcher"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/network/consensus/adapters"
)

type conveyorDispatcher struct {
	conveyor *conveyor.PulseConveyor
}

var _ dispatcher.Dispatcher = &conveyorDispatcher{}

func (c *conveyorDispatcher) BeginPulse(ctx context.Context, pulseObject insolar.Pulse) {
	pulseRange := adapters.NewPulseData(pulseObject).AsRange()
	if err := c.conveyor.CommitPulseChange(pulseRange); err != nil {
		panic(err.Error())
	}
}

func (c *conveyorDispatcher) ClosePulse(ctx context.Context, pulseObject insolar.Pulse) {
}

type DispatcherMessage struct {
	MessageMeta message.Metadata
	PayloadMeta *payload.Meta
}

func (c *conveyorDispatcher) Process(msg *message.Message) error {
	pl, err := payload.Unmarshal(msg.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload.Meta")
	}
	meta, ok := pl.(*payload.Meta)
	if !ok {
		return errors.Errorf("unexpected type: %T (expected payload.Meta)", pl)
	}

	return c.conveyor.AddInput(context.Background(), meta.Pulse, &DispatcherMessage{
		MessageMeta: msg.Metadata,
		PayloadMeta: meta,
	})
}

func NewConveyorDispatcher(conveyor *conveyor.PulseConveyor) dispatcher.Dispatcher {
	return &conveyorDispatcher{conveyor: conveyor}
}
