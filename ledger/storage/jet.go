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

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/jetdrop"
)

// GetDrop returns jet drop for a given pulse number.
func (db *DB) GetDrop(ctx context.Context, pulse core.PulseNumber) (*jetdrop.JetDrop, error) {
	k := prefixkey(scopeIDJetDrop, pulse.Bytes())
	buf, err := db.get(ctx, k)
	if err != nil {
		return nil, err
	}
	drop, err := jetdrop.Decode(buf)
	if err != nil {
		return nil, err
	}
	return drop, nil
}

// CreateDrop creates and stores jet drop for given pulse number.
//
// Previous JetDrop hash should be provided. On success returns saved drop and slot records.
func (db *DB) CreateDrop(ctx context.Context, pulse core.PulseNumber, prevHash []byte) (
	*jetdrop.JetDrop,
	[][]byte,
	error,
) {
	var err error
	db.waitinflight()

	hw := db.PlatformCryptographyScheme.ReferenceHasher()
	_, err = hw.Write(prevHash)
	if err != nil {
		return nil, nil, err
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	var messages [][]byte
	var messagesError error

	go func() {
		messagesPrefix := make([]byte, core.PulseNumberSize+1)
		messagesPrefix[0] = scopeIDMessage
		copy(messagesPrefix[1:], pulse.Bytes())

		messagesError = db.db.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			for it.Seek(messagesPrefix); it.ValidForPrefix(messagesPrefix); it.Next() {
				val, err := it.Item().ValueCopy(nil)
				if err != nil {
					return err
				}
				messages = append(messages, val)
			}
			return nil
		})

		wg.Done()
	}()

	var jetDropHashError error

	go func() {
		recordPrefix := make([]byte, core.PulseNumberSize+1)
		recordPrefix[0] = scopeIDRecord
		copy(recordPrefix[1:], pulse.Bytes())

		jetDropHashError = db.db.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			for it.Seek(recordPrefix); it.ValidForPrefix(recordPrefix); it.Next() {
				val, err := it.Item().ValueCopy(nil)
				if err != nil {
					return err
				}
				hw.Sum(val)
			}
			return nil
		})

		wg.Done()
	}()

	wg.Wait()

	if messagesError != nil {
		return nil, nil, messagesError
	}
	if jetDropHashError != nil {
		return nil, nil, jetDropHashError
	}

	drop := jetdrop.JetDrop{
		Pulse:    pulse,
		PrevHash: prevHash,
		Hash:     hw.Sum(nil),
	}
	return &drop, messages, nil
}

// SetDrop saves provided JetDrop in db.
func (db *DB) SetDrop(ctx context.Context, drop *jetdrop.JetDrop) error {
	k := prefixkey(scopeIDJetDrop, drop.Pulse.Bytes())
	_, err := db.get(ctx, k)
	if err == nil {
		return ErrOverride
	}

	encoded, err := jetdrop.Encode(drop)
	if err != nil {
		return err
	}
	return db.set(ctx, k, encoded)
}
