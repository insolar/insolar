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
	"fmt"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

// PulseStorage implements core.PulseStorage
type PulseStorage struct {
	db           *DB
	rwLock       sync.RWMutex
	currentPulse *core.Pulse
}

// NewPulseStorage creates new pulse storage
func NewPulseStorage(db *DB) *PulseStorage {
	return &PulseStorage{db: db}
}

// Current returns current pulse of the system
func (ps *PulseStorage) Current(ctx context.Context) (*core.Pulse, error) {
	pulse, err := ps.pulseFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if pulse != nil {
		fmt.Println("*********FROMCTX Err is nil. Pulse is: ", pulse)
		return pulse, nil
	}

	currentPulse := ps.getCachedPulse()
	if currentPulse != nil {
		return currentPulse, nil
	}

	currentPulse, err = ps.reloadPulse(ctx)
	if err != nil {
		return nil, err
	}

	return currentPulse, nil
}

func (ps *PulseStorage) getCachedPulse() *core.Pulse {
	ps.rwLock.RLock()
	defer ps.rwLock.RUnlock()

	return ps.currentPulse
}

func (ps *PulseStorage) reloadPulse(ctx context.Context) (*core.Pulse, error) {
	ps.rwLock.Lock()
	defer ps.rwLock.Unlock()

	if ps.currentPulse == nil {
		currentPulse, err := ps.db.GetLatestPulse(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "[ PulseStorage.reloadPulse ] Can't GetLatestPulse")
		}
		ps.currentPulse = &currentPulse.Pulse
	}

	return ps.currentPulse, nil
}

func (ps *PulseStorage) pulseFromContext(ctx context.Context) (*core.Pulse, error) {
	pulseNumber, err := core.NewPulseNumberFromContext(ctx)
	if err != nil {
		if err == core.ErrNoPulseInContext {
			return nil, nil
		}
		return nil, err
	}

	currentPulse := ps.getCachedPulse()
	inslogger.FromContext(ctx).Debugf("[ PulseStorage.pulseFromContext ] Getting pulse %d from context", pulseNumber)
	if currentPulse != nil {
		if currentPulse.PulseNumber == pulseNumber {
			return currentPulse, nil
		}
		inslogger.FromContext(ctx).Warnf(
			"[ PulseStorage.pulseFromContext ] Current pulse (%d) differs from context pulse (%d)",
			currentPulse.PulseNumber,
			pulseNumber,
		)
	}

	pulse, err := ps.db.GetPulse(ctx, pulseNumber)
	if err != nil {
		return nil, errors.Wrapf(err, "[ PulseStorage.pulseFromContext ] Can't GetPulse %d from context", pulseNumber)
	}

	return &pulse.Pulse, nil
}

func (ps *PulseStorage) Set(pulse *core.Pulse) {
	ps.rwLock.Lock()
	defer ps.rwLock.Unlock()

	ps.currentPulse = pulse
}

// Lock takes lock on parent's pulse storage
func (ps *PulseStorage) Lock() {
	ps.rwLock.Lock()
}

// Unlock takes unlock on parent's pulse storage
func (ps *PulseStorage) Unlock() {
	ps.rwLock.Unlock()
}
