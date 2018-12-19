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

package storage

import (
	"context"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

// PulseStorage implements core.PulseStorage
type PulseStorage struct {
	db     *DB
	rwLock sync.RWMutex
}

// NewPulseStorage creates new pulse storage
func NewPulseStorage(db *DB) *PulseStorage {
	return &PulseStorage{db: db}
}

// Current returns current pulse of the system
func (ps *PulseStorage) Current(ctx context.Context) (*core.Pulse, error) {
	ps.rwLock.RLock()
	defer ps.rwLock.RUnlock()

	latestPulse, err := ps.db.GetLatestPulse(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "pulse manager failed to get current pulse")
	}

	return &latestPulse.Pulse, nil
}

// Lock takes lock on parent's pulse storage
func (ps *PulseStorage) Lock() {
	ps.rwLock.Lock()
}

// Unlock takes unlock on parent's pulse storage
func (ps *PulseStorage) Unlock() {
	ps.rwLock.Unlock()
}
