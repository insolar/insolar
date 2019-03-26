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

package insolar

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

// PulseManager provides Ledger's methods related to Pulse.
//go:generate minimock -i github.com/insolar/insolar/insolar.PulseManager -o ../testutils -s _mock.go
type PulseManager interface {
	// Set set's new pulse and closes current jet drop. If dry is true, nothing will be saved to storage.
	Set(ctx context.Context, pulse Pulse, persist bool) error
}

// JetCoordinator provides methods for calculating Jet affinity
// (e.g. to which Jet a message should be sent).
//go:generate minimock -i github.com/insolar/insolar/insolar.JetCoordinator -o ../testutils -s _mock.go
type JetCoordinator interface {
	// Me returns current node.
	Me() Reference

	// IsAuthorized checks for role on concrete pulse for the address.
	IsAuthorized(ctx context.Context, role DynamicRole, obj ID, pulse PulseNumber, node Reference) (bool, error)

	// QueryRole returns node refs responsible for role bound operations for given object and pulse.
	QueryRole(ctx context.Context, role DynamicRole, obj ID, pulse PulseNumber) ([]Reference, error)

	VirtualExecutorForObject(ctx context.Context, objID ID, pulse PulseNumber) (*Reference, error)
	VirtualValidatorsForObject(ctx context.Context, objID ID, pulse PulseNumber) ([]Reference, error)

	LightExecutorForObject(ctx context.Context, objID ID, pulse PulseNumber) (*Reference, error)
	LightValidatorsForObject(ctx context.Context, objID ID, pulse PulseNumber) ([]Reference, error)
	// LightExecutorForJet calculates light material executor for provided jet.
	LightExecutorForJet(ctx context.Context, jetID ID, pulse PulseNumber) (*Reference, error)
	LightValidatorsForJet(ctx context.Context, jetID ID, pulse PulseNumber) ([]Reference, error)

	Heavy(ctx context.Context, pulse PulseNumber) (*Reference, error)

	IsBeyondLimit(ctx context.Context, currentPN, targetPN PulseNumber) (bool, error)
	NodeForJet(ctx context.Context, jetID ID, rootPN, targetPN PulseNumber) (*Reference, error)

	// NodeForObject calculates a node (LME or heavy) for a specific jet for a specific pulseNumber
	NodeForObject(ctx context.Context, objectID ID, rootPN, targetPN PulseNumber) (*Reference, error)
}

// ArtifactManager is a high level storage interface.
//go:generate minimock -i github.com/insolar/insolar/insolar.ArtifactManager -o ../testutils -s _mock.go
type ArtifactManager interface {
	// GenesisRef returns the root record reference.
	//
	// Root record is the parent for all top-level records.
	GenesisRef() *Reference

	// RegisterRequest creates request record in storage.
	RegisterRequest(ctx context.Context, object Reference, parcel Parcel) (*ID, error)

	// RegisterValidation marks provided object state as approved or disapproved.
	//
	// When fetching object, validity can be specified.
	RegisterValidation(ctx context.Context, object Reference, state ID, isValid bool, validationMessages []Message) error

	// RegisterResult saves VM method call result.
	RegisterResult(ctx context.Context, object, request Reference, payload []byte) (*ID, error)

	// GetCode returns code from code record by provided reference according to provided machine preference.
	//
	// This method is used by VM to fetch code for execution.
	GetCode(ctx context.Context, ref Reference) (CodeDescriptor, error)

	// GetObject returns descriptor for provided state.
	//
	// If provided state is nil, the latest state will be returned (with deactivation check). Returned descriptor will
	// provide methods for fetching all related data.
	GetObject(ctx context.Context, head Reference, state *ID, approved bool) (ObjectDescriptor, error)

	// GetPendingRequest returns a pending request for object.
	GetPendingRequest(ctx context.Context, objectID ID) (Parcel, error)

	// HasPendingRequests returns true if object has unclosed requests.
	HasPendingRequests(ctx context.Context, object Reference) (bool, error)

	// GetDelegate returns provided object's delegate reference for provided type.
	//
	// Object delegate should be previously created for this object. If object delegate does not exist, an error will
	// be returned.
	GetDelegate(ctx context.Context, head, asType Reference) (*Reference, error)

	// GetChildren returns children iterator.
	//
	// During iteration children refs will be fetched from remote source (parent object).
	GetChildren(ctx context.Context, parent Reference, pulse *PulseNumber) (RefIterator, error)

	// DeclareType creates new type record in storage.
	//
	// Type is a contract interface. It contains one method signature.
	DeclareType(ctx context.Context, domain, request Reference, typeDec []byte) (*ID, error)

	// DeployCode creates new code record in storage.
	//
	// Code records are used to activate prototype.
	DeployCode(ctx context.Context, domain, request Reference, code []byte, machineType MachineType) (*ID, error)

	// ActivatePrototype creates activate object record in storage. Provided prototype reference will be used as objects prototype
	// memory as memory of created object. If memory is not provided, the prototype default memory will be used.
	//
	// Request reference will be this object's identifier and referred as "object head".
	ActivatePrototype(
		ctx context.Context,
		domain, request, parent, code Reference,
		memory []byte,
	) (ObjectDescriptor, error)

	// ActivateObject creates activate object record in storage. If memory is not provided, the prototype default
	// memory will be used.
	//
	// Request reference will be this object's identifier and referred as "object head".
	ActivateObject(
		ctx context.Context,
		domain, request, parent, prototype Reference,
		asDelegate bool,
		memory []byte,
	) (ObjectDescriptor, error)

	// UpdatePrototype creates amend object record in storage. Provided reference should be a reference to the head of
	// the prototype. Provided memory well be the new object memory.
	//
	// Returned reference will be the latest object state (exact) reference.
	UpdatePrototype(
		ctx context.Context,
		domain, request Reference,
		obj ObjectDescriptor,
		memory []byte,
		code *Reference,
	) (ObjectDescriptor, error)

	// UpdateObject creates amend object record in storage. Provided reference should be a reference to the head of the
	// object. Provided memory well be the new object memory.
	//
	// Returned reference will be the latest object state (exact) reference.
	UpdateObject(
		ctx context.Context,
		domain, request Reference,
		obj ObjectDescriptor,
		memory []byte,
	) (ObjectDescriptor, error)

	// DeactivateObject creates deactivate object record in storage. Provided reference should be a reference to the head
	// of the object. If object is already deactivated, an error should be returned.
	//
	// Deactivated object cannot be changed.
	DeactivateObject(ctx context.Context, domain, request Reference, obj ObjectDescriptor) (*ID, error)

	// State returns hash state for artifact manager.
	State() ([]byte, error)
}

// CodeDescriptor represents meta info required to fetch all code data.
//go:generate minimock -i github.com/insolar/insolar/insolar.CodeDescriptor -o ../testutils -s _mock.go
type CodeDescriptor interface {
	// Ref returns reference to represented code record.
	Ref() *Reference

	// MachineType returns code machine type for represented code.
	MachineType() MachineType

	// Code returns code data.
	Code() ([]byte, error)
}

// ObjectDescriptor represents meta info required to fetch all object data.
//go:generate minimock -i github.com/insolar/insolar/insolar.ObjectDescriptor -o ../testutils -s _mock.go
type ObjectDescriptor interface {
	// HeadRef returns head reference to represented object record.
	HeadRef() *Reference

	// StateID returns reference to object state record.
	StateID() *ID

	// Memory fetches object memory from storage.
	Memory() []byte

	// IsPrototype determines if the object is a prototype.
	IsPrototype() bool

	// Code returns code reference.
	Code() (*Reference, error)

	// Prototype returns prototype reference.
	Prototype() (*Reference, error)

	// Children returns object's children references.
	Children(pulse *PulseNumber) (RefIterator, error)

	// ChildPointer returns the latest child for this object.
	ChildPointer() *ID

	// Parent returns object's parent.
	Parent() *Reference
}

// RefIterator is used for iteration over affined children(parts) of container.
type RefIterator interface {
	Next() (*Reference, error)
	HasNext() bool
}

// KV is a generic key/value struct.
type KV struct {
	K []byte
	V []byte
}

// KVSize returns size of key/value array in bytes.
func KVSize(kvs []KV) (amount int64) {
	for _, kv := range kvs {
		amount += int64(len(kv.K) + len(kv.V))
	}
	return
}

// StorageExportResult represents storage data view.
type StorageExportResult struct {
	Data     map[string]interface{}
	NextFrom *PulseNumber
	Size     int
}

var (
	// TODOJetID temporary stub for passing jet ID in ledger functions
	// on period Jet ID full implementation
	// TODO: remove it after jets support readyness - @nordicdyno 5.Dec.2018
	TODOJetID = *NewID(PulseNumberJet, nil)
	DomainID  = *NewID(0, nil)
)

// PulseStorage provides the interface for fetching current pulse of the system
//go:generate minimock -i github.com/insolar/insolar/insolar.PulseStorage -o ../testutils -s _mock.go
type PulseStorage interface {
	Current(ctx context.Context) (*Pulse, error)
}
