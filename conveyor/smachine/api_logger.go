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

import (
	"context"
	"time"
)

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

type StepLoggerFlags uint8

const (
	StepLoggerUpdateErrorDefault StepLoggerFlags = iota
	StepLoggerUpdateErrorMuted
	StepLoggerUpdateErrorRecovered
	StepLoggerUpdateErrorRecoveryDenied
)
const StepLoggerErrorMask StepLoggerFlags = 3

const (
	StepLoggerDetached StepLoggerFlags = 1 << (2 + iota)
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
	Flags       StepLoggerFlags
}

type StepLoggerUpdateData struct {
	UpdateType string
	PrevStepNo uint32
	NextStep   StepDeclaration

	InactivityNano time.Duration // zero or negative - means that value is not applicable / not valid
	ActivityNano   time.Duration // zero or negative - means that value is not applicable / not valid
}

type SlotMachineLogger interface {
	CreateStepLogger(context.Context, StateMachine, TracerId) StepLogger
	LogMachineInternal(data SlotMachineData, msg string)
	LogMachineCritical(data SlotMachineData, msg string)
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

	CreateAsyncLogger(context.Context, *StepLoggerData) (context.Context, StepLogger)
}

type StepLoggerFunc func(*StepLoggerData, *StepLoggerUpdateData)

type TracerId = string

type Logger struct { // we use an explicit struct here to enable compiler optimizations when logging is not needed
	ctx      context.Context
	loggerFn interface {
		getStepLogger() (StepLogger, StepLogLevel, uint32)
	}
	stepData StepLoggerData
}

func (p Logger) getStepLogger() (StepLogger, StepLogLevel, uint32) {
	if p.loggerFn != nil {
		return p.loggerFn.getStepLogger()
	}
	return nil, 0, 0
}

func (p Logger) GetContext() context.Context {
	_, _, _ = p.getStepLogger() // check context availability
	return p.ctx
}

func (p Logger) GetTracerId() TracerId {
	if stepLogger, _, _ := p.getStepLogger(); stepLogger != nil {
		return stepLogger.GetTracerId()
	}
	return ""
}

func (p Logger) _checkLog(eventType StepLoggerEvent) (StepLogger, uint32, StepLoggerEvent) {
	if stepLogger, stepLevel, stepUpdate := p.getStepLogger(); stepLogger != nil {
		if stepLogger.CanLogEvent(eventType, stepLevel) {
			if stepLevel == StepLogLevelTracing && eventType == StepLoggerTrace {
				eventType = StepLoggerActiveTrace
			}
			return stepLogger, stepUpdate, eventType
		}
	}
	return nil, 0, 0
}

func (p Logger) _doLog(stepLogger StepLogger, stepUpdate uint32, eventType StepLoggerEvent, msg interface{}, err error) {
	stepData := p.stepData
	stepData.EventType = eventType
	stepData.Error = err
	if stepUpdate != 0 {
		stepData.StepNo.step = stepUpdate
	}
	stepLogger.LogEvent(stepData, msg)
}

// NB! keep method simple to ensure inlining
func (p Logger) Trace(msg interface{}) {
	if stepLogger, stepUpdate, eventType := p._checkLog(StepLoggerTrace); stepLogger != nil {
		p._doLog(stepLogger, stepUpdate, eventType, msg, nil)
	}
}

// NB! keep method simple to ensure inlining
func (p Logger) Warn(msg interface{}) {
	if stepLogger, stepUpdate, eventType := p._checkLog(StepLoggerWarn); stepLogger != nil {
		p._doLog(stepLogger, stepUpdate, eventType, msg, nil)
	}
}

// NB! keep method simple to ensure inlining
func (p Logger) Error(msg interface{}, err error) {
	if stepLogger, stepUpdate, eventType := p._checkLog(StepLoggerError); stepLogger != nil {
		p._doLog(stepLogger, stepUpdate, eventType, msg, err)
	}
}

// NB! keep method simple to ensure inlining
func (p Logger) Fatal(msg interface{}) {
	if stepLogger, stepUpdate, eventType := p._checkLog(StepLoggerFatal); stepLogger != nil {
		p._doLog(stepLogger, stepUpdate, eventType, msg, nil)
	}
}
