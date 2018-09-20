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
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/pkg/errors"
)

// CodeDescriptor represents meta info required to fetch all code data.
type CodeDescriptor struct {
	cache struct {
		machineType core.MachineType
		code        []byte
	}

	db          *storage.DB
	ref         *record.Reference
	machinePref []core.MachineType
}

// Ref returns reference to represented code record.
func (d *CodeDescriptor) Ref() *core.RecordRef {
	return d.ref.CoreRef()
}

// MachineType fetches code from storage and returns first available machine type according to architecture
// preferences.
//
// Code for returned machine type will be fetched by Code method.
func (d *CodeDescriptor) MachineType() (core.MachineType, error) {
	_, mt, err := d.code()
	if err != nil {
		return core.MachineTypeNotExist, err
	}
	return mt, nil
}

// Code fetches code from storage. Code will be fetched according to architecture preferences
// set via SetArchPref in artifact manager. If preferences are not provided, an error will be returned.
func (d *CodeDescriptor) Code() ([]byte, error) {
	code, _, err := d.code()
	if err != nil {
		return nil, err
	}
	return code, nil
}

// Validate checks code record integrity.
func (d *CodeDescriptor) Validate() error {
	_, _, err := d.code()
	return err
}

func NewCodeDescriptor(db *storage.DB, ref record.Reference, machinePref []core.MachineType) (*CodeDescriptor, error) {
	desc := CodeDescriptor{
		db:          db,
		ref:         &ref,
		machinePref: machinePref,
	}

	return &desc, nil
}

func (d *CodeDescriptor) code() ([]byte, core.MachineType, error) {
	if d.cache.code != nil && d.cache.machineType != core.MachineTypeNotExist {
		return d.cache.code, d.cache.machineType, nil
	}
	// TODO: local / non-local check
	return d.codeLocal()
}

func (d *CodeDescriptor) codeLocal() ([]byte, core.MachineType, error) {
	rec, err := d.db.GetRecord(d.ref)
	if err != nil {
		return nil, core.MachineTypeNotExist, errors.Wrap(err, "failed to retrieve code record")
	}
	codeRec, ok := rec.(*record.CodeRecord)
	if !ok {
		return nil, core.MachineTypeNotExist, errors.Wrap(ErrInvalidRef, "failed to retrieve code record")
	}
	code, mt, err := codeRec.GetCode(d.machinePref)
	if err != nil {
		return nil, mt, errors.Wrap(err, "failed to retrieve code from record")
	}

	return code, mt, nil
}

// ClassDescriptor represents meta info required to fetch all class data.
type ClassDescriptor struct {
	cache struct {
		lifelineIndex  *index.ClassLifeline
		stateRecord    record.ClassState
		codeDescriptor *CodeDescriptor
	}

	db       *storage.DB
	headRef  *record.Reference
	stateRef *record.Reference
}

// HeadRef returns head reference to represented class record.
func (d *ClassDescriptor) HeadRef() *core.RecordRef {
	return d.headRef.CoreRef()
}

// StateRef returns reference to represented class state record.
func (d *ClassDescriptor) StateRef() (*core.RecordRef, error) {
	idx, err := d.index()
	if err != nil {
		return nil, err
	}
	return idx.LatestStateRef.CoreRef(), nil
}

// CodeDescriptor returns descriptor for fetching object's code data.
func (d *ClassDescriptor) CodeDescriptor(machinePref []core.MachineType) (core.CodeDescriptor, error) {
	if d.cache.codeDescriptor != nil {
		return d.cache.codeDescriptor, nil
	}

	state, err := d.state()
	if err != nil {
		return nil, err
	}

	codeRef := state.GetCode()
	if codeRef == nil {
		return nil, errors.New("class has no code")
	}

	desc, err := NewCodeDescriptor(d.db, *codeRef, machinePref)
	if err != nil {
		return nil, err
	}

	return desc, nil
}

// IsActive checks if class is active.
func (d *ClassDescriptor) IsActive() (bool, error) {
	state, err := d.state()
	if err != nil {
		return false, err
	}
	return !state.IsDeactivation(), nil
}

func NewClassDescriptor(db *storage.DB, head record.Reference, state *record.Reference) (*ClassDescriptor, error) {
	desc := ClassDescriptor{
		db:       db,
		headRef:  &head,
		stateRef: state,
	}
	return &desc, nil
}

func (d *ClassDescriptor) index() (*index.ClassLifeline, error) {
	if d.cache.lifelineIndex != nil {
		return d.cache.lifelineIndex, nil
	}
	// TODO: local / non-local check
	return d.indexLocal()
}

func (d *ClassDescriptor) indexLocal() (*index.ClassLifeline, error) {
	idx, err := d.db.GetClassIndex(d.headRef)
	if err != nil {
		return nil, errors.Wrap(err, "inconsistent class index")
	}
	return idx, nil
}

func (d *ClassDescriptor) state() (record.ClassState, error) {
	if d.cache.stateRecord != nil {
		return d.cache.stateRecord, nil
	}
	// TODO: local / non-local check
	return d.stateLocal()
}

func (d *ClassDescriptor) stateLocal() (record.ClassState, error) {
	if d.cache.stateRecord == nil {
		var stateRef *record.Reference
		if d.stateRef == nil {
			idx, err := d.index()
			if err != nil {
				return nil, err
			}
			stateRef = &idx.LatestStateRef
		} else {
			stateRef = d.stateRef
		}
		state, err := d.db.GetRecord(stateRef)
		if err != nil {
			return nil, err
		}
		stateRec, ok := state.(record.ClassState)
		if !ok {
			return nil, errors.New("invalid class record")
		}
		d.cache.stateRecord = stateRec
	}

	return d.cache.stateRecord, nil
}

// ObjectDescriptor represents meta info required to fetch all object data.
type ObjectDescriptor struct {
	cache struct {
		lifelineIndex   *index.ObjectLifeline
		stateRecord     record.ObjectState
		classDescriptor *ClassDescriptor
	}
	db       *storage.DB
	headRef  *record.Reference
	stateRef *record.Reference
}

// HeadRef returns reference to represented object record.
func (d *ObjectDescriptor) HeadRef() *core.RecordRef {
	return d.headRef.CoreRef()
}

// StateRef returns reference to object state record.
func (d *ObjectDescriptor) StateRef() (*core.RecordRef, error) {
	idx, err := d.index()
	if err != nil {
		return nil, err
	}
	return idx.LatestStateRef.CoreRef(), nil
}

// Memory fetches latest memory of the object known to storage.
func (d *ObjectDescriptor) Memory() ([]byte, error) {
	state, err := d.state()
	if err != nil {
		return nil, err
	}
	return state.GetMemory(), nil
}

// IsActive checks if object is active.
func (d *ObjectDescriptor) IsActive() (bool, error) {
	state, err := d.state()
	if err != nil {
		return false, err
	}
	return !state.IsDeactivation(), nil
}

// ClassDescriptor returns descriptor for fetching object's class data.
func (d *ObjectDescriptor) ClassDescriptor(state *core.RecordRef) (core.ClassDescriptor, error) {
	var (
		class *ClassDescriptor
		err   error
	)
	if d.cache.classDescriptor != nil {
		return d.cache.classDescriptor, nil
	}

	idx, err := d.index()
	if err != nil {
		return nil, err
	}
	if state != nil {
		classRef := record.Core2Reference(*state)
		class, err = NewClassDescriptor(d.db, idx.ClassRef, &classRef)
	} else {
		class, err = NewClassDescriptor(d.db, idx.ClassRef, nil)
	}
	if err != nil {
		return nil, err
	}

	return class, nil
}

func NewObjectDescriptor(db *storage.DB, head record.Reference, state *record.Reference) (*ObjectDescriptor, error) {
	desc := ObjectDescriptor{
		db:       db,
		headRef:  &head,
		stateRef: state,
	}
	return &desc, nil
}

func (d *ObjectDescriptor) index() (*index.ObjectLifeline, error) {
	if d.cache.lifelineIndex != nil {
		return d.cache.lifelineIndex, nil
	}
	// TODO: local / non-local check
	return d.indexLocal()
}

func (d *ObjectDescriptor) indexLocal() (*index.ObjectLifeline, error) {
	idx, err := d.db.GetObjectIndex(d.headRef)
	if err != nil {
		return nil, errors.Wrap(err, "inconsistent object index")
	}
	return idx, nil
}

func (d *ObjectDescriptor) state() (record.ObjectState, error) {
	if d.cache.stateRecord != nil {
		return d.cache.stateRecord, nil
	}
	// TODO: local / non-local check
	return d.stateLocal()
}

func (d *ObjectDescriptor) stateLocal() (record.ObjectState, error) {
	if d.cache.stateRecord == nil {
		var stateRef *record.Reference
		if d.stateRef == nil {
			idx, err := d.index()
			if err != nil {
				return nil, err
			}
			stateRef = &idx.LatestStateRef
		} else {
			stateRef = d.stateRef
		}
		state, err := d.db.GetRecord(stateRef)
		if err != nil {
			return nil, err
		}
		stateRec, ok := state.(record.ObjectState)
		if !ok {
			return nil, errors.New("invalid class record")
		}
		d.cache.stateRecord = stateRec
	}

	return d.cache.stateRecord, nil
}

type RefIterator struct {
	elements     []record.Reference
	currentIndex int
}

func (i *RefIterator) HasNext() bool {
	return len(i.elements) > i.currentIndex
}

func (i *RefIterator) Next() (core.RecordRef, error) {
	el := i.elements[i.currentIndex]
	i.currentIndex++
	return *el.CoreRef(), nil
}
