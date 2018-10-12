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
	"sync"

	"github.com/insolar/insolar/log"
)

//go:generate stringer -type=State
type State int

const (
	Failed State = iota
	WaitingForStart
	GenerateEntropy
	WaitingForEntropySigns
	SendingEntropy
	WaitingForEntropy
	SendingVector
	WaitingForVectors
	Verifying
	SendingPulseSign
	WaitingForPulseSigns
	SendingPulse
)

// StateSwitcher is a base for pulsar's state machine
type StateSwitcher interface {
	switchToState(state State, args interface{})
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

func (switcher *StateSwitcherImpl) switchToState(state State, args interface{}) {
	log.Debugf("Switch state from %v to %v", switcher.GetState().String(), state.String())
	if state < switcher.GetState() && (state != WaitingForStart && state != Failed) {
		panic("Attempt to set a backward step")
	}

	switcher.setState(state)

	switch state {
	case WaitingForStart:
		switcher.pulsar.clearState()
	case WaitingForEntropySigns:
		switcher.pulsar.waitForEntropySigns()
	case SendingEntropy:
		switcher.pulsar.sendEntropy()
	case WaitingForEntropy:
		switcher.pulsar.waitForEntropy()
	case SendingVector:
		switcher.pulsar.sendVector()
	case WaitingForVectors:
		switcher.pulsar.receiveVectors()
	case Verifying:
		switcher.pulsar.verify()
	case WaitingForPulseSigns:
		switcher.pulsar.waitForPulseSigns()
	case SendingPulseSign:
		switcher.pulsar.sendPulseSign()
	case SendingPulse:
		switcher.pulsar.sendPulse()
	case Failed:
		switcher.pulsar.handleErrorState(args.(error))
		switcher.setState(WaitingForStart)
	}
}
