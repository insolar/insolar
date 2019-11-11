package statemachine

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/insolar/insolar/conveyor"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type ConveyorLogger struct {
	smachine.StepLoggerStub

	logger insolar.Logger
}

func (c ConveyorLogger) CanLogEvent(eventType smachine.StepLoggerEvent, stepLevel smachine.StepLogLevel) bool {
	return true
}

type LogStepMessage struct {
	*insolar.LogObjectTemplate

	Message   string
	Component string `txt:"sm"`
	TraceID   string

	MachineName interface{} `fmt:"%T"`
	MachineID   string
	SlotStep    string
	From        string
	To          string `opt:""`
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

func prepareStepName(sd *smachine.StepDeclaration) {
	if !sd.IsNameless() {
		return
	}
	sd.Name = getStepName(sd.Transition)
}

func (c ConveyorLogger) LogEvent(data smachine.StepLoggerData, msg interface{}) {
	c.logger.Error(msg)
}

func (c ConveyorLogger) LogUpdate(stepLoggerData smachine.StepLoggerData, stepLoggerUpdateData smachine.StepLoggerUpdateData) {
	special := ""

	switch stepLoggerData.EventType {
	case smachine.StepLoggerUpdate:
	case smachine.StepLoggerMigrate:
		special = "migrate "
	default:
		panic("illegal value")
	}

	prepareStepName(&stepLoggerData.CurrentStep)
	prepareStepName(&stepLoggerUpdateData.NextStep)

	detached := ""
	if stepLoggerUpdateData.Flags&smachine.StepLoggerDetached != 0 {
		detached = " (detached)"
	}

	if _, ok := stepLoggerData.Declaration.(*conveyor.PulseSlotMachine); ok {
		return
	}

	c.logger.Error(LogStepMessage{
		Message: special + stepLoggerUpdateData.UpdateType + detached,

		MachineName: stepLoggerData.Declaration,
		MachineID:   fmt.Sprintf("%s[%3d]", stepLoggerData.StepNo.MachineId(), stepLoggerData.CycleNo),
		SlotStep:    fmt.Sprintf("%03d @ %03d", stepLoggerData.StepNo.StepNo(), stepLoggerData.StepNo.StepNo()),
		From:        stepLoggerData.CurrentStep.GetStepName(),
		To:          stepLoggerUpdateData.NextStep.GetStepName(),
	})
}

func NewConveyorLogger(ctx context.Context, _ smachine.StateMachine, traceID smachine.TracerId) smachine.StepLogger {
	_, logger := inslogger.WithTraceField(context.Background(), traceID)
	return &ConveyorLogger{
		StepLoggerStub: smachine.StepLoggerStub{TracerId: traceID},
		logger:         logger,
	}
}
