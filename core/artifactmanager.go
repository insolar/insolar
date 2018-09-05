/*
 *    Copyright 2018 INS Ecosystem
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

// ArtifactManager is a high level storage interface.
type ArtifactManager interface {
	// SetArchPref stores a list of preferred VM architectures memory.
	//
	// When returning classes storage will return compiled code according to this preferences. VM is responsible for
	// calling this method before fetching object in a new process. If preference is not provided, object getters will
	// return an error.
	SetArchPref(pref []MachineType)

	// GetCode returns code from code record by provided reference.
	//
	// This method is used by VM to fetch code for execution.
	GetCode(code RecordRef) (CodeDescriptor, error)

	// GetLatestObj returns descriptors for latest known state of the object/class known to the storage.
	// If the object or the class is deactivated, an error should be returned.
	//
	// Returned descriptors will provide methods for fetching migrations and appends relative to the provided states.
	GetLatestObj(head RecordRef) (ObjectDescriptor, error)

	// DeclareType creates new type record in storage.
	//
	// Type is a contract interface. It contains one method signature.
	DeclareType(domain, request RecordRef, typeDec []byte) (*RecordRef, error)

	// DeployCode creates new code record in storage.
	//
	// Code records are used to activate class or as migration code for an object.
	DeployCode(domain, request RecordRef, types []RecordRef, codeMap map[MachineType][]byte) (*RecordRef, error)

	// ActivateClass creates activate class record in storage. Provided code reference will be used as a class code
	// and memory as the default memory for class objects.
	//
	// Activation reference will be this class'es identifier and referred as "class head".
	ActivateClass(domain, request, code RecordRef, memory []byte) (*RecordRef, error)

	// DeactivateClass creates deactivate record in storage. Provided reference should be a reference to the head of
	// the class. If class is already deactivated, an error should be returned.
	//
	// Deactivated class cannot be changed or instantiate objects.
	DeactivateClass(domain, request, class RecordRef) (*RecordRef, error)

	// UpdateClass creates amend class record in storage. Provided reference should be a reference to the head of
	// the class. Migrations are references to code records.
	//
	// Migration code will be executed by VM to migrate objects memory in the order they appear in provided slice.
	UpdateClass(domain, request, class, code RecordRef, migrationRefs []RecordRef) (*RecordRef, error)

	// ActivateObj creates activate object record in storage. Provided class reference will be used as objects class
	// memory as memory of crated object. If memory is not provided, the class default memory will be used.
	//
	// Activation reference will be this object's identifier and referred as "object head".
	ActivateObj(domain, request, class RecordRef, memory []byte) (*RecordRef, error)

	// DeactivateObj creates deactivate object record in storage. Provided reference should be a reference to the head
	// of the object. If object is already deactivated, an error should be returned.
	//
	// Deactivated object cannot be changed.
	DeactivateObj(domain, request, obj RecordRef) (*RecordRef, error)

	// UpdateObj creates amend object record in storage. Provided reference should be a reference to the head of the
	// object. Provided memory well be the new object memory.
	//
	// This will nullify all the object's append delegates. VM is responsible for collecting all appends and adding
	// them to the new memory manually if its required.
	UpdateObj(domain, request, obj RecordRef, memory []byte) (*RecordRef, error)

	// AppendObjDelegate creates append object record in storage. Provided reference should be a reference to the head
	// of the object. Provided memory well be used as append delegate memory.
	//
	// Object's delegates will be provided by GetLatestObj. Any object update will nullify all the object's append
	// delegates. VM is responsible for collecting all appends and adding them to the new memory manually if its
	// required.
	AppendObjDelegate(domain, request, obj RecordRef, memory []byte) (*RecordRef, error)
}

// CodeDescriptor represents meta info required to fetch all code data.
type CodeDescriptor interface {
	// Ref returns reference to represented code record.
	Ref() *RecordRef

	// Code fetches code from storage. Code will be fetched according to architecture preferences
	// set via SetArchPref in artifact manager. If preferences are not provided, an error will be returned.
	Code() ([]byte, error)
}

// ClassDescriptor represents meta info required to fetch all object data.
type ClassDescriptor interface {
	// HeadRef returns head reference to represented class record.
	HeadRef() *RecordRef

	// StateRef returns reference to represented class state record.
	StateRef() *RecordRef

	// CodeDescriptor returns descriptor for fetching class's code data.
	CodeDescriptor() (CodeDescriptor, error)
}

// ObjectDescriptor represents meta info required to fetch all object data.
type ObjectDescriptor interface {
	// HeadRef returns head reference to represented object record.
	HeadRef() *RecordRef

	// StateRef returns reference to object state record.
	StateRef() *RecordRef

	// Memory fetches object memory from storage.
	Memory() ([]byte, error)

	// CodeDescriptor returns descriptor for fetching object's code data.
	CodeDescriptor() (CodeDescriptor, error)

	// ClassDescriptor returns descriptor for fetching object's class data.
	ClassDescriptor() (ClassDescriptor, error)
}
