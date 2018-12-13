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

// HeavyPayload carries Key/Value records and pulse number
// that replicates to Heavy Material node.
type HeavyPayload struct {
	JetID    core.RecordID
	PulseNum core.PulseNumber
	Records  []core.KV
}

// AllowedSenderObjectAndRole implements interface method
func (*HeavyPayload) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return nil, 0
}

// DefaultRole returns role for this event
func (*HeavyPayload) DefaultRole() core.DynamicRole {
	return core.DynamicRoleHeavyExecutor
}

// DefaultTarget returns of target of this event.
func (hp *HeavyPayload) DefaultTarget() *core.RecordRef {
	return &core.RecordRef{}
}

// GetCaller implementation of Message interface.
func (HeavyPayload) GetCaller() *core.RecordRef {
	return nil
}

// Type implementation of Message interface.
func (e *HeavyPayload) Type() core.MessageType {
	return core.TypeHeavyPayload
}

// HeavyStartStop carries heavy replication start/stop signal with pulse number.
type HeavyStartStop struct {
	JetID    core.RecordID
	PulseNum core.PulseNumber
	Finished bool
}

// AllowedSenderObjectAndRole implements interface method
func (*HeavyStartStop) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return nil, 0
}

// DefaultTarget returns of target of this event.
func (*HeavyStartStop) DefaultTarget() *core.RecordRef {
	return &core.RecordRef{}
}

// DefaultRole returns role for this event
func (*HeavyStartStop) DefaultRole() core.DynamicRole {
	return core.DynamicRoleHeavyExecutor
}

// GetCaller implementation of Message interface.
func (HeavyStartStop) GetCaller() *core.RecordRef {
	return nil
}

// Type implementation of Message interface.
func (e *HeavyStartStop) Type() core.MessageType {
	return core.TypeHeavyStartStop
}

// HeavyReset carries heavy replication start/stop signal with pulse number.
type HeavyReset struct {
	JetID    core.RecordID
	PulseNum core.PulseNumber
}

// AllowedSenderObjectAndRole implements interface method
func (*HeavyReset) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return nil, 0
}

// DefaultTarget returns of target of this event.
func (*HeavyReset) DefaultTarget() *core.RecordRef {
	return &core.RecordRef{}
}

// DefaultRole returns role for this event
func (*HeavyReset) DefaultRole() core.DynamicRole {
	return core.DynamicRoleHeavyExecutor
}

// GetCaller implementation of Message interface.
func (HeavyReset) GetCaller() *core.RecordRef {
	return nil
}

// Type implementation of Message interface.
func (e *HeavyReset) Type() core.MessageType {
	return core.TypeHeavyStartStop
}
