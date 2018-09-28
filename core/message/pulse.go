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

// Pulse is a message type for pulsar.
type Pulse struct {
	Pulse core.Pulse
}

// Type returns message type.
func (p *Pulse) Type() core.MessageType {
	return TypePulse
}

// Target returns nil to send for all actors for the role returned by Role method.
func (p *Pulse) Target() *core.RecordRef {
	return nil
}

// TargetRole returns jet role to actors of which Message should be sent.
func (p *Pulse) TargetRole() core.JetRole {
	return core.RoleAllRoles
}
