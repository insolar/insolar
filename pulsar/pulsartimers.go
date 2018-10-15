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
			if currentPulsar.IsStateFailed() || currentPulsar.StateSwitcher.GetState() == SendingPulse {
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
			if currentPulsar.IsStateFailed() || currentPulsar.StateSwitcher.GetState() == SendingVector {
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
			if currentPulsar.IsStateFailed() || currentPulsar.StateSwitcher.GetState() == SendingEntropy {
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
			if currentPulsar.IsStateFailed() || currentPulsar.StateSwitcher.GetState() == Verifying {
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
