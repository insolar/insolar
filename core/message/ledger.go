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

// ExtractAllowedSenderObjectAndRole implements interface method
func (sr *SetRecord) ExtractAllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &sr.TargetRef, core.DynamicRoleVirtualExecutor
}

// ExtractRole returns role for this event
func (*SetRecord) ExtractRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// ExtractTarget returns of target of this event.
func (sr *SetRecord) ExtractTarget() *core.RecordRef {
	return &sr.TargetRef
}

// Type implementation of Message interface.
func (e *SetRecord) Type() core.MessageType {
	return core.TypeSetRecord
}

// GetCode retrieves code From storage.
type GetCode struct {
	ledgerMessage
	Code core.RecordRef
}

// ExtractAllowedSenderObjectAndRole implements interface method
func (gc *GetCode) ExtractAllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &gc.Code, core.DynamicRoleVirtualExecutor
}

// ExtractRole returns role for this event
func (*GetCode) ExtractRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// ExtractTarget returns of target of this event.
func (gc *GetCode) ExtractTarget() *core.RecordRef {
	return &gc.Code
}

// Type implementation of Message interface.
func (e *GetCode) Type() core.MessageType {
	return core.TypeGetCode
}

// GetObject retrieves object From storage.
type GetObject struct {
	ledgerMessage
	Head     core.RecordRef
	State    *core.RecordID // If nil, will fetch the latest state.
	Approved bool
}

// ExtractAllowedSenderObjectAndRole implements interface method
func (getObj *GetObject) ExtractAllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &getObj.Head, core.DynamicRoleVirtualExecutor
}

// ExtractRole returns role for this event
func (*GetObject) ExtractRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// ExtractTarget returns of target of this event.
func (getObj *GetObject) ExtractTarget() *core.RecordRef {
	return &getObj.Head
}

// Type implementation of Message interface.
func (getObj *GetObject) Type() core.MessageType {
	return core.TypeGetObject
}

// GetDelegate retrieves object represented as provided type.
type GetDelegate struct {
	ledgerMessage
	Head   core.RecordRef
	AsType core.RecordRef
}

// ExtractAllowedSenderObjectAndRole implements interface method
func (gd *GetDelegate) ExtractAllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &gd.Head, core.DynamicRoleVirtualExecutor
}

// ExtractRole returns role for this event
func (*GetDelegate) ExtractRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// ExtractTarget returns of target of this event.
func (gd *GetDelegate) ExtractTarget() *core.RecordRef {
	return &gd.Head
}

// Type implementation of Message interface.
func (e *GetDelegate) Type() core.MessageType {
	return core.TypeGetDelegate
}

// UpdateObject amends object.
type UpdateObject struct {
	ledgerMessage

	Record []byte
	Object core.RecordRef
}

// ExtractAllowedSenderObjectAndRole implements interface method
func (uo *UpdateObject) ExtractAllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &uo.Object, core.DynamicRoleVirtualExecutor
}

// ExtractRole returns role for this event
func (*UpdateObject) ExtractRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// ExtractTarget returns of target of this event.
func (uo *UpdateObject) ExtractTarget() *core.RecordRef {
	return &uo.Object
}

// Type implementation of Message interface.
func (e *UpdateObject) Type() core.MessageType {
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

// ExtractAllowedSenderObjectAndRole implements interface method
func (rc *RegisterChild) ExtractAllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &rc.Child, core.DynamicRoleVirtualExecutor
}

// ExtractRole returns role for this event
func (*RegisterChild) ExtractRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// ExtractTarget returns of target of this event.
func (rc *RegisterChild) ExtractTarget() *core.RecordRef {
	return &rc.Parent
}

// Type implementation of Message interface.
func (rc *RegisterChild) Type() core.MessageType {
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

// ExtractAllowedSenderObjectAndRole implements interface method
func (gc *GetChildren) ExtractAllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &gc.Parent, core.DynamicRoleVirtualExecutor
}

// ExtractRole returns role for this event
func (*GetChildren) ExtractRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// ExtractTarget returns of target of this event.
func (gc *GetChildren) ExtractTarget() *core.RecordRef {
	return &gc.Parent
}

// Type implementation of Message interface.
func (e *GetChildren) Type() core.MessageType {
	return core.TypeGetChildren
}

// JetDrop spreads jet drop
type JetDrop struct {
	ledgerMessage
	Jet         core.RecordRef
	Drop        []byte
	Messages    [][]byte
	PulseNumber core.PulseNumber
}

// ExtractAllowedSenderObjectAndRole implements interface method
func (jd *JetDrop) ExtractAllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	// This check is not needed, because JetDrop sender is explicitly checked in handler.
	return nil, core.DynamicRoleUndefined
}

// ExtractRole returns role for this event
func (*JetDrop) ExtractRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// ExtractTarget returns of target of this event.
func (jd *JetDrop) ExtractTarget() *core.RecordRef {
	return &jd.Jet
}

// Type implementation of Message interface.
func (e *JetDrop) Type() core.MessageType {
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

// ExtractAllowedSenderObjectAndRole implements interface method
func (vr *ValidateRecord) ExtractAllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &vr.Object, core.DynamicRoleVirtualExecutor
}

// ExtractRole returns role for this event
func (*ValidateRecord) ExtractRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// ExtractTarget returns of target of this event.
func (vr *ValidateRecord) ExtractTarget() *core.RecordRef {
	return &vr.Object
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

// ExtractAllowedSenderObjectAndRole implements interface method
func (sb *SetBlob) ExtractAllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &sb.TargetRef, core.DynamicRoleVirtualExecutor
}

// ExtractRole returns role for this event
func (*SetBlob) ExtractRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// ExtractTarget returns of target of this event.
func (sb *SetBlob) ExtractTarget() *core.RecordRef {
	return &sb.TargetRef
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

// ExtractAllowedSenderObjectAndRole implements interface method
func (getObjectIndex *GetObjectIndex) ExtractAllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &getObjectIndex.Object, core.DynamicRoleLightExecutor
}

// ExtractRole returns role for this event
func (*GetObjectIndex) ExtractRole() core.DynamicRole {
	return core.DynamicRoleHeavyExecutor
}

// ExtractTarget returns of target of this event.
func (getObjectIndex *GetObjectIndex) ExtractTarget() *core.RecordRef {
	return &getObjectIndex.Object
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

// ExtractTarget returns of target of this event.
func (vc *ValidationCheck) ExtractTarget() *core.RecordRef {
	// TODO: return smth real
	return nil
}

// ExtractRole returns role for this event
func (vc *ValidationCheck) ExtractRole() core.DynamicRole {
	return core.DynamicRoleUndefined
}

// ExtractAllowedSenderObjectAndRole implements interface method
func (vc *ValidationCheck) ExtractAllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
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
	Jet             core.RecordRef
	Drop            jet.JetDrop
	RecentObjects   map[core.RecordID]*HotIndex
	PendingRequests map[core.RecordID][]byte
	PulseNumber     core.PulseNumber
}

// ExtractAllowedSenderObjectAndRole implements interface method
func (hd *HotData) ExtractAllowedSenderObjectAndRole() (*core.RecordRef, core.DynamicRole) {
	return &hd.Jet, core.DynamicRoleLightExecutor
}

// ExtractRole returns role for this event
func (*HotData) ExtractRole() core.DynamicRole {
	return core.DynamicRoleLightExecutor
}

// ExtractTarget returns of target of this event.
func (hd *HotData) ExtractTarget() *core.RecordRef {
	return &hd.Jet
}

// HotIndex contains meat about hot-data
type HotIndex struct {
	TTL   int
	Index []byte
}

// Type implementation of Message interface.
func (*HotData) Type() core.MessageType {
	return core.TypeHotRecords
}
