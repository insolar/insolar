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

package reply

import (
	"github.com/insolar/insolar/core"
)

// Code is code from storage.
type Code struct {
	Code        []byte
	MachineType core.MachineType
}

// Type implementation of Reply interface.
func (e *Code) Type() core.ReplyType {
	return TypeCode
}

// Class is class from storage.
type Class struct {
	Head  core.RecordRef
	State core.RecordRef
	Code  *core.RecordRef // Can be nil.
}

// Type implementation of Reply interface.
func (e *Class) Type() core.ReplyType {
	return TypeClass
}

// Object is object from storage.
type Object struct {
	Head     core.RecordRef
	State    core.RecordRef
	Class    core.RecordRef
	Memory   []byte
	Children []core.RecordRef
}

// Type implementation of Reply interface.
func (e *Object) Type() core.ReplyType {
	return TypeObject
}

// Delegate is delegate reference from storage.
type Delegate struct {
	Head core.RecordRef
}

// Type implementation of Reply interface.
func (e *Delegate) Type() core.ReplyType {
	return TypeDelegate
}

// Reference is common reaction for methods returning reference to created records.
type Reference struct {
	Ref core.RecordRef
}

// Type implementation of Reply interface.
func (e *Reference) Type() core.ReplyType {
	return TypeReference
}

// ID is common reaction for methods returning id to lifeline states.
type ID struct {
	ID core.RecordID
}

// Type implementation of Reply interface.
func (e *ID) Type() core.ReplyType {
	return TypeReference
}
