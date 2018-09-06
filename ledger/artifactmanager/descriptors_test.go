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

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"

	"github.com/insolar/insolar/ledger/storage/storagetest"
)

type preparedCodeDescriptorTestData struct {
	db      *storage.DB
	manager *LedgerArtifactManager
	rec     *record.CodeRecord
	ref     *record.Reference
}

func prepareCodeDescriptorTestData(t *testing.T) (preparedCodeDescriptorTestData, func()) {
	db, cleaner := storagetest.TmpDB(t, "")

	rec := record.CodeRecord{TargetedCode: map[core.MachineType][]byte{1: {1, 2, 3}}}
	ref, err := db.SetRecord(&rec)
	assert.NoError(t, err)

	return preparedCodeDescriptorTestData{
		db: db,
		manager: &LedgerArtifactManager{
			db:       db,
			archPref: []core.MachineType{1},
		},
		rec: &rec,
		ref: ref,
	}, cleaner
}

func TestCodeDescriptor_MachineType(t *testing.T) {
	td, cleaner := prepareCodeDescriptorTestData(t)
	defer cleaner()

	desc := CodeDescriptor{
		manager: td.manager,
		ref:     td.ref,
	}

	mt, err := desc.MachineType()
	assert.NoError(t, err)
	assert.Equal(t, core.MachineType(1), mt)
}

func TestCodeDescriptor_Code(t *testing.T) {
	td, cleaner := prepareCodeDescriptorTestData(t)
	defer cleaner()

	desc := CodeDescriptor{
		manager: td.manager,
		ref:     td.ref,
	}

	code, err := desc.Code()
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, code)
}

type preparedClassDescriptorTestData struct {
	ledger   *storage.DB
	manager  *LedgerArtifactManager
	classRec *record.ClassActivateRecord
	classRef *record.Reference
}

func prepareClassDescriptorTestData(t *testing.T) (preparedClassDescriptorTestData, func()) {
	ledger, cleaner := storagetest.TmpDB(t, "")

	rec := record.ClassActivateRecord{}
	ref, err := ledger.SetRecord(&rec)
	assert.NoError(t, err)

	return preparedClassDescriptorTestData{
		ledger: ledger,
		manager: &LedgerArtifactManager{
			db:       ledger,
			archPref: []core.MachineType{1},
		},
		classRec: &rec,
		classRef: ref,
	}, cleaner
}

func TestClassDescriptor_GetMigrations(t *testing.T) {
	td, cleaner := prepareClassDescriptorTestData(t)
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
		manager:       td.manager,
		stateRef:      amendRef1,
		headRecord:    td.classRec,
		stateRecord:   &amendRec3,
		lifelineIndex: &idx,
	}

	migrations, err := desc.GetMigrations()
	assert.NoError(t, err)
	assert.Equal(t, [][]byte{{2}, {3}, {4}}, migrations)
}

type preparedObjectDescriptorTestData struct {
	ledger  *storage.DB
	manager *LedgerArtifactManager
	objRec  *record.ObjectActivateRecord
	objRef  *record.Reference
}

func prepareObjectDescriptorTestData(t *testing.T) (preparedObjectDescriptorTestData, func()) {
	ledger, cleaner := storagetest.TmpDB(t, "")

	rec := record.ObjectActivateRecord{Memory: []byte{1}}
	ref, err := ledger.SetRecord(&rec)
	assert.NoError(t, err)

	return preparedObjectDescriptorTestData{
		ledger: ledger,
		manager: &LedgerArtifactManager{
			db:       ledger,
			archPref: []core.MachineType{1},
		},
		objRec: &rec,
		objRef: ref,
	}, cleaner
}

func TestObjectDescriptor_GetMemory(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareObjectDescriptorTestData(t)
	defer cleaner()

	amendRec := record.ObjectAmendRecord{NewMemory: []byte{2}}
	amendRef, _ := td.ledger.SetRecord(&amendRec)
	idx := index.ObjectLifeline{
		LatestStateRef: *amendRef,
	}
	td.ledger.SetObjectIndex(td.objRef, &idx)

	desc := ObjectDescriptor{
		manager:       td.manager,
		headRecord:    td.objRec,
		stateRecord:   nil,
		lifelineIndex: &idx,
	}
	mem, err := desc.Memory()
	assert.NoError(t, err)
	assert.Equal(t, []byte{1}, mem)

	desc = ObjectDescriptor{
		manager:       td.manager,
		headRecord:    td.objRec,
		stateRecord:   &amendRec,
		lifelineIndex: &idx,
	}
	mem, err = desc.Memory()
	assert.NoError(t, err)
	assert.Equal(t, []byte{2}, mem)
}

func TestObjectDescriptor_GetDelegates(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareObjectDescriptorTestData(t)
	defer cleaner()

	appendRec1 := record.ObjectAppendRecord{AppendMemory: []byte{2}}
	appendRec2 := record.ObjectAppendRecord{AppendMemory: []byte{3}}
	appendRef1, _ := td.ledger.SetRecord(&appendRec1)
	appendRef2, _ := td.ledger.SetRecord(&appendRec2)
	idx := index.ObjectLifeline{
		LatestStateRef: *td.objRef,
		AppendRefs:     []record.Reference{*appendRef1, *appendRef2},
	}
	td.ledger.SetObjectIndex(td.objRef, &idx)

	desc := ObjectDescriptor{
		manager:       td.manager,
		headRecord:    td.objRec,
		stateRecord:   nil,
		lifelineIndex: &idx,
	}

	appends, err := desc.GetDelegates()
	assert.NoError(t, err)
	assert.Equal(t, [][]byte{{2}, {3}}, appends)
}
