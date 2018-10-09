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
	failed State = iota
	waitingForStart
	generateEntropy
	waitingForEntropySigns
	sendingEntropy
	waitingForEntropy
	sendingVector
	waitingForVectors
	verifying
	sendingPulseSign
	waitingForPulseSigns
	sendingPulse
)

// StateSwitcher is a base for pulsar's state machine
type StateSwitcher interface {
	switchToState(state State, args interface{})
	getState() State
	setState(state State)
	SetPulsar(pulsar *Pulsar)
}

// StateSwitcherImpl is a base implementation of the pulsar's state machine
type StateSwitcherImpl struct {
	pulsar *Pulsar
	state  State
	lock   sync.RWMutex
}

func (switcher *StateSwitcherImpl) getState() State {
	switcher.lock.RLock()
	defer switcher.lock.RUnlock()
	return switcher.state
}

func (switcher *StateSwitcherImpl) setState(state State) {
	switcher.lock.Lock()
	defer switcher.lock.Unlock()
	switcher.state = state
}

func (switcher *StateSwitcherImpl) SetPulsar(pulsar *Pulsar) {
	switcher.pulsar = pulsar
}

func (switcher *StateSwitcherImpl) switchToState(state State, args interface{}) {
	log.Debugf("Switch state from %v to %v", switcher.getState().String(), state.String())
	if state < switcher.getState() && (state != waitingForStart && state != failed) {
		panic("Attempt to set a backward step")
	}

	switcher.setState(state)

	switch state {
	case waitingForStart:
		switcher.pulsar.clearState()
	case waitingForEntropySigns:
		switcher.pulsar.waitForEntropySigns()
	case sendingEntropy:
		switcher.pulsar.sendEntropy()
	case waitingForEntropy:
		switcher.pulsar.waitForEntropy()
	case sendingVector:
		switcher.pulsar.sendVector()
	case waitingForVectors:
		switcher.pulsar.receiveVectors()
	case verifying:
		switcher.pulsar.verify()
	case waitingForPulseSigns:
		switcher.pulsar.waitForPulseSigns()
	case sendingPulseSign:
		switcher.pulsar.sendPulseSign()
	case sendingPulse:
		switcher.pulsar.sendPulse()
	case failed:
		switcher.pulsar.handleErrorState(args.(error))
		switcher.setState(waitingForStart)
	}
}
