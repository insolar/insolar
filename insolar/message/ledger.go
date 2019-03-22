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
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage/drop"
)

// FIXME: @andreyromancev. 21.12.18. Remove this and create 'LogicRunnerMessage' interface to get rid of 'GetCaller' in ledger.
type ledgerMessage struct {
}

// GetCaller implementation of Message interface.
func (ledgerMessage) GetCaller() *insolar.RecordRef {
	return nil
}

// SetRecord saves record in storage.
type SetRecord struct {
	ledgerMessage

	Record    []byte
	TargetRef insolar.RecordRef
}

// AllowedSenderObjectAndRole implements interface method
func (m *SetRecord) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return &m.TargetRef, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*SetRecord) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *SetRecord) DefaultTarget() *insolar.RecordRef {
	return &m.TargetRef
}

// Type implementation of Message interface.
func (m *SetRecord) Type() insolar.MessageType {
	return insolar.TypeSetRecord
}

// GetCode retrieves code From storage.
type GetCode struct {
	ledgerMessage
	Code insolar.RecordRef
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetCode) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return &m.Code, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*GetCode) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetCode) DefaultTarget() *insolar.RecordRef {
	return &m.Code
}

// Type implementation of Message interface.
func (*GetCode) Type() insolar.MessageType {
	return insolar.TypeGetCode
}

// GetObject retrieves object From storage.
type GetObject struct {
	ledgerMessage
	Head     insolar.RecordRef
	State    *insolar.RecordID // If nil, will fetch the latest state.
	Approved bool
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetObject) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return &m.Head, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*GetObject) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetObject) DefaultTarget() *insolar.RecordRef {
	return &m.Head
}

// Type implementation of Message interface.
func (*GetObject) Type() insolar.MessageType {
	return insolar.TypeGetObject
}

// GetDelegate retrieves object represented as provided type.
type GetDelegate struct {
	ledgerMessage
	Head   insolar.RecordRef
	AsType insolar.RecordRef
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetDelegate) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return &m.Head, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*GetDelegate) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetDelegate) DefaultTarget() *insolar.RecordRef {
	return &m.Head
}

// Type implementation of Message interface.
func (*GetDelegate) Type() insolar.MessageType {
	return insolar.TypeGetDelegate
}

// UpdateObject amends object.
type UpdateObject struct {
	ledgerMessage

	Record []byte
	Object insolar.RecordRef
	Memory []byte
}

// AllowedSenderObjectAndRole implements interface method
func (m *UpdateObject) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return &m.Object, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*UpdateObject) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *UpdateObject) DefaultTarget() *insolar.RecordRef {
	return &m.Object
}

// Type implementation of Message interface.
func (*UpdateObject) Type() insolar.MessageType {
	return insolar.TypeUpdateObject
}

// RegisterChild amends object.
type RegisterChild struct {
	ledgerMessage
	Record []byte
	Parent insolar.RecordRef
	Child  insolar.RecordRef
	AsType *insolar.RecordRef // If not nil, considered as delegate.
}

// AllowedSenderObjectAndRole implements interface method
func (m *RegisterChild) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return &m.Child, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*RegisterChild) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *RegisterChild) DefaultTarget() *insolar.RecordRef {
	return &m.Parent
}

// Type implementation of Message interface.
func (*RegisterChild) Type() insolar.MessageType {
	return insolar.TypeRegisterChild
}

// GetChildren retrieves a chunk of children references.
type GetChildren struct {
	ledgerMessage
	Parent    insolar.RecordRef
	FromChild *insolar.RecordID
	FromPulse *insolar.PulseNumber
	Amount    int
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetChildren) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return &m.Parent, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*GetChildren) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetChildren) DefaultTarget() *insolar.RecordRef {
	return &m.Parent
}

// Type implementation of Message interface.
func (*GetChildren) Type() insolar.MessageType {
	return insolar.TypeGetChildren
}

// Drop spreads jet drop
type JetDrop struct {
	ledgerMessage

	JetID insolar.RecordID

	Drop        []byte
	Messages    [][]byte
	PulseNumber insolar.PulseNumber
}

// AllowedSenderObjectAndRole implements interface method
func (m *JetDrop) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	// This check is not needed, because Drop sender is explicitly checked in handler.
	return nil, insolar.DynamicRoleUndefined
}

// DefaultRole returns role for this event
func (*JetDrop) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *JetDrop) DefaultTarget() *insolar.RecordRef {
	return insolar.NewRecordRef(insolar.RecordID{}, m.JetID)
}

// Type implementation of Message interface.
func (*JetDrop) Type() insolar.MessageType {
	return insolar.TypeJetDrop
}

// ValidateRecord creates VM validation for specific object record.
type ValidateRecord struct {
	ledgerMessage

	Object             insolar.RecordRef
	State              insolar.RecordID
	IsValid            bool
	ValidationMessages []insolar.Message
}

// AllowedSenderObjectAndRole implements interface method
func (m *ValidateRecord) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return &m.Object, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*ValidateRecord) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *ValidateRecord) DefaultTarget() *insolar.RecordRef {
	return &m.Object
}

// Type implementation of Message interface.
func (*ValidateRecord) Type() insolar.MessageType {
	return insolar.TypeValidateRecord
}

// SetBlob saves blob in storage.
type SetBlob struct {
	ledgerMessage

	TargetRef insolar.RecordRef
	Memory    []byte
}

// AllowedSenderObjectAndRole implements interface method
func (m *SetBlob) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return &m.TargetRef, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*SetBlob) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *SetBlob) DefaultTarget() *insolar.RecordRef {
	return &m.TargetRef
}

// Type implementation of Message interface.
func (*SetBlob) Type() insolar.MessageType {
	return insolar.TypeSetBlob
}

// GetObjectIndex fetches objects index.
type GetObjectIndex struct {
	ledgerMessage

	Object insolar.RecordRef
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetObjectIndex) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return &m.Object, insolar.DynamicRoleLightExecutor
}

// DefaultRole returns role for this event
func (*GetObjectIndex) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleHeavyExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetObjectIndex) DefaultTarget() *insolar.RecordRef {
	return &m.Object
}

// Type implementation of Message interface.
func (*GetObjectIndex) Type() insolar.MessageType {
	return insolar.TypeGetObjectIndex
}

// ValidationCheck checks if validation of a particular record can be performed.
type ValidationCheck struct {
	ledgerMessage

	Object              insolar.RecordRef
	ValidatedState      insolar.RecordID
	LatestStateApproved *insolar.RecordID
}

// DefaultTarget returns of target of this event.
func (m *ValidationCheck) DefaultTarget() *insolar.RecordRef {
	return &m.Object
}

// DefaultRole returns role for this event
func (m *ValidationCheck) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// AllowedSenderObjectAndRole implements interface method
func (m *ValidationCheck) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	// TODO: return smth real
	return nil, 0
}

// Type implementation of Message interface.
func (*ValidationCheck) Type() insolar.MessageType {
	return insolar.TypeValidationCheck
}

// HotData contains hot-data
type HotData struct {
	ledgerMessage
	Jet             insolar.RecordRef
	Drop            drop.Drop
	RecentObjects   map[insolar.RecordID]HotIndex
	PendingRequests map[insolar.RecordID]recentstorage.PendingObjectContext
	PulseNumber     insolar.PulseNumber
}

// AllowedSenderObjectAndRole implements interface method
func (m *HotData) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return nil, insolar.DynamicRoleUndefined
}

// DefaultRole returns role for this event
func (*HotData) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *HotData) DefaultTarget() *insolar.RecordRef {
	return &m.Jet
}

// Type implementation of Message interface.
func (*HotData) Type() insolar.MessageType {
	return insolar.TypeHotRecords
}

// HotIndex contains meat about hot-data
type HotIndex struct {
	TTL   int
	Index []byte
}

// GetPendingRequests fetches pending requests for object.
type GetPendingRequests struct {
	ledgerMessage

	Object insolar.RecordRef
}

// Type implementation of Message interface.
func (*GetPendingRequests) Type() insolar.MessageType {
	return insolar.TypeGetPendingRequests
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetPendingRequests) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return &m.Object, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*GetPendingRequests) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetPendingRequests) DefaultTarget() *insolar.RecordRef {
	return &m.Object
}

// GetJet requests to calculate a jet for provided object.
type GetJet struct {
	ledgerMessage

	Object insolar.RecordID
	Pulse  insolar.PulseNumber
}

// Type implementation of Message interface.
func (*GetJet) Type() insolar.MessageType {
	return insolar.TypeGetJet
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetJet) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return insolar.NewRecordRef(insolar.DomainID, m.Object), insolar.DynamicRoleLightExecutor
}

// DefaultRole returns role for this event
func (*GetJet) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetJet) DefaultTarget() *insolar.RecordRef {
	return insolar.NewRecordRef(insolar.DomainID, m.Object)
}

// AbandonedRequestsNotification informs virtual node about unclosed requests.
type AbandonedRequestsNotification struct {
	ledgerMessage

	Object insolar.RecordID
}

// Type implementation of Message interface.
func (*AbandonedRequestsNotification) Type() insolar.MessageType {
	return insolar.TypeAbandonedRequestsNotification
}

// AllowedSenderObjectAndRole implements interface method
func (m *AbandonedRequestsNotification) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return nil, insolar.DynamicRoleUndefined
}

// DefaultRole returns role for this event
func (*AbandonedRequestsNotification) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleVirtualExecutor
}

// DefaultTarget returns of target of this event.
func (m *AbandonedRequestsNotification) DefaultTarget() *insolar.RecordRef {
	return insolar.NewRecordRef(insolar.DomainID, m.Object)
}

// GetRequest fetches request from ledger.
type GetRequest struct {
	ledgerMessage

	Request insolar.RecordID
}

// Type implementation of Message interface.
func (*GetRequest) Type() insolar.MessageType {
	return insolar.TypeGetRequest
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetRequest) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return nil, insolar.DynamicRoleUndefined
}

// DefaultRole returns role for this event
func (*GetRequest) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetRequest) DefaultTarget() *insolar.RecordRef {
	return insolar.NewRecordRef(insolar.DomainID, m.Request)
}

// GetPendingRequestID fetches a pending request id for an object from current LME
type GetPendingRequestID struct {
	ledgerMessage

	ObjectID insolar.RecordID
}

// Type implementation of Message interface.
func (*GetPendingRequestID) Type() insolar.MessageType {
	return insolar.TypeGetPendingRequestID
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetPendingRequestID) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return nil, insolar.DynamicRoleUndefined
}

// DefaultRole returns role for this event
func (*GetPendingRequestID) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetPendingRequestID) DefaultTarget() *insolar.RecordRef {
	return insolar.NewRecordRef(insolar.DomainID, m.ObjectID)
}
