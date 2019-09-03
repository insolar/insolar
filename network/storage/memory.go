//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package storage

import (
	"context"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/node"
	"sync"
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
