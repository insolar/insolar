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

type conveyorSlotMachineLogger struct {
}

func (conveyorSlotMachineLogger) LogMachineInternal(data smachine.SlotMachineData, msg string) {
	fmt.Printf("[MACHINE][LOG] %s[%3d]: %03d @ %03d: internal %s err=%v\n", data.StepNo.MachineId(), data.CycleNo,
		data.StepNo.SlotID(), data.StepNo.StepNo(), msg, data.Error)
}

func (conveyorSlotMachineLogger) LogMachineCritical(data smachine.SlotMachineData, msg string) {
	fmt.Printf("[MACHINE][ERR] %s[%3d]: %03d @ %03d: internal %s err=%v\n", data.StepNo.MachineId(), data.CycleNo,
		data.StepNo.SlotID(), data.StepNo.StepNo(), msg, data.Error)
}

func (conveyorSlotMachineLogger) CreateStepLogger(ctx context.Context, sm smachine.StateMachine, tracer smachine.TracerId) smachine.StepLogger {
	return conveyorStepLogger{ctx, sm, tracer}
}

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

func (v conveyorStepLogger) CreateAsyncLogger(ctx context.Context, _ *smachine.StepLoggerData) (context.Context, smachine.StepLogger) {
	return ctx, v
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
	if data.Flags&smachine.StepLoggerDetached != 0 {
		detached = "(detached)"
	}

	durations := ""
	if upd.InactivityNano > 0 || upd.ActivityNano > 0 {
		durations = fmt.Sprintf(" timing=%s/%s", upd.InactivityNano, upd.ActivityNano)
	}

	if data.Error == nil {
		fmt.Printf("[LOG] %s[%3d]: %03d @ %03d: %s%s%s%s current=%v next=%v payload=%T tracer=%v\n", data.StepNo.MachineId(), data.CycleNo,
			data.StepNo.SlotID(), data.StepNo.StepNo(),
			special, upd.UpdateType, detached, durations,
			data.CurrentStep.GetStepName(), upd.NextStep.GetStepName(), v.sm, v.tracer)
		return
	}

	errSpecial := ""
	switch data.Flags & smachine.StepLoggerErrorMask {
	case smachine.StepLoggerUpdateErrorMuted:
		errSpecial = "muted "
	case smachine.StepLoggerUpdateErrorRecovered:
		errSpecial = "recovered "
	case smachine.StepLoggerUpdateErrorRecoveryDenied:
		errSpecial = "recover-denied "
	}

	fmt.Printf("[ERR] %s[%3d]: %03d @ %03d: %s%s%s%s current=%v next=%v payload=%T tracer=%v err=%v\n", data.StepNo.MachineId(), data.CycleNo,
		data.StepNo.SlotID(), data.StepNo.StepNo(),
		special, errSpecial, upd.UpdateType, detached, data.CurrentStep.GetStepName(), upd.NextStep.GetStepName(), v.sm, v.tracer, data.Error)
}

func (v conveyorStepLogger) LogInternal(data smachine.StepLoggerData, updateType string) {
	v.prepareStepName(&data.CurrentStep)

	if data.Error == nil {
		fmt.Printf("[LOG] %s[%3d]: %03d @ %03d: internal %s current=%v payload=%T tracer=%v\n", data.StepNo.MachineId(), data.CycleNo,
			data.StepNo.SlotID(), data.StepNo.StepNo(),
			updateType, data.CurrentStep.GetStepName(), v.sm, v.tracer)
	} else {
		fmt.Printf("[ERR] %s[%3d]: %03d @ %03d: internal %s current=%v payload=%T tracer=%v err=%v\n", data.StepNo.MachineId(), data.CycleNo,
			data.StepNo.SlotID(), data.StepNo.StepNo(),
			updateType, data.CurrentStep.GetStepName(), v.sm, v.tracer, data.Error)
	}
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
		fmt.Printf("[U%d] %s[%3d]: %03d @ %03d: current=%v payload=%T tracer=%v custom=%v\n", data.EventType, data.StepNo.MachineId(), data.CycleNo,
			data.StepNo.SlotID(), data.StepNo.StepNo(),
			data.CurrentStep.GetStepName(), v.sm, v.tracer, customEvent)
		return
	}
	fmt.Printf("[%s] %s[%3d]: %03d @ %03d: current=%v payload=%T tracer=%v custom=%v\n", special, data.StepNo.MachineId(), data.CycleNo,
		data.StepNo.SlotID(), data.StepNo.StepNo(),
		data.CurrentStep.GetStepName(), v.sm, v.tracer, customEvent)

	if data.EventType == smachine.StepLoggerFatal {
		panic("os.Exit(1)")
	}
}

func (v conveyorStepLogger) LogAdapter(data smachine.StepLoggerData, adapterId smachine.AdapterId, callId uint64) {
	//case smachine.StepLoggerAdapterCall:
	s := "?"
	switch data.Flags & smachine.StepLoggerAdapterMask {
	case smachine.StepLoggerAdapterSyncCall:
		s = "sync-call"
	case smachine.StepLoggerAdapterAsyncCall:
		s = "async-call"
	case smachine.StepLoggerAdapterAsyncResult:
		s = "async-result"
	case smachine.StepLoggerAdapterAsyncCancel:
		s = "async-cancel"
	}
	fmt.Printf("[ADP] %s %s[%3d]: %03d @ %03d: current=%v payload=%T tracer=%v adapter=%v/%v\n", s, data.StepNo.MachineId(), data.CycleNo,
		data.StepNo.SlotID(), data.StepNo.StepNo(),
		data.CurrentStep.GetStepName(), v.sm, v.tracer, adapterId, callId)
}
