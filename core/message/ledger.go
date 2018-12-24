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
	"github.com/insolar/insolar/ledger/storage/jet"
)

// FIXME: @andreyromancev. 21.12.18. Remove this and create 'LogicRunnerMessage' interface to get rid of 'GetCaller' in ledger.
type ledgerMessage struct {
}

// GetCaller implementation of Message interface.
func (ledgerMessage) GetCaller() *core.RecordRef {
	return nil
}

// SetRecord saves record in storage.
type SetRecord struct {
	ledgerMessage

	Record    []byte
	TargetRef core.RecordRef
}

// AllowedSenderObjectAndRole implements interface method
func (m *SetRecord) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &m.TargetRef, core.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*SetRecord) DefaultRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *SetRecord) DefaultTarget() *core.RecordRef {
	return &m.TargetRef
}

// Type implementation of Message interface.
func (m *SetRecord) Type() core.MessageType {
	return core.TypeSetRecord
}

// GetCode retrieves code From storage.
type GetCode struct {
	ledgerMessage
	Code core.RecordRef
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetCode) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &m.Code, core.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*GetCode) DefaultRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetCode) DefaultTarget() *core.RecordRef {
	return &m.Code
}

// Type implementation of Message interface.
func (*GetCode) Type() core.MessageType {
	return core.TypeGetCode
}

// GetObject retrieves object From storage.
type GetObject struct {
	ledgerMessage
	Head     core.RecordRef
	State    *core.RecordID // If nil, will fetch the latest state.
	Approved bool
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetObject) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &m.Head, core.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*GetObject) DefaultRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetObject) DefaultTarget() *core.RecordRef {
	return &m.Head
}

// Type implementation of Message interface.
func (*GetObject) Type() core.MessageType {
	return core.TypeGetObject
}

// GetDelegate retrieves object represented as provided type.
type GetDelegate struct {
	ledgerMessage
	Head   core.RecordRef
	AsType core.RecordRef
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetDelegate) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &m.Head, core.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*GetDelegate) DefaultRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetDelegate) DefaultTarget() *core.RecordRef {
	return &m.Head
}

// Type implementation of Message interface.
func (*GetDelegate) Type() core.MessageType {
	return core.TypeGetDelegate
}

// UpdateObject amends object.
type UpdateObject struct {
	ledgerMessage

	Record []byte
	Object core.RecordRef
}

// AllowedSenderObjectAndRole implements interface method
func (m *UpdateObject) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &m.Object, core.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*UpdateObject) DefaultRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *UpdateObject) DefaultTarget() *core.RecordRef {
	return &m.Object
}

// Type implementation of Message interface.
func (*UpdateObject) Type() core.MessageType {
	return core.TypeUpdateObject
}

// RegisterChild amends object.
type RegisterChild struct {
	ledgerMessage
	Record []byte
	Parent core.RecordRef
	Child  core.RecordRef
	AsType *core.RecordRef // If not nil, considered as delegate.
}

// AllowedSenderObjectAndRole implements interface method
func (m *RegisterChild) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &m.Child, core.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*RegisterChild) DefaultRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *RegisterChild) DefaultTarget() *core.RecordRef {
	return &m.Parent
}

// Type implementation of Message interface.
func (*RegisterChild) Type() core.MessageType {
	return core.TypeRegisterChild
}

// GetChildren retrieves a chunk of children references.
type GetChildren struct {
	ledgerMessage
	Parent    core.RecordRef
	FromChild *core.RecordID
	FromPulse *core.PulseNumber
	Amount    int
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetChildren) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &m.Parent, core.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*GetChildren) DefaultRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetChildren) DefaultTarget() *core.RecordRef {
	return &m.Parent
}

// Type implementation of Message interface.
func (*GetChildren) Type() core.MessageType {
	return core.TypeGetChildren
}

// JetDrop spreads jet drop
type JetDrop struct {
	ledgerMessage

	JetID core.RecordID

	Drop        []byte
	Messages    [][]byte
	PulseNumber core.PulseNumber
}

// AllowedSenderObjectAndRole implements interface method
func (m *JetDrop) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	// This check is not needed, because JetDrop sender is explicitly checked in handler.
	return nil, core.DynamicRoleUndefined
}

// DefaultRole returns role for this event
func (*JetDrop) DefaultRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *JetDrop) DefaultTarget() *core.RecordRef {
	return core.NewRecordRef(core.RecordID{}, m.JetID)
}

// Type implementation of Message interface.
func (*JetDrop) Type() core.MessageType {
	return core.TypeJetDrop
}

// ValidateRecord creates VM validation for specific object record.
type ValidateRecord struct {
	ledgerMessage

	Object             core.RecordRef
	State              core.RecordID
	IsValid            bool
	ValidationMessages []core.Message
}

// AllowedSenderObjectAndRole implements interface method
func (m *ValidateRecord) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &m.Object, core.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*ValidateRecord) DefaultRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *ValidateRecord) DefaultTarget() *core.RecordRef {
	return &m.Object
}

// Type implementation of Message interface.
func (*ValidateRecord) Type() core.MessageType {
	return core.TypeValidateRecord
}

// SetBlob saves blob in storage.
type SetBlob struct {
	ledgerMessage

	TargetRef core.RecordRef
	Memory    []byte
}

// AllowedSenderObjectAndRole implements interface method
func (m *SetBlob) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &m.TargetRef, core.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*SetBlob) DefaultRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *SetBlob) DefaultTarget() *core.RecordRef {
	return &m.TargetRef
}

// Type implementation of Message interface.
func (*SetBlob) Type() core.MessageType {
	return core.TypeSetBlob
}

// GetObjectIndex fetches objects index.
type GetObjectIndex struct {
	ledgerMessage

	Object core.RecordRef
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetObjectIndex) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &m.Object, core.DynamicRoleLightExecutor
}

// DefaultRole returns role for this event
func (*GetObjectIndex) DefaultRole() core.DynamicRole {
	return core.DynamicRoleHeavyExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetObjectIndex) DefaultTarget() *core.RecordRef {
	return &m.Object
}

// Type implementation of Message interface.
func (*GetObjectIndex) Type() core.MessageType {
	return core.TypeGetObjectIndex
}

// ValidationCheck checks if validation of a particular record can be performed.
type ValidationCheck struct {
	ledgerMessage

	Object              core.RecordRef
	ValidatedState      core.RecordID
	LatestStateApproved *core.RecordID
}

// DefaultTarget returns of target of this event.
func (m *ValidationCheck) DefaultTarget() *core.RecordRef {
	return &m.Object
}

// DefaultRole returns role for this event
func (m *ValidationCheck) DefaultRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// AllowedSenderObjectAndRole implements interface method
func (m *ValidationCheck) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	// TODO: return smth real
	return nil, 0
}

// Type implementation of Message interface.
func (*ValidationCheck) Type() core.MessageType {
	return core.TypeValidationCheck
}

// HotData contains hot-data
type HotData struct {
	ledgerMessage
	Jet                core.RecordRef
	Drop               jet.JetDrop
	RecentObjects      map[core.RecordID]*HotIndex
	PendingRequests    map[core.RecordID]map[core.RecordID][]byte
	PulseNumber        core.PulseNumber
	JetDropSizeHistory jet.DropSizeHistory
}

// AllowedSenderObjectAndRole implements interface method
func (m *HotData) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &m.Jet, core.DynamicRoleLightExecutor
}

// DefaultRole returns role for this event
func (*HotData) DefaultRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *HotData) DefaultTarget() *core.RecordRef {
	return &m.Jet
}

// Type implementation of Message interface.
func (*HotData) Type() core.MessageType {
	return core.TypeHotRecords
}

// HotIndex contains meat about hot-data
type HotIndex struct {
	TTL   int
	Index []byte
}

// GetPendingRequests fetches pending requests for object.
type GetPendingRequests struct {
	ledgerMessage

	Object core.RecordRef
}

// Type implementation of Message interface.
func (*GetPendingRequests) Type() core.MessageType {
	return core.TypeGetPendingRequests
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetPendingRequests) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &m.Object, core.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*GetPendingRequests) DefaultRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetPendingRequests) DefaultTarget() *core.RecordRef {
	return &m.Object
}

// GetJet requests to calculate a jet for provided object.
type GetJet struct {
	ledgerMessage

	Object core.RecordID
}

// Type implementation of Message interface.
func (*GetJet) Type() core.MessageType {
	return core.TypeGetJet
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetJet) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return core.NewRecordRef(core.DomainID, m.Object), core.DynamicRoleLightExecutor
}

// DefaultRole returns role for this event
func (*GetJet) DefaultRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetJet) DefaultTarget() *core.RecordRef {
	return core.NewRecordRef(core.DomainID, m.Object)
}

// AbandonedRequestsNotification informs virtual node about unclosed requests.
type AbandonedRequestsNotification struct {
	ledgerMessage

	Object   core.RecordID
	Requests []core.RecordID
}

// Type implementation of Message interface.
func (*AbandonedRequestsNotification) Type() core.MessageType {
	return core.TypeAbandonedRequestsNotification
}

// AllowedSenderObjectAndRole implements interface method
func (m *AbandonedRequestsNotification) AllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return nil, core.DynamicRoleUndefined
}

// DefaultRole returns role for this event
func (*AbandonedRequestsNotification) DefaultRole() core.DynamicRole {
	return core.DynamicRoleVirtualExecutor
}

// DefaultTarget returns of target of this event.
func (m *AbandonedRequestsNotification) DefaultTarget() *core.RecordRef {
	return core.NewRecordRef(core.DomainID, m.Object)
}
