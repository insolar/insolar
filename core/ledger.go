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

// JetRole is number representing a node role.
type JetRole int

const (
	// RoleVirtualExecutor is responsible for current pulse CPU operations.
	RoleVirtualExecutor = JetRole(iota + 1)
	// RoleVirtualValidator is responsible for previous pulse CPU operations.
	RoleVirtualValidator
	// RoleLightExecutor is responsible for current pulse Disk operations.
	RoleLightExecutor
	// RoleLightValidator is responsible for previous pulse Disk operations.
	RoleLightValidator
	// RoleHeavyExecutor is responsible for permanent Disk operations.
	RoleHeavyExecutor
)

// Ledger is the global ledger handler. Other system parts communicate with ledger through it.
type Ledger interface {
	// GetArtifactManager returns artifact manager to work with.
	GetArtifactManager() ArtifactManager

	// GetJetCoordinator returns jet coordinator to work with.
	GetJetCoordinator() JetCoordinator

	// GetPulseManager returns pulse manager to work with.
	GetPulseManager() PulseManager
}

// PulseManager provides Ledger's methods related to Pulse.
type PulseManager interface {
	// Current returns current pulse structure.
	Current() (*Pulse, error)

	// Set set's new pulse and closes current jet drop.
	Set(Pulse) error
}

// JetCoordinator provides methods for calculating Jet affinity
// (e.g. to which Jet a message should be sent).
type JetCoordinator interface {
	// IsAuthorized checks for role on concrete pulse for the address.
	IsAuthorized(role JetRole, obj RecordRef, pulse PulseNumber, node RecordRef) (bool, error)

	// QueryRole returns node refs responsible for role bound operations for given object and pulse.
	QueryRole(role JetRole, obj RecordRef, pulse PulseNumber) ([]RecordRef, error)
}

// ArtifactManager is a high level storage interface.
type ArtifactManager interface {
	// GenesisRef returns the root record reference.
	//
	// Root record is the parent for all top-level records.
	GenesisRef() *RecordRef

	// RegisterRequest creates or check call request record and returns it RecordRef.
	// (used by VM on executing side)
	RegisterRequest(ctx Context, message Message) (*RecordRef, error)

	// GetCode returns code from code record by provided reference according to provided machine preference.
	//
	// This method is used by VM to fetch code for execution.
	GetCode(ctx Context, ref RecordRef) (CodeDescriptor, error)

	// GetClass returns descriptor for provided state.
	//
	// If provided state is nil, the latest state will be returned (with deactivation check). Returned descriptor will
	// provide methods for fetching all related data.
	GetClass(ctx Context, head RecordRef, state *RecordRef) (ClassDescriptor, error)

	// GetObject returns descriptor for provided state.
	//
	// If provided state is nil, the latest state will be returned (with deactivation check). Returned descriptor will
	// provide methods for fetching all related data.
	GetObject(ctx Context, head RecordRef, state *RecordRef) (ObjectDescriptor, error)

	// GetDelegate returns provided object's delegate reference for provided class.
	//
	// Object delegate should be previously created for this object. If object delegate does not exist, an error will
	// be returned.
	GetDelegate(ctx Context, head, asClass RecordRef) (*RecordRef, error)

	// GetChildren returns children iterator.
	//
	// During iteration children refs will be fetched from remote source (parent object).
	GetChildren(ctx Context, parent RecordRef, pulse *PulseNumber) (RefIterator, error)

	// DeclareType creates new type record in storage.
	//
	// Type is a contract interface. It contains one method signature.
	DeclareType(ctx Context, domain, request RecordRef, typeDec []byte) (*RecordRef, error)

	// DeployCode creates new code record in storage.
	//
	// Code records are used to activate class or as migration code for an object.
	DeployCode(ctx Context, domain, request RecordRef, code []byte, machineType MachineType) (*RecordRef, error)

	// ActivateClass creates activate class record in storage. Provided code reference will be used as a class code.
	//
	// Request reference will be this class'es identifier and referred as "class head".
	ActivateClass(ctx Context, domain, request, code RecordRef) (*RecordID, error)

	// DeactivateClass creates deactivate record in storage. Provided reference should be a reference to the head of
	// the class. If class is already deactivated, an error should be returned.
	//
	// Deactivated class cannot be changed or instantiate objects.
	DeactivateClass(domain, request, class RecordRef) (*RecordID, error)

	// UpdateClass creates amend class record in storage. Provided reference should be a reference to the head of
	// the class. Migrations are references to code records.
	//
	// Returned reference will be the latest class state (exact) reference. Migration code will be executed by VM to
	// migrate objects memory in the order they appear in provided slice.
	UpdateClass(domain, request, class, code RecordRef, migrationRefs []RecordRef) (*RecordID, error)

	// ActivateObject creates activate object record in storage. Provided class reference will be used as object's class.
	// If memory is not provided, the class default memory will be used.
	//
	// Request reference will be this object's identifier and referred as "object head".
	ActivateObject(domain, request, class, parent RecordRef, memory []byte) (*RecordID, error)

	// ActivateObjectDelegate is similar to ActivateObject but it created object will be parent's delegate of provided class.
	ActivateObjectDelegate(domain, request, class, parent RecordRef, memory []byte) (*RecordID, error)

	// DeactivateObject creates deactivate object record in storage. Provided reference should be a reference to the head
	// of the object. If object is already deactivated, an error should be returned.
	//
	// Deactivated object cannot be changed.
	DeactivateObject(domain, request, obj RecordRef) (*RecordID, error)

	// UpdateObject creates amend object record in storage. Provided reference should be a reference to the head of the
	// object. Provided memory well be the new object memory.
	//
	// Returned reference will be the latest object state (exact) reference.
	UpdateObject(domain, request, obj RecordRef, memory []byte) (*RecordID, error)
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

// ClassDescriptor represents meta info required to fetch all object data.
type ClassDescriptor interface {
	// HeadRef returns head reference to represented class record.
	HeadRef() *RecordRef

	// StateID returns reference to represented class state record.
	StateID() *RecordID

	// CodeDescriptor returns descriptor for fetching class's code data.
	CodeDescriptor() CodeDescriptor
}

// ObjectDescriptor represents meta info required to fetch all object data.
type ObjectDescriptor interface {
	// HeadRef returns head reference to represented object record.
	HeadRef() *RecordRef

	// StateID returns reference to object state record.
	StateID() *RecordID

	// Memory fetches object memory from storage.
	Memory() []byte

	// ClassDescriptor returns descriptor for fetching object's class data.
	ClassDescriptor(state *RecordRef) (ClassDescriptor, error)

	// Children returns object's children references.
	Children(pulse *PulseNumber) (RefIterator, error)
}

// RefIterator is used for iteration over affined children(parts) of container.
type RefIterator interface {
	Next() (*RecordRef, error)
	HasNext() bool
}
