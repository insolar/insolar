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
	"os"
	"testing"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/leveldb"
	"github.com/stretchr/testify/assert"
)

func prepareClassDescriptorTest() (
	storage.LedgerStorer, *LedgerArtifactManager, *record.ClassActivateRecord, *record.Reference,
) {
	if err := leveldb.DropDB(); err != nil {
		os.Exit(1)
	}
	ledger, _ := leveldb.InitDB()
	manager := LedgerArtifactManager{
		storer:   ledger,
		archPref: []record.ArchType{1},
	}
	rec := record.ClassActivateRecord{}
	ref, _ := ledger.SetRecord(&rec)

	return ledger, &manager, &rec, ref
}

func TestClassDescriptor_GetCode(t *testing.T) {
	ledger, manager, classRec, classRef := prepareClassDescriptorTest()
	codeRef, _ := ledger.SetRecord(&record.CodeRecord{TargetedCode: map[record.ArchType][]byte{
		1: {1, 2, 3},
	}})
	amendRec := record.ClassAmendRecord{NewCode: *codeRef}
	amendRef, _ := ledger.SetRecord(&amendRec)
	idx := index.ClassLifeline{
		LatestStateRef: *amendRef,
	}
	ledger.SetClassIndex(classRef, &idx)

	desc := ClassDescriptor{
		manager:           manager,
		activateRecord:    classRec,
		latestAmendRecord: &amendRec,
		lifelineIndex:     &idx,
	}

	code, err := desc.GetCode()
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, code)
}

func TestClassDescriptor_GetMigrations(t *testing.T) {
	ledger, manager, classRec, classRef := prepareClassDescriptorTest()
	codeRef1, _ := ledger.SetRecord(&record.CodeRecord{TargetedCode: map[record.ArchType][]byte{
		record.ArchType(1): {1},
	}})
	codeRef2, _ := ledger.SetRecord(&record.CodeRecord{TargetedCode: map[record.ArchType][]byte{
		record.ArchType(1): {2},
	}})
	codeRef3, _ := ledger.SetRecord(&record.CodeRecord{TargetedCode: map[record.ArchType][]byte{
		record.ArchType(1): {3},
	}})
	codeRef4, _ := ledger.SetRecord(&record.CodeRecord{TargetedCode: map[record.ArchType][]byte{
		record.ArchType(1): {4},
	}})

	amendRec3 := record.ClassAmendRecord{Migrations: []record.Reference{*codeRef4}}
	amendRef1, _ := ledger.SetRecord(&record.ClassAmendRecord{Migrations: []record.Reference{*codeRef1}})
	amendRef2, _ := ledger.SetRecord(&record.ClassAmendRecord{Migrations: []record.Reference{*codeRef2, *codeRef3}})
	amendRef3, _ := ledger.SetRecord(&amendRec3)
	idx := index.ClassLifeline{
		LatestStateRef: *amendRef2,
		AmendRefs:      []record.Reference{*amendRef1, *amendRef2, *amendRef3},
	}
	ledger.SetClassIndex(classRef, &idx)

	desc := ClassDescriptor{
		manager:           manager,
		fromState:         *amendRef1,
		activateRecord:    classRec,
		latestAmendRecord: &amendRec3,
		lifelineIndex:     &idx,
	}

	migrations, err := desc.GetMigrations()
	assert.NoError(t, err)
	assert.Equal(t, [][]byte{{2}, {3}, {4}}, migrations)
}
