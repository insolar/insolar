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

// GetObject retrieves object from storage.
type GetObject struct {
	ledgerMessage
	Head     core.RecordRef
	State    *core.RecordID // If nil, will fetch the latest state.
	Approved bool
}

// Type implementation of Message interface.
func (e *GetObject) Type() core.MessageType {
	return core.TypeGetObject
}

// Target implementation of Message interface.
func (e *GetObject) Target() *core.RecordRef {
	return &e.Head
}

// GetDelegate retrieves object represented as provided type.
type GetDelegate struct {
	ledgerMessage
	Head   core.RecordRef
	AsType core.RecordRef
}

// Type implementation of Message interface.
func (e *GetDelegate) Type() core.MessageType {
	return core.TypeGetDelegate
}

// Target implementation of Message interface.
func (e *GetDelegate) Target() *core.RecordRef {
	return &e.Head
}

// UpdateObject amends object.
type UpdateObject struct {
	ledgerMessage

	Record []byte
	Object core.RecordRef
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
	Record []byte
	Parent core.RecordRef
	Child  core.RecordRef
	AsType *core.RecordRef // If not nil, considered as delegate.
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

// GetHistory retrieves a chunk of history references.
type GetHistory struct {
	ledgerMessage
	Object core.RecordRef
	From   *core.RecordID
	Pulse  *core.PulseNumber
	Amount int
}

// Type implementation of Message interface.
func (e *GetHistory) Type() core.MessageType {
	return core.TypeGetHistory
}

// Target implementation of Message interface.
func (e *GetHistory) Target() *core.RecordRef {
	return &e.Object
}

// JetDrop spreads jet drop
type JetDrop struct {
	ledgerMessage
	Jet         core.RecordRef
	Drop        []byte
	Messages    [][]byte
	PulseNumber core.PulseNumber
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

// ValidateRecord creates VM validation for specific object record.
type ValidateRecord struct {
	ledgerMessage

	Object             core.RecordRef
	State              core.RecordID
	IsValid            bool
	ValidationMessages []core.Message
}

// Type implementation of Message interface.
func (*ValidateRecord) Type() core.MessageType {
	return core.TypeValidateRecord
}

// Target implementation of Message interface.
func (m *ValidateRecord) Target() *core.RecordRef {
	return &m.Object
}

// TargetRole implementation of Message interface.
func (*ValidateRecord) TargetRole() core.JetRole {
	return core.RoleLightExecutor
}

// SetBlob saves blob in storage.
type SetBlob struct {
	ledgerMessage

	TargetRef core.RecordRef
	Memory    []byte
}

// Type implementation of Message interface.
func (*SetBlob) Type() core.MessageType {
	return core.TypeSetBlob
}

// Target implementation of Message interface.
func (m *SetBlob) Target() *core.RecordRef {
	return &m.TargetRef
}
