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
	StepLoggerDetached StepLoggerUpdateFlags = 1 << iota
)

type StepLoggerData struct {
	CycleNo     uint32
	StepNo      StepLink
	CurrentStep SlotStep
	Declaration StateMachineDeclaration
	EventType   StepLoggerEvent
	Error       error
}

type StepLoggerUpdateData struct {
	UpdateType string
	NextStep   SlotStep
	Flags      StepLoggerUpdateFlags
}

type StepLoggerFactoryFunc func(context.Context, StateMachine, TracerId) StepLogger

type StepLogger interface {
	CanLogEvent(eventType StepLoggerEvent, isTracing bool) bool
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
		isTracing := p.c.s.isTracing()
		if !stepLogger.CanLogEvent(eventType, isTracing) {
			return
		}
		s := p.c.s

		if isTracing && eventType == StepLoggerTrace {
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

func (StepLoggerStub) CanLogEvent(StepLoggerEvent, bool) bool {
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
