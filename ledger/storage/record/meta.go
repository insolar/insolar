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

package record

import (
	"io"

	"github.com/insolar/insolar/core"
)

// GenesisRecord is the first record created on storage. It's used to link root objects to it.
type GenesisRecord struct {
}

// PrevStateID returns previous state id.
func (r *GenesisRecord) PrevStateID() *core.RecordID {
	return nil
}

// State returns state id.
func (r *GenesisRecord) State() State {
	return StateActivation
}

// Type implementation of Record interface.
func (r *GenesisRecord) Type() TypeID { return typeGenesis }

// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
func (r *GenesisRecord) WriteHashData(w io.Writer) (int, error) {
	return w.Write(SerializeRecord(r))
}

// GetMemory returns state memory.
func (*GenesisRecord) GetMemory() *core.RecordID {
	return nil
}

// GetImage returns state code.
func (*GenesisRecord) GetImage() *core.RecordRef {
	return nil
}

// GetIsPrototype returns state code.
func (*GenesisRecord) GetIsPrototype() bool {
	return false
}

// ChildRecord is a child activation record. Its used for children iterating.
type ChildRecord struct {
	PrevChild *core.RecordID

	Ref core.RecordRef // Reference to the child's head.
}

// Type implementation of Record interface.
func (r *ChildRecord) Type() TypeID { return typeChild }

// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
func (r *ChildRecord) WriteHashData(w io.Writer) (int, error) {
	return w.Write(SerializeRecord(r))
}

// JetRecord represents Jet.
type JetRecord struct {
	// TODO: should contain prefix.
}

// Type implementation of Record interface.
func (r *JetRecord) Type() TypeID { return typeJet }

// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
func (r *JetRecord) WriteHashData(w io.Writer) (int, error) {
	return w.Write(SerializeRecord(r))
}
