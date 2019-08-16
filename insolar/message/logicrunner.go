//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package message

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/platformpolicy"
)

// ReturnResults - push results of methods
type ReturnResults struct {
	Target     insolar.Reference
	RequestRef insolar.Reference
	Reason     insolar.Reference
	Reply      insolar.Reply
	Error      string
}

func (rr *ReturnResults) Type() insolar.MessageType {
	return insolar.TypeReturnResults
}

func (rr *ReturnResults) GetCaller() *insolar.Reference {
	return nil
}

func (rr *ReturnResults) DefaultTarget() *insolar.Reference {
	return &rr.Target
}

func (rr *ReturnResults) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleVirtualExecutor
}

func (rr *ReturnResults) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return nil, insolar.DynamicRoleVirtualExecutor
}

// CallMethod - Simply call method and return result
type CallMethod struct {
	record.IncomingRequest

	PulseNum insolar.PulseNumber // DIRTY: EVIL: HACK
}

func (cm *CallMethod) GetCaller() *insolar.Reference {
	return &cm.Caller
}

// AllowedSenderObjectAndRole implements interface method
func (cm *CallMethod) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	c := cm.Caller
	if c.IsEmpty() {
		return nil, 0
	}
	return &c, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*CallMethod) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleVirtualExecutor
}

var pcs = platformpolicy.NewPlatformCryptographyScheme() // TODO: create message factory

// DefaultTarget returns of target of this event.
func (cm *CallMethod) DefaultTarget() *insolar.Reference {
	return record.CalculateRequestAffinityRef(&cm.IncomingRequest, cm.PulseNum, pcs)
}

func (cm *CallMethod) GetReference() insolar.Reference {
	return *record.CalculateRequestAffinityRef(&cm.IncomingRequest, cm.PulseNum, pcs)
}

// Type returns TypeCallMethod.
func (cm *CallMethod) Type() insolar.MessageType {
	return insolar.TypeCallMethod
}

type ValidationResults struct {
	Caller           insolar.Reference
	RecordRef        insolar.Reference
	PassedStepsCount int
	Error            string
}

// AllowedSenderObjectAndRole implements interface method
func (vr *ValidationResults) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return &vr.RecordRef, insolar.DynamicRoleVirtualValidator
}

// DefaultRole returns role for this event
func (*ValidationResults) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleVirtualExecutor
}

// DefaultTarget returns of target of this event.
func (vr *ValidationResults) DefaultTarget() *insolar.Reference {
	return &vr.RecordRef
}

func (vr *ValidationResults) Type() insolar.MessageType {
	return insolar.TypeValidationResults
}

// TODO change after changing pulsar
func (vr *ValidationResults) GetCaller() *insolar.Reference {
	return &vr.Caller // TODO actually it's not right. There is no caller.
}

func (vr *ValidationResults) GetReference() insolar.Reference {
	return vr.RecordRef
}

// AdditionalCallFromPreviousExecutor is sent to the current executor
// by previous executor when Flow cancels after registering the request
// but before adding the request to the execution queue. For this reason
// this one request may be invisible by OnPulse handler. See HandleCall
// for more details.
type AdditionalCallFromPreviousExecutor struct {
	ObjectReference insolar.Reference
	// pending in this is message can be used to say: "this is the only
	// incomplete request on the object, you can start execution".
	// However, we lost this ability.
	Pending     insolar.PendingState
	RequestRef  insolar.Reference
	Request     record.IncomingRequest
	ServiceData ServiceData
}

func (m *AdditionalCallFromPreviousExecutor) GetCaller() *insolar.Reference {
	// Contract that initiated this call
	return &m.ObjectReference
}

func (m *AdditionalCallFromPreviousExecutor) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	// This type of message currently can be send from any node
	return nil, 0
}

func (m *AdditionalCallFromPreviousExecutor) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleVirtualExecutor
}

func (m *AdditionalCallFromPreviousExecutor) DefaultTarget() *insolar.Reference {
	return &m.ObjectReference
}

func (m *AdditionalCallFromPreviousExecutor) Type() insolar.MessageType {
	return insolar.TypeAdditionalCallFromPreviousExecutor
}
