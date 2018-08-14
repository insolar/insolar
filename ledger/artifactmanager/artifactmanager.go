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
	"github.com/pkg/errors"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
)

// ArtifactManager is a high level storage interface.
type ArtifactManager interface {
	// SetArchPref allows to set a list of preferred VM architectures. When returning classes storage will return
	// compiled code according to this preferences.
	SetArchPref(pref []record.ArchType)

	// GetExactObj returns code and memory of provided object/class state. Deactivation records should be ignored
	// (e.g. object considered to be active).
	//
	// This method is used by validator to fetch the exact state of the object that was used by the executor.
	GetExactObj(classRef, objectRef record.Reference) ([]byte, record.Memory, error)

	// GetLatestObj returns descriptors for latest known state of the object/class known to the storage. The caller
	// should provide latest known states of the object/class known to it. If the object or the class is deactivated,
	// an error should be returned.
	//
	// Returned descriptors will provide methods for fetching migrations and appends relative to the provided states.
	GetLatestObj(objectRef, storedClassState, storedObjState record.Reference) (
		*ClassDescriptor, *ObjectDescriptor, error,
	)

	// DeployCode creates new code record in storage.
	DeployCode(requestRef record.Reference, codeMap map[record.ArchType][]byte) (*record.Reference, error)

	// ActivateClass creates activate class record in storage. Provided code reference will be used as a class code
	// and memory as the default memory for class objects.
	//
	// Activation reference will be this class'es identifier and referred as "class head".
	ActivateClass(requestRef, codeRef record.Reference, memory record.Memory) (*record.Reference, error)

	// DeactivateClass creates deactivate record in storage. Provided reference should be a reference to the head of
	// the class. If class is already deactivated, an error should be returned.
	//
	// Deactivated class cannot be changed or instantiate objects.
	DeactivateClass(requestRef, classRef record.Reference) (*record.Reference, error)

	// UpdateClass creates amend class record in storage. Provided reference should be a reference to the head of
	// the class. Migrations are references to code records.
	//
	// Migration code will be executed by VM to migrate objects memory in the order they appear in provided slice.
	UpdateClass(requestRef, classRef, codeRef record.Reference, migrationRefs []record.Reference) (
		*record.Reference, error,
	)

	// ActivateObj creates activate object record in storage. Provided class reference will be used as objects class
	// memory as memory of crated object. If memory is not provided, the class default memory will be used.
	//
	// Activation reference will be this object's identifier and referred as "object head".
	ActivateObj(requestRef, classRef record.Reference, memory record.Memory) (*record.Reference, error)

	// DeactivateObj creates deactivate object record in storage. Provided reference should be a reference to the head
	// of the object. If object is already deactivated, an error should be returned.
	//
	// Deactivated object cannot be changed.
	DeactivateObj(requestRef, objRef record.Reference) (*record.Reference, error)

	// UpdateObj creates amend object record in storage. Provided reference should be a reference to the head of the
	// object. Provided memory well be the new object memory.
	//
	// This will nullify all the object's append delegates. VM is responsible for collecting all appends and adding
	// them to the new memory manually if its required.
	UpdateObj(requestRef, objRef record.Reference, memory record.Memory) (*record.Reference, error)

	// AppendObjDelegate creates append object record in storage. Provided reference should be a reference to the head
	// of the object. Provided memory well be used as append delegate memory.
	//
	// Object's delegates will be provided by GetLatestObj. Any object update will nullify all the object's append
	// delegates. VM is responsible for collecting all appends and adding them to the new memory manually if its
	// required.
	AppendObjDelegate(requestRef, objRef record.Reference, memory record.Memory) (*record.Reference, error)
}

// LedgerArtifactManager provides concrete API to storage for virtual processing module
type LedgerArtifactManager struct {
	storer   storage.LedgerStorer
	archPref []record.ArchType
}

func (m *LedgerArtifactManager) checkRequestRecord(requestRef *record.Reference) error {
	// TODO: implement request check
	return nil
}

func (m *LedgerArtifactManager) getCodeRecord(codeRef record.Reference) (*record.CodeRecord, error) {
	rec, err := m.storer.GetRecord(&codeRef)
	if err != nil {
		return nil, errors.Wrap(err, "code record is not found")
	}
	codeRec, ok := rec.(*record.CodeRecord)
	if !ok {
		return nil, errors.New("provided reference is not a code reference")
	}
	return codeRec, nil
}

func (m *LedgerArtifactManager) getCodeRecordCode(codeRef record.Reference) ([]byte, error) {
	codeRec, err := m.getCodeRecord(codeRef)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve code record")
	}
	code, err := codeRec.GetCode(m.archPref)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve code")
	}

	return code, nil
}

func (m *LedgerArtifactManager) storeRecord(rec record.Record) (*record.Reference, error) {
	ref, err := m.storer.SetRecord(rec)
	if err != nil {
		return nil, errors.Wrap(err, "record store failed")
	}
	return ref, nil
}

func (m *LedgerArtifactManager) getActiveClass(classRef record.Reference) (
	*record.ClassActivateRecord, *record.ClassAmendRecord, *index.ClassLifeline, error,
) {
	classRecord, err := m.storer.GetRecord(&classRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "class record is not found")
	}
	activateRec, ok := classRecord.(*record.ClassActivateRecord)
	if !ok {
		return nil, nil, nil, errors.New("provided reference is not a class record")
	}
	classIndex, err := m.storer.GetClassIndex(&classRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "inconsistent class index")
	}
	latestClassRecord, err := m.storer.GetRecord(&classIndex.LatestStateRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "latest class record is not found")
	}
	if _, ok := latestClassRecord.(*record.DeactivationRecord); ok {
		return nil, nil, nil, errors.New("class is deactivated")
	}
	amendRecord, ok := latestClassRecord.(*record.ClassAmendRecord)
	if classRef.IsNotEqual(classIndex.LatestStateRef) && !ok {
		return nil, nil, nil, errors.New("wrong index record")
	}

	return activateRec, amendRecord, classIndex, nil
}

func (m *LedgerArtifactManager) getActiveObject(objRef record.Reference) (
	*record.ObjectActivateRecord, *record.ObjectAmendRecord, *index.ObjectLifeline, error,
) {
	objRecord, err := m.storer.GetRecord(&objRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "object record is not found")
	}
	activateRec, ok := objRecord.(*record.ObjectActivateRecord)
	if !ok {
		return nil, nil, nil, errors.New("provided reference is not an object record")
	}

	objIndex, err := m.storer.GetObjectIndex(&objRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "inconsistent object index")
	}
	latestObjRecord, err := m.storer.GetRecord(&objIndex.LatestStateRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "latest object record is not found")
	}
	if _, ok := latestObjRecord.(*record.DeactivationRecord); ok {
		return nil, nil, nil, errors.New("object is deactivated")
	}
	amendRecord, ok := latestObjRecord.(*record.ObjectAmendRecord)
	if objRef.IsNotEqual(objIndex.LatestStateRef) && !ok {
		return nil, nil, nil, errors.New("wrong index record")
	}

	return activateRec, amendRecord, objIndex, nil
}

func (m *LedgerArtifactManager) SetArchPref(pref []record.ArchType) {
	m.archPref = pref
}

// DeployCode deploys new code to storage (CodeRecord).
func (m *LedgerArtifactManager) DeployCode(
	requestRef record.Reference, codeMap map[record.ArchType][]byte,
) (*record.Reference, error) {
	err := m.checkRequestRecord(&requestRef)
	if err != nil {
		return nil, err
	}

	rec := record.CodeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
		},
		TargetedCode: codeMap,
	}
	return m.storeRecord(&rec)
}

// ActivateClass activates class from given code (ClassActivateRecord).
func (m *LedgerArtifactManager) ActivateClass(
	requestRef, codeRef record.Reference, memory record.Memory,
) (*record.Reference, error) {
	err := m.checkRequestRecord(&requestRef)
	if err != nil {
		return nil, err
	}
	_, err = m.getCodeRecord(codeRef)
	if err != nil {
		return nil, err
	}

	rec := record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
		},
		CodeRecord:    codeRef,
		DefaultMemory: memory,
	}
	classRef, err := m.storeRecord(&rec)
	if err != nil {
		return nil, err
	}
	err = m.storer.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *classRef,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to store lifeline index")
	}

	return classRef, nil
}

// DeactivateClass deactivates class (DeactivationRecord)
func (m *LedgerArtifactManager) DeactivateClass(
	requestRef, classRef record.Reference,
) (*record.Reference, error) {
	err := m.checkRequestRecord(&requestRef)
	if err != nil {
		return nil, err
	}

	_, _, classIndex, err := m.getActiveClass(classRef)
	if err != nil {
		return nil, err
	}

	rec := record.DeactivationRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
			HeadRecord: classRef,
			AmendedRecord: record.Reference{
				Domain: classRef.Domain,
				Record: classIndex.LatestStateRef.Record,
			},
		},
	}
	deactivationRef, err := m.storeRecord(&rec)
	if err != nil {
		return nil, errors.New("failed to store deactivation record")
	}
	classIndex.LatestStateRef = *deactivationRef
	err = m.storer.SetClassIndex(&classRef, classIndex)
	if err != nil {
		// TODO: add transaction
		return nil, errors.New("failed to store lifeline index")
	}

	return deactivationRef, nil
}

// UpdateClass allows to change class code etc. (ClassAmendRecord).
func (m *LedgerArtifactManager) UpdateClass(
	requestRef, classRef, codeRef record.Reference, migrationRefs []record.Reference,
) (*record.Reference, error) {
	err := m.checkRequestRecord(&requestRef)
	if err != nil {
		return nil, err
	}

	_, _, classIndex, err := m.getActiveClass(classRef)
	if err != nil {
		return nil, err
	}

	_, err = m.getCodeRecord(codeRef)
	if err != nil {
		return nil, err
	}
	for _, migrationRef := range migrationRefs {
		_, err = m.getCodeRecord(migrationRef)
		if err != nil {
			return nil, errors.Wrap(err, "invalid migrations")
		}
	}

	rec := record.ClassAmendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
			HeadRecord: classRef,
			AmendedRecord: record.Reference{
				Domain: classRef.Domain,
				Record: classIndex.LatestStateRef.Record,
			},
		},
		NewCode:    codeRef,
		Migrations: migrationRefs,
	}

	amendRef, err := m.storeRecord(&rec)
	if err != nil {
		return nil, errors.New("failed to store amend record")
	}
	classIndex.LatestStateRef = *amendRef
	classIndex.AmendRefs = append(classIndex.AmendRefs, *amendRef)
	err = m.storer.SetClassIndex(&classRef, classIndex)
	if err != nil {
		// TODO: add transaction
		return nil, errors.New("failed to store lifeline index")
	}

	return amendRef, nil
}

// ActivateObj creates and activates new object from given class (ObjectActivateRecord).
func (m *LedgerArtifactManager) ActivateObj(
	requestRef, classRef record.Reference, memory record.Memory,
) (*record.Reference, error) {
	err := m.checkRequestRecord(&requestRef)
	if err != nil {
		return nil, err
	}

	_, _, _, err = m.getActiveClass(classRef)
	if err != nil {
		return nil, err
	}

	rec := record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
		},
		ClassActivateRecord: classRef,
		Memory:              memory,
	}

	objRef, err := m.storeRecord(&rec)
	if err != nil {
		return nil, err
	}
	err = m.storer.SetObjectIndex(objRef, &index.ObjectLifeline{
		ClassRef:       classRef,
		LatestStateRef: *objRef,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to store lifeline index")
	}

	return objRef, nil
}

// DeactivateObj deactivates object (DeactivationRecord).
func (m *LedgerArtifactManager) DeactivateObj(requestRef, objRef record.Reference) (*record.Reference, error) {
	err := m.checkRequestRecord(&requestRef)
	if err != nil {
		return nil, err
	}

	_, _, objIndex, err := m.getActiveObject(objRef)
	if err != nil {
		return nil, err
	}

	rec := record.DeactivationRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
			HeadRecord: objRef,
			AmendedRecord: record.Reference{
				Domain: objRef.Domain,
				Record: objIndex.LatestStateRef.Record,
			},
		},
	}
	deactivationRef, err := m.storeRecord(&rec)
	if err != nil {
		return nil, errors.New("failed to store deactivation record")
	}
	objIndex.LatestStateRef = *deactivationRef
	err = m.storer.SetObjectIndex(&objRef, objIndex)
	if err != nil {
		// TODO: add transaction
		return nil, errors.New("failed to store lifeline index")
	}
	return deactivationRef, nil
}

// UpdateObj allows to change object state (ObjectAmendRecord).
func (m *LedgerArtifactManager) UpdateObj(
	requestRef, objRef record.Reference, memory record.Memory,
) (*record.Reference, error) {
	err := m.checkRequestRecord(&requestRef)
	if err != nil {
		return nil, err
	}

	_, _, objIndex, err := m.getActiveObject(objRef)
	if err != nil {
		return nil, err
	}

	rec := record.ObjectAmendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
			HeadRecord: objRef,
			AmendedRecord: record.Reference{
				Domain: objRef.Domain,
				Record: objIndex.LatestStateRef.Record,
			},
		},
		NewMemory: memory,
	}

	amendRef, err := m.storeRecord(&rec)
	if err != nil {
		return nil, errors.New("failed to store amend record")
	}
	objIndex.LatestStateRef = *amendRef
	objIndex.AppendRefs = []record.Reference{}
	err = m.storer.SetObjectIndex(&objRef, objIndex)
	if err != nil {
		// TODO: add transaction
		return nil, errors.New("failed to store lifeline index")
	}
	return amendRef, nil
}

func (m *LedgerArtifactManager) AppendObjDelegate(
	requestRef, objRef record.Reference, memory record.Memory,
) (*record.Reference, error) {
	err := m.checkRequestRecord(&requestRef)
	if err != nil {
		return nil, err
	}

	_, _, objIndex, err := m.getActiveObject(objRef)
	if err != nil {
		return nil, err
	}

	rec := record.ObjectAppendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
			HeadRecord: objRef,
			AmendedRecord: record.Reference{
				Domain: objRef.Domain,
				Record: objIndex.LatestStateRef.Record,
			},
		},
		AppendMemory: memory,
	}

	appendRef, err := m.storeRecord(&rec)
	if err != nil {
		return nil, errors.New("failed to store append record")
	}
	objIndex.AppendRefs = append(objIndex.AppendRefs, *appendRef)
	err = m.storer.SetObjectIndex(&objRef, objIndex)
	if err != nil {
		// TODO: add transaction
		return nil, errors.New("failed to store lifeline index")
	}
	return appendRef, nil
}

func (m *LedgerArtifactManager) GetExactObj(
	classState, objectState record.Reference,
) ([]byte, record.Memory, error) {
	classRec, err := m.storer.GetRecord(&classState)
	if err != nil {
		return nil, nil, errors.Wrap(err, "class record not found")
	}
	var codeRef record.Reference
	var classHeadRef record.Reference
	switch rec := classRec.(type) {
	case *record.ClassActivateRecord:
		codeRef = rec.CodeRecord
		classHeadRef = classState
	case *record.ClassAmendRecord:
		codeRef = rec.NewCode
		classHeadRef = rec.HeadRecord
	default:
		return nil, nil, errors.New("wrong class reference")
	}
	code, err := m.getCodeRecordCode(codeRef)
	if err != nil {
		return nil, nil, err
	}

	objectRec, err := m.storer.GetRecord(&objectState)
	if err != nil {
		return nil, nil, errors.Wrap(err, "object record not found")
	}
	var memory record.Memory
	var objectHeadRef record.Reference
	switch rec := objectRec.(type) {
	case *record.ObjectActivateRecord:
		memory = rec.Memory
		objectHeadRef = objectState
	case *record.ObjectAmendRecord:
		memory = rec.NewMemory
		objectHeadRef = rec.HeadRecord
	default:
		return nil, nil, errors.New("wrong object reference")
	}
	objectIndex, err := m.storer.GetObjectIndex(&objectHeadRef)
	if err != nil {
		return nil, nil, errors.Wrap(err, "object index not found")
	}

	if objectIndex.ClassRef.IsNotEqual(classHeadRef) {
		return nil, nil, errors.New("the object does not belong to the class")
	}

	return code, memory, nil
}

func (m *LedgerArtifactManager) GetLatestObj(
	objectRef, storedClassState, storedObjState record.Reference,
) (*ClassDescriptor, *ObjectDescriptor, error) {
	var (
		class  *ClassDescriptor
		object *ObjectDescriptor
	)

	objActivateRec, objStateRec, objIndex, err := m.getActiveObject(objectRef)
	if err != nil {
		return nil, nil, err
	}
	classActivateRec, classStateRec, classIndex, err := m.getActiveClass(objIndex.ClassRef)
	if err != nil {
		return nil, nil, err
	}

	if storedClassState.IsNotEqual(classIndex.LatestStateRef) {
		class = &ClassDescriptor{
			StateRef: record.Reference{
				Domain: storedClassState.Domain,
				Record: classIndex.LatestStateRef.Record,
			},

			manager:           m,
			fromState:         storedClassState,
			activateRecord:    classActivateRec,
			latestAmendRecord: classStateRec,
			lifelineIndex:     classIndex,
		}
	}

	if storedObjState.IsNotEqual(objIndex.LatestStateRef) {
		object = &ObjectDescriptor{
			StateRef: record.Reference{
				Domain: storedObjState.Domain,
				Record: objIndex.LatestStateRef.Record,
			},

			manager:           m,
			activateRecord:    objActivateRec,
			latestAmendRecord: objStateRec,
			lifelineIndex:     objIndex,
		}
	}

	return class, object, nil
}
