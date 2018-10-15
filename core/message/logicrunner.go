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

type IBaseLogicMessage interface {
	core.Message
	GetReference() core.RecordRef
}

// BaseLogicMessage base of event class family, do not use it standalone
type BaseLogicMessage struct {
	Caller core.RecordRef
	Nonce  uint64
	sign   []byte
}

// SetSign sets a signature to message.
func (b *BaseLogicMessage) SetSign(sign []byte) {
	b.sign = sign
}

// GetSign returns a sign.
func (b *BaseLogicMessage) GetSign() []byte {
	return b.sign
}

func (e *BaseLogicMessage) GetCaller() *core.RecordRef {
	return &e.Caller
}

// TargetRole returns RoleVirtualExecutor as routing target role.
func (e *BaseLogicMessage) TargetRole() core.JetRole {
	return core.RoleVirtualExecutor
}

// CallMethod - Simply call method and return result
type CallMethod struct {
	BaseLogicMessage
	ReturnMode MethodReturnMode
	ObjectRef  core.RecordRef
	Method     string
	Arguments  core.Arguments
}

func (e *CallMethod) GetReference() core.RecordRef {
	return e.ObjectRef
}

// Type returns TypeCallMethod.
func (e *CallMethod) Type() core.MessageType {
	return core.TypeCallMethod
}

// Target returns ObjectRef as routing target.
func (e *CallMethod) Target() *core.RecordRef {
	return &e.ObjectRef
}

type SaveAs int

const (
	Child SaveAs = iota
	Delegate
)

// CallConstructor is a message for calling constructor and obtain its reply
type CallConstructor struct {
	BaseLogicMessage
	ParentRef core.RecordRef
	SaveAs    SaveAs
	ClassRef  core.RecordRef
	Name      string
	Arguments core.Arguments
	PulseNum  core.PulseNumber
}

func (e *CallConstructor) GetReference() core.RecordRef {
	return e.ClassRef
}

// Type returns TypeCallConstructor.
func (e *CallConstructor) Type() core.MessageType {
	return core.TypeCallConstructor
}

// Target returns request ref as routing target.
func (e *CallConstructor) Target() *core.RecordRef {
	if e.SaveAs == Delegate {
		return &e.ParentRef
	}
	return core.GenRequest(e.PulseNum, MustSerializeBytes(e))
}

type ExecutorResults struct {
	RecordRef   core.RecordRef
	CaseRecords []core.CaseRecord
	sign        []byte
}

func (e *ExecutorResults) Type() core.MessageType {
	return core.TypeExecutorResults
}

func (e *ExecutorResults) TargetRole() core.JetRole {
	return core.RoleVirtualExecutor
}

func (e *ExecutorResults) Target() *core.RecordRef {
	return &e.RecordRef
}

// TODO change after changing pulsar
func (e *ExecutorResults) GetCaller() *core.RecordRef {
	return &core.RecordRef{}
}

func (e *ExecutorResults) GetReference() core.RecordRef {
	return e.RecordRef
}

func (e *ExecutorResults) GetSign() []byte {
	return e.sign
}

func (e *ExecutorResults) SetSign(sign []byte) {
	e.sign = sign
}

type IValidateCaseBind interface {
	GetReference() core.RecordRef
	GetCaseRecords() []core.CaseRecord
	GetPulse() core.Pulse
	GetSign() []byte
	SetSign(sign []byte)
}

type ValidateCaseBind struct {
	RecordRef   core.RecordRef
	CaseRecords []core.CaseRecord
	Pulse       core.Pulse
	sign        []byte
}

func (e *ValidateCaseBind) Type() core.MessageType {
	return core.TypeValidateCaseBind
}

func (e *ValidateCaseBind) TargetRole() core.JetRole {
	return core.RoleVirtualValidator
}

// TODO change after changing pulsar
func (e *ValidateCaseBind) Target() *core.RecordRef {
	return &e.RecordRef
}

// TODO change after changing pulsar
// need to implement core.Message interface
func (e *ValidateCaseBind) GetCaller() *core.RecordRef {
	return &e.RecordRef // TODO actually it's not right. There is no caller.
}

// TODO change after changing pulsar
func (e *ValidateCaseBind) GetReference() core.RecordRef {
	return e.RecordRef
}

func (e *ValidateCaseBind) GetCaseRecords() []core.CaseRecord {
	return e.CaseRecords
}

func (e *ValidateCaseBind) GetPulse() core.Pulse {
	return e.Pulse
}

func (e *ValidateCaseBind) GetSign() []byte {
	return e.sign
}

func (e *ValidateCaseBind) SetSign(sign []byte) {
	e.sign = sign
}

type ValidationResults struct {
	RecordRef        core.RecordRef
	PassedStepsCount int
	Error            error
	sign             []byte
}

func (e *ValidationResults) Type() core.MessageType {
	return core.TypeValidationResults
}

func (e ValidationResults) TargetRole() core.JetRole {
	return core.RoleVirtualExecutor
}

func (e *ValidationResults) Target() *core.RecordRef {
	return &e.RecordRef
}

// TODO change after changing pulsar
// need to implement core.Message interface
func (e *ValidationResults) GetCaller() *core.RecordRef {
	return &e.RecordRef // TODO actually it's not right. There is no caller.
}

func (e *ValidationResults) GetReference() core.RecordRef {
	return e.RecordRef
}

func (e *ValidationResults) GetSign() []byte {
	return e.sign
}

func (e *ValidationResults) SetSign(sign []byte) {
	e.sign = sign
}
