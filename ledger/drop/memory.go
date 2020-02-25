// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package drop

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"go.opencensus.io/stats"
)

type dropKey struct {
	pulse insolar.PulseNumber
	jetID insolar.JetID
}

type dropStorageMemory struct {
	lock  sync.RWMutex
	drops map[dropKey]Drop
}

// NewStorageMemory creates a new storage, that holds data in a memory.
func NewStorageMemory() *dropStorageMemory { // nolint: golint
	return &dropStorageMemory{
		drops: map[dropKey]Drop{},
	}
}

// ForPulse returns a Drop for a provided pulse, that is stored in a memory
func (m *dropStorageMemory) ForPulse(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber) (Drop, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	key := dropKey{jetID: jetID, pulse: pulse}
	d, ok := m.drops[key]
	if !ok {
		return Drop{}, ErrNotFound
	}

	return d, nil
}

// Set saves a provided Drop to a memory
func (m *dropStorageMemory) Set(ctx context.Context, drop Drop) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	key := dropKey{jetID: drop.JetID, pulse: drop.Pulse}
	_, ok := m.drops[key]
	if ok {
		return ErrOverride
	}
	m.drops[key] = drop

	stats.Record(ctx,
		statDropInMemoryAddedCount.M(1),
	)

	return nil
}

// DeleteForPN methods removes a drop from a memory storage.
func (m *dropStorageMemory) DeleteForPN(ctx context.Context, pulse insolar.PulseNumber) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for key := range m.drops {
		if key.pulse == pulse {
			delete(m.drops, key)
			stats.Record(ctx,
				statDropInMemoryRemovedCount.M(1),
			)
		}
	}
}
