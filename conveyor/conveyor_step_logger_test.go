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

	"github.com/insolar/insolar/conveyor/smachine"
)

type conveyorStepLogger struct {
	ctx    context.Context
	sm     smachine.StateMachine
	tracer smachine.TracerId
}

func (conveyorStepLogger) CanLogEvent(eventType smachine.StepLoggerEvent, stepLevel smachine.StepLogLevel) bool {
	return true
}

func (v conveyorStepLogger) GetTracerId() smachine.TracerId {
	return v.tracer
}

func (v conveyorStepLogger) LogUpdate(data smachine.StepLoggerData, upd smachine.StepLoggerUpdateData) {
	special := ""

	switch data.EventType {
	case smachine.StepLoggerUpdate:
	case smachine.StepLoggerMigrate:
		special = "migrate "
	default:
		panic("illegal value")
	}

	detached := ""
	if upd.Flags&smachine.StepLoggerDetached != 0 {
		detached = "(detached)"
	}
	fmt.Printf("%s[%3d]: %03d @ %03d: %s%s%s current=%p next=%p payload=%T tracer=%v\n", data.StepNo.MachineId(), data.CycleNo,
		data.StepNo.SlotID(), data.StepNo.StepNo(),
		special, upd.UpdateType, detached, data.CurrentStep.Transition, upd.NextStep.Transition, v.sm, v.tracer)
}

func (v conveyorStepLogger) LogInternal(data smachine.StepLoggerData, updateType string) {
	fmt.Printf("%s[%3d]: %03d @ %03d: internal %s current=%p payload=%T tracer=%v\n", data.StepNo.MachineId(), data.CycleNo,
		data.StepNo.SlotID(), data.StepNo.StepNo(),
		updateType, data.CurrentStep.Transition, v.sm, v.tracer)
}

func (v conveyorStepLogger) LogEvent(data smachine.StepLoggerData, customEvent interface{}) {
	special := ""

	switch data.EventType {
	case smachine.StepLoggerTrace:
		special = "TRC"
	case smachine.StepLoggerActiveTrace:
		special = "TRA"
	case smachine.StepLoggerWarn:
		special = "WRN"
	case smachine.StepLoggerError:
		special = "ERR"
	case smachine.StepLoggerFatal:
		special = "FTL"
	default:
		fmt.Printf("%s[%3d]: %03d @ %03d: unknown (%s) current=%p payload=%T tracer=%v\n", data.StepNo.MachineId(), data.CycleNo,
			data.StepNo.SlotID(), data.StepNo.StepNo(),
			customEvent, data.CurrentStep.Transition, v.sm, v.tracer)
		return
	}
	fmt.Printf("%s[%3d]: %03d @ %03d: custom %s current=%p event=%v payload=%T tracer=%v\n", data.StepNo.MachineId(), data.CycleNo,
		data.StepNo.SlotID(), data.StepNo.StepNo(),
		special, data.CurrentStep.Transition, customEvent, v.sm, v.tracer)

	if data.EventType == smachine.StepLoggerFatal {
		panic("os.Exit(1)")
	}
}
