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

package pulsar

import (
	"context"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

func (currentPulsar *Pulsar) waitForPulseSigns(ctx context.Context) {
	ctx, span := instracer.StartSpan(ctx, "Pulsar.waitForPulseSigns")
	defer span.End()

	inslogger.FromContext(ctx).Debug("[waitForPulseSigns]")
	ticker := time.NewTicker(10 * time.Millisecond)
	currentTimeOut := time.Now().Add(time.Duration(currentPulsar.Config.ReceivingSignsForChosenTimeout) * time.Millisecond)
	go func() {
		ctx, span := instracer.StartSpan(ctx, "goroutine")
		defer span.End()

		for range ticker.C {
			if currentPulsar.IsStateFailed() || currentPulsar.StateSwitcher.GetState() > WaitingForPulseSigns {
				ticker.Stop()
				return
			}

			if time.Now().After(currentTimeOut) {
				ticker.Stop()
				go currentPulsar.StateSwitcher.SwitchToState(ctx, SendingPulse, nil)
				return
			}
		}
	}()
}

func (currentPulsar *Pulsar) waitForEntropy(ctx context.Context) {
	ctx, span := instracer.StartSpan(ctx, "Pulsar.waitForEntropy")
	defer span.End()

	inslogger.FromContext(ctx).Debug("[waitForEntropy]")
	ticker := time.NewTicker(10 * time.Millisecond)
	timeout := time.Now().Add(time.Duration(currentPulsar.Config.ReceivingNumberTimeout) * time.Millisecond)
	go func() {
		ctx, span := instracer.StartSpan(ctx, "goroutine")
		defer span.End()

		for range ticker.C {
			if currentPulsar.IsStateFailed() || currentPulsar.StateSwitcher.GetState() > WaitingForEntropy {
				ticker.Stop()
				return
			}

			if time.Now().After(timeout) {
				ticker.Stop()
				go currentPulsar.StateSwitcher.SwitchToState(ctx, SendingVector, nil)
				return
			}
		}
	}()
}

func (currentPulsar *Pulsar) waitForEntropySigns(ctx context.Context) {
	ctx, span := instracer.StartSpan(ctx, "Pulsar.waitForEntropySigns")
	defer span.End()

	inslogger.FromContext(ctx).Debug("[waitForEntropySigns]")
	ticker := time.NewTicker(10 * time.Millisecond)
	currentTimeOut := time.Now().Add(time.Duration(currentPulsar.Config.ReceivingSignTimeout) * time.Millisecond)
	go func() {
		ctx, span := instracer.StartSpan(ctx, "goroutine")
		defer span.End()

		for range ticker.C {
			if currentPulsar.IsStateFailed() || currentPulsar.StateSwitcher.GetState() > WaitingForEntropySigns {
				ticker.Stop()
				return
			}

			if time.Now().After(currentTimeOut) {
				ticker.Stop()
				go currentPulsar.StateSwitcher.SwitchToState(ctx, SendingEntropy, nil)
				return
			}
		}
	}()
}

func (currentPulsar *Pulsar) waitForVectors(ctx context.Context) {
	ctx, span := instracer.StartSpan(ctx, "Pulsar.waitForVectors")
	defer span.End()

	inslogger.FromContext(ctx).Debug("[waitForVectors]")
	ticker := time.NewTicker(10 * time.Millisecond)
	currentTimeOut := time.Now().Add(time.Duration(currentPulsar.Config.ReceivingVectorTimeout) * time.Millisecond)
	go func() {
		ctx, span := instracer.StartSpan(ctx, "goroutine")
		defer span.End()

		for range ticker.C {
			if currentPulsar.IsStateFailed() || currentPulsar.StateSwitcher.GetState() > WaitingForVectors {
				ticker.Stop()
				return
			}

			if time.Now().After(currentTimeOut) {
				ticker.Stop()
				go currentPulsar.StateSwitcher.SwitchToState(ctx, Verifying, nil)
				return
			}
		}
	}()
}
