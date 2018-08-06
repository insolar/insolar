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
	"errors"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
)

// LedgerArtifactManager provides concrete API to storage for virtual processing module
type LedgerArtifactManager struct {
	storer   storage.LedgerStorer
	archPref []record.ArchType
}

func (m *LedgerArtifactManager) storeRecord(rec record.Record) (record.Reference, error) {
	ref, err := m.storer.AddRecord(rec)
	if err != nil {
		return record.Reference{}, errors.New("record store failed")
	}

	return ref, nil
}

func (m *LedgerArtifactManager) getActiveClassIndex(classRef record.Reference) (*index.Lifeline, error) {
	classRecord, isFound := m.storer.GetRecord(record.ID2Key(classRef.Record))
	if !isFound {
		return nil, errors.New("class record is not found")
	}
	if _, ok := classRecord.(*record.ClassActivateRecord); !ok {
		return nil, errors.New("provided reference is not a class record")
	}

	classIndex, isFound := m.storer.GetIndex(classRef.Record)
	if !isFound {
		return nil, errors.New("class is not activated")
	}
	latestClassRecord, isFound := m.storer.GetRecord(record.ID2Key(classIndex.LatestStateID))
	if !isFound {
		return nil, errors.New("latest class record is not found")
	}
	if _, ok := latestClassRecord.(*record.DeactivationRecord); ok {
		return nil, errors.New("class is deactivated")
	}

	return classIndex, nil
}

func (m *LedgerArtifactManager) SetArchPref(pref []record.ArchType) {
	m.archPref = pref
}

// DeployCode deploys new code to storage (CodeRecord).
func (m *LedgerArtifactManager) DeployCode(requestRef record.Reference) (record.Reference, error) {
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
	codeRecord, isFound := m.storer.GetRecord(record.ID2Key(codeRef.Record))
	if !isFound {
		return record.Reference{}, errors.New("code reference is not found")
	}
	if _, ok := codeRecord.(*record.CodeRecord); !ok {
		return record.Reference{}, errors.New("provided reference is not a code reference")
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
	return m.storeRecord(&rec)
}

// DeactivateClass deactivates class (DeactivationRecord)
func (m *LedgerArtifactManager) DeactivateClass(
	requestRef, classRef record.Reference,
) (record.Reference, error) {
	classIndex, err := m.getActiveClassIndex(classRef)
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
	err = m.storer.SetIndex(classRef.Record, classIndex)
	if err != nil {
		// TODO: add transaction
		return record.Reference{}, errors.New("failed to store lifeline index")
	}

	return deactivationRef, nil
}

// UpdateClass allows to change class code etc. (ClassAmendRecord).
func (m *LedgerArtifactManager) UpdateClass(
	requestRef, classRef record.Reference, migrations []record.MemoryMigrationCode,
) (record.Reference, error) {
	classIndex, err := m.getActiveClassIndex(classRef)
	if err != nil {
		return record.Reference{}, err
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
		// TODO: require and fill code and migration params
	}

	amendRef, err := m.storeRecord(&rec)
	if err != nil {
		return record.Reference{}, errors.New("failed to store amend record")
	}
	classIndex.LatestStateID = amendRef.Record
	err = m.storer.SetIndex(classRef.Record, classIndex)
	if err != nil {
		// TODO: add transaction
		return record.Reference{}, errors.New("failed to store lifeline index")
	}

	return amendRef, nil
}
