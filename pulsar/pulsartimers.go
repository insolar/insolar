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
	"time"

	"github.com/insolar/insolar/log"
)

func (currentPulsar *Pulsar) waitForPulseSigns() {
	log.Debug("[waitForPulseSigns]")
	ticker := time.NewTicker(10 * time.Millisecond)
	currentTimeOut := time.Now().Add(time.Duration(currentPulsar.Config.ReceivingSignsForChosenTimeout) * time.Millisecond)
	go func() {
		for range ticker.C {
			if currentPulsar.IsStateFailed() || currentPulsar.StateSwitcher.GetState() > WaitingForPulseSigns {
				ticker.Stop()
				return
			}

			if time.Now().After(currentTimeOut) {
				ticker.Stop()
				currentPulsar.StateSwitcher.SwitchToState(SendingPulse, nil)
			}
		}
	}()
}

func (currentPulsar *Pulsar) waitForEntropy() {
	log.Debug("[waitForEntropy]")
	ticker := time.NewTicker(10 * time.Millisecond)
	timeout := time.Now().Add(time.Duration(currentPulsar.Config.ReceivingNumberTimeout) * time.Millisecond)
	go func() {
		for range ticker.C {
			if currentPulsar.IsStateFailed() || currentPulsar.StateSwitcher.GetState() > WaitingForEntropy {
				ticker.Stop()
				return
			}

			if time.Now().After(timeout) {
				ticker.Stop()
				currentPulsar.StateSwitcher.SwitchToState(SendingVector, nil)
			}
		}
	}()
}

func (currentPulsar *Pulsar) waitForEntropySigns() {
	log.Debug("[waitForEntropySigns]")
	ticker := time.NewTicker(10 * time.Millisecond)
	currentTimeOut := time.Now().Add(time.Duration(currentPulsar.Config.ReceivingSignTimeout) * time.Millisecond)
	go func() {
		for range ticker.C {
			if currentPulsar.IsStateFailed() || currentPulsar.StateSwitcher.GetState() > WaitingForEntropySigns {
				ticker.Stop()
				return
			}

			if time.Now().After(currentTimeOut) {
				ticker.Stop()
				currentPulsar.StateSwitcher.SwitchToState(SendingEntropy, nil)
			}
		}
	}()
}

func (currentPulsar *Pulsar) waitForVectors() {
	log.Debug("[waitForVectors]")
	ticker := time.NewTicker(10 * time.Millisecond)
	currentTimeOut := time.Now().Add(time.Duration(currentPulsar.Config.ReceivingVectorTimeout) * time.Millisecond)
	go func() {
		for range ticker.C {
			if currentPulsar.IsStateFailed() || currentPulsar.StateSwitcher.GetState() > WaitingForVectors {
				ticker.Stop()
				return
			}

			if time.Now().After(currentTimeOut) {
				ticker.Stop()
				currentPulsar.StateSwitcher.SwitchToState(Verifying, nil)
			}
		}
	}()
}
