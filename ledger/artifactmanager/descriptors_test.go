/*
 *    Copyright 2018 Insolar
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

func TestCodeDescriptor_Ref(t *testing.T) {
	td, cleaner := prepareCodeDescriptorTestData(t)
	defer cleaner()

	desc, err := NewCodeDescriptor(td.db, *td.ref, td.manager.archPref)
	assert.NoError(t, err)
	ref := desc.Ref()
	assert.Equal(t, *td.ref.CoreRef(), *ref)
}

func TestCodeDescriptor_MachineType(t *testing.T) {
	td, cleaner := prepareCodeDescriptorTestData(t)
	defer cleaner()

	desc, err := NewCodeDescriptor(td.db, *td.ref, td.manager.archPref)
	assert.NoError(t, err)
	mt, err := desc.MachineType()
	assert.NoError(t, err)
	assert.Equal(t, core.MachineType(1), mt)
}

func TestCodeDescriptor_Code(t *testing.T) {
	td, cleaner := prepareCodeDescriptorTestData(t)
	defer cleaner()

	desc, err := NewCodeDescriptor(td.db, *td.ref, td.manager.archPref)
	assert.NoError(t, err)
	code, err := desc.Code()
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, code)
}

func TestCodeDescriptor_Validate(t *testing.T) {
	td, cleaner := prepareCodeDescriptorTestData(t)
	defer cleaner()

	desc, err := NewCodeDescriptor(td.db, *genRandomRef(), td.manager.archPref)
	assert.NoError(t, err)
	err = desc.Validate()
	assert.Error(t, err)

	desc, err = NewCodeDescriptor(td.db, *td.ref, td.manager.archPref)
	assert.NoError(t, err)
	err = desc.Validate()
	assert.NoError(t, err)
}

type preparedClassDescriptorTestData struct {
	db       *storage.DB
	manager  *LedgerArtifactManager
	classRec *record.ClassActivateRecord
	classRef *record.Reference
}

func prepareClassDescriptorTestData(t *testing.T) (preparedClassDescriptorTestData, func()) {
	db, cleaner := storagetest.TmpDB(t, "")

	rec := record.ClassActivateRecord{}
	ref, err := db.SetRecord(&rec)
	assert.NoError(t, err)

	return preparedClassDescriptorTestData{
		db: db,
		manager: &LedgerArtifactManager{
			db:       db,
			archPref: []core.MachineType{1},
		},
		classRec: &rec,
		classRef: ref,
	}, cleaner
}

func TestClassDescriptor_HeadRef(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareClassDescriptorTestData(t)
	defer cleaner()

	desc, err := NewClassDescriptor(td.db, *td.classRef, nil)
	assert.NoError(t, err)
	assert.Equal(t, *td.classRef.CoreRef(), *desc.HeadRef())
}

func TestClassDescriptor_StateRef(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareClassDescriptorTestData(t)
	defer cleaner()

	stateRef := genRandomRef()
	err := td.db.SetClassIndex(td.classRef, &index.ClassLifeline{LatestStateRef: *stateRef})

	desc, err := NewClassDescriptor(td.db, *td.classRef, nil)
	assert.NoError(t, err)
	descStateRef, err := desc.StateRef()
	assert.NoError(t, err)
	assert.Equal(t, *stateRef.CoreRef(), *descStateRef)

	desc, err = NewClassDescriptor(td.db, *td.classRef, stateRef)
	descStateRef, err = desc.StateRef()
	assert.NoError(t, err)
	assert.Equal(t, *stateRef.CoreRef(), *descStateRef)
}

func TestClassDescriptor_IsActive(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareClassDescriptorTestData(t)
	defer cleaner()

	err := td.db.SetClassIndex(td.classRef, &index.ClassLifeline{LatestStateRef: *td.classRef})
	assert.NoError(t, err)

	desc, err := NewClassDescriptor(td.db, *td.classRef, nil)
	assert.NoError(t, err)
	active, err := desc.IsActive()
	assert.NoError(t, err)
	assert.Equal(t, true, active)

	deactivateRef, err := td.db.SetRecord(&record.DeactivationRecord{})
	err = td.db.SetClassIndex(td.classRef, &index.ClassLifeline{LatestStateRef: *deactivateRef})
	desc, err = NewClassDescriptor(td.db, *td.classRef, nil)
	assert.NoError(t, err)
	active, err = desc.IsActive()
	assert.NoError(t, err)
	assert.Equal(t, false, active)
}

func TestClassDescriptor_CodeDescriptor(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareClassDescriptorTestData(t)
	defer cleaner()

	codeRef := genRandomRef()
	amendRef, err := td.db.SetRecord(&record.ClassAmendRecord{NewCode: *codeRef})
	err = td.db.SetClassIndex(td.classRef, &index.ClassLifeline{LatestStateRef: *amendRef})
	assert.NoError(t, err)

	desc, err := NewClassDescriptor(td.db, *td.classRef, nil)
	assert.NoError(t, err)
	_, err = desc.CodeDescriptor(nil)
	assert.NoError(t, err)
}

type preparedObjectDescriptorTestData struct {
	db      *storage.DB
	manager *LedgerArtifactManager
	objRec  *record.ObjectActivateRecord
	objRef  *record.Reference
}

func prepareObjectDescriptorTestData(t *testing.T) (preparedObjectDescriptorTestData, func()) {
	db, cleaner := storagetest.TmpDB(t, "")

	rec := record.ObjectActivateRecord{Memory: []byte{1}}
	ref, err := db.SetRecord(&rec)
	assert.NoError(t, err)

	return preparedObjectDescriptorTestData{
		db: db,
		manager: &LedgerArtifactManager{
			db:       db,
			archPref: []core.MachineType{1},
		},
		objRec: &rec,
		objRef: ref,
	}, cleaner
}

func TestObjectDescriptor_HeadRef(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareObjectDescriptorTestData(t)
	defer cleaner()

	desc, err := NewObjectDescriptor(td.db, *td.objRef, nil)
	assert.NoError(t, err)
	assert.Equal(t, *td.objRef.CoreRef(), *desc.HeadRef())
}

func TestObjectDescriptor_StateRef(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareObjectDescriptorTestData(t)
	defer cleaner()

	stateRef := genRandomRef()
	err := td.db.SetObjectIndex(td.objRef, &index.ObjectLifeline{LatestStateRef: *stateRef})

	desc, err := NewObjectDescriptor(td.db, *td.objRef, nil)
	assert.NoError(t, err)
	descStateRef, err := desc.StateRef()
	assert.NoError(t, err)
	assert.Equal(t, *stateRef.CoreRef(), *descStateRef)

	desc, err = NewObjectDescriptor(td.db, *td.objRef, stateRef)
	descStateRef, err = desc.StateRef()
	assert.NoError(t, err)
	assert.Equal(t, *stateRef.CoreRef(), *descStateRef)
}

func TestObjectDescriptor_IsActive(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareObjectDescriptorTestData(t)
	defer cleaner()

	err := td.db.SetObjectIndex(td.objRef, &index.ObjectLifeline{LatestStateRef: *td.objRef})
	assert.NoError(t, err)

	desc, err := NewObjectDescriptor(td.db, *td.objRef, nil)
	assert.NoError(t, err)
	active, err := desc.IsActive()
	assert.NoError(t, err)
	assert.Equal(t, true, active)

	deactivateRef, err := td.db.SetRecord(&record.DeactivationRecord{})
	err = td.db.SetObjectIndex(td.objRef, &index.ObjectLifeline{LatestStateRef: *deactivateRef})
	desc, err = NewObjectDescriptor(td.db, *td.objRef, nil)
	assert.NoError(t, err)
	active, err = desc.IsActive()
	assert.NoError(t, err)
	assert.Equal(t, false, active)
}

func TestObjectDescriptor_ClassDescriptor(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareObjectDescriptorTestData(t)
	defer cleaner()

	classRef := genRandomRef()
	err := td.db.SetObjectIndex(td.objRef, &index.ObjectLifeline{LatestStateRef: *td.objRef, ClassRef: *classRef})
	assert.NoError(t, err)

	desc, err := NewObjectDescriptor(td.db, *td.objRef, nil)
	assert.NoError(t, err)
	_, err = desc.ClassDescriptor(nil)
	assert.NoError(t, err)
}

func TestObjectDescriptor_GetMemory(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareObjectDescriptorTestData(t)
	defer cleaner()

	amendRec := record.ObjectAmendRecord{NewMemory: []byte{2}}
	amendRef, _ := td.db.SetRecord(&amendRec)
	idx := index.ObjectLifeline{
		LatestStateRef: *amendRef,
	}
	td.db.SetObjectIndex(td.objRef, &idx)

	desc, err := NewObjectDescriptor(td.db, *td.objRef, nil)
	assert.NoError(t, err)
	mem, err := desc.Memory()
	assert.NoError(t, err)
	assert.Equal(t, []byte{2}, mem)

	desc, err = NewObjectDescriptor(td.db, *td.objRef, td.objRef)
	assert.NoError(t, err)
	mem, err = desc.Memory()
	assert.NoError(t, err)
	assert.Equal(t, []byte{1}, mem)

	desc, err = NewObjectDescriptor(td.db, *td.objRef, amendRef)
	assert.NoError(t, err)
	mem, err = desc.Memory()
	assert.NoError(t, err)
	assert.Equal(t, []byte{2}, mem)
}
