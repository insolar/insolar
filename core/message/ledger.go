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

type ledgerMessage struct {
}

// GetCaller implementation of Message interface.
func (ledgerMessage) GetCaller() *core.RecordRef {
	return nil
}

// TargetRole implementation of Message interface.
func (ledgerMessage) TargetRole() core.JetRole {
	return core.RoleLightExecutor
}

// SetRecord saves record in storage.
type SetRecord struct {
	ledgerMessage

	Record    []byte
	TargetRef core.RecordRef
}

// Type implementation of Message interface.
func (e *SetRecord) Type() core.MessageType {
	return core.TypeSetRecord
}

// Target implementation of Message interface.
func (e *SetRecord) Target() *core.RecordRef {
	return &e.TargetRef
}

// GetCode retrieves code from storage.
type GetCode struct {
	ledgerMessage
	Code core.RecordRef
}

// Type implementation of Message interface.
func (e *GetCode) Type() core.MessageType {
	return core.TypeGetCode
}

// Target implementation of Message interface.
func (e *GetCode) Target() *core.RecordRef {
	return &e.Code
}

// GetClass retrieves class from storage.
type GetClass struct {
	ledgerMessage
	Head  core.RecordRef
	State *core.RecordRef // If nil, will fetch the latest state.
}

// Type implementation of Message interface.
func (e *GetClass) Type() core.MessageType {
	return core.TypeGetClass
}

// Target implementation of Message interface.
func (e *GetClass) Target() *core.RecordRef {
	return &e.Head
}

// GetObject retrieves object from storage.
type GetObject struct {
	ledgerMessage
	Head  core.RecordRef
	State *core.RecordRef // If nil, will fetch the latest state.
}

// Type implementation of Message interface.
func (e *GetObject) Type() core.MessageType {
	return core.TypeGetObject
}

// Target implementation of Message interface.
func (e *GetObject) Target() *core.RecordRef {
	return &e.Head
}

// GetDelegate retrieves object represented as provided class.
type GetDelegate struct {
	ledgerMessage
	Head    core.RecordRef
	AsClass core.RecordRef
}

// Type implementation of Message interface.
func (e *GetDelegate) Type() core.MessageType {
	return core.TypeGetDelegate
}

// Target implementation of Message interface.
func (e *GetDelegate) Target() *core.RecordRef {
	return &e.Head
}

// DeclareType creates new type.
type DeclareType struct {
	ledgerMessage
	Domain  core.RecordRef
	Request core.RecordRef
	TypeDec []byte
}

// Type implementation of Message interface.
func (e *DeclareType) Type() core.MessageType {
	return core.TypeDeclareType
}

// Target implementation of Message interface.
func (e *DeclareType) Target() *core.RecordRef {
	return &e.Request
}

// DeployCode creates new code.
type DeployCode struct {
	ledgerMessage
	Domain      core.RecordRef
	Request     core.RecordRef
	Code        []byte
	MachineType core.MachineType
}

// Type implementation of Message interface.
func (e *DeployCode) Type() core.MessageType {
	return core.TypeDeployCode
}

// Target implementation of Message interface.
func (e *DeployCode) Target() *core.RecordRef {
	return &e.Request
}

// ActivateClass activates class.
type ActivateClass struct {
	ledgerMessage
	Domain  core.RecordRef
	Request core.RecordRef
	Code    core.RecordRef
}

// Type implementation of Message interface.
func (e *ActivateClass) Type() core.MessageType {
	return core.TypeActivateClass
}

// Target implementation of Message interface.
func (e *ActivateClass) Target() *core.RecordRef {
	return &e.Code
}

// DeactivateClass deactivates class.
type DeactivateClass struct {
	ledgerMessage
	Domain  core.RecordRef
	Request core.RecordRef
	Class   core.RecordRef
}

// Type implementation of Message interface.
func (e *DeactivateClass) Type() core.MessageType {
	return core.TypeDeactivateClass
}

// Target implementation of Message interface.
func (e *DeactivateClass) Target() *core.RecordRef {
	return &e.Class
}

// UpdateClass amends class.
type UpdateClass struct {
	ledgerMessage

	Record []byte
	Class  core.RecordRef
}

// Type implementation of Message interface.
func (e *UpdateClass) Type() core.MessageType {
	return core.TypeUpdateClass
}

// Target implementation of Message interface.
func (e *UpdateClass) Target() *core.RecordRef {
	return &e.Class
}

// ActivateObject activates object.
type ActivateObject struct {
	ledgerMessage
	Domain  core.RecordRef
	Request core.RecordRef
	Class   core.RecordRef
	Parent  core.RecordRef
	Memory  []byte
}

// Type implementation of Message interface.
func (e *ActivateObject) Type() core.MessageType {
	return core.TypeActivateObject
}

// Target implementation of Message interface.
func (e *ActivateObject) Target() *core.RecordRef {
	return &e.Class
}

// ActivateObjectDelegate similar to ActivateObjType but it creates object as parent's delegate of provided class.
type ActivateObjectDelegate struct {
	ledgerMessage
	Domain  core.RecordRef
	Request core.RecordRef
	Class   core.RecordRef
	Parent  core.RecordRef
	Memory  []byte
}

// Type implementation of Message interface.
func (e *ActivateObjectDelegate) Type() core.MessageType {
	return core.TypeActivateObjectDelegate
}

// Target implementation of Message interface.
func (e *ActivateObjectDelegate) Target() *core.RecordRef {
	return &e.Class
}

// DeactivateObject deactivates object.
type DeactivateObject struct {
	ledgerMessage
	Domain  core.RecordRef
	Request core.RecordRef
	Object  core.RecordRef
}

// Type implementation of Message interface.
func (e *DeactivateObject) Type() core.MessageType {
	return core.TypeDeactivateObject
}

// Target implementation of Message interface.
func (e *DeactivateObject) Target() *core.RecordRef {
	return &e.Object
}

// UpdateObject amends object.
type UpdateObject struct {
	ledgerMessage
	Domain  core.RecordRef
	Request core.RecordRef
	Object  core.RecordRef
	Memory  []byte
}

// Type implementation of Message interface.
func (e *UpdateObject) Type() core.MessageType {
	return core.TypeUpdateObject
}

// Target implementation of Message interface.
func (e *UpdateObject) Target() *core.RecordRef {
	return &e.Object
}

// RegisterChild amends object.
type RegisterChild struct {
	ledgerMessage
	Parent core.RecordRef
	Child  core.RecordRef
}

// Type implementation of Message interface.
func (e *RegisterChild) Type() core.MessageType {
	return core.TypeRegisterChild
}

// Target implementation of Message interface.
func (e *RegisterChild) Target() *core.RecordRef {
	return &e.Parent
}

// RequestCall is a Ledger's message wrapping logicrunner's Call messages.
type RequestCall struct {
	core.Message
}

// TargetRole implementation of Message interface.
func (*RequestCall) TargetRole() core.JetRole {
	return core.RoleLightExecutor
}

// Type implementation of Message interface.
func (*RequestCall) Type() core.MessageType {
	return core.TypeRequestCall
}

// GetChildren retrieves a chunk of children references.
type GetChildren struct {
	ledgerMessage
	Parent    core.RecordRef
	FromChild *core.RecordID
	FromPulse *core.PulseNumber
	Amount    int
}

// Type implementation of Message interface.
func (e *GetChildren) Type() core.MessageType {
	return core.TypeGetChildren
}

// Target implementation of Message interface.
func (e *GetChildren) Target() *core.RecordRef {
	return &e.Parent
}

// JetDrop spreads jet drop
type JetDrop struct {
	ledgerMessage
	Jet     core.RecordRef
	Drop    []byte
	Records [][2][]byte
	Indexes [][2][]byte
}

// Type implementation of Message interface.
func (e *JetDrop) Type() core.MessageType {
	return core.TypeJetDrop
}

// Target implementation of Message interface.
func (e *JetDrop) Target() *core.RecordRef {
	return &e.Jet
}

// TargetRole implementation of Message interface.
func (JetDrop) TargetRole() core.JetRole {
	return core.RoleLightValidator
}
