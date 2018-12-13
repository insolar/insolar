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

// Object is object from storage.
type Object struct {
	Head         core.RecordRef
	State        core.RecordID
	Prototype    *core.RecordRef
	IsPrototype  bool
	ChildPointer *core.RecordID
	Memory       []byte
	Parent       core.RecordRef
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

// ID is common reaction for methods returning id to lifeline states.
type ID struct {
	ID core.RecordID
}

// Type implementation of Reply interface.
func (e *ID) Type() core.ReplyType {
	return TypeID
}

// Children is common reaction for methods returning id to lifeline states.
type Children struct {
	Refs     []core.RecordRef
	NextFrom *core.RecordID
}

// Type implementation of Reply interface.
func (e *Children) Type() core.ReplyType {
	return TypeChildren
}

// ObjectIndex contains serialized object index. It can be stored in DB without processing.
type ObjectIndex struct {
	Index []byte
}

// Type implementation of Reply interface.
func (e *ObjectIndex) Type() core.ReplyType {
	return TypeObjectIndex
}

// JetMiss is returned for miscalculated jets due to incomplete jet tree.
type JetMiss struct {
	JetID core.RecordID
}

// Type implementation of Reply interface.
func (e *JetMiss) Type() core.ReplyType {
	return TypeJetMiss
}
