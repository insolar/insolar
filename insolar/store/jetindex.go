// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package store

import (
	"sync"

	"github.com/insolar/insolar/insolar"
)

//go:generate minimock -i github.com/insolar/insolar/insolar/store.JetIndexModifier -o ./ -s _mock.go -g -g

// JetIndexModifier is an interface for modifying index records.
type JetIndexModifier interface {
	Add(id insolar.ID, jetID insolar.JetID)
	Delete(id insolar.ID, jetID insolar.JetID)
}

//go:generate minimock -i github.com/insolar/insolar/insolar/store.JetIndexAccessor -o ./ -s _mock.go -g

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

// For returns a collection of ids, that are stored for a specific jetID
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
