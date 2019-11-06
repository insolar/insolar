//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package conveyor

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/pulse"
)

func TestConveyor(t *testing.T) {
	machineConfig := smachine.SlotMachineConfig{
		SlotPageSize:    1000,
		PollingPeriod:   500 * time.Millisecond,
		PollingTruncate: 1 * time.Millisecond,
		ScanCountLimit:  100000,
	}

	factoryFn := func(pn pulse.Number, v InputEvent) smachine.CreateFunc {
		return func(ctx smachine.ConstructionContext) smachine.StateMachine {
			sm := &AppEventSM{eventValue: v, pn: pn}
			return sm
		}
	}
	machineConfig.StepLoggerFactoryFn = func(ctx context.Context) smachine.StepLoggerFunc {
		return func(data *smachine.StepLoggerData) {
			stepLogger(ctx, data)
		}
	}

	conveyor := NewPulseConveyor(context.Background(), machineConfig, 100*time.Millisecond, factoryFn, machineConfig, nil)

	pd := pulse.NewFirstPulsarData(10, longbits.Bits256{})

	conveyor.StartWorker(nil, nil)

	require.NoError(t, conveyor.CommitPulseChange(pd.AsRange()))
	eventCount := 0

	time.Sleep(100 * time.Millisecond)

	for i := 0; i < 100; i++ {
		pd = pd.CreateNextPulsarPulse(10, func() longbits.Bits256 {
			return longbits.Bits256{}
		})
		fmt.Println("================================== ", pd, " ====================================")
		require.NoError(t, conveyor.CommitPulseChange(pd.AsRange()))

		if eventCount < math.MaxInt32 {
			eventCount++
			require.NoError(t, conveyor.AddInput(context.Background(), pd.NextPulseNumber(), fmt.Sprintf("event-%d-future", eventCount)))
			eventCount++
			require.NoError(t, conveyor.AddInput(context.Background(), pd.PrevPulseNumber(), fmt.Sprintf("event-%d-past", eventCount)))

			for j := 0; j < 1; j++ {
				eventCount++
				require.NoError(t, conveyor.AddInput(context.Background(), pd.PulseNumber, fmt.Sprintf("event-%d-present", eventCount)))
			}
		}

		time.Sleep(time.Second)
	}
	time.Sleep(100 * time.Millisecond)
	fmt.Println("======================")
	conveyor.StopNoWait()
	time.Sleep(time.Hour)
}

func stepLogger(_ context.Context, data *smachine.StepLoggerData) {
	if data.Flags&smachine.StepLoggerInternal != 0 {
		fmt.Printf("%s: %03d @ %03d: internal %s current=%p payload=%T:%+v\n", data.StepNo.MachineId(), data.StepNo.SlotID(), data.StepNo.StepNo(),
			data.UpdateType, data.CurrentStep.Transition, data.SM, data.SM)
		return
	}

	special := ""
	if data.Flags&smachine.StepLoggerMigrate != 0 {
		special = "migrate "
	}
	detached := ""
	if data.Flags&smachine.StepLoggerDetached != 0 {
		detached = "(detached)"
	}
	fmt.Printf("%s: %03d @ %03d: %s%s%s current=%p next=%p payload=%T:%+v\n", data.StepNo.MachineId(), data.StepNo.SlotID(), data.StepNo.StepNo(),
		special, data.UpdateType, detached, data.CurrentStep.Transition, data.NextStep.Transition, data.SM, data.SM)
}
