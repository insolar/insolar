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

package event

import (
	"io"

	"github.com/insolar/insolar/core"
)

// CallMethodEvent - Simply call method and return result
type CallMethodEvent struct {
	baseEvent
	ObjectRef core.RecordRef
	Request   core.RecordRef
	Method    string
	Arguments core.Arguments
}

// React handles event and returns associated reaction.
func (e *CallMethodEvent) React(core.Components) (core.Reaction, error) {
	panic("implement me")
}

// GetOperatingRole returns operating jet role for given event type.
func (e *CallMethodEvent) GetOperatingRole() core.JetRole {
	return core.RoleVirtualExecutor
}

// GetReference returns referenced object.
func (e *CallMethodEvent) GetReference() core.RecordRef {
	return e.ObjectRef
}

// Serialize serializes event.
func (e *CallMethodEvent) Serialize() (io.Reader, error) {
	return serialize(e, CallMethodEventType)
}
