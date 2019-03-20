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
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/platformpolicy"
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

type PendingState int

const (
	PendingUnknown PendingState = iota
	NotPending
	InPending
)

type IBaseLogicMessage interface {
	core.Message
	GetBaseLogicMessage() *BaseLogicMessage
	GetReference() core.RecordRef
	GetRequest() core.RecordRef
	GetCallerPrototype() *core.RecordRef
}

// BaseLogicMessage base of event class family, do not use it standalone
type BaseLogicMessage struct {
	Caller          core.RecordRef
	Request         core.RecordRef
	CallerPrototype core.RecordRef
	Nonce           uint64
	Sequence        uint64
}

func (m *BaseLogicMessage) GetBaseLogicMessage() *BaseLogicMessage {
	return m
}

func (m *BaseLogicMessage) Type() core.MessageType {
	panic("Virtual")
}

func (m *BaseLogicMessage) DefaultTarget() *core.RecordRef {
	panic("Virtual")
}

func (m *BaseLogicMessage) DefaultRole() core.DynamicRole {
	panic("implement me")
}

func (m *BaseLogicMessage) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	panic("implement me")
}

func (m *BaseLogicMessage) GetReference() core.RecordRef {
	panic("implement me")
}

func (m *BaseLogicMessage) GetCaller() *core.RecordRef {
	return &m.Caller
}

func (m *BaseLogicMessage) GetCallerPrototype() *core.RecordRef {
	return &m.CallerPrototype
}

// GetRequest returns DynamicRoleVirtualExecutor as routing target role.
func (m *BaseLogicMessage) GetRequest() core.RecordRef {
	return m.Request
}

// ReturnResults - push results of methods
type ReturnResults struct {
	Target   core.RecordRef
	Caller   core.RecordRef
	Sequence uint64
	Reply    core.Reply
	Error    string
}

func (rr *ReturnResults) Type() core.MessageType {
	return core.TypeReturnResults
}

func (rr *ReturnResults) GetCaller() *core.RecordRef {
	return &rr.Caller
}

func (rr *ReturnResults) DefaultTarget() *core.RecordRef {
	return &rr.Target
}

func (rr *ReturnResults) DefaultRole() core.DynamicRole {
	return core.DynamicRoleVirtualExecutor
}

func (rr *ReturnResults) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return nil, core.DynamicRoleVirtualExecutor
}

// CallMethod - Simply call method and return result
type CallMethod struct {
	BaseLogicMessage
	ReturnMode     MethodReturnMode
	ObjectRef      core.RecordRef
	Method         string
	Arguments      core.Arguments
	ProxyPrototype core.RecordRef
}

// ToMap returns map representation of CallMethod.
// Temporary until ledger.exporter api response reorganization
func (cm *CallMethod) ToMap() (map[string]interface{}, error) {
	msg := make(map[string]interface{})

	// BaseLogicMessage fields
	msg["Caller"] = cm.BaseLogicMessage.Caller.String()
	msg["Request"] = cm.BaseLogicMessage.Request.String()
	msg["CallerPrototype"] = cm.BaseLogicMessage.CallerPrototype.String()
	msg["Nonce"] = cm.BaseLogicMessage.Nonce
	msg["Sequence"] = cm.BaseLogicMessage.Sequence

	// CallMethod fields
	msg["ReturnMode"] = cm.ReturnMode
	msg["ObjectRef"] = cm.ObjectRef.String()
	msg["Method"] = cm.Method
	msg["ProxyPrototype"] = cm.ProxyPrototype.String()
	args, err := cm.Arguments.MarshalJSON()
	if err != nil {
		msg["Arguments"] = cm.Arguments
	} else {
		msg["Arguments"] = string(args)
	}

	return msg, nil
}

// AllowedSenderObjectAndRole implements interface method
func (cm *CallMethod) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	c := cm.GetCaller()
	if c.IsEmpty() {
		return nil, 0
	}
	return c, core.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*CallMethod) DefaultRole() core.DynamicRole {
	return core.DynamicRoleVirtualExecutor
}

// DefaultTarget returns of target of this event.
func (cm *CallMethod) DefaultTarget() *core.RecordRef {
	return &cm.ObjectRef
}

func (cm *CallMethod) GetReference() core.RecordRef {
	return cm.ObjectRef
}

// Type returns TypeCallMethod.
func (cm *CallMethod) Type() core.MessageType {
	return core.TypeCallMethod
}

type SaveAs int

const (
	Child SaveAs = iota
	Delegate
)

// CallConstructor is a message for calling constructor and obtain its reply
type CallConstructor struct {
	BaseLogicMessage
	ParentRef    core.RecordRef
	SaveAs       SaveAs
	PrototypeRef core.RecordRef
	Method       string
	Arguments    core.Arguments
	PulseNum     core.PulseNumber
}

// ToMap returns map representation of CallConstructor.
// Temporary until ledger.exporter api response reorganization
func (cc *CallConstructor) ToMap() (map[string]interface{}, error) {
	msg := make(map[string]interface{})

	// BaseLogicMessage fields
	msg["Caller"] = cc.BaseLogicMessage.Caller.String()
	msg["Request"] = cc.BaseLogicMessage.Request.String()
	msg["CallerPrototype"] = cc.BaseLogicMessage.CallerPrototype.String()
	msg["Nonce"] = cc.BaseLogicMessage.Nonce
	msg["Sequence"] = cc.BaseLogicMessage.Sequence

	// CallConstructor fields
	msg["ParentRef"] = cc.ParentRef.String()
	msg["SaveAs"] = cc.SaveAs
	msg["PrototypeRef"] = cc.PrototypeRef.String()
	msg["Method"] = cc.Method
	msg["PulseNum"] = cc.PulseNum
	args, err := cc.Arguments.MarshalJSON()
	if err != nil {
		msg["Arguments"] = cc.Arguments
	} else {
		msg["Arguments"] = string(args)
	}

	return msg, nil
}

//
func (cc *CallConstructor) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	c := cc.GetCaller()
	if c.IsEmpty() {
		return nil, 0
	}
	return c, core.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*CallConstructor) DefaultRole() core.DynamicRole {
	return core.DynamicRoleVirtualExecutor
}

// DefaultTarget returns of target of this event.
func (cc *CallConstructor) DefaultTarget() *core.RecordRef {
	if cc.SaveAs == Delegate {
		return &cc.ParentRef
	}
	return genRequest(cc.PulseNum, MustSerializeBytes(cc), cc.Request.Domain())
}

func (cc *CallConstructor) GetReference() core.RecordRef {
	return *genRequest(cc.PulseNum, MustSerializeBytes(cc), cc.Request.Domain())
}

// Type returns TypeCallConstructor.
func (cc *CallConstructor) Type() core.MessageType {
	return core.TypeCallConstructor
}

// TODO rename to executorObjectResult (results?)
type ExecutorResults struct {
	Caller                core.RecordRef
	RecordRef             core.RecordRef
	Requests              []CaseBindRequest
	Queue                 []ExecutionQueueElement
	LedgerHasMoreRequests bool
	Pending               PendingState
}

type ExecutionQueueElement struct {
	Parcel  core.Parcel
	Request *core.RecordRef
}

// AllowedSenderObjectAndRole implements interface method
func (er *ExecutorResults) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	// TODO need to think - this message can send only Executor of Previous Pulse, this function
	return nil, 0
}

// DefaultRole returns role for this event
func (er *ExecutorResults) DefaultRole() core.DynamicRole {
	return core.DynamicRoleVirtualExecutor
}

// DefaultTarget returns of target of this event.
func (er *ExecutorResults) DefaultTarget() *core.RecordRef {
	return &er.RecordRef
}

func (er *ExecutorResults) Type() core.MessageType {
	return core.TypeExecutorResults
}

// TODO change after changing pulsar
func (er *ExecutorResults) GetCaller() *core.RecordRef {
	return &er.Caller
}

func (er *ExecutorResults) GetReference() core.RecordRef {
	return er.RecordRef
}

type ValidateCaseBind struct {
	Caller    core.RecordRef
	RecordRef core.RecordRef
	Requests  []CaseBindRequest
	Pulse     core.Pulse
}

type CaseBindRequest struct {
	Parcel         core.Parcel
	Request        core.RecordRef
	MessageBusTape []byte
	Reply          core.Reply
	Error          string
}

// AllowedSenderObjectAndRole implements interface method
func (vcb *ValidateCaseBind) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &vcb.RecordRef, core.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*ValidateCaseBind) DefaultRole() core.DynamicRole {
	return core.DynamicRoleVirtualValidator
}

// DefaultTarget returns of target of this event.
func (vcb *ValidateCaseBind) DefaultTarget() *core.RecordRef {
	return &vcb.RecordRef
}

func (vcb *ValidateCaseBind) Type() core.MessageType {
	return core.TypeValidateCaseBind
}

// TODO change after changing pulsar
func (vcb *ValidateCaseBind) GetCaller() *core.RecordRef {
	return &vcb.Caller // TODO actually it's not right. There is no caller.
}

func (vcb *ValidateCaseBind) GetReference() core.RecordRef {
	return vcb.RecordRef
}

func (vcb *ValidateCaseBind) GetPulse() core.Pulse {
	return vcb.Pulse
}

type ValidationResults struct {
	Caller           core.RecordRef
	RecordRef        core.RecordRef
	PassedStepsCount int
	Error            string
}

// AllowedSenderObjectAndRole implements interface method
func (vr *ValidationResults) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &vr.RecordRef, core.DynamicRoleVirtualValidator
}

// DefaultRole returns role for this event
func (*ValidationResults) DefaultRole() core.DynamicRole {
	return core.DynamicRoleVirtualExecutor
}

// DefaultTarget returns of target of this event.
func (vr *ValidationResults) DefaultTarget() *core.RecordRef {
	return &vr.RecordRef
}

func (vr *ValidationResults) Type() core.MessageType {
	return core.TypeValidationResults
}

// TODO change after changing pulsar
func (vr *ValidationResults) GetCaller() *core.RecordRef {
	return &vr.Caller // TODO actually it's not right. There is no caller.
}

func (vr *ValidationResults) GetReference() core.RecordRef {
	return vr.RecordRef
}

var hasher = platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher() // TODO: create message factory

// GenRequest calculates RecordRef for request message from pulse number and request's payload.
func genRequest(pn core.PulseNumber, payload []byte, domain *core.RecordID) *core.RecordRef {
	ref := core.NewRecordRef(
		*domain,
		*core.NewRecordID(pn, hasher.Hash(payload)),
	)
	return ref
}

// PendingFinished is sent by the old executor to the current executor
// when pending execution finishes.
type PendingFinished struct {
	Reference core.RecordRef // object pended in executor
}

func (pf *PendingFinished) GetCaller() *core.RecordRef {
	// Contract that initiated this call
	return &pf.Reference
}

func (pf *PendingFinished) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	// This type of message currently can be send from any node todo: rethink it
	return nil, 0
}

func (pf *PendingFinished) DefaultRole() core.DynamicRole {
	return core.DynamicRoleVirtualExecutor
}

func (pf *PendingFinished) DefaultTarget() *core.RecordRef {
	return &pf.Reference
}

func (pf *PendingFinished) Type() core.MessageType {
	return core.TypePendingFinished
}

// StillExecuting
type StillExecuting struct {
	Reference core.RecordRef // object we still executing
}

func (se *StillExecuting) GetCaller() *core.RecordRef {
	return &se.Reference
}

func (se *StillExecuting) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return nil, 0
}

func (se *StillExecuting) DefaultRole() core.DynamicRole {
	return core.DynamicRoleVirtualExecutor
}

func (se *StillExecuting) DefaultTarget() *core.RecordRef {
	return &se.Reference
}

func (se *StillExecuting) Type() core.MessageType {
	return core.TypeStillExecuting
}
