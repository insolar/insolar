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

type JetRole int

const (
	RoleVirtualExecutor  = JetRole(iota + 1) // Role responsible for current pulse CPU operations.
	RoleVirtualValidator                     // Role responsible for previous pulse CPU operations.
	RoleLightExecutor                        // Role responsible for current pulse Disk operations.
	RoleLightValidator                       // Role responsible for previous pulse Disk operations.
	RoleHeavyExecutor                        // Role responsible for permanent Disk operations.
)

// Ledger is the global ledger handler. Other system parts communicate with ledger through it.
type Ledger interface {
	// GetArtifactManager returns artifact manager to work with.
	GetArtifactManager() ArtifactManager

	// GetJetCoordinator returns jet coordinator to work with.
	GetJetCoordinator() JetCoordinator

	// GetPulseManager returns pulse manager to work with.
	GetPulseManager() PulseManager

	HandleEvent(Event) (Reaction, error)
}

type JetCoordinator interface {
	// IsAuthorized checks for role on concrete pulse for the address.
	IsAuthorized(role JetRole, obj RecordRef, pulse PulseNumber, node RecordRef) (bool, error)

	// QueryRole returns node refs responsible for role bound operations for given object and pulse.
	QueryRole(role JetRole, obj RecordRef, pulse PulseNumber) ([]RecordRef, error)
}

// ArtifactManager is a high level storage interface.
type ArtifactManager interface {
	// RootRef returns the root record reference.
	//
	// Root record is the parent for all top-level records.
	RootRef() *RecordRef

	// GetCode returns code from code record by provided reference.
	//
	// This method is used by VM to fetch code for execution.
	GetCode(ref RecordRef, machinePref []MachineType) (CodeDescriptor, error)

	// GetClass returns descriptor for latest state of the class known to storage.
	// If the class is deactivated, an error should be returned.
	//
	// Returned descriptor will provide methods for fetching all related data.
	GetClass(head RecordRef, state *RecordRef) (ClassDescriptor, error)

	// GetObject returns descriptor for latest state of the object known to storage.
	// If the object or the class is deactivated, an error should be returned.
	//
	// Returned descriptor will provide methods for fetching all related data.
	GetObject(head RecordRef, state *RecordRef) (ObjectDescriptor, error)

	// GetDelegate returns provided object's delegate reference for provided class.
	//
	// Object delegate should be previously created for this object. If object delegate does not exist, an error will
	// be returned.
	GetDelegate(head, asClass RecordRef) (*RecordRef, error)

	// DeclareType creates new type record in storage.
	//
	// Type is a contract interface. It contains one method signature.
	DeclareType(domain, request RecordRef, typeDec []byte) (*RecordRef, error)

	// DeployCode creates new code record in storage.
	//
	// Code records are used to activate class or as migration code for an object.
	DeployCode(domain, request RecordRef, codeMap map[MachineType][]byte) (*RecordRef, error)

	// ActivateClass creates activate class record in storage. Provided code reference will be used as a class code.
	//
	// Activation reference will be this class'es identifier and referred as "class head".
	ActivateClass(domain, request RecordRef) (*RecordRef, error)

	// DeactivateClass creates deactivate record in storage. Provided reference should be a reference to the head of
	// the class. If class is already deactivated, an error should be returned.
	//
	// Deactivated class cannot be changed or instantiate objects.
	DeactivateClass(domain, request, class RecordRef) (*RecordRef, error)

	// UpdateClass creates amend class record in storage. Provided reference should be a reference to the head of
	// the class. Migrations are references to code records.
	//
	// Returned reference will be the latest class state (exact) reference. Migration code will be executed by VM to
	// migrate objects memory in the order they appear in provided slice.
	UpdateClass(domain, request, class, code RecordRef, migrationRefs []RecordRef) (*RecordRef, error)

	// ActivateObject creates activate object record in storage. Provided class reference will be used as object's class.
	// If memory is not provided, the class default memory will be used.
	//
	// Activation reference will be this object's identifier and referred as "object head".
	ActivateObject(domain, request, class, parent RecordRef, memory []byte) (*RecordRef, error)

	// ActivateObjectDelegate is similar to ActivateObj but it created object will be parent's delegate of provided class.
	ActivateObjectDelegate(domain, request, class, parent RecordRef, memory []byte) (*RecordRef, error)

	// DeactivateObject creates deactivate object record in storage. Provided reference should be a reference to the head
	// of the object. If object is already deactivated, an error should be returned.
	//
	// Deactivated object cannot be changed.
	DeactivateObject(domain, request, obj RecordRef) (*RecordRef, error)

	// UpdateObject creates amend object record in storage. Provided reference should be a reference to the head of the
	// object. Provided memory well be the new object memory.
	//
	// Returned reference will be the latest object state (exact) reference.
	UpdateObject(domain, request, obj RecordRef, memory []byte) (*RecordRef, error)
}

// CodeDescriptor represents meta info required to fetch all code data.
type CodeDescriptor interface {
	// Ref returns reference to represented code record.
	Ref() *RecordRef

	// MachineType fetches code from storage and returns first available machine type according to architecture
	// preferences.
	//
	// Code for returned machine type will be fetched by Code method.
	MachineType() MachineType

	// Code fetches code from storage. Code will be fetched according to architecture preferences
	// set via SetArchPref in artifact manager. If preferences are not provided, an error will be returned.
	Code() []byte
}

// ClassDescriptor represents meta info required to fetch all object data.
type ClassDescriptor interface {
	// HeadRef returns head reference to represented class record.
	HeadRef() *RecordRef

	// StateRef returns reference to represented class state record.
	StateRef() *RecordRef

	// CodeDescriptor returns descriptor for fetching class's code data.
	CodeDescriptor(machinePref []MachineType) (CodeDescriptor, error)
}

// ObjectDescriptor represents meta info required to fetch all object data.
type ObjectDescriptor interface {
	// HeadRef returns head reference to represented object record.
	HeadRef() *RecordRef

	// StateRef returns reference to object state record.
	StateRef() *RecordRef

	// Memory fetches object memory from storage.
	Memory() []byte

	// ClassDescriptor returns descriptor for fetching object's class data.
	ClassDescriptor(state *RecordRef) (ClassDescriptor, error)

	// Children returns object's children references.
	Children() RefIterator
}

type RefIterator interface {
	Next() (RecordRef, error)
	HasNext() bool
}
