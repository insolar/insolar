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
	"github.com/insolar/insolar/insolar/bus/meta"
	"github.com/insolar/insolar/insolar/flow/dispatcher"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	logger "github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/pulse"
)

type dispatcherInitializationState int8

const (
	InitializationStarted dispatcherInitializationState = iota
	FirstPulseClosed
	InitializationDone // SecondPulseOpened
)

type conveyorDispatcher struct {
	conveyor      *conveyor.PulseConveyor
	state         dispatcherInitializationState
	previousPulse insolar.PulseNumber
}

var _ dispatcher.Dispatcher = &conveyorDispatcher{}

func (c *conveyorDispatcher) BeginPulse(ctx context.Context, pulseObject insolar.Pulse) {
	logger := inslogger.FromContext(ctx)
	var (
		pulseData  = adapters.NewPulseData(pulseObject)
		pulseRange pulse.Range
	)
	switch c.state {
	case InitializationDone:
		pulseRange = pulseData.AsRange()
	case FirstPulseClosed:
		pulseRange = pulse.NewLeftGapRange(c.previousPulse, 0, pulseData)
		c.state = InitializationDone
	case InitializationStarted:
		fallthrough
	default:
		panic("unreachable")
	}

	logger.Errorf("BeginPulse -> [%d, %d]", c.previousPulse, pulseData.PulseNumber)
	if err := c.conveyor.CommitPulseChange(pulseRange); err != nil {
		panic(err)
	}
}

func (c *conveyorDispatcher) ClosePulse(ctx context.Context, pulseObject insolar.Pulse) {
	logger.Errorf("ClosePulse -> [%d]", pulseObject.PulseNumber)
	c.previousPulse = pulseObject.PulseNumber

	switch c.state {
	case InitializationDone:
		// channel := make(conveyor.PreparePulseChangeChannel, 1)
		channel := conveyor.PreparePulseChangeChannel(nil)
		if err := c.conveyor.PreparePulseChange(channel); err != nil {
			panic(err)
		}
	case InitializationStarted:
		c.state = FirstPulseClosed
		return
	case FirstPulseClosed:
		fallthrough
	default:
		panic("unreachable")
	}
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
	plMeta, ok := pl.(*payload.Meta)
	if !ok {
		return errors.Errorf("unexpected type: %T (expected payload.Meta)", pl)
	}

	ctx, _ := inslogger.WithTraceField(context.Background(), msg.Metadata.Get(meta.TraceID))
	return c.conveyor.AddInput(ctx, plMeta.Pulse, &DispatcherMessage{
		MessageMeta: msg.Metadata,
		PayloadMeta: plMeta,
	})
}

func NewConveyorDispatcher(conveyor *conveyor.PulseConveyor) dispatcher.Dispatcher {
	return &conveyorDispatcher{conveyor: conveyor}
}
