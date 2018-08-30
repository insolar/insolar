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
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/pkg/errors"
)

// ObjectDescriptor represents meta info required to fetch all object data.
type ObjectDescriptor struct {
	StateRef record.Reference

	manager           *LedgerArtifactManager
	activateRecord    *record.ObjectActivateRecord
	latestAmendRecord *record.ObjectAmendRecord
	lifelineIndex     *index.ObjectLifeline
}

// GetMemory fetches latest memory of the object known to storage.
func (d *ObjectDescriptor) GetMemory() ([]byte, error) {
	if d.latestAmendRecord != nil {
		return d.latestAmendRecord.NewMemory, nil
	}

	return d.activateRecord.Memory, nil
}

// GetDelegates fetches unamended delegates from storage.
//
// VM is responsible for collecting all delegates and adding them to the object memory manually if its required.
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
