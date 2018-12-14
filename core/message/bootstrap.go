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

// GenesisRequest is used for genesis records generation.
type GenesisRequest struct {
	// Name should be unique for each genesis record.
	Name string
}

// AllowedSenderObjectAndRole implements interface method
func (*GenesisRequest) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return nil, 0
}

// DefaultRole returns role for this event
func (*GenesisRequest) DefaultRole() core.DynamicRole {
	return core.DynamicRoleVirtualExecutor
}

// DefaultTarget returns of target of this event.
func (gr *GenesisRequest) DefaultTarget() *core.RecordRef {
	return &core.RecordRef{}
}

// Type implementation for genesis request.
func (*GenesisRequest) Type() core.MessageType {
	return core.TypeBootstrapRequest
}

// GetCaller implementation for genesis request.
func (*GenesisRequest) GetCaller() *core.RecordRef {
	return nil
}
