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
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/leveldb/leveltestutils"
)

type preparedDCTestData struct {
	ledger   storage.Store
	manager  *LedgerArtifactManager
	classRec *record.ClassActivateRecord
	classRef *record.Reference
}

func prepareDCTestData(t *testing.T) (preparedDCTestData, func()) {
	ledger, cleaner := leveltestutils.TmpDB(t, "")

	rec := record.ClassActivateRecord{}
	ref, err := ledger.SetRecord(&rec)
	assert.NoError(t, err)

	return preparedDCTestData{
		ledger: ledger,
		manager: &LedgerArtifactManager{
			store:    ledger,
			archPref: []core.MachineType{1},
		},
		classRec: &rec,
		classRef: ref,
	}, cleaner
}

func TestClassDescriptor_GetCode(t *testing.T) {
	td, cleaner := prepareDCTestData(t)
	defer cleaner()

	codeRef, _ := td.ledger.SetRecord(&record.CodeRecord{
		TargetedCode: map[core.MachineType][]byte{
			1: {1, 2, 3},
		},
	})
	amendRec := record.ClassAmendRecord{NewCode: *codeRef}
	amendRef, _ := td.ledger.SetRecord(&amendRec)
	idx := index.ClassLifeline{
		LatestStateRef: *amendRef,
	}
	td.ledger.SetClassIndex(td.classRef, &idx)

	desc := ClassDescriptor{
		manager:           td.manager,
		activateRecord:    td.classRec,
		latestAmendRecord: &amendRec,
		lifelineIndex:     &idx,
	}

	code, err := desc.GetCode()
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, code)
}

func TestClassDescriptor_GetMigrations(t *testing.T) {
	td, cleaner := prepareDCTestData(t)
	defer cleaner()

	codeRef1, _ := td.ledger.SetRecord(&record.CodeRecord{
		TargetedCode: map[core.MachineType][]byte{
			core.MachineType(1): {1},
		},
	})
	codeRef2, _ := td.ledger.SetRecord(&record.CodeRecord{
		TargetedCode: map[core.MachineType][]byte{
			core.MachineType(1): {2},
		},
	})
	codeRef3, _ := td.ledger.SetRecord(&record.CodeRecord{
		TargetedCode: map[core.MachineType][]byte{
			core.MachineType(1): {3},
		},
	})
	codeRef4, _ := td.ledger.SetRecord(&record.CodeRecord{
		TargetedCode: map[core.MachineType][]byte{
			core.MachineType(1): {4},
		},
	})

	amendRec3 := record.ClassAmendRecord{Migrations: []record.Reference{*codeRef4}}
	amendRef1, _ := td.ledger.SetRecord(&record.ClassAmendRecord{
		Migrations: []record.Reference{*codeRef1},
	})
	amendRef2, _ := td.ledger.SetRecord(&record.ClassAmendRecord{
		Migrations: []record.Reference{*codeRef2, *codeRef3},
	})
	amendRef3, _ := td.ledger.SetRecord(&amendRec3)
	idx := index.ClassLifeline{
		LatestStateRef: *amendRef2,
		AmendRefs:      []record.Reference{*amendRef1, *amendRef2, *amendRef3},
	}
	td.ledger.SetClassIndex(td.classRef, &idx)

	desc := ClassDescriptor{
		manager:           td.manager,
		fromState:         *amendRef1,
		activateRecord:    td.classRec,
		latestAmendRecord: &amendRec3,
		lifelineIndex:     &idx,
	}

	migrations, err := desc.GetMigrations()
	assert.NoError(t, err)
	assert.Equal(t, [][]byte{{2}, {3}, {4}}, migrations)
}
