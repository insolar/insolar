//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package object

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
)

// GenesisRecord is the first record created on storage. It's used to link root objects to it.
// this type is a wrapper for insolar.GenesisRecord, which makes records storage work for it.
type GenesisRecord struct {
	record.VirtualRecord
}

// methods below implements State interface (required in some places of logic runner code)

var _ State = &GenesisRecord{}

// PrevStateID returns previous state id.
func (r *GenesisRecord) PrevStateID() *insolar.ID {
	return nil
}

// StateID returns state id.
func (r *GenesisRecord) ID() StateID {
	return StateActivation
}

// GetMemory returns state memory.
func (*GenesisRecord) GetMemory() *insolar.ID {
	return nil
}

// GetImage returns state code.
func (*GenesisRecord) GetImage() *insolar.Reference {
	return nil
}

// GetIsPrototype returns state code.
func (*GenesisRecord) GetIsPrototype() bool {
	return false
}
