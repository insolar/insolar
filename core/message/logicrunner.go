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

type IBaseLogicMessage interface {
	core.Message
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

// CallMethod - Simply call method and return result
type CallMethod struct {
	BaseLogicMessage
	ReturnMode MethodReturnMode
	ObjectRef  core.RecordRef
	Method     string
	Arguments  core.Arguments
}

func (m *CallMethod) GetReference() core.RecordRef {
	return m.ObjectRef
}

// Type returns TypeCallMethod.
func (m *CallMethod) Type() core.MessageType {
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
	Name         string
	Arguments    core.Arguments
	PulseNum     core.PulseNumber
}

func (m *CallConstructor) GetReference() core.RecordRef {
	return *genRequest(m.PulseNum, MustSerializeBytes(m))
}

// Type returns TypeCallConstructor.
func (m *CallConstructor) Type() core.MessageType {
	return core.TypeCallConstructor
}

type ExecutorResults struct {
	Caller    core.RecordRef
	RecordRef core.RecordRef
	CaseBind  core.CaseBind
}

func (m *ExecutorResults) Type() core.MessageType {
	return core.TypeExecutorResults
}

// TODO change after changing pulsar
func (m *ExecutorResults) GetCaller() *core.RecordRef {
	return &m.Caller
}

func (m *ExecutorResults) GetReference() core.RecordRef {
	return m.RecordRef
}

type ValidateCaseBind struct {
	Caller    core.RecordRef
	RecordRef core.RecordRef
	CaseBind  core.CaseBind
	Pulse     core.Pulse
}

func (m *ValidateCaseBind) Type() core.MessageType {
	return core.TypeValidateCaseBind
}

// TODO change after changing pulsar
func (m *ValidateCaseBind) GetCaller() *core.RecordRef {
	return &m.Caller // TODO actually it's not right. There is no caller.
}

func (m *ValidateCaseBind) GetReference() core.RecordRef {
	return m.RecordRef
}

func (m *ValidateCaseBind) GetPulse() core.Pulse {
	return m.Pulse
}

type ValidationResults struct {
	Caller           core.RecordRef
	RecordRef        core.RecordRef
	PassedStepsCount int
	Error            string
}

func (m *ValidationResults) Type() core.MessageType {
	return core.TypeValidationResults
}

// TODO change after changing pulsar
func (m *ValidationResults) GetCaller() *core.RecordRef {
	return &m.Caller // TODO actually it's not right. There is no caller.
}

func (m *ValidationResults) GetReference() core.RecordRef {
	return m.RecordRef
}

var hasher = platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher() // TODO: create message factory

// GenRequest calculates RecordRef for request message from pulse number and request's payload.
func genRequest(pn core.PulseNumber, payload []byte) *core.RecordRef {
	ref := core.NewRecordRef(
		core.RecordID{},
		*core.NewRecordID(pn, hasher.Hash(payload)),
	)
	return ref
}
