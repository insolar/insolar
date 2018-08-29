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
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
)

// LedgerArtifactManager provides concrete API to storage for processing module
type LedgerArtifactManager struct {
	storer   storage.LedgerStorer
	archPref []core.MachineType
}

func (m *LedgerArtifactManager) checkRequestRecord(requestRef *record.Reference) error {
	// TODO: implement request check
	return nil
}

func (m *LedgerArtifactManager) getCodeRecord(codeRef record.Reference) (*record.CodeRecord, error) {
	rec, err := m.storer.GetRecord(&codeRef)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve code record")
	}
	codeRec, ok := rec.(*record.CodeRecord)
	if !ok {
		return nil, errors.Wrap(ErrInvalidRef, "failed to retrieve code record")
	}
	return codeRec, nil
}

func (m *LedgerArtifactManager) getCodeRecordCode(codeRef record.Reference) ([]byte, error) {
	codeRec, err := m.getCodeRecord(codeRef)
	if err != nil {
		return nil, err
	}
	code, err := codeRec.GetCode(m.archPref)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve code")
	}

	return code, nil
}

func (m *LedgerArtifactManager) getActiveClass(classRef record.Reference) (
	*record.ClassActivateRecord, *record.ClassAmendRecord, *index.ClassLifeline, error,
) {
	classRecord, err := m.storer.GetRecord(&classRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to retrieve class record")
	}
	activateRec, isClassRec := classRecord.(*record.ClassActivateRecord)
	if !isClassRec {
		return nil, nil, nil, errors.Wrap(ErrInvalidRef, "failed to retrieve class record")
	}
	classIndex, err := m.storer.GetClassIndex(&classRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "inconsistent class index")
	}
	latestClassRecord, err := m.storer.GetRecord(&classIndex.LatestStateRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "inconsistent class index")
	}
	if _, isDeactivated := latestClassRecord.(*record.DeactivationRecord); isDeactivated {
		return nil, nil, nil, ErrClassDeactivated
	}
	amendRecord, isLatestStateAmend := latestClassRecord.(*record.ClassAmendRecord)
	if classRef.IsNotEqual(classIndex.LatestStateRef) && !isLatestStateAmend {
		return nil, nil, nil, errors.Wrap(ErrInconsistentIndex, "inconsistent class index")
	}

	return activateRec, amendRecord, classIndex, nil
}

func (m *LedgerArtifactManager) getActiveObject(objRef record.Reference) (
	*record.ObjectActivateRecord, *record.ObjectAmendRecord, *index.ObjectLifeline, error,
) {
	objRecord, err := m.storer.GetRecord(&objRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to retrieve object record")
	}
	activateRec, isObjectRec := objRecord.(*record.ObjectActivateRecord)
	if !isObjectRec {
		return nil, nil, nil, errors.Wrap(ErrInvalidRef, "failed to retrieve object record")
	}

	objIndex, err := m.storer.GetObjectIndex(&objRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "inconsistent object index")
	}
	latestObjRecord, err := m.storer.GetRecord(&objIndex.LatestStateRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "inconsistent object index")
	}
	if _, isDeactivated := latestObjRecord.(*record.DeactivationRecord); isDeactivated {
		return nil, nil, nil, ErrObjectDeactivated
	}
	amendRecord, isLatestAmend := latestObjRecord.(*record.ObjectAmendRecord)
	if objRef.IsNotEqual(objIndex.LatestStateRef) && !isLatestAmend {
		return nil, nil, nil, errors.Wrap(ErrInconsistentIndex, "inconsistent object index")
	}

	return activateRec, amendRecord, objIndex, nil
}

// SetArchPref stores a list of preferred VM architectures memory.
//
// When returning classes storage will return compiled code according to this preferences. VM is responsible for
// calling this method before fetching object in a new process. If preference is not provided, object getters will
// return an error.
func (m *LedgerArtifactManager) SetArchPref(pref []core.MachineType) {
	m.archPref = pref
}

// DeployCode creates new code record in storage.
//
// Code records are used to activate class or as migration code for an object.
func (m *LedgerArtifactManager) DeployCode(
	domain, request core.RecordRef, codeMap map[core.MachineType][]byte,
) (core.RecordRef, error) {
	domainRef := record.Bytes2Reference(domain)
	requestRef := record.Bytes2Reference(request)

	err := m.checkRequestRecord(&requestRef)
	if err != nil {
		return nil, err
	}

	rec := record.CodeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  domainRef,
					RequestRecord: requestRef,
				},
			},
		},
		TargetedCode: codeMap,
	}
	codeRef, err := m.storer.SetRecord(&rec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store record")
	}
	return codeRef.Bytes(), nil
}

// ActivateClass creates activate class record in storage. Provided code reference will be used as a class code
// and memory as the default memory for class objects.
//
// Activation reference will be this class'es identifier and referred as "class head".
func (m *LedgerArtifactManager) ActivateClass(
	domain, request, code core.RecordRef, memory []byte,
) (core.RecordRef, error) {
	domainRef := record.Bytes2Reference(domain)
	requestRef := record.Bytes2Reference(request)
	codeRef := record.Bytes2Reference(code)

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
					DomainRecord:  domainRef,
					RequestRecord: requestRef,
				},
			},
		},
		CodeRecord:    codeRef,
		DefaultMemory: memory,
	}
	classRef, err := m.storer.SetRecord(&rec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store record")
	}
	err = m.storer.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *classRef,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to store lifeline index")
	}

	return classRef.Bytes(), nil
}

// DeactivateClass creates deactivate record in storage. Provided reference should be a reference to the head of
// the class. If class is already deactivated, an error should be returned.
//
// Deactivated class cannot be changed or instantiate objects.
func (m *LedgerArtifactManager) DeactivateClass(
	domain, request, class core.RecordRef,
) (core.RecordRef, error) {
	domainRef := record.Bytes2Reference(domain)
	requestRef := record.Bytes2Reference(request)
	classRef := record.Bytes2Reference(class)

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
					DomainRecord:  domainRef,
					RequestRecord: requestRef,
				},
			},
			HeadRecord:    classRef,
			AmendedRecord: classIndex.LatestStateRef,
		},
	}
	deactivationRef, err := m.storer.SetRecord(&rec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store record")
	}
	classIndex.LatestStateRef = *deactivationRef
	err = m.storer.SetClassIndex(&classRef, classIndex)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store lifeline index")
	}

	return deactivationRef.Bytes(), nil
}

// UpdateClass creates amend class record in storage. Provided reference should be a reference to the head of
// the class. Migrations are references to code records.
//
// Migration code will be executed by VM to migrate objects memory in the order they appear in provided slice.
func (m *LedgerArtifactManager) UpdateClass(
	domain, request, class, code core.RecordRef, migrations []core.RecordRef,
) (core.RecordRef, error) {
	domainRef := record.Bytes2Reference(domain)
	requestRef := record.Bytes2Reference(request)
	classRef := record.Bytes2Reference(class)
	codeRef := record.Bytes2Reference(code)
	migrationRefs := make([]record.Reference, 0, len(migrations))
	for _, migration := range migrations {
		migrationRefs = append(migrationRefs, record.Bytes2Reference(migration))
	}

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
					DomainRecord:  domainRef,
					RequestRecord: requestRef,
				},
			},
			HeadRecord:    classRef,
			AmendedRecord: classIndex.LatestStateRef,
		},
		NewCode:    codeRef,
		Migrations: migrationRefs,
	}

	amendRef, err := m.storer.SetRecord(&rec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store record")
	}
	classIndex.LatestStateRef = *amendRef
	classIndex.AmendRefs = append(classIndex.AmendRefs, *amendRef)
	err = m.storer.SetClassIndex(&classRef, classIndex)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store lifeline index")
	}

	return amendRef.Bytes(), nil
}

// ActivateObj creates activate object record in storage. Provided class reference will be used as objects class
// memory as memory of crated object. If memory is not provided, the class default memory will be used.
//
// Activation reference will be this object's identifier and referred as "object head".
func (m *LedgerArtifactManager) ActivateObj(
	domain, request, class core.RecordRef, memory []byte,
) (core.RecordRef, error) {
	domainRef := record.Bytes2Reference(domain)
	requestRef := record.Bytes2Reference(request)
	classRef := record.Bytes2Reference(class)

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
					DomainRecord:  domainRef,
					RequestRecord: requestRef,
				},
			},
		},
		ClassActivateRecord: classRef,
		Memory:              memory,
	}

	objRef, err := m.storer.SetRecord(&rec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store record")
	}
	err = m.storer.SetObjectIndex(objRef, &index.ObjectLifeline{
		ClassRef:       classRef,
		LatestStateRef: *objRef,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to store lifeline index")
	}

	return objRef.Bytes(), nil
}

// DeactivateObj creates deactivate object record in storage. Provided reference should be a reference to the head
// of the object. If object is already deactivated, an error should be returned.
//
// Deactivated object cannot be changed.
func (m *LedgerArtifactManager) DeactivateObj(
	domain, request, obj core.RecordRef,
) (core.RecordRef, error) {
	domainRef := record.Bytes2Reference(domain)
	requestRef := record.Bytes2Reference(request)
	objRef := record.Bytes2Reference(obj)

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
					DomainRecord:  domainRef,
					RequestRecord: requestRef,
				},
			},
			HeadRecord:    objRef,
			AmendedRecord: objIndex.LatestStateRef,
		},
	}
	deactivationRef, err := m.storer.SetRecord(&rec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store record")
	}
	objIndex.LatestStateRef = *deactivationRef
	err = m.storer.SetObjectIndex(&objRef, objIndex)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store lifeline index")
	}
	return deactivationRef.Bytes(), nil
}

// UpdateObj creates amend object record in storage. Provided reference should be a reference to the head of the
// object. Provided memory well be the new object memory.
//
// This will nullify all the object's append delegates. VM is responsible for collecting all appends and adding
// them to the new memory manually if its required.
func (m *LedgerArtifactManager) UpdateObj(
	domain, request, obj core.RecordRef, memory []byte,
) (core.RecordRef, error) {
	domainRef := record.Bytes2Reference(domain)
	requestRef := record.Bytes2Reference(request)
	objRef := record.Bytes2Reference(obj)

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
					DomainRecord:  domainRef,
					RequestRecord: requestRef,
				},
			},
			HeadRecord:    objRef,
			AmendedRecord: objIndex.LatestStateRef,
		},
		NewMemory: memory,
	}

	amendRef, err := m.storer.SetRecord(&rec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store record")
	}
	objIndex.LatestStateRef = *amendRef
	objIndex.AppendRefs = []record.Reference{}
	err = m.storer.SetObjectIndex(&objRef, objIndex)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store lifeline index")
	}
	return amendRef.Bytes(), nil
}

// AppendObjDelegate creates append object record in storage. Provided reference should be a reference to the head
// of the object. Provided memory well be used as append delegate memory.
//
// Object's delegates will be provided by GetLatestObj. Any object update will nullify all the object's append
// delegates. VM is responsible for collecting all appends and adding them to the new memory manually if its
// required.
func (m *LedgerArtifactManager) AppendObjDelegate(
	domain, request, obj core.RecordRef, memory []byte,
) (core.RecordRef, error) {
	domainRef := record.Bytes2Reference(domain)
	requestRef := record.Bytes2Reference(request)
	objRef := record.Bytes2Reference(obj)

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
					DomainRecord:  domainRef,
					RequestRecord: requestRef,
				},
			},
			HeadRecord:    objRef,
			AmendedRecord: objIndex.LatestStateRef,
		},
		AppendMemory: memory,
	}

	appendRef, err := m.storer.SetRecord(&rec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store record")
	}
	objIndex.AppendRefs = append(objIndex.AppendRefs, *appendRef)
	err = m.storer.SetObjectIndex(&objRef, objIndex)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store lifeline index")
	}
	return appendRef.Bytes(), nil
}

// GetExactObj returns code and memory of provided object/class state. Deactivation records should be ignored
// (e.g. object considered to be active).
//
// This method is used by validator to fetch the exact state of the object that was used by the executor.
func (m *LedgerArtifactManager) GetExactObj( // nolint: gocyclo
	classState, objectState core.RecordRef,
) ([]byte, []byte, error) {
	classRef := record.Bytes2Reference(classState)
	objRef := record.Bytes2Reference(objectState)

	// Fetching class data
	classRec, err := m.storer.GetRecord(&classRef)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to retrieve class record")
	}

	var codeRef record.Reference
	var classHeadRef record.Reference
	switch rec := classRec.(type) {
	case *record.ClassActivateRecord:
		codeRef = rec.CodeRecord
		classHeadRef = classRef
	case *record.ClassAmendRecord:
		codeRef = rec.NewCode
		classHeadRef = rec.HeadRecord
	default:
		return nil, nil, errors.Wrap(ErrInvalidRef, "failed to retrieve class record")
	}

	code, err := m.getCodeRecordCode(codeRef)
	if err != nil {
		return nil, nil, err
	}

	// Fetching object data
	objectRec, err := m.storer.GetRecord(&objRef)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to retrieve object record")
	}

	var memory []byte
	var objectHeadRef record.Reference
	switch rec := objectRec.(type) {
	case *record.ObjectActivateRecord:
		memory = rec.Memory
		objectHeadRef = objRef
	case *record.ObjectAmendRecord:
		memory = rec.NewMemory
		objectHeadRef = rec.HeadRecord
	default:
		return nil, nil, errors.Wrap(ErrInvalidRef, "failed to retrieve object record")
	}

	objectIndex, err := m.storer.GetObjectIndex(&objectHeadRef)
	if err != nil {
		return nil, nil, errors.Wrap(err, "inconsistent object index")
	}

	// Checking if the object belongs to the class
	if objectIndex.ClassRef.IsNotEqual(classHeadRef) {
		return nil, nil, ErrWrongObject
	}

	return code, memory, nil
}

// GetLatestObj returns descriptors for latest known state of the object/class known to the storage. The caller
// should provide latest known states of the object/class known to it. If the object or the class is deactivated,
// an error should be returned.
//
// Returned descriptors will provide methods for fetching migrations and appends relative to the provided states.
func (m *LedgerArtifactManager) GetLatestObj(
	obj, storedClassState, storedObjState core.RecordRef,
) (core.ClassDescriptor, core.ObjectDescriptor, error) {
	objRef := record.Bytes2Reference(obj)
	storedClassRef := record.Bytes2Reference(storedClassState)
	storedObjRef := record.Bytes2Reference(storedObjState)

	var (
		class  *ClassDescriptor
		object *ObjectDescriptor
	)

	objActivateRec, objStateRec, objIndex, err := m.getActiveObject(objRef)
	if err != nil {
		return nil, nil, err
	}
	classActivateRec, classStateRec, classIndex, err := m.getActiveClass(objIndex.ClassRef)
	if err != nil {
		return nil, nil, err
	}

	// if provided reference is the last reference in the lifeline, we can return nil
	if storedClassRef.IsNotEqual(classIndex.LatestStateRef) {
		class = &ClassDescriptor{
			StateRef: classIndex.LatestStateRef,

			manager:           m,
			fromState:         storedClassRef,
			activateRecord:    classActivateRec,
			latestAmendRecord: classStateRec,
			lifelineIndex:     classIndex,
		}
	}

	// if provided reference is the last reference in the lifeline, we can return nil
	if storedObjRef.IsNotEqual(objIndex.LatestStateRef) {
		object = &ObjectDescriptor{
			StateRef: objIndex.LatestStateRef,

			manager:           m,
			activateRecord:    objActivateRec,
			latestAmendRecord: objStateRec,
			lifelineIndex:     objIndex,
		}
	}

	return class, object, nil
}
