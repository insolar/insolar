///
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
///

package smachine

type StepFlags uint16

const (
	StepResetAllFlags StepFlags = 1 << iota
	StepWeak
	StepPriority
	//StepIgnoreAsyncWakeup
	//StepForceAsyncWakeup
	//StepIgnoreAsyncPanic
)

type SlotStep struct {
	Transition StateFunc
	Migration  MigrateFunc
	Flags      StepFlags
	Handler    ErrorHandlerFunc
}

func (s *SlotStep) IsZero() bool {
	return s.Transition == nil && s.Flags == 0 && s.Migration == nil && s.Handler == nil
}

func (s *SlotStep) ensureTransition() {
	if s.Transition == nil {
		panic("illegal value")
	}
}
