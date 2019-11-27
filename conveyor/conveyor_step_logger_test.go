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
	"reflect"
	"runtime"
	"strings"

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

func getStepName(step interface{}) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(step).Pointer()).Name()
	if lastIndex := strings.LastIndex(fullName, "."); lastIndex >= 0 {
		fullName = fullName[lastIndex+1:]
	}
	if lastIndex := strings.LastIndex(fullName, "-"); lastIndex >= 0 {
		fullName = fullName[:lastIndex]
	}

	return fullName
}

func (v conveyorStepLogger) prepareStepName(sd *smachine.StepDeclaration) {
	if !sd.IsNameless() {
		return
	}
	sd.Name = getStepName(sd.Transition)
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

	v.prepareStepName(&data.CurrentStep)
	v.prepareStepName(&upd.NextStep)

	detached := ""
	if upd.Flags&smachine.StepLoggerDetached != 0 {
		detached = "(detached)"
	}
	fmt.Printf("%s[%3d]: %03d @ %03d: %s%s%s current=%v next=%v payload=%T tracer=%v\n", data.StepNo.MachineId(), data.CycleNo,
		data.StepNo.SlotID(), data.StepNo.StepNo(),
		special, upd.UpdateType, detached, data.CurrentStep.GetStepName(), upd.NextStep.GetStepName(), v.sm, v.tracer)
}

func (v conveyorStepLogger) LogInternal(data smachine.StepLoggerData, updateType string) {
	v.prepareStepName(&data.CurrentStep)

	fmt.Printf("%s[%3d]: %03d @ %03d: internal %s current=%v payload=%T tracer=%v\n", data.StepNo.MachineId(), data.CycleNo,
		data.StepNo.SlotID(), data.StepNo.StepNo(),
		updateType, data.CurrentStep.GetStepName(), v.sm, v.tracer)
}

func (v conveyorStepLogger) LogEvent(data smachine.StepLoggerData, customEvent interface{}) {
	special := ""

	v.prepareStepName(&data.CurrentStep)

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
		fmt.Printf("%s[%3d]: %03d @ %03d: unknown (%s) current=%v payload=%T tracer=%v\n", data.StepNo.MachineId(), data.CycleNo,
			data.StepNo.SlotID(), data.StepNo.StepNo(),
			customEvent, data.CurrentStep.GetStepName(), v.sm, v.tracer)
		return
	}
	fmt.Printf("%s[%3d]: %03d @ %03d: custom %s current=%v event=%v payload=%T tracer=%v\n", data.StepNo.MachineId(), data.CycleNo,
		data.StepNo.SlotID(), data.StepNo.StepNo(),
		special, data.CurrentStep.GetStepName(), customEvent, v.sm, v.tracer)

	if data.EventType == smachine.StepLoggerFatal {
		panic("os.Exit(1)")
	}
}