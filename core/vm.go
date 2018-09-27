/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package core

import (
	"time"
)

// MachineType is a type of virtual machine
type MachineType int

// Real constants of MachineType
const (
	MachineTypeNotExist             = 0
	MachineTypeBuiltin  MachineType = iota + 1
	MachineTypeGoPlugin

	MachineTypesLastID
)

// MachineLogicExecutor is an interface for implementers of one particular machine type
type MachineLogicExecutor interface {
	CallMethod(ctx *LogicCallContext, code RecordRef, data []byte, method string, args Arguments) (newObjectState []byte, methodResults Arguments, err error)
	CallConstructor(ctx *LogicCallContext, code RecordRef, name string, args Arguments) (objectState []byte, err error)
	Stop() error
}

// LogicRunner is an interface that should satisfy logic executor
type LogicRunner interface {
	Execute(event LogicRunnerEvent) (res Reaction, err error)
	OnPulse(pulse Pulse) error
}

// LogicCallContext is a context of contract execution
type LogicCallContext struct {
	Callee *RecordRef // Contract that was called
	Class  *RecordRef // Class of the callee
	Parent *RecordRef // Parent of the callee
	Caller *RecordRef // Contract that made the call
	Time   time.Time  // Time when call was made
	Pulse  Pulse      // Number of the pulse
}
