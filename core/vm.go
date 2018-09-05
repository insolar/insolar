/*
 *    Copyright 2018 INS Ecosystem
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

// MachineType is a type of virtual machine
type MachineType int

// Real constants of MachineType
const (
	MachineTypeBuiltin MachineType = iota + 1
	MachineTypeGoPlugin

	MachineTypesLastID
)

// MachineLogicExecutor is an interface for implementers of one particular machine type
type MachineLogicExecutor interface {
	CallMethod(codeRef RecordRef, data []byte, method string, args Arguments) (newObjectState []byte, methodResults Arguments, err error)
	CallConstructor(codeRef RecordRef, name string, args Arguments) (objectState []byte, err error)
	Stop() error
}

// LogicRunner is an interface that should satisfy logic executor
type LogicRunner interface {
	Component
	Execute(msg Message) (res *Response)
}
