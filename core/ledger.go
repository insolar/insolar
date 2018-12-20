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

package core

import (
	"context"
)

// DynamicRole is number representing a node role.
type DynamicRole int

const (
	// DynamicRoleUndefined is used for special cases.
	DynamicRoleUndefined = DynamicRole(iota)
	// DynamicRoleVirtualExecutor is responsible for current pulse CPU operations.
	DynamicRoleVirtualExecutor
	// DynamicRoleVirtualValidator is responsible for previous pulse CPU operations.
	DynamicRoleVirtualValidator
	// DynamicRoleLightExecutor is responsible for current pulse Disk operations.
	DynamicRoleLightExecutor
	// DynamicRoleLightValidator is responsible for previous pulse Disk operations.
	DynamicRoleLightValidator
	// DynamicRoleHeavyExecutor is responsible for permanent Disk operations.
	DynamicRoleHeavyExecutor
)

// IsVirtualRole checks if node role is virtual (validator or executor).
func (r DynamicRole) IsVirtualRole() bool {
	switch r {
	case DynamicRoleVirtualExecutor:
		return true
	case DynamicRoleVirtualValidator:
		return true
	}
	return false
}

// Ledger is the global ledger handler. Other system parts communicate with ledger through it.
// FIXME: THIS INTERFACE IS DEPRECATED. USE DI.
type Ledger interface {
	// GetArtifactManager returns artifact manager to work with.
	GetArtifactManager() ArtifactManager

	// GetJetCoordinator returns jet coordinator to work with.
	GetJetCoordinator() JetCoordinator

	// GetPulseManager returns pulse manager to work with.
	GetPulseManager() PulseManager

	// GetLocalStorage returns local storage to work with.
	GetLocalStorage() LocalStorage
}

// PulseManager provides Ledger's methods related to Pulse.
//go:generate minimock -i github.com/insolar/insolar/core.PulseManager -o ../testutils -s _mock.go
type PulseManager interface {
	// Set set's new pulse and closes current jet drop. If dry is true, nothing will be saved to storage.
	Set(ctx context.Context, pulse Pulse, persist bool) error
}

// JetCoordinator provides methods for calculating Jet affinity
// (e.g. to which Jet a message should be sent).
//go:generate minimock -i github.com/insolar/insolar/core.JetCoordinator -o ../testutils -s _mock.go
type JetCoordinator interface {
	// Me returns current node.
	Me() RecordRef

	// IsAuthorized checks for role on concrete pulse for the address.
	IsAuthorized(ctx context.Context, role DynamicRole, obj RecordID, pulse PulseNumber, node RecordRef) (bool, error)

	// QueryRole returns node refs responsible for role bound operations for given object and pulse.
	QueryRole(ctx context.Context, role DynamicRole, obj RecordID, pulse PulseNumber) ([]RecordRef, error)

	VirtualExecutorForObject(ctx context.Context, objID RecordID, pulse PulseNumber) (*RecordRef, error)
	VirtualValidatorsForObject(ctx context.Context, objID RecordID, pulse PulseNumber) ([]RecordRef, error)

	LightExecutorForObject(ctx context.Context, objID RecordID, pulse PulseNumber) (*RecordRef, error)
	LightValidatorsForObject(ctx context.Context, objID RecordID, pulse PulseNumber) ([]RecordRef, error)
	// LightExecutorForJet calculates light material executor for provided jet.
	LightExecutorForJet(ctx context.Context, jetID RecordID, pulse PulseNumber) (*RecordRef, error)
	LightValidatorsForJet(ctx context.Context, jetID RecordID, pulse PulseNumber) ([]RecordRef, error)

	Heavy(ctx context.Context, pulse PulseNumber) (*RecordRef, error)
}

// ArtifactManager is a high level storage interface.
//go:generate minimock -i github.com/insolar/insolar/core.ArtifactManager -o ../testutils -s _mock.go
type ArtifactManager interface {
	// GenesisRef returns the root record reference.
	//
	// Root record is the parent for all top-level records.
	GenesisRef() *RecordRef

	// RegisterRequest creates request record in storage.
	RegisterRequest(ctx context.Context, object RecordRef, parcel Parcel) (*RecordID, error)

	// RegisterValidation marks provided object state as approved or disapproved.
	//
	// When fetching object, validity can be specified.
	RegisterValidation(ctx context.Context, object RecordRef, state RecordID, isValid bool, validationMessages []Message) error

	// RegisterResult saves VM method call result.
	RegisterResult(ctx context.Context, object, request RecordRef, payload []byte) (*RecordID, error)

	// GetCode returns code from code record by provided reference according to provided machine preference.
	//
	// This method is used by VM to fetch code for execution.
	GetCode(ctx context.Context, ref RecordRef) (CodeDescriptor, error)

	// GetObject returns descriptor for provided state.
	//
	// If provided state is nil, the latest state will be returned (with deactivation check). Returned descriptor will
	// provide methods for fetching all related data.
	GetObject(ctx context.Context, head RecordRef, state *RecordID, approved bool) (ObjectDescriptor, error)

	// HasPendingRequests returns true if object has unclosed requests.
	HasPendingRequests(ctx context.Context, object RecordRef) (bool, error)

	// GetDelegate returns provided object's delegate reference for provided type.
	//
	// Object delegate should be previously created for this object. If object delegate does not exist, an error will
	// be returned.
	GetDelegate(ctx context.Context, head, asType RecordRef) (*RecordRef, error)

	// GetChildren returns children iterator.
	//
	// During iteration children refs will be fetched from remote source (parent object).
	GetChildren(ctx context.Context, parent RecordRef, pulse *PulseNumber) (RefIterator, error)

	// DeclareType creates new type record in storage.
	//
	// Type is a contract interface. It contains one method signature.
	DeclareType(ctx context.Context, domain, request RecordRef, typeDec []byte) (*RecordID, error)

	// DeployCode creates new code record in storage.
	//
	// Code records are used to activate prototype.
	DeployCode(ctx context.Context, domain, request RecordRef, code []byte, machineType MachineType) (*RecordID, error)

	// ActivatePrototype creates activate object record in storage. Provided prototype reference will be used as objects prototype
	// memory as memory of created object. If memory is not provided, the prototype default memory will be used.
	//
	// Request reference will be this object's identifier and referred as "object head".
	ActivatePrototype(
		ctx context.Context,
		domain, request, parent, code RecordRef,
		memory []byte,
	) (ObjectDescriptor, error)

	// ActivateObject creates activate object record in storage. If memory is not provided, the prototype default
	// memory will be used.
	//
	// Request reference will be this object's identifier and referred as "object head".
	ActivateObject(
		ctx context.Context,
		domain, request, parent, prototype RecordRef,
		asDelegate bool,
		memory []byte,
	) (ObjectDescriptor, error)

	// UpdatePrototype creates amend object record in storage. Provided reference should be a reference to the head of
	// the prototype. Provided memory well be the new object memory.
	//
	// Returned reference will be the latest object state (exact) reference.
	UpdatePrototype(
		ctx context.Context,
		domain, request RecordRef,
		obj ObjectDescriptor,
		memory []byte,
		code *RecordRef,
	) (ObjectDescriptor, error)

	// UpdateObject creates amend object record in storage. Provided reference should be a reference to the head of the
	// object. Provided memory well be the new object memory.
	//
	// Returned reference will be the latest object state (exact) reference.
	UpdateObject(
		ctx context.Context,
		domain, request RecordRef,
		obj ObjectDescriptor,
		memory []byte,
	) (ObjectDescriptor, error)

	// DeactivateObject creates deactivate object record in storage. Provided reference should be a reference to the head
	// of the object. If object is already deactivated, an error should be returned.
	//
	// Deactivated object cannot be changed.
	DeactivateObject(ctx context.Context, domain, request RecordRef, obj ObjectDescriptor) (*RecordID, error)

	// State returns hash state for artifact manager.
	State() ([]byte, error)
}

// CodeDescriptor represents meta info required to fetch all code data.
type CodeDescriptor interface {
	// Ref returns reference to represented code record.
	Ref() *RecordRef

	// MachineType returns code machine type for represented code.
	MachineType() MachineType

	// Code returns code data.
	Code() ([]byte, error)
}

// ObjectDescriptor represents meta info required to fetch all object data.
//go:generate minimock -i github.com/insolar/insolar/core.ObjectDescriptor -o ../testutils -s _mock.go
type ObjectDescriptor interface {
	// HeadRef returns head reference to represented object record.
	HeadRef() *RecordRef

	// StateID returns reference to object state record.
	StateID() *RecordID

	// Memory fetches object memory from storage.
	Memory() []byte

	// IsPrototype determines if the object is a prototype.
	IsPrototype() bool

	// Code returns code reference.
	Code() (*RecordRef, error)

	// Prototype returns prototype reference.
	Prototype() (*RecordRef, error)

	// Children returns object's children references.
	Children(pulse *PulseNumber) (RefIterator, error)

	// ChildPointer returns the latest child for this object.
	ChildPointer() *RecordID

	// Parent returns object's parent.
	Parent() *RecordRef
}

// RefIterator is used for iteration over affined children(parts) of container.
type RefIterator interface {
	Next() (*RecordRef, error)
	HasNext() bool
}

// LocalStorage allows a node to save local data.
//go:generate minimock -i github.com/insolar/insolar/core.LocalStorage -o ../testutils -s _mock.go
type LocalStorage interface {
	// Set saves data in storage.
	Set(ctx context.Context, pulse PulseNumber, key []byte, data []byte) error
	// Get retrieves data from storage.
	Get(ctx context.Context, pulse PulseNumber, key []byte) ([]byte, error)
	// Iterate iterates over all record with specified prefix and calls handler with key and value of that record.
	//
	// The key will be returned without prefix (e.g. the remaining slice) and value will be returned as it was saved.
	Iterate(ctx context.Context, pulse PulseNumber, prefix []byte, handler func(k, v []byte) error) error
}

// KV is a generic key/value struct.
type KV struct {
	K []byte
	V []byte
}

// StorageExportResult represents storage data view.
type StorageExportResult struct {
	Data     map[string]interface{}
	NextFrom *PulseNumber
	Size     int
}

// StorageExporter provides methods for fetching data view from storage.
type StorageExporter interface {
	// Export returns data view from storage.
	Export(ctx context.Context, fromPulse PulseNumber, size int) (*StorageExportResult, error)
}

var (
	// TODOJetID temporary stub for passing jet ID in ledger functions
	// on period Jet ID full implementation
	// TODO: remove it after jets support readyness - @nordicdyno 5.Dec.2018
	TODOJetID = *NewRecordID(PulseNumberJet, nil)
	DomainID  = *NewRecordID(0, nil)
)

// PulseStorage provides the interface for fetching current pulse of the system
//go:generate minimock -i github.com/insolar/insolar/core.PulseStorage -o ../testutils -s _mock.go
type PulseStorage interface {
	Current(ctx context.Context) (*Pulse, error)
}
