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

import (
	"github.com/insolar/insolar/ledger/record"
)

// ArtifactManager is a high level storage interface.
type ArtifactManager interface {
	// SetArchPref allows to set list of preferred VM architectures.
	// When returning classes storage will return compiled code
	// according to this preferences.
	SetArchPref(pref []record.ArchType)

	// GetExactObj returns exact object data and code ref without calculating last state
	GetExactObj(classRef, objectRef record.Reference) ([]byte, record.Memory, error)

	// GetLatestObj returns object by reference.
	GetLatestObj(
		objectRef record.Reference,
		storedClassState record.Reference,
		storedObjState record.Reference,
	) (
		*ClassDescriptor,
		*ObjectDescriptor,
		error,
	)

	// DeployCode deploys new code to storage (CodeRecord).
	DeployCode(requestRef record.Reference, codeMap map[record.ArchType][]byte) (record.Reference, error)

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
	ActivateObj(requestRef, classRef record.Reference, memory record.Memory) (record.Reference, error)

	// DeactivateObj deactivates object (DeactivationRecord).
	DeactivateObj(requestRef, objRef record.Reference) (record.Reference, error)

	// UpdateObj allows to change object state (ObjectAmendRecord).
	UpdateObj(requestRef, objRef record.Reference, memory record.Memory) (record.Reference, error)

	// AppendDelegate allows to append some class'es delegate to object (ObjectAppendRecord).
	AppendObjDelegate(requestRef, objRef record.Reference, memory record.Memory) (record.Reference, error)
}
