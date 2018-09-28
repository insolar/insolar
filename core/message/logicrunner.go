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

package message

import (
	"github.com/insolar/insolar/core"
)

// MethodReturnMode ENUM to set when method returns its result
type MethodReturnMode int

const (
	// ReturnResult - return result as soon as it is ready
	ReturnResult MethodReturnMode = iota
	// ReturnNoWait - call method and return without results
	ReturnNoWait
	// ReturnValidated (not yet) - return result only when it's validated
	// ReturnValidated
)

// CallMethod - Simply call method and return result
type CallMethod struct {
	BaseLogicEvent
	ReturnMode MethodReturnMode
	ObjectRef  core.RecordRef
	Request    core.RecordRef
	Method     string
	Arguments  core.Arguments
}

func (e *CallMethod) Type() core.MessageType {
	return TypeCallMethod
}

func (e *CallMethod) TargetRole() core.JetRole {
	return core.RoleVirtualExecutor
}

func (e *CallMethod) Target() *core.RecordRef {
	return &e.ObjectRef
}

// CallConstructor is a message for calling constructor and obtain its reply
type CallConstructor struct {
	BaseLogicEvent
	ClassRef  core.RecordRef
	Name      string
	Arguments core.Arguments
}

func (e *CallConstructor) Type() core.MessageType {
	return TypeCallConstructor
}

func (e *CallConstructor) Target() *core.RecordRef {
	return &e.ClassRef
}

func (e *CallConstructor) TargetRole() core.JetRole {
	return core.RoleVirtualExecutor
}
