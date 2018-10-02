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
	"github.com/insolar/insolar/log"
)

//go:generate stringer -type=State
type State int

const (
	WaitingForTheStart State = iota + 1
	WaitingForTheSigns
	SendingEntropy
	WaitingForTheEntropy
	SendingVector
	WaitingForTheVectors
	Verifying
	SendingSignForChosen
	WaitingForChosenSigns
	SendingEntropyToNodes
	Failed
)

// StateSwitcher is a base for pulsar's state machine
type StateSwitcher interface {
	switchToState(state State, args interface{})
}

// StateSwitcherImpl is a base implementation of the pulsar's state machine
type StateSwitcherImpl struct {
	pulsar *Pulsar
}

func (switcher *StateSwitcherImpl) SetPulsar(pulsar *Pulsar) {
	switcher.pulsar = pulsar
}

func (switcher *StateSwitcherImpl) switchToState(state State, args interface{}) {
	log.Debugf("Switch state from %v to %v", switcher.pulsar.State.String(), state.String())
	if state < switcher.pulsar.State && state != WaitingForTheStart {
		panic("Attempt to set a backward step")
	}

	switcher.pulsar.State = state
	switch state {
	case WaitingForTheStart:
		log.Info("Switch to start")
	case WaitingForTheSigns:
		switcher.pulsar.stateSwitchedToWaitingForSigns()
	case SendingEntropy:
		switcher.pulsar.stateSwitchedToSendingEntropy()
	case WaitingForTheEntropy:
		switcher.pulsar.stateSwitchedWaitingForTheEntropy()
	case SendingVector:
		switcher.pulsar.stateSwitchedToSendingVector()
	case WaitingForTheVectors:
		switcher.pulsar.stateSwitchedToReceivingVector()
	case Verifying:
		switcher.pulsar.stateSwitchedToVerifying()
	case WaitingForChosenSigns:
		switcher.pulsar.stateSwitchedToWaitingForChosenSigns()
	case SendingSignForChosen:
		switcher.pulsar.stateSwitchedToSendingSignForChosen()
	case SendingEntropyToNodes:
		switcher.pulsar.stateSwitchedToSendingEntropyToNodes()
	case Failed:
		switcher.pulsar.stateSwitchedToFailed(args.(error))
	}
}
