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
	"github.com/insolar/insolar/conveyor/injector"
	smachine "github.com/insolar/insolar/conveyor/smachinev2"
	"github.com/insolar/insolar/pulse"
)

type PulseSlotState uint8

const (
	Uninitialized PulseSlotState = iota
	Future
	Present
	Past
)

type PulseSlotConfig struct {
	config         smachine.SlotMachineConfig
	signalCallback func()
	parentRegistry injector.DependencyRegistry
}

func newFuturePulseSlot(pd pulse.Data, config PulseSlotConfig) *PulseSlotMachine {
	return &PulseSlotMachine{
		SlotMachine: smachine.NewSlotMachine(config.config, config.signalCallback, config.parentRegistry),
		pulseSlot: PulseSlot{&futurePulseDataHolder{
			pd: pd,
		}},
	}
}

func newPresentPulseSlot(pd pulse.Data, config PulseSlotConfig) *PulseSlotMachine {
	return &PulseSlotMachine{
		SlotMachine: smachine.NewSlotMachine(config.config, config.signalCallback, config.parentRegistry),
		pulseSlot: PulseSlot{&presentPulseDataHolder{
			pd: pd,
		}},
	}
}

func newPastPulseSlot(pd pulse.Data, config PulseSlotConfig) *PulseSlotMachine {
	return &PulseSlotMachine{
		SlotMachine: smachine.NewSlotMachine(config.config, config.signalCallback, config.parentRegistry),
		pulseSlot: PulseSlot{&pastPulseDataHolder{
			pd: pd,
		}},
	}
}

type PulseSlotMachine struct {
	smachine.SlotMachine
	pulseSlot PulseSlot
	isAntique bool
}

type PulseSlot struct {
	pulseData pulseDataHolder
}

func (p *PulseSlot) State() PulseSlotState {
	return p.pulseData.State()
}

func (p *PulseSlot) PulseData() pulse.Data {
	return p.pulseData.PulseData()
}
