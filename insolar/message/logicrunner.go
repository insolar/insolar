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
	insolar.Message
	GetBaseLogicMessage() *BaseLogicMessage
	GetReference() insolar.RecordRef
	GetRequest() insolar.RecordRef
	GetCallerPrototype() *insolar.RecordRef
}

// BaseLogicMessage base of event class family, do not use it standalone
type BaseLogicMessage struct {
	Caller          insolar.RecordRef
	Request         insolar.RecordRef
	CallerPrototype insolar.RecordRef
	Nonce           uint64
	Sequence        uint64
}

func (m *BaseLogicMessage) GetBaseLogicMessage() *BaseLogicMessage {
	return m
}

func (m *BaseLogicMessage) Type() insolar.MessageType {
	panic("Virtual")
}

func (m *BaseLogicMessage) DefaultTarget() *insolar.RecordRef {
	panic("Virtual")
}

func (m *BaseLogicMessage) DefaultRole() insolar.DynamicRole {
	panic("implement me")
}

func (m *BaseLogicMessage) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	panic("implement me")
}

func (m *BaseLogicMessage) GetReference() insolar.RecordRef {
	panic("implement me")
}

func (m *BaseLogicMessage) GetCaller() *insolar.RecordRef {
	return &m.Caller
}

func (m *BaseLogicMessage) GetCallerPrototype() *insolar.RecordRef {
	return &m.CallerPrototype
}

// GetRequest returns DynamicRoleVirtualExecutor as routing target role.
func (m *BaseLogicMessage) GetRequest() insolar.RecordRef {
	return m.Request
}

// ReturnResults - push results of methods
type ReturnResults struct {
	Target   insolar.RecordRef
	Caller   insolar.RecordRef
	Sequence uint64
	Reply    insolar.Reply
	Error    string
}

func (rr *ReturnResults) Type() insolar.MessageType {
	return insolar.TypeReturnResults
}

func (rr *ReturnResults) GetCaller() *insolar.RecordRef {
	return &rr.Caller
}

func (rr *ReturnResults) DefaultTarget() *insolar.RecordRef {
	return &rr.Target
}

func (rr *ReturnResults) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleVirtualExecutor
}

func (rr *ReturnResults) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return nil, insolar.DynamicRoleVirtualExecutor
}

// CallMethod - Simply call method and return result
type CallMethod struct {
	BaseLogicMessage
	ReturnMode     MethodReturnMode
	ObjectRef      insolar.RecordRef
	Method         string
	Arguments      insolar.Arguments
	ProxyPrototype insolar.RecordRef
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
func (cm *CallMethod) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	c := cm.GetCaller()
	if c.IsEmpty() {
		return nil, 0
	}
	return c, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*CallMethod) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleVirtualExecutor
}

// DefaultTarget returns of target of this event.
func (cm *CallMethod) DefaultTarget() *insolar.RecordRef {
	return &cm.ObjectRef
}

func (cm *CallMethod) GetReference() insolar.RecordRef {
	return cm.ObjectRef
}

// Type returns TypeCallMethod.
func (cm *CallMethod) Type() insolar.MessageType {
	return insolar.TypeCallMethod
}

type SaveAs int

const (
	Child SaveAs = iota
	Delegate
)

// CallConstructor is a message for calling constructor and obtain its reply
type CallConstructor struct {
	BaseLogicMessage
	ParentRef    insolar.RecordRef
	SaveAs       SaveAs
	PrototypeRef insolar.RecordRef
	Method       string
	Arguments    insolar.Arguments
	PulseNum     insolar.PulseNumber
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
func (cc *CallConstructor) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	c := cc.GetCaller()
	if c.IsEmpty() {
		return nil, 0
	}
	return c, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*CallConstructor) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleVirtualExecutor
}

// DefaultTarget returns of target of this event.
func (cc *CallConstructor) DefaultTarget() *insolar.RecordRef {
	if cc.SaveAs == Delegate {
		return &cc.ParentRef
	}
	return genRequest(cc.PulseNum, MustSerializeBytes(cc), cc.Request.Domain())
}

func (cc *CallConstructor) GetReference() insolar.RecordRef {
	return *genRequest(cc.PulseNum, MustSerializeBytes(cc), cc.Request.Domain())
}

// Type returns TypeCallConstructor.
func (cc *CallConstructor) Type() insolar.MessageType {
	return insolar.TypeCallConstructor
}

// TODO rename to executorObjectResult (results?)
type ExecutorResults struct {
	Caller                insolar.RecordRef
	RecordRef             insolar.RecordRef
	Requests              []CaseBindRequest
	Queue                 []ExecutionQueueElement
	LedgerHasMoreRequests bool
	Pending               PendingState
}

type ExecutionQueueElement struct {
	Parcel  insolar.Parcel
	Request *insolar.RecordRef
}

// AllowedSenderObjectAndRole implements interface method
func (er *ExecutorResults) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	// TODO need to think - this message can send only Executor of Previous Pulse, this function
	return nil, 0
}

// DefaultRole returns role for this event
func (er *ExecutorResults) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleVirtualExecutor
}

// DefaultTarget returns of target of this event.
func (er *ExecutorResults) DefaultTarget() *insolar.RecordRef {
	return &er.RecordRef
}

func (er *ExecutorResults) Type() insolar.MessageType {
	return insolar.TypeExecutorResults
}

// TODO change after changing pulsar
func (er *ExecutorResults) GetCaller() *insolar.RecordRef {
	return &er.Caller
}

func (er *ExecutorResults) GetReference() insolar.RecordRef {
	return er.RecordRef
}

type ValidateCaseBind struct {
	Caller    insolar.RecordRef
	RecordRef insolar.RecordRef
	Requests  []CaseBindRequest
	Pulse     insolar.Pulse
}

type CaseBindRequest struct {
	Parcel         insolar.Parcel
	Request        insolar.RecordRef
	MessageBusTape []byte
	Reply          insolar.Reply
	Error          string
}

// AllowedSenderObjectAndRole implements interface method
func (vcb *ValidateCaseBind) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return &vcb.RecordRef, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*ValidateCaseBind) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleVirtualValidator
}

// DefaultTarget returns of target of this event.
func (vcb *ValidateCaseBind) DefaultTarget() *insolar.RecordRef {
	return &vcb.RecordRef
}

func (vcb *ValidateCaseBind) Type() insolar.MessageType {
	return insolar.TypeValidateCaseBind
}

// TODO change after changing pulsar
func (vcb *ValidateCaseBind) GetCaller() *insolar.RecordRef {
	return &vcb.Caller // TODO actually it's not right. There is no caller.
}

func (vcb *ValidateCaseBind) GetReference() insolar.RecordRef {
	return vcb.RecordRef
}

func (vcb *ValidateCaseBind) GetPulse() insolar.Pulse {
	return vcb.Pulse
}

type ValidationResults struct {
	Caller           insolar.RecordRef
	RecordRef        insolar.RecordRef
	PassedStepsCount int
	Error            string
}

// AllowedSenderObjectAndRole implements interface method
func (vr *ValidationResults) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return &vr.RecordRef, insolar.DynamicRoleVirtualValidator
}

// DefaultRole returns role for this event
func (*ValidationResults) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleVirtualExecutor
}

// DefaultTarget returns of target of this event.
func (vr *ValidationResults) DefaultTarget() *insolar.RecordRef {
	return &vr.RecordRef
}

func (vr *ValidationResults) Type() insolar.MessageType {
	return insolar.TypeValidationResults
}

// TODO change after changing pulsar
func (vr *ValidationResults) GetCaller() *insolar.RecordRef {
	return &vr.Caller // TODO actually it's not right. There is no caller.
}

func (vr *ValidationResults) GetReference() insolar.RecordRef {
	return vr.RecordRef
}

var hasher = platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher() // TODO: create message factory

// GenRequest calculates RecordRef for request message from pulse number and request's payload.
func genRequest(pn insolar.PulseNumber, payload []byte, domain *insolar.RecordID) *insolar.RecordRef {
	ref := insolar.NewRecordRef(
		*domain,
		*insolar.NewRecordID(pn, hasher.Hash(payload)),
	)
	return ref
}

// PendingFinished is sent by the old executor to the current executor
// when pending execution finishes.
type PendingFinished struct {
	Reference insolar.RecordRef // object pended in executor
}

func (pf *PendingFinished) GetCaller() *insolar.RecordRef {
	// Contract that initiated this call
	return &pf.Reference
}

func (pf *PendingFinished) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	// This type of message currently can be send from any node todo: rethink it
	return nil, 0
}

func (pf *PendingFinished) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleVirtualExecutor
}

func (pf *PendingFinished) DefaultTarget() *insolar.RecordRef {
	return &pf.Reference
}

func (pf *PendingFinished) Type() insolar.MessageType {
	return insolar.TypePendingFinished
}

// StillExecuting
type StillExecuting struct {
	Reference insolar.RecordRef // object we still executing
}

func (se *StillExecuting) GetCaller() *insolar.RecordRef {
	return &se.Reference
}

func (se *StillExecuting) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return nil, 0
}

func (se *StillExecuting) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleVirtualExecutor
}

func (se *StillExecuting) DefaultTarget() *insolar.RecordRef {
	return &se.Reference
}

func (se *StillExecuting) Type() insolar.MessageType {
	return insolar.TypeStillExecuting
}
