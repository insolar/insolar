/*
 *    Copyright 2018 Insolar
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

package pulsar

import (
	"context"
	"fmt"
	"sync"

	"github.com/insolar/insolar/instrumentation/inslogger"
)

//go:generate stringer -type=State
type State int

const (
	// Failed means that current iteration is broken
	Failed State = iota

	// WaitingForStart means that state machine is waiting for the start
	WaitingForStart

	// GenerateEntropy means that state machine is generating entropy for a current slot
	GenerateEntropy

	// WaitingForEntropySigns means that state machine is waiting for other pulsars' signs of entropy
	WaitingForEntropySigns

	// SendingEntropy means that state machine is sending entropy to other pulsars
	SendingEntropy

	// WaitingForEntropy means that state machine is waiting for the entropy for other pulsars
	WaitingForEntropy

	// SendingVector means that state machine is sending verctor to other pulsars
	SendingVector

	// WaitingForVectors means that state machine is waiting for other pulsars' vectors
	WaitingForVectors

	// Verifying means that state machine is verifying bft-table
	Verifying

	// SendingPulseSign means that state machine is sending sign to chosen pulsar
	SendingPulseSign

	// WaitingForPulseSigns means that state machine is waiting for signs to chosen pulsar
	WaitingForPulseSigns

	// SendingPulseSign means that state machine is sending pulse to network
	SendingPulse
)

// StateSwitcher is a base for pulsar's state machine
type StateSwitcher interface {
	SwitchToState(ctx context.Context, state State, args interface{})
	GetState() State
	setState(state State)
	SetPulsar(pulsar *Pulsar)
}

// StateSwitcherImpl is a base implementation of the pulsar's state machine
type StateSwitcherImpl struct {
	pulsar *Pulsar
	state  State
	lock   sync.RWMutex
}

func (switcher *StateSwitcherImpl) GetState() State {
	switcher.lock.RLock()
	defer switcher.lock.RUnlock()
	return switcher.state
}

func (switcher *StateSwitcherImpl) setState(state State) {
	switcher.lock.Lock()
	defer switcher.lock.Unlock()
	switcher.state = state
}

// SetPulsar sets pulsar of the current instance
func (switcher *StateSwitcherImpl) SetPulsar(pulsar *Pulsar) {
	switcher.setState(WaitingForStart)
	switcher.pulsar = pulsar
}

// SwitchToState switches the state-machine to another step
func (switcher *StateSwitcherImpl) SwitchToState(ctx context.Context, state State, args interface{}) {
	logger := inslogger.FromContext(ctx)
	logger.Debugf("Switch state from %v to %v, node - %v", switcher.GetState().String(), state.String(), switcher.pulsar.Config.MainListenerAddress)
	if state < switcher.GetState() && (state != WaitingForStart && state != Failed) {
		panic(fmt.Sprintf("Attempt to set a backward step. %v", switcher.pulsar.Config.MainListenerAddress))
	}

	switcher.setState(state)

	switch state {
	case WaitingForStart:
		switcher.pulsar.clearState()
	case WaitingForEntropySigns:
		switcher.pulsar.waitForEntropySigns(ctx)
	case SendingEntropy:
		switcher.pulsar.sendEntropy(ctx)
	case WaitingForEntropy:
		switcher.pulsar.waitForEntropy(ctx)
	case SendingVector:
		switcher.pulsar.sendVector(ctx)
	case WaitingForVectors:
		switcher.pulsar.waitForVectors(ctx)
	case Verifying:
		switcher.pulsar.verify(ctx)
	case WaitingForPulseSigns:
		switcher.pulsar.waitForPulseSigns(ctx)
	case SendingPulseSign:
		switcher.pulsar.sendPulseSign(ctx)
	case SendingPulse:
		switcher.pulsar.sendPulseToNodesAndPulsars(ctx)
	case Failed:
		switcher.pulsar.handleErrorState(ctx, args.(error))
		switcher.setState(WaitingForStart)
	}
}
