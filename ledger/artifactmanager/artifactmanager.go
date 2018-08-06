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

// LedgerArtifactManager provides concrete API to storage for virtual processing module
type LedgerArtifactManager struct {
	storer   storage.LedgerStorer
	archPref []record.ArchType
}

func (m *LedgerArtifactManager) checkRequestRecord(requestRef record.Reference) error {
	// TODO: implement request check
	return nil
}

func (m *LedgerArtifactManager) checkCodeRecord(codeRef record.Reference) error {
	codeRecord, err := m.storer.GetRecord(codeRef.Record)
	if err != nil {
		return errors.Wrap(err, "code record is not found")
	}
	if _, ok := codeRecord.(*record.CodeRecord); !ok {
		return errors.New("provided reference is not a code reference")
	}
	return nil
}

func (m *LedgerArtifactManager) storeRecord(rec record.Record) (record.Reference, error) {
	id, err := m.storer.SetRecord(rec)
	if err != nil {
		return record.Reference{}, errors.Wrap(err, "record store failed")
	}
	return record.Reference{Domain: rec.Domain(), Record: id}, nil
}

func (m *LedgerArtifactManager) getActiveClassIndex(classID record.ID) (*index.ClassLifeline, error) {
	classRecord, err := m.storer.GetRecord(classID)
	if err != nil {
		return nil, errors.Wrap(err, "class record is not found")
	}
	if _, ok := classRecord.(*record.ClassActivateRecord); !ok {
		return nil, errors.New("provided reference is not a class record")
	}

	classIndex, isFound := m.storer.GetClassIndex(classID)
	if !isFound {
		return nil, errors.New("inconsistent class index")
	}
	latestClassRecord, err := m.storer.GetRecord(classIndex.LatestStateID)
	if err != nil {
		return nil, errors.Wrap(err, "latest class record is not found")
	}
	if _, ok := latestClassRecord.(*record.DeactivationRecord); ok {
		return nil, errors.New("class is deactivated")
	}

	return classIndex, nil
}

func (m *LedgerArtifactManager) getActiveObjectIndex(objID record.ID) (*index.ObjectLifeline, error) {
	objRecord, err := m.storer.GetRecord(objID)
	if err != nil {
		return nil, errors.Wrap(err, "object record is not found")
	}
	if _, ok := objRecord.(*record.ObjectActivateRecord); !ok {
		return nil, errors.New("provided reference is not an object record")
	}

	objIndex, isFound := m.storer.GetObjectIndex(objID)
	if !isFound {
		return nil, errors.New("inconsistent object index")
	}
	latestObjRecord, err := m.storer.GetRecord(objIndex.LatestStateID)
	if err != nil {
		return nil, errors.Wrap(err, "latest object record is not found")
	}
	if _, ok := latestObjRecord.(*record.DeactivationRecord); ok {
		return nil, errors.New("object is deactivated")
	}

	return objIndex, nil
}

func (m *LedgerArtifactManager) SetArchPref(pref []record.ArchType) {
	m.archPref = pref
}

// DeployCode deploys new code to storage (CodeRecord).
func (m *LedgerArtifactManager) DeployCode(requestRef record.Reference) (record.Reference, error) {
	err := m.checkRequestRecord(requestRef)
	if err != nil {
		return record.Reference{}, nil
	}

	rec := record.CodeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
		},
		// TODO: require and fill code params
	}
	return m.storeRecord(&rec)
}

// ActivateClass activates class from given code (ClassActivateRecord).
func (m *LedgerArtifactManager) ActivateClass(
	requestRef, codeRef record.Reference, memory record.Memory,
) (record.Reference, error) {
	err := m.checkRequestRecord(requestRef)
	if err != nil {
		return record.Reference{}, nil
	}
	err = m.checkCodeRecord(codeRef)
	if err != nil {
		return record.Reference{}, err
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
		return record.Reference{}, err
	}
	err = m.storer.SetClassIndex(classRef.Record, &index.ClassLifeline{
		LatestStateID: classRef.Record,
	})
	if err != nil {
		return record.Reference{}, errors.Wrap(err, "failed to store lifeline index")
	}

	return classRef, nil
}

// DeactivateClass deactivates class (DeactivationRecord)
func (m *LedgerArtifactManager) DeactivateClass(
	requestRef, classRef record.Reference,
) (record.Reference, error) {
	err := m.checkRequestRecord(requestRef)
	if err != nil {
		return record.Reference{}, nil
	}

	classIndex, err := m.getActiveClassIndex(classRef.Record)
	if err != nil {
		return record.Reference{}, err
	}

	rec := record.DeactivationRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
			AmendedRecord: record.Reference{
				Domain: classRef.Domain,
				Record: classIndex.LatestStateID,
			},
		},
	}
	deactivationRef, err := m.storeRecord(&rec)
	if err != nil {
		return record.Reference{}, errors.New("failed to store deactivation record")
	}
	classIndex.LatestStateID = deactivationRef.Record
	err = m.storer.SetClassIndex(classRef.Record, classIndex)
	if err != nil {
		// TODO: add transaction
		return record.Reference{}, errors.New("failed to store lifeline index")
	}

	return deactivationRef, nil
}

// UpdateClass allows to change class code etc. (ClassAmendRecord).
func (m *LedgerArtifactManager) UpdateClass(
	requestRef, classRef, codeRef record.Reference, migrationRefs []record.Reference,
) (record.Reference, error) {
	err := m.checkRequestRecord(requestRef)
	if err != nil {
		return record.Reference{}, nil
	}

	classIndex, err := m.getActiveClassIndex(classRef.Record)
	if err != nil {
		return record.Reference{}, err
	}

	err = m.checkCodeRecord(codeRef)
	if err != nil {
		return record.Reference{}, err
	}
	for _, migrationRef := range migrationRefs {
		err = m.checkCodeRecord(migrationRef)
		if err != nil {
			return record.Reference{}, errors.Wrap(err, "invalid migrations")
		}
	}

	rec := record.ClassAmendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
			AmendedRecord: record.Reference{
				Domain: classRef.Domain,
				Record: classIndex.LatestStateID,
			},
		},
		NewCode:    codeRef,
		Migrations: migrationRefs,
	}

	amendRef, err := m.storeRecord(&rec)
	if err != nil {
		return record.Reference{}, errors.New("failed to store amend record")
	}
	classIndex.LatestStateID = amendRef.Record
	classIndex.MigrationIDs = append(classIndex.MigrationIDs, amendRef.Record)
	err = m.storer.SetClassIndex(classRef.Record, classIndex)
	if err != nil {
		// TODO: add transaction
		return record.Reference{}, errors.New("failed to store lifeline index")
	}

	return amendRef, nil
}

// ActivateObj creates and activates new object from given class (ObjectActivateRecord).
func (m *LedgerArtifactManager) ActivateObj(
	requestRef, classRef record.Reference, memory record.Memory,
) (record.Reference, error) {
	err := m.checkRequestRecord(requestRef)
	if err != nil {
		return record.Reference{}, nil
	}

	_, err = m.getActiveClassIndex(classRef.Record)
	if err != nil {
		return record.Reference{}, err
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
		return record.Reference{}, err
	}
	err = m.storer.SetObjectIndex(objRef.Record, &index.ObjectLifeline{
		ClassID:       classRef.Record,
		LatestStateID: objRef.Record,
	})
	if err != nil {
		return record.Reference{}, errors.Wrap(err, "failed to store lifeline index")
	}

	return objRef, nil
}

// DeactivateObj deactivates object (DeactivationRecord).
func (m *LedgerArtifactManager) DeactivateObj(requestRef, objRef record.Reference) (record.Reference, error) {
	err := m.checkRequestRecord(requestRef)
	if err != nil {
		return record.Reference{}, nil
	}

	objIndex, err := m.getActiveObjectIndex(objRef.Record)
	if err != nil {
		return record.Reference{}, err
	}

	rec := record.DeactivationRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
			AmendedRecord: record.Reference{
				Domain: objRef.Domain,
				Record: objIndex.LatestStateID,
			},
		},
	}
	deactivationRef, err := m.storeRecord(&rec)
	if err != nil {
		return record.Reference{}, errors.New("failed to store deactivation record")
	}
	objIndex.LatestStateID = deactivationRef.Record
	err = m.storer.SetObjectIndex(objRef.Record, objIndex)
	if err != nil {
		// TODO: add transaction
		return record.Reference{}, errors.New("failed to store lifeline index")
	}
	return deactivationRef, nil
}

// UpdateObj allows to change object state (ObjectAmendRecord).
func (m *LedgerArtifactManager) UpdateObj(
	requestRef, objRef record.Reference, memory record.Memory,
) (record.Reference, error) {
	err := m.checkRequestRecord(requestRef)
	if err != nil {
		return record.Reference{}, nil
	}

	objIndex, err := m.getActiveObjectIndex(objRef.Record)
	if err != nil {
		return record.Reference{}, err
	}

	rec := record.ObjectAmendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
			AmendedRecord: record.Reference{
				Domain: objRef.Domain,
				Record: objIndex.LatestStateID,
			},
		},
		NewMemory: memory,
	}

	amendRef, err := m.storeRecord(&rec)
	if err != nil {
		return record.Reference{}, errors.New("failed to store amend record")
	}
	objIndex.LatestStateID = amendRef.Record
	objIndex.AppendIDs = []record.ID{}
	err = m.storer.SetObjectIndex(objRef.Record, objIndex)
	if err != nil {
		// TODO: add transaction
		return record.Reference{}, errors.New("failed to store lifeline index")
	}
	return amendRef, nil
}

func (m *LedgerArtifactManager) AppendObjDelegate(
	requestRef, objRef record.Reference, memory record.Memory,
) (record.Reference, error) {
	err := m.checkRequestRecord(requestRef)
	if err != nil {
		return record.Reference{}, nil
	}

	objIndex, err := m.getActiveObjectIndex(objRef.Record)
	if err != nil {
		return record.Reference{}, err
	}

	rec := record.ObjectAppendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
			AmendedRecord: record.Reference{
				Domain: objRef.Domain,
				Record: objIndex.LatestStateID,
			},
		},
		AppendMemory: memory,
	}

	appendRef, err := m.storeRecord(&rec)
	if err != nil {
		return record.Reference{}, errors.New("failed to store append record")
	}
	objIndex.LatestStateID = appendRef.Record
	objIndex.AppendIDs = append(objIndex.AppendIDs, appendRef.Record)
	err = m.storer.SetObjectIndex(objRef.Record, objIndex)
	if err != nil {
		// TODO: add transaction
		return record.Reference{}, errors.New("failed to store lifeline index")
	}
	return appendRef, nil
}

func (m *LedgerArtifactManager) GetObj(
	ref, storedClassRef, storedObjRef record.Reference,
) (*ClassDescr, *ObjDescr, error) {
	getMigrations := func(classIndex *index.ClassLifeline) ([][]byte, error) {
		result := make([][]byte, len(classIndex.MigrationIDs))
		for i := len(classIndex.MigrationIDs); i >= 0; i-- {
			if storedClassRef.Record == classIndex.MigrationIDs[i] {
				break
			}
			rec, err := m.storer.GetRecord(classIndex.MigrationIDs[i])
			if err != nil {
				return nil, err
			}
			codeRec, ok := rec.(*record.CodeRecord)
			if !ok {
				return nil, errors.New("migration reference is not a code record")
			}
			code, err := codeRec.GetCode(m.archPref)
			if err != nil {
				return nil, err
			}
			result = append(result, code)
		}
		sortedResult := make([][]byte, len(result))
		for i := len(result); i >= 0; i-- {
			sortedResult = append(sortedResult, result[i])
		}
		return sortedResult, nil
	}

	objIndex, err := m.getActiveObjectIndex(ref.Record)
	if err != nil {
		return nil, nil, err
	}
	classIndex, err := m.getActiveClassIndex(objIndex.ClassID)
	if err != nil {
		return nil, nil, err
	}

	var (
		class  *ClassDescr = nil
		object *ObjDescr   = nil
	)

	if storedClassRef.Record != classIndex.LatestStateID {
		migrations, err := getMigrations(classIndex)
		if err != nil {
			return nil, nil, errors.Wrap(err, "invalid migrations")
		}
		class = &ClassDescr{
			Migrations: migrations,
		}
	}

	// TODO: add object processing
	return class, object, nil
}
