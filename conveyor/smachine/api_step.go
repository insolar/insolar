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
	// Indicates that Slot's default flags (set by SetDefaultFlags()) will be ignored, otherwise ORed.
	StepResetAllFlags StepFlags = 1 << iota

	// When SM is at a step that StepWeak flag, then SM is considered as "weak".
	// SlotMachine will delete all "weak" SMs when there are no "non-weak" or working SMs left.
	StepWeak

	// A step with StepPriority flag will be executed before other steps in a cycle.
	StepPriority

	// A marker for logger to log this step without tracing
	StepElevatedLog

	//StepIgnoreAsyncWakeup
	//StepForceAsyncWakeup
)

// Describes a step of a SM
type SlotStep struct {
	// Function to be called when the step is executed. MUST NOT be nil
	Transition StateFunc

	// Function to be called for migration of this step. Overrides SetDefaultMigration() when not nil.
	Migration MigrateFunc

	// Step will be executed with the given flags. When StepResetAllFlags is specified, then SetDefaultFlags() is ignored, otherwise ORed.
	Flags StepFlags

	// Function to be called to handler errors of this step. Overrides SetDefaultErrorHandler() when not nil.
	Handler ErrorHandlerFunc
}

func (s *SlotStep) IsZero() bool {
	return s.Transition == nil && s.Flags == 0 && s.Migration == nil && s.Handler == nil
}

func (s *SlotStep) ensureTransition() {
	if s.Transition == nil {
		panic("illegal value")
	}
}
