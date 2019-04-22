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
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/recentstorage"
)

// FIXME: @andreyromancev. 21.12.18. Remove this and create 'LogicRunnerMessage' interface to get rid of 'GetCaller' in ledger.
type ledgerMessage struct {
}

// GetCaller implementation of Message interface.
func (ledgerMessage) GetCaller() *insolar.Reference {
	return nil
}

// SetRecord saves record in storage.
type SetRecord struct {
	ledgerMessage

	Record    []byte
	TargetRef insolar.Reference
}

// AllowedSenderObjectAndRole implements interface method
func (m *SetRecord) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return &m.TargetRef, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*SetRecord) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *SetRecord) DefaultTarget() *insolar.Reference {
	return &m.TargetRef
}

// Type implementation of Message interface.
func (m *SetRecord) Type() insolar.MessageType {
	return insolar.TypeSetRecord
}

// GetCode retrieves code From storage.
type GetCode struct {
	ledgerMessage
	Code insolar.Reference
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetCode) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return &m.Code, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*GetCode) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetCode) DefaultTarget() *insolar.Reference {
	return &m.Code
}

// Type implementation of Message interface.
func (*GetCode) Type() insolar.MessageType {
	return insolar.TypeGetCode
}

// GetObject retrieves object From storage.
type GetObject struct {
	ledgerMessage
	Head     insolar.Reference
	State    *insolar.ID // If nil, will fetch the latest state.
	Approved bool
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetObject) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return &m.Head, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*GetObject) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetObject) DefaultTarget() *insolar.Reference {
	return &m.Head
}

// Type implementation of Message interface.
func (*GetObject) Type() insolar.MessageType {
	return insolar.TypeGetObject
}

// GetDelegate retrieves object represented as provided type.
type GetDelegate struct {
	ledgerMessage
	Head   insolar.Reference
	AsType insolar.Reference
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetDelegate) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return &m.Head, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*GetDelegate) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetDelegate) DefaultTarget() *insolar.Reference {
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
	Object insolar.Reference
	Memory []byte
}

// AllowedSenderObjectAndRole implements interface method
func (m *UpdateObject) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return &m.Object, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*UpdateObject) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *UpdateObject) DefaultTarget() *insolar.Reference {
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
	Parent insolar.Reference
	Child  insolar.Reference
	AsType *insolar.Reference // If not nil, considered as delegate.
}

// AllowedSenderObjectAndRole implements interface method
func (m *RegisterChild) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return &m.Child, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*RegisterChild) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *RegisterChild) DefaultTarget() *insolar.Reference {
	return &m.Parent
}

// Type implementation of Message interface.
func (*RegisterChild) Type() insolar.MessageType {
	return insolar.TypeRegisterChild
}

// GetChildren retrieves a chunk of children references.
type GetChildren struct {
	ledgerMessage
	Parent    insolar.Reference
	FromChild *insolar.ID
	FromPulse *insolar.PulseNumber
	Amount    int
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetChildren) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return &m.Parent, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*GetChildren) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetChildren) DefaultTarget() *insolar.Reference {
	return &m.Parent
}

// Type implementation of Message interface.
func (*GetChildren) Type() insolar.MessageType {
	return insolar.TypeGetChildren
}

// ValidateRecord creates VM validation for specific object record.
type ValidateRecord struct {
	ledgerMessage

	Object             insolar.Reference
	State              insolar.ID
	IsValid            bool
	ValidationMessages []insolar.Message
}

// AllowedSenderObjectAndRole implements interface method
func (m *ValidateRecord) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return &m.Object, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*ValidateRecord) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *ValidateRecord) DefaultTarget() *insolar.Reference {
	return &m.Object
}

// Type implementation of Message interface.
func (*ValidateRecord) Type() insolar.MessageType {
	return insolar.TypeValidateRecord
}

// SetBlob saves blob in storage.
type SetBlob struct {
	ledgerMessage

	TargetRef insolar.Reference
	Memory    []byte
}

// AllowedSenderObjectAndRole implements interface method
func (m *SetBlob) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return &m.TargetRef, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*SetBlob) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *SetBlob) DefaultTarget() *insolar.Reference {
	return &m.TargetRef
}

// Type implementation of Message interface.
func (*SetBlob) Type() insolar.MessageType {
	return insolar.TypeSetBlob
}

// GetObjectIndex fetches objects index.
type GetObjectIndex struct {
	ledgerMessage

	Object insolar.Reference
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetObjectIndex) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return &m.Object, insolar.DynamicRoleLightExecutor
}

// DefaultRole returns role for this event
func (*GetObjectIndex) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleHeavyExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetObjectIndex) DefaultTarget() *insolar.Reference {
	return &m.Object
}

// Type implementation of Message interface.
func (*GetObjectIndex) Type() insolar.MessageType {
	return insolar.TypeGetObjectIndex
}

// HotData contains hot-data
type HotData struct {
	ledgerMessage
	Jet             insolar.Reference
	Drop            drop.Drop
	HotIndexes      map[insolar.ID]HotIndex
	PendingRequests map[insolar.ID]recentstorage.PendingObjectContext
	PulseNumber     insolar.PulseNumber
}

// AllowedSenderObjectAndRole implements interface method
func (m *HotData) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return nil, insolar.DynamicRoleUndefined
}

// DefaultRole returns role for this event
func (*HotData) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *HotData) DefaultTarget() *insolar.Reference {
	return &m.Jet
}

// Type implementation of Message interface.
func (*HotData) Type() insolar.MessageType {
	return insolar.TypeHotRecords
}

// HotIndex contains meat about hot-data
type HotIndex struct {
	LastUsed insolar.PulseNumber
	Index    []byte
}

// GetPendingRequests fetches pending requests for object.
type GetPendingRequests struct {
	ledgerMessage

	Object insolar.Reference
}

// Type implementation of Message interface.
func (*GetPendingRequests) Type() insolar.MessageType {
	return insolar.TypeGetPendingRequests
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetPendingRequests) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return &m.Object, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*GetPendingRequests) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetPendingRequests) DefaultTarget() *insolar.Reference {
	return &m.Object
}

// GetJet requests to calculate a jet for provided object.
type GetJet struct {
	ledgerMessage

	Object insolar.ID
	Pulse  insolar.PulseNumber
}

// Type implementation of Message interface.
func (*GetJet) Type() insolar.MessageType {
	return insolar.TypeGetJet
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetJet) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return insolar.NewReference(insolar.DomainID, m.Object), insolar.DynamicRoleLightExecutor
}

// DefaultRole returns role for this event
func (*GetJet) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetJet) DefaultTarget() *insolar.Reference {
	return insolar.NewReference(insolar.DomainID, m.Object)
}

// AbandonedRequestsNotification informs virtual node about unclosed requests.
type AbandonedRequestsNotification struct {
	ledgerMessage

	Object insolar.ID
}

// Type implementation of Message interface.
func (*AbandonedRequestsNotification) Type() insolar.MessageType {
	return insolar.TypeAbandonedRequestsNotification
}

// AllowedSenderObjectAndRole implements interface method
func (m *AbandonedRequestsNotification) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return nil, insolar.DynamicRoleUndefined
}

// DefaultRole returns role for this event
func (*AbandonedRequestsNotification) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleVirtualExecutor
}

// DefaultTarget returns of target of this event.
func (m *AbandonedRequestsNotification) DefaultTarget() *insolar.Reference {
	return insolar.NewReference(insolar.DomainID, m.Object)
}

// GetRequest fetches request from ledger.
type GetRequest struct {
	ledgerMessage

	Request insolar.ID
}

// Type implementation of Message interface.
func (*GetRequest) Type() insolar.MessageType {
	return insolar.TypeGetRequest
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetRequest) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return nil, insolar.DynamicRoleUndefined
}

// DefaultRole returns role for this event
func (*GetRequest) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetRequest) DefaultTarget() *insolar.Reference {
	return insolar.NewReference(insolar.DomainID, m.Request)
}

// GetPendingRequestID fetches a pending request id for an object from current LME
type GetPendingRequestID struct {
	ledgerMessage

	ObjectID insolar.ID
}

// Type implementation of Message interface.
func (*GetPendingRequestID) Type() insolar.MessageType {
	return insolar.TypeGetPendingRequestID
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetPendingRequestID) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return nil, insolar.DynamicRoleUndefined
}

// DefaultRole returns role for this event
func (*GetPendingRequestID) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetPendingRequestID) DefaultTarget() *insolar.Reference {
	return insolar.NewReference(insolar.DomainID, m.ObjectID)
}
