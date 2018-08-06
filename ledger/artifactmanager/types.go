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

package artifactmanager

import "github.com/insolar/insolar/ledger/record"

// ArtifactManager is a high level storage interface.
type ArtifactManager interface {
	// SetArchPref allows to set list of preferred VM architectures.
	// When returning classes storage will return compiled code
	// according to this preferences.
	SetArchPref(pref []record.ArchType)

	// GetObj returns object by reference.
	GetObj(
		// ref is target object reference
		ref record.Reference,
		// lastClassRef is reference to class that is already deployed to VM. Can be nil.
		lastClassRef record.Reference,
		// lastObjRef is reference to object that is already deployed to VM. Can be nil.
		lastObjRef record.Reference,
	) (
		ClassDescr,
		ObjDescr,
		error,
	)

	// DeployCode deploys new code to storage (CodeRecord).
	DeployCode(requestRef record.Reference) (record.Reference, error)

	// ActivateClass activates class from given code (ClassActivateRecord).
	ActivateClass(requestRef, codeRef record.Reference, memory record.Memory) (record.Reference, error)

	// DeactivateClass deactivates class (DeactivationRecord)
	DeactivateClass(requestRef, classRef record.Reference) (record.Reference, error)

	// UpdateClass allows to change class code etc. (ClassAmendRecord).
	UpdateClass(
		requestRef,
		classRef record.Reference,
		migrationRefs []record.Reference,
	) (record.Reference, error)

	// ActivateObj creates and activates new object from given class (ObjectActivateRecord).
	ActivateObj(record record.ObjectActivateRecord) (record.Reference, error)

	// DeactivateObj deactivates object (DeactivationRecord).
	DeactivateObj(ref record.Reference) (record.Reference, error)

	// UpdateObj allows to change object state (ObjectAmendRecord).
	UpdateObj(ref record.Reference, memory record.Memory) (record.Reference, error)

	// AppendDelegate allows to append some class'es delegate to object (ObjectAppendRecord).
	AppendDelegate(ref record.Reference, delegate record.Memory) (record.Reference, error)
}

// ClassDescr contains class code and migration procedures if any.
type ClassDescr struct {
	// TODO: implement MachineBinaryCode
	Code       []byte
	Migrations [][]byte
}

// ObjDescr contains object memory and delegate appends if any.
type ObjDescr struct {
	ObjectMemory record.Memory   // nil if LastObjRef is actual
	Appends      []record.Memory // can be empty
}
