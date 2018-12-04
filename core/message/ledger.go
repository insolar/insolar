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
	"github.com/insolar/insolar/ledger/jetdrop"
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

// Type implementation of Message interface.
func (e *SetRecord) Type() core.MessageType {
	return core.TypeSetRecord
}

// GetCode retrieves code From storage.
type GetCode struct {
	ledgerMessage
	Code core.RecordRef
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

// Type implementation of Message interface.
func (e *GetObject) Type() core.MessageType {
	return core.TypeGetObject
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

// GetObjectIndex fetches objects index.
type GetObjectIndex struct {
	ledgerMessage

	Object core.RecordRef
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

// Type implementation of Message interface.
func (*ValidationCheck) Type() core.MessageType {
	return core.TypeValidationCheck
}

// HotData contains hot-data
type HotData struct {
	ledgerMessage
	Jet             core.RecordRef
	Drop            jetdrop.JetDrop
	RecentObjects   map[core.RecordID]*HotIndex
	PendingRequests map[core.RecordID][]byte
	PulseNumber     core.PulseNumber
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
