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
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/pkg/errors"
)

// CodeDescriptor represents meta info required to fetch all code data.
type CodeDescriptor struct {
	ref     *record.Reference
	manager *LedgerArtifactManager
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
	_, mt, err := d.manager.getCodeRecordCode(*d.ref)
	return mt, err
}

// Code fetches code from storage. Code will be fetched according to architecture preferences
// set via SetArchPref in artifact manager. If preferences are not provided, an error will be returned.
func (d *CodeDescriptor) Code() ([]byte, error) {
	code, _, err := d.manager.getCodeRecordCode(*d.ref)
	if err != nil {
		return nil, err
	}

	return code, nil
}

// ClassDescriptor represents meta info required to fetch all class data.
type ClassDescriptor struct {
	manager *LedgerArtifactManager

	headRef       *record.Reference
	stateRef      *record.Reference
	headRecord    *record.ClassActivateRecord
	stateRecord   *record.ClassAmendRecord
	lifelineIndex *index.ClassLifeline

	codeDescriptor *CodeDescriptor
}

// HeadRef returns head reference to represented class record.
func (d *ClassDescriptor) HeadRef() *core.RecordRef {
	return d.headRef.CoreRef()
}

// StateRef returns reference to represented class state record.
func (d *ClassDescriptor) StateRef() *core.RecordRef {
	return d.stateRef.CoreRef()
}

// CodeDescriptor returns descriptor for fetching object's code data.
func (d *ClassDescriptor) CodeDescriptor() (core.CodeDescriptor, error) {
	if d.codeDescriptor == nil {
		codeRef := d.headRecord.CodeRecord
		if d.stateRecord != nil {
			codeRef = d.stateRecord.NewCode
		}
		d.codeDescriptor = &CodeDescriptor{
			ref:     &codeRef,
			manager: d.manager,
		}
	}

	return d.codeDescriptor, nil
}

// GetMigrations fetches all migrations from provided to artifact manager state to the last state known to storage. VM
// is responsible for applying these migrations and updating objects.
// TODO: not used for now
func (d *ClassDescriptor) GetMigrations() ([][]byte, error) {
	var amends []*record.ClassAmendRecord
	// Search for provided state in class amends from the end of the list.
	// Record keys are hashes and are not incremental, so we can't say if provided state should be before or after.
	for i := len(d.lifelineIndex.AmendRefs) - 1; i >= 0; i-- {
		amendRef := d.lifelineIndex.AmendRefs[i]
		if d.stateRef.IsEqual(amendRef) {
			break // Provided state is found. It means we now have all the amends we need.
		}
		rec, err := d.manager.store.GetRecord(&amendRef)
		if err != nil {
			return nil, errors.Wrap(err, "inconsistent class index")
		}
		amendRec, ok := rec.(*record.ClassAmendRecord)
		if !ok {
			return nil, errors.Wrap(ErrInvalidRef, "inconsistent class index")
		}
		amends = append(amends, amendRec)
	}
	// Reverse found amends again (we appended them from the end) so they'll have the original order.
	sortedAmends := make([]*record.ClassAmendRecord, len(amends))
	for i, amend := range amends {
		sortedAmends[len(amends)-i-1] = amend
	}

	// Flatten the migrations list from amends.
	var migrations [][]byte
	for _, amendRec := range sortedAmends {
		for _, codeRef := range amendRec.Migrations {
			code, _, err := d.manager.getCodeRecordCode(codeRef)
			if err != nil {
				return nil, errors.Wrap(err, "invalid migration reference in amend record")
			}
			migrations = append(migrations, code)
		}
	}

	return migrations, nil
}

// ObjectDescriptor represents meta info required to fetch all object data.
type ObjectDescriptor struct {
	manager *LedgerArtifactManager

	headRef       *record.Reference
	stateRef      *record.Reference
	headRecord    *record.ObjectActivateRecord
	stateRecord   *record.ObjectAmendRecord
	lifelineIndex *index.ObjectLifeline

	classDescriptor *ClassDescriptor
}

// HeadRef returns reference to represented object record.
func (d *ObjectDescriptor) HeadRef() *core.RecordRef {
	return d.headRef.CoreRef()
}

// StateRef returns reference to object state record.
func (d *ObjectDescriptor) StateRef() *core.RecordRef {
	return d.stateRef.CoreRef()
}

// Memory fetches latest memory of the object known to storage.
func (d *ObjectDescriptor) Memory() ([]byte, error) {
	if d.stateRecord != nil {
		return d.stateRecord.NewMemory, nil
	}

	return d.headRecord.Memory, nil
}

// CodeDescriptor returns descriptor for fetching object's code data.
func (d *ObjectDescriptor) CodeDescriptor() (core.CodeDescriptor, error) {
	return d.classDescriptor.codeDescriptor, nil
}

// ClassDescriptor returns descriptor for fetching object's class data.
func (d *ObjectDescriptor) ClassDescriptor() (core.ClassDescriptor, error) {
	return d.classDescriptor, nil
}

// GetDelegates fetches unamended delegates from storage.
//
// VM is responsible for collecting all delegates and adding them to the object memory manually if its required.
// TODO: not used for now
func (d *ObjectDescriptor) GetDelegates() ([][]byte, error) {
	var delegates [][]byte
	for _, appendRef := range d.lifelineIndex.AppendRefs {
		rec, err := d.manager.store.GetRecord(&appendRef)
		if err != nil {
			return nil, errors.Wrap(err, "inconsistent object index")
		}
		appendRec, ok := rec.(*record.ObjectAppendRecord)
		if !ok {
			return nil, errors.Wrap(ErrInvalidRef, "inconsistent object index")
		}
		delegates = append(delegates, appendRec.AppendMemory)
	}

	return delegates, nil
}
