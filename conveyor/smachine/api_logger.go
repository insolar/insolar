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

package smachine

import "context"

type StepLoggerEvent uint8

const (
	StepLoggerUpdate StepLoggerEvent = iota
	StepLoggerMigrate
	StepLoggerInternal

	StepLoggerTrace
	StepLoggerActiveTrace
	StepLoggerWarn
	StepLoggerError
	StepLoggerFatal
)

type StepLoggerUpdateFlags uint8

const (
	StepLoggerUpdateErrorDefault StepLoggerUpdateFlags = iota
	StepLoggerUpdateErrorMuted
	StepLoggerUpdateErrorRecovered
	StepLoggerUpdateErrorRecoveryDenied
)
const StepLoggerErrorMask StepLoggerUpdateFlags = 3

const (
	StepLoggerDetached StepLoggerUpdateFlags = 1 << (2 + iota)
)

type SlotMachineData struct {
	CycleNo uint32
	StepNo  StepLink
	Error   error
}

type StepLoggerData struct {
	CycleNo     uint32
	StepNo      StepLink
	CurrentStep StepDeclaration
	Declaration StateMachineDeclaration
	EventType   StepLoggerEvent
	Error       error
}

type StepLoggerUpdateData struct {
	UpdateType string
	PrevStepNo uint32
	NextStep   StepDeclaration
	Flags      StepLoggerUpdateFlags
}

type SlotMachineLogger interface {
	CreateStepLogger(context.Context, StateMachine, TracerId) StepLogger
	LogInternal(data SlotMachineData, msg string)
	LogCritical(data SlotMachineData, msg string)
}

type StepLoggerFactoryFunc func(context.Context, StateMachine, TracerId) StepLogger

type StepLogLevel uint8

const (
	StepLogLevelDefault StepLogLevel = iota
	StepLogLevelElevated
	StepLogLevelTracing
)

type StepLogger interface {
	CanLogEvent(eventType StepLoggerEvent, stepLevel StepLogLevel) bool
	//LogMetric()
	LogUpdate(StepLoggerData, StepLoggerUpdateData)
	LogInternal(data StepLoggerData, updateType string)
	LogEvent(data StepLoggerData, customEvent interface{})

	GetTracerId() TracerId
}

type StepLoggerFunc func(*StepLoggerData, *StepLoggerUpdateData)

type TracerId = string

type Logger = slotLogger

type slotLogger struct { // we use an explicit struct here to enable compiler optimizations when logging is not needed
	c *slotContext
}

func (p slotLogger) IsTracing() bool {
	p.c.ensureValid()
	return p.c.s.isTracing()
}

func (p slotLogger) getStepLogger() StepLogger {
	p.c.ensureValid()
	return p.c.s.stepLogger
}

func (p slotLogger) GetTracerId() TracerId {
	if stepLogger := p.getStepLogger(); stepLogger != nil {
		return stepLogger.GetTracerId()
	}
	return ""
}

func (p slotLogger) _logCustom(eventType StepLoggerEvent, msg interface{}, err error) {
	if stepLogger := p.getStepLogger(); stepLogger != nil {
		stepLevel := p.c.s.getStepLogLevel()
		if !stepLogger.CanLogEvent(eventType, stepLevel) {
			return
		}
		s := p.c.s

		if stepLevel == StepLogLevelTracing && eventType == StepLoggerTrace {
			eventType = StepLoggerActiveTrace
		}

		stepData := s.newStepLoggerData(eventType, s.NewStepLink())
		stepData.Error = err
		stepLogger.LogEvent(stepData, msg)
	}
}

func (p slotLogger) Trace(msg interface{}) {
	p._logCustom(StepLoggerTrace, msg, nil)
}

func (p slotLogger) Warn(msg interface{}) {
	p._logCustom(StepLoggerWarn, msg, nil)
}

func (p slotLogger) Error(msg interface{}, err error) {
	p._logCustom(StepLoggerError, msg, err)
}

func (p slotLogger) Fatal(msg interface{}) {
	p._logCustom(StepLoggerFatal, msg, nil)
}

type StepLoggerStub struct {
	TracerId TracerId
}

func (StepLoggerStub) CanLogEvent(StepLoggerEvent, StepLogLevel) bool {
	return false
}

func (StepLoggerStub) LogUpdate(StepLoggerData, StepLoggerUpdateData) {
}

func (StepLoggerStub) LogInternal(StepLoggerData, string) {
}

func (StepLoggerStub) LogEvent(StepLoggerData, interface{}) {
}

func (v StepLoggerStub) GetTracerId() TracerId {
	return v.TracerId
}
