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

type NodeSignPayloadInt interface {
	core.Message
	GetNodeRef() *core.RecordRef
}

type NodeSignPayload struct {
	NodeRef *core.RecordRef
}

// GetCaller implementation of Message interface.
func (NodeSignPayload) GetCaller() *core.RecordRef {
	return nil
}

// Type implementation of Message interface.
func (e *NodeSignPayload) Type() core.MessageType {
	return core.NetworkCoordinatorNodeSignRequest
}

func (e *NodeSignPayload) GetNodeRef() *core.RecordRef {
	return e.NodeRef
}
