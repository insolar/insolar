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
	"path/filepath"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is returned when value was not found.
	ErrNotFound = errors.New("value not found")
	ErrBadPulse = errors.New("pulse should be bigger than latest")
)

type pulseKey insolar.PulseNumber

func (k pulseKey) Scope() Scope {
	return ScopePulse
}

func (k pulseKey) ID() []byte {
	return append([]byte{prefixPulse}, insolar.PulseNumber(k).Bytes()...)
}

// DB provides a simple key-value store interface for persisting data.
type DB interface {
	Get(key Key) (value []byte, err error)
	Set(key Key, value []byte) error
}

// Key represents a key for the key-value store. Scope is required to separate different DB clients and should be
// unique.
type Key interface {
	Scope() Scope
	ID() []byte
}

// Scope separates DB clients.
type Scope byte

// Bytes returns binary scope representation.
func (s Scope) Bytes() []byte {
	return []byte{byte(s)}
}

const (
	// ScopePulse is the scope for pulse storage.
	ScopePulse Scope = 1
)

type BadgerDB struct {
	db *badger.DB
}

// NewTestBadgerDB creates new badger DB instance for tests.
// Flag DoNotCompact lets close db immediately
func NewTestBadgerDB(conf configuration.ServiceNetwork) (*BadgerDB, error) {
	dir, err := filepath.Abs(conf.CacheDirectory)
	if err != nil {
		return nil, err
	}

	ops := badger.DefaultOptions
	ops.SyncWrites = false
	ops.DoNotCompact = true
	ops.ValueDir = dir
	ops.Dir = dir
	bdb, err := badger.Open(ops)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open badger")
	}

	db := &BadgerDB{
		db: bdb,
	}
	return db, nil
}

// NewBadgerDB creates new badger DB instance. Configuration should contain DataDirectory option. Badger will create
// files there.
func NewBadgerDB(conf configuration.ServiceNetwork) (*BadgerDB, error) {
	dir, err := filepath.Abs(conf.CacheDirectory)
	if err != nil {
		return nil, err
	}

	ops := badger.DefaultOptions(dir)
	bdb, err := badger.Open(ops)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open badger")
	}

	db := &BadgerDB{
		db: bdb,
	}
	return db, nil
}

// Get returns value for specified key or an error. A copy of a value will be returned (i.e. getting large value can be
// long).
func (b *BadgerDB) Get(key Key) (value []byte, err error) {
	fullKey := append(key.Scope().Bytes(), key.ID()...)

	err = b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(fullKey)
		if err != nil {
			return err
		}
		value, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return
}

// Set stores value for a key.
func (b *BadgerDB) Set(key Key, value []byte) error {
	fullKey := append(key.Scope().Bytes(), key.ID()...)

	err := b.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(fullKey, value)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Stop gracefully stops all disk writes. After calling this, it's safe to kill the process without losing data.
func (b *BadgerDB) Stop(ctx context.Context) error {
	return b.db.Close()
}
