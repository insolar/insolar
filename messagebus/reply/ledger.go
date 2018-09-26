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
	"io"

	"github.com/insolar/insolar/core"
)

// Code is code from storage.
type Code struct {
	Code        []byte
	MachineType core.MachineType
}

// Serialize serializes to gob.
func (e *Code) Serialize() (io.Reader, error) {
	return serialize(e, TypeCode)
}

// Class is class from storage.
type Class struct {
	Head  core.RecordRef
	State core.RecordRef
	Code  *core.RecordRef // Can be nil.
}

// Serialize serializes to gob.
func (e *Class) Serialize() (io.Reader, error) {
	return serialize(e, TypeClass)
}

// Object is object from storage.
type Object struct {
	Head     core.RecordRef
	State    core.RecordRef
	Class    core.RecordRef
	Memory   []byte
	Children []core.RecordRef
}

// Serialize serializes to gob.
func (e *Object) Serialize() (io.Reader, error) {
	return serialize(e, TypeObject)
}

// Delegate is delegate reference from storage.
type Delegate struct {
	Head core.RecordRef
}

// Serialize serializes to gob.
func (e *Delegate) Serialize() (io.Reader, error) {
	return serialize(e, TypeDelegate)
}

// Reference is common reaction for methods returning reference to created records.
type Reference struct {
	Ref core.RecordRef
}

// Serialize serializes to gob.
func (e *Reference) Serialize() (io.Reader, error) {
	return serialize(e, TypeReference)
}
