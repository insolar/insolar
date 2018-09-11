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
	db       *storage.DB
	archPref []core.MachineType
}

func (m *LedgerArtifactManager) checkRequestRecord(s storage.Store, requestRef *record.Reference) error {
	// TODO: implement request check
	return nil
}

func (m *LedgerArtifactManager) getCodeRecord(s storage.Store, codeRef record.Reference) (*record.CodeRecord, error) {
	rec, err := s.GetRecord(&codeRef)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve code record")
	}
	codeRec, ok := rec.(*record.CodeRecord)
	if !ok {
		return nil, errors.Wrap(ErrInvalidRef, "failed to retrieve code record")
	}
	return codeRec, nil
}

func (m *LedgerArtifactManager) getCodeRecordCode(s storage.Store, codeRef record.Reference) ([]byte, core.MachineType, error) {
	codeRec, err := m.getCodeRecord(s, codeRef)
	if err != nil {
		return nil, core.MachineTypeNotExist, err
	}
	code, mt, err := codeRec.GetCode(m.archPref)
	if err != nil {
		return nil, mt, errors.Wrap(err, "failed to retrieve code")
	}

	return code, mt, nil
}

func (m *LedgerArtifactManager) getActiveClass(s storage.Store, classRef record.Reference) (
	*record.ClassActivateRecord, *record.ClassAmendRecord, *index.ClassLifeline, error,
) {
	classRecord, err := s.GetRecord(&classRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to retrieve class record")
	}
	activateRec, isClassRec := classRecord.(*record.ClassActivateRecord)
	if !isClassRec {
		return nil, nil, nil, errors.Wrap(ErrInvalidRef, "failed to retrieve class record")
	}
	classIndex, err := s.GetClassIndex(&classRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "inconsistent class index")
	}
	latestClassRecord, err := s.GetRecord(&classIndex.LatestStateRef)
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

func (m *LedgerArtifactManager) getActiveObject(s storage.Store, objRef record.Reference) (
	*record.ObjectActivateRecord, *record.ObjectAmendRecord, *index.ObjectLifeline, error,
) {
	objRecord, err := s.GetRecord(&objRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to retrieve object record")
	}
	activateRec, isObjectRec := objRecord.(*record.ObjectActivateRecord)
	if !isObjectRec {
		return nil, nil, nil, errors.Wrap(ErrInvalidRef, "failed to retrieve object record")
	}

	objIndex, err := s.GetObjectIndex(&objRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "inconsistent object index")
	}
	latestObjRecord, err := s.GetRecord(&objIndex.LatestStateRef)
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

// RootRef returns the root record reference.
//
// Root record is the parent for all top-level records.
func (m *LedgerArtifactManager) RootRef() *core.RecordRef {
	return m.db.RootRef().CoreRef()
}

// SetArchPref Stores a list of preferred VM architectures memory.
//
// When returning classes storage will return compiled code according to this preferences. VM is responsible for
// calling this method before fetching object in a new process. If preference is not provided, object getters will
// return an error.
func (m *LedgerArtifactManager) SetArchPref(pref []core.MachineType) {
	m.archPref = pref
}

// DeclareType creates new type record in storage.
//
// Type is a contract interface. It contains one method signature.
func (m *LedgerArtifactManager) DeclareType(
	domain, request core.RecordRef, typeDec []byte,
) (*core.RecordRef, error) {
	domainRef := record.Core2Reference(domain)
	requestRef := record.Core2Reference(request)

	err := m.checkRequestRecord(m.db, &requestRef)
	if err != nil {
		return nil, err
	}

	rec := record.TypeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  domainRef,
					RequestRecord: requestRef,
				},
			},
		},
		TypeDeclaration: typeDec,
	}
	codeRef, err := m.db.SetRecord(&rec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store record")
	}
	return codeRef.CoreRef(), nil
}

// DeployCode creates new code record in storage.
//
// Code records are used to activate class or as migration code for an object.
func (m *LedgerArtifactManager) DeployCode(
	domain, request core.RecordRef, types []core.RecordRef, codeMap map[core.MachineType][]byte,
) (*core.RecordRef, error) {
	domainRef := record.Core2Reference(domain)
	requestRef := record.Core2Reference(request)

	err := m.checkRequestRecord(m.db, &requestRef)
	if err != nil {
		return nil, err
	}

	typeRefs := make([]record.Reference, 0, len(types))
	for _, tp := range types {
		ref := record.Core2Reference(tp)
		rec, tErr := m.db.GetRecord(&ref)
		if tErr != nil {
			return nil, errors.Wrap(tErr, "failed to retrieve type record")
		}
		if _, ok := rec.(*record.TypeRecord); !ok {
			return nil, errors.Wrap(ErrInvalidRef, "failed to retrieve type record")
		}
		typeRefs = append(typeRefs, ref)
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
		Types:        typeRefs,
		TargetedCode: codeMap,
	}
	codeRef, err := m.db.SetRecord(&rec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store record")
	}
	return codeRef.CoreRef(), nil
}

// ActivateClass creates activate class record in storage. Provided code reference will be used as a class code.
//
// Activation reference will be this class'es identifier and referred as "class head".
func (m *LedgerArtifactManager) ActivateClass(
	domain, request, code core.RecordRef,
) (*core.RecordRef, error) {
	domainRef := record.Core2Reference(domain)
	requestRef := record.Core2Reference(request)
	codeRef := record.Core2Reference(code)

	err := m.checkRequestRecord(m.db, &requestRef)
	if err != nil {
		return nil, err
	}
	_, err = m.getCodeRecord(m.db, codeRef)
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
		CodeRecord: codeRef,
	}

	var classRef *record.Reference
	err = m.db.Update(func(tx *storage.TransactionManager) error {
		classRef, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		err = tx.SetClassIndex(classRef, &index.ClassLifeline{
			LatestStateRef: *classRef,
		})
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return classRef.CoreRef(), nil
}

// DeactivateClass creates deactivate record in storage. Provided reference should be a reference to the head of
// the class. If class is already deactivated, an error should be returned.
//
// Deactivated class cannot be changed or instantiate objects.
func (m *LedgerArtifactManager) DeactivateClass(
	domain, request, class core.RecordRef,
) (*core.RecordRef, error) {
	var err error
	domainRef := record.Core2Reference(domain)
	requestRef := record.Core2Reference(request)
	classRef := record.Core2Reference(class)

	err = m.checkRequestRecord(m.db, &requestRef)
	if err != nil {
		return nil, err
	}

	var deactivationRef *record.Reference
	err = m.db.Update(func(tx *storage.TransactionManager) error {
		_, _, classIndex, err := m.getActiveClass(tx, classRef)
		if err != nil {
			return err
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

		deactivationRef, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		classIndex.LatestStateRef = *deactivationRef
		err = tx.SetClassIndex(&classRef, classIndex)
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return deactivationRef.CoreRef(), nil
}

// UpdateClass creates amend class record in storage. Provided reference should be a reference to the head of
// the class. Migrations are references to code records.
//
// Returned reference will be the latest class state (exact) reference. Migration code will be executed by VM to
// migrate objects memory in the order they appear in provided slice.
func (m *LedgerArtifactManager) UpdateClass(
	domain, request, class, code core.RecordRef, migrations []core.RecordRef,
) (*core.RecordRef, error) {
	var err error
	domainRef := record.Core2Reference(domain)
	requestRef := record.Core2Reference(request)
	classRef := record.Core2Reference(class)
	codeRef := record.Core2Reference(code)
	migrationRefs := make([]record.Reference, 0, len(migrations))
	for _, migration := range migrations {
		migrationRefs = append(migrationRefs, record.Core2Reference(migration))
	}

	err = m.checkRequestRecord(m.db, &requestRef)
	if err != nil {
		return nil, err
	}

	var amendRef *record.Reference
	err = m.db.Update(func(tx *storage.TransactionManager) error {
		_, _, classIndex, err := m.getActiveClass(tx, classRef)
		if err != nil {
			return err
		}

		_, err = m.getCodeRecord(tx, codeRef)
		if err != nil {
			return err
		}
		for _, migrationRef := range migrationRefs {
			_, err = m.getCodeRecord(tx, migrationRef)
			if err != nil {
				return errors.Wrap(err, "invalid migrations")
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

		amendRef, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		classIndex.LatestStateRef = *amendRef
		classIndex.AmendRefs = append(classIndex.AmendRefs, *amendRef)
		err = tx.SetClassIndex(&classRef, classIndex)
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return amendRef.CoreRef(), nil
}

// ActivateObj creates activate object record in storage. Provided class reference will be used as objects class
// memory as memory of crated object. If memory is not provided, the class default memory will be used.
//
// Activation reference will be this object's identifier and referred as "object head".
func (m *LedgerArtifactManager) ActivateObj(
	domain, request, class, parent core.RecordRef, memory []byte,
) (*core.RecordRef, error) {
	var err error
	domainRef := record.Core2Reference(domain)
	requestRef := record.Core2Reference(request)
	classRef := record.Core2Reference(class)
	parentRef := record.Core2Reference(parent)

	err = m.checkRequestRecord(m.db, &requestRef)
	if err != nil {
		return nil, err
	}

	var objRef *record.Reference
	err = m.db.Update(func(tx *storage.TransactionManager) error {
		_, _, _, err = m.getActiveClass(tx, classRef)
		if err != nil {
			return err
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
			Parent:              parentRef,
		}

		// save new record and it's index
		objRef, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		err = tx.SetObjectIndex(objRef, &index.ObjectLifeline{
			ClassRef:       classRef,
			LatestStateRef: *objRef,
			Delegates:      map[core.RecordRef]record.Reference{},
		})
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		// append new record parent's children
		parentIdx, err := tx.GetObjectIndex(&parentRef)
		if err != nil {
			return errors.Wrap(err, "inconsistent index")
		}
		parentIdx.Children = append(parentIdx.Children, *objRef)
		err = tx.SetObjectIndex(&parentRef, parentIdx)
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return objRef.CoreRef(), nil
}

// ActivateObjDelegate is similar to ActivateObj but it created object will be parent's delegate of provided class.
func (m *LedgerArtifactManager) ActivateObjDelegate(
	domain, request, class, parent core.RecordRef, memory []byte,
) (*core.RecordRef, error) {
	var err error
	domainRef := record.Core2Reference(domain)
	requestRef := record.Core2Reference(request)
	classRef := record.Core2Reference(class)
	parentRef := record.Core2Reference(parent)

	err = m.checkRequestRecord(m.db, &requestRef)
	if err != nil {
		return nil, err
	}

	var objRef *record.Reference
	err = m.db.Update(func(tx *storage.TransactionManager) error {
		_, _, _, err = m.getActiveClass(tx, classRef)
		if err != nil {
			return err
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
			Parent:              parentRef,
			Delegate:            true,
		}

		// save new record and it's index
		objRef, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		err = tx.SetObjectIndex(objRef, &index.ObjectLifeline{
			ClassRef:       classRef,
			LatestStateRef: *objRef,
			Delegates:      map[core.RecordRef]record.Reference{},
		})
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		// append new record parent's delegates
		parentIdx, err := tx.GetObjectIndex(&parentRef)
		if err != nil {
			return errors.Wrap(err, "inconsistent index")
		}
		if _, ok := parentIdx.Delegates[class]; ok {
			return errors.New("delegate for this class already exists")
		}
		parentIdx.Delegates[class] = *objRef
		err = tx.SetObjectIndex(&parentRef, parentIdx)
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return objRef.CoreRef(), nil
}

// DeactivateObj creates deactivate object record in storage. Provided reference should be a reference to the head
// of the object. If object is already deactivated, an error should be returned.
//
// Deactivated object cannot be changed.
func (m *LedgerArtifactManager) DeactivateObj(
	domain, request, obj core.RecordRef,
) (*core.RecordRef, error) {
	var err error
	domainRef := record.Core2Reference(domain)
	requestRef := record.Core2Reference(request)
	objRef := record.Core2Reference(obj)

	err = m.checkRequestRecord(m.db, &requestRef)
	if err != nil {
		return nil, err
	}

	var deactivationRef *record.Reference
	err = m.db.Update(func(tx *storage.TransactionManager) error {
		_, _, objIndex, err := m.getActiveObject(tx, objRef)
		if err != nil {
			return err
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
		deactivationRef, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		objIndex.LatestStateRef = *deactivationRef
		err = tx.SetObjectIndex(&objRef, objIndex)
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return deactivationRef.CoreRef(), nil
}

// UpdateObj creates amend object record in storage. Provided reference should be a reference to the head of the
// object. Provided memory well be the new object memory.
//
// Returned reference will be the latest object state (exact) reference.
func (m *LedgerArtifactManager) UpdateObj(
	domain, request, obj core.RecordRef, memory []byte,
) (*core.RecordRef, error) {
	var err error
	domainRef := record.Core2Reference(domain)
	requestRef := record.Core2Reference(request)
	objRef := record.Core2Reference(obj)

	err = m.checkRequestRecord(m.db, &requestRef)
	if err != nil {
		return nil, err
	}

	var amendRef *record.Reference
	err = m.db.Update(func(tx *storage.TransactionManager) error {
		_, _, objIndex, err := m.getActiveObject(tx, objRef)
		if err != nil {
			return err
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

		amendRef, err = tx.SetRecord(&rec)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		objIndex.LatestStateRef = *amendRef
		err = tx.SetObjectIndex(&objRef, objIndex)
		if err != nil {
			return errors.Wrap(err, "failed to store lifeline index")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return amendRef.CoreRef(), nil
}

// GetExactObj returns code and memory of provided object/class state. Deactivation records should be ignored
// (e.g. object considered to be active).
//
// This method is used by validator to fetch the exact state of the object that was used by the executor.
// TODO: not used for now
func (m *LedgerArtifactManager) GetExactObj( // nolint: gocyclo
	classState, objectState core.RecordRef,
) ([]byte, []byte, error) {
	classRef := record.Core2Reference(classState)
	objRef := record.Core2Reference(objectState)

	// Fetching class data
	classRec, err := m.db.GetRecord(&classRef)
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

	code, _, err := m.getCodeRecordCode(m.db, codeRef)
	if err != nil {
		return nil, nil, err
	}

	// Fetching object data
	objectRec, err := m.db.GetRecord(&objRef)
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

	objectIndex, err := m.db.GetObjectIndex(&objectHeadRef)
	if err != nil {
		return nil, nil, errors.Wrap(err, "inconsistent object index")
	}

	// Checking if the object belongs to the class
	if objectIndex.ClassRef.IsNotEqual(classHeadRef) {
		return nil, nil, ErrWrongObject
	}

	return code, memory, nil
}

// GetCode returns code from code record by provided reference.
//
// This method is used by VM to fetch code for execution.
func (m *LedgerArtifactManager) GetCode(code core.RecordRef) (core.CodeDescriptor, error) {
	codeRef := record.Core2Reference(code)

	desc := CodeDescriptor{
		ref:     &codeRef,
		manager: m,
	}

	return &desc, nil
}

// GetLatestClass returns descriptor for latest state of the class known to storage.
// If the class is deactivated, an error should be returned.
//
// Returned descriptor will provide methods for fetching all related data.
func (m *LedgerArtifactManager) GetLatestClass(head core.RecordRef) (core.ClassDescriptor, error) {
	classRef := record.Core2Reference(head)

	classActivateRec, classStateRec, classIndex, err := m.getActiveClass(m.db, classRef)
	if err != nil {
		return nil, err
	}

	class := &ClassDescriptor{
		manager: m,

		headRef:       &classRef,
		stateRef:      &classIndex.LatestStateRef,
		headRecord:    classActivateRec,
		stateRecord:   classStateRec,
		lifelineIndex: classIndex,
	}

	return class, nil
}

// GetLatestObj returns descriptor for latest state of the object known to storage.
// If the object or the class is deactivated, an error should be returned.
//
// Returned descriptor will provide methods for fetching all related data.
func (m *LedgerArtifactManager) GetLatestObj(head core.RecordRef) (core.ObjectDescriptor, error) {
	objRef := record.Core2Reference(head)

	objActivateRec, objStateRec, objIndex, err := m.getActiveObject(m.db, objRef)
	if err != nil {
		return nil, err
	}
	class, err := m.GetLatestClass(*objIndex.ClassRef.CoreRef())
	if err != nil {
		return nil, err
	}

	object := &ObjectDescriptor{
		manager: m,

		headRef:       &objRef,
		stateRef:      &objIndex.LatestStateRef,
		headRecord:    objActivateRec,
		stateRecord:   objStateRec,
		lifelineIndex: objIndex,

		classDescriptor: class.(*ClassDescriptor),
	}

	return object, nil
}

// GetObjChildren returns provided object's children references.
func (m *LedgerArtifactManager) GetObjChildren(head core.RecordRef) ([]core.RecordRef, error) {
	objRef := record.Core2Reference(head)
	_, _, objIndex, err := m.getActiveObject(m.db, objRef)
	if err != nil {
		return nil, err
	}

	childRefs := make([]core.RecordRef, 0, len(objIndex.Children))
	for _, ch := range objIndex.Children {
		childRefs = append(childRefs, *ch.CoreRef())
	}

	return childRefs, nil
}

// GetObjDelegate returns provided object's delegate reference for provided class.
//
// Object delegate should be previously created for this object. If object delegate does not exist, an error will
// be returned.
func (m *LedgerArtifactManager) GetObjDelegate(head, asClass core.RecordRef) (*core.RecordRef, error) {
	objRef := record.Core2Reference(head)

	_, _, objIndex, err := m.getActiveObject(m.db, objRef)
	if err != nil {
		return nil, err
	}

	delegateRef, ok := objIndex.Delegates[asClass]
	if !ok {
		return nil, ErrNotFound
	}

	return delegateRef.CoreRef(), nil
}

// NewArtifactManger creates new manager instance.
func NewArtifactManger(db *storage.DB) (*LedgerArtifactManager, error) {
	return &LedgerArtifactManager{db: db}, nil
}
