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

package store

import (
	"sync"

	"github.com/insolar/insolar/insolar"
)

//go:generate minimock -i github.com/insolar/insolar/internal/ledger/store.JetIndexModifier -o ./ -s _mock.go

// JetIndexModifier is an interface for modifying index records.
type JetIndexModifier interface {
	Add(id insolar.ID, jetID insolar.JetID)
	Delete(id insolar.ID, jetID insolar.JetID)
}

//go:generate minimock -i github.com/insolar/insolar/internal/ledger/store.JetIndexAccesssor -o ./ -s _mock.go

// JetIndexAccessor is an interface for modifying index records.
type JetIndexAccessor interface {
	For(jetID insolar.JetID) map[insolar.ID]struct{}
}

// JetIndex contains methods to implement quick access to data by jet. Indexes are stored in memory. Consider disk
// implementation for large collections.
type JetIndex struct {
	lock    sync.Mutex
	storage map[insolar.JetID]recordSet
}

type recordSet map[insolar.ID]struct{}

// NewJetIndex creates new index instance.
func NewJetIndex() *JetIndex {
	return &JetIndex{storage: map[insolar.JetID]recordSet{}}
}

// Add creates index record for specified id and jet. To remove clean up index, use "Delete" method.
func (i *JetIndex) Add(id insolar.ID, jetID insolar.JetID) {
	i.lock.Lock()
	defer i.lock.Unlock()

	jet, ok := i.storage[jetID]
	if !ok {
		jet = recordSet{}
		i.storage[jetID] = jet
	}
	jet[id] = struct{}{}
}

// Delete removes specified id - jet record from index.
func (i *JetIndex) Delete(id insolar.ID, jetID insolar.JetID) {
	i.lock.Lock()
	defer i.lock.Unlock()

	jet, ok := i.storage[jetID]
	if !ok {
		return
	}

	delete(jet, id)
	if len(jet) == 0 {
		delete(i.storage, jetID)
	}
}

// For returns a collection of ids, that are stored for a specific jetID and a pulse number
func (i *JetIndex) For(jetID insolar.JetID) map[insolar.ID]struct{} {
	i.lock.Lock()
	defer i.lock.Unlock()

	ids, ok := i.storage[jetID]
	if !ok {
		return nil
	}

	res := map[insolar.ID]struct{}{}
	for id := range ids {
		res[id] = struct{}{}
	}

	return res
}
