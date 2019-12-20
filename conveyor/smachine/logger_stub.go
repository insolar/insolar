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

type StepLoggerStub struct {
	TracerId TracerId
}

func (StepLoggerStub) CanLogEvent(StepLoggerEvent, StepLogLevel) bool {
	return false
}

func (StepLoggerStub) LogUpdate(StepLoggerData, StepLoggerUpdateData) {}
func (StepLoggerStub) LogInternal(StepLoggerData, string)             {}
func (StepLoggerStub) LogEvent(StepLoggerData, interface{})           {}
func (StepLoggerStub) LogAdapter(StepLoggerData, AdapterId, uint64)   {}

func (StepLoggerStub) CreateAsyncLogger(ctx context.Context, data *StepLoggerData) (context.Context, StepLogger) {
	return ctx, nil
}

func (v StepLoggerStub) GetTracerId() TracerId {
	return v.TracerId
}

type fixedSlotLogger struct {
	logger StepLogger
	level  StepLogLevel
	data   StepLoggerData
}

func (v fixedSlotLogger) getStepLogger() (StepLogger, StepLogLevel, uint32) {
	if step, ok := v.data.StepNo.SlotLink.GetStepLink(); ok {
		return v.logger, v.level, step.StepNo()
	}
	return nil, 0, 0
}

func (v fixedSlotLogger) getStepLoggerData() StepLoggerData {
	return v.data
}
