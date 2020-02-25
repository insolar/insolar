// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package storage

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/node"
)

const entriesCount = 10

// NewMemoryStorage constructor creates MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		entries:         make([]insolar.Pulse, 0),
		snapshotEntries: make(map[insolar.PulseNumber]*node.Snapshot),
	}
}

type MemoryStorage struct {
	lock            sync.RWMutex
	entries         []insolar.Pulse
	snapshotEntries map[insolar.PulseNumber]*node.Snapshot
}

// truncate deletes all entries except Count
func (m *MemoryStorage) truncate(count int) {
	if len(m.entries) <= count {
		return
	}

	truncatePulses := m.entries[:len(m.entries)-count]
	m.entries = m.entries[len(truncatePulses):]
	for _, p := range truncatePulses {
		delete(m.snapshotEntries, p.PulseNumber)
	}
}

func (m *MemoryStorage) AppendPulse(ctx context.Context, pulse insolar.Pulse) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.entries = append(m.entries, pulse)
	m.truncate(entriesCount)
	return nil
}

func (m *MemoryStorage) GetPulse(ctx context.Context, number insolar.PulseNumber) (insolar.Pulse, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for _, p := range m.entries {
		if p.PulseNumber == number {
			return p, nil
		}
	}

	return *insolar.GenesisPulse, ErrNotFound
}

func (m *MemoryStorage) GetLatestPulse(ctx context.Context) (insolar.Pulse, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if len(m.entries) == 0 {
		return *insolar.GenesisPulse, ErrNotFound
	}
	return m.entries[len(m.entries)-1], nil
}

func (m *MemoryStorage) Append(pulse insolar.PulseNumber, snapshot *node.Snapshot) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.snapshotEntries[pulse] = snapshot
	return nil
}

func (m *MemoryStorage) ForPulseNumber(pulse insolar.PulseNumber) (*node.Snapshot, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if s, ok := m.snapshotEntries[pulse]; ok {
		return s, nil
	}
	return nil, ErrNotFound
}
