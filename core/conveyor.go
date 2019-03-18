/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package core

import (
	"time"

	"github.com/insolar/insolar/conveyor/queue"
)

// EventSink allows to push events to conveyor
type EventSink interface {
	// SinkPush adds event to conveyor
	SinkPush(pulseNumber PulseNumber, data interface{}) error
	// SinkPushAll adds several events to conveyor
	SinkPushAll(pulseNumber PulseNumber, data []interface{}) error
}

// ConveyorState is the states of conveyor
type ConveyorState int

//go:generate stringer -type=ConveyorState
const (
	ConveyorActive = ConveyorState(iota)
	ConveyorPreparingPulse
	ConveyorShuttingDown
	ConveyorInactive
)

// Control allows to control conveyor and pulse
type Control interface {
	// PreparePulse is preparing conveyor for working with provided pulse
	PreparePulse(pulse Pulse, callback queue.SyncDone) error
	// ActivatePulse is activate conveyor with prepared pulse
	ActivatePulse() error
	// GetState returns current state of conveyor
	GetState() ConveyorState
	// IsOperational shows if conveyor is ready for work
	IsOperational() bool
	// InitiateShutdown shutting conveyor down and cancels tasks in adapters if force param set
	InitiateShutdown(force bool)
}

// Conveyor is responsible for all pulse-dependent processing logic
//go:generate minimock -i github.com/insolar/insolar/core.Conveyor -o ../testutils -s _mock.go
type Conveyor interface {
	EventSink
	Control
}

// ConveyorFuture is pending for response from conveyor
type ConveyorFuture interface {
	// ID returns number.
	ID() uint64
	// Result is a channel to listen for future result.
	Result() <-chan Reply
	// SetResult makes packet to appear in result channel.
	SetResult(res Reply)
	// GetResult gets the future result from Result() channel with a timeout set to `duration`.
	GetResult(duration time.Duration) (Reply, error)
	// Cancel closes all channels and cleans up underlying structures.
	Cancel()
}
