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

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/leveldb/leveltestutils"
)

type preparedDOTestData struct {
	ledger  storage.Store
	manager *LedgerArtifactManager
	objRec  *record.ObjectActivateRecord
	objRef  *record.Reference
}

func prepareDOTestData(t *testing.T) (preparedDOTestData, func()) {
	ledger, cleaner := leveltestutils.TmpDB(t, "")

	rec := record.ObjectActivateRecord{Memory: record.Memory{1}}
	ref, err := ledger.SetRecord(&rec)
	assert.NoError(t, err)

	return preparedDOTestData{
		ledger: ledger,
		manager: &LedgerArtifactManager{
			storer:   ledger,
			archPref: []record.ArchType{1},
		},
		objRec: &rec,
		objRef: ref,
	}, cleaner
}

func TestObjectDescriptor_GetMemory(t *testing.T) {
	td, cleaner := prepareDOTestData(t)
	defer cleaner()

	amendRec := record.ObjectAmendRecord{NewMemory: record.Memory{2}}
	amendRef, _ := td.ledger.SetRecord(&amendRec)
	idx := index.ObjectLifeline{
		LatestStateRef: *amendRef,
	}
	td.ledger.SetObjectIndex(td.objRef, &idx)

	desc := ObjectDescriptor{
		manager:           td.manager,
		activateRecord:    td.objRec,
		latestAmendRecord: nil,
		lifelineIndex:     &idx,
	}
	mem, err := desc.GetMemory()
	assert.NoError(t, err)
	assert.Equal(t, record.Memory{1}, mem)

	desc = ObjectDescriptor{
		manager:           td.manager,
		activateRecord:    td.objRec,
		latestAmendRecord: &amendRec,
		lifelineIndex:     &idx,
	}
	mem, err = desc.GetMemory()
	assert.NoError(t, err)
	assert.Equal(t, record.Memory{2}, mem)
}

func TestObjectDescriptor_GetDelegates(t *testing.T) {
	td, cleaner := prepareDOTestData(t)
	defer cleaner()

	appendRec1 := record.ObjectAppendRecord{AppendMemory: record.Memory{2}}
	appendRec2 := record.ObjectAppendRecord{AppendMemory: record.Memory{3}}
	appendRef1, _ := td.ledger.SetRecord(&appendRec1)
	appendRef2, _ := td.ledger.SetRecord(&appendRec2)
	idx := index.ObjectLifeline{
		LatestStateRef: *td.objRef,
		AppendRefs:     []record.Reference{*appendRef1, *appendRef2},
	}
	td.ledger.SetObjectIndex(td.objRef, &idx)

	desc := ObjectDescriptor{
		manager:           td.manager,
		activateRecord:    td.objRec,
		latestAmendRecord: nil,
		lifelineIndex:     &idx,
	}

	appends, err := desc.GetDelegates()
	assert.NoError(t, err)
	assert.Equal(t, []record.Memory{{2}, {3}}, appends)
}
