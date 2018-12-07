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
	"bytes"
	"context"
	"sync"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/ugorji/go/codec"
)

// GetDrop returns jet drop for a given pulse number and jet id.
func (db *DB) GetDrop(ctx context.Context, jetID core.RecordID, pulse core.PulseNumber) (*jet.JetDrop, error) {
	k := prefixkeyany(scopeIDJetDrop, jetID[:], pulse.Bytes())

	buf, err := db.get(ctx, k)
	if err != nil {
		return nil, err
	}
	drop, err := jet.Decode(buf)
	if err != nil {
		return nil, err
	}
	return drop, nil
}

// CreateDrop creates and stores jet drop for given pulse number.
//
// Previous JetDrop hash should be provided. On success returns saved drop and slot records.
func (db *DB) CreateDrop(ctx context.Context, jetID core.RecordID, pulse core.PulseNumber, prevHash []byte) (
	*jet.JetDrop,
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
		messagesPrefix := prefixkeyany(scopeIDMessage, jetID[:], pulse.Bytes())

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
		recordPrefix := prefixkeyany(scopeIDRecord, jetID[:], pulse.Bytes())

		jetDropHashError = db.db.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			for it.Seek(recordPrefix); it.ValidForPrefix(recordPrefix); it.Next() {
				val, err := it.Item().ValueCopy(nil)
				if err != nil {
					return err
				}
				_, err = hw.Write(val)
				if err != nil {
					return err
				}
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

	drop := jet.JetDrop{
		Pulse:    pulse,
		PrevHash: prevHash,
		Hash:     hw.Sum(nil),
	}
	return &drop, messages, nil
}

// SetDrop saves provided JetDrop in db.
func (db *DB) SetDrop(ctx context.Context, jetID core.RecordID, drop *jet.JetDrop) error {
	k := prefixkeyany(scopeIDJetDrop, jetID[:], drop.Pulse.Bytes())

	_, err := db.get(ctx, k)
	if err == nil {
		return ErrOverride
	}

	encoded, err := jet.Encode(drop)
	if err != nil {
		return err
	}
	return db.set(ctx, k, encoded)
}

// SetJetTree stores jet tree for specified pulse.
func (db *DB) SetJetTree(ctx context.Context, pulse core.PulseNumber, tree *jet.Tree) error {
	k := prefixkey(scopeIDSystem, append([]byte{sysJetTree}, pulse.Bytes()...))
	_, err := db.get(ctx, k)
	if err == nil {
		return ErrOverride
	}

	return db.set(ctx, k, tree.Bytes())
}

// GetJetTree fetches tree for specified pulse.
func (db *DB) GetJetTree(ctx context.Context, pulse core.PulseNumber) (*jet.Tree, error) {
	k := prefixkey(scopeIDSystem, append([]byte{sysJetTree}, pulse.Bytes()...))
	buff, err := db.get(ctx, k)
	if err != nil {
		return nil, err
	}

	dec := codec.NewDecoder(bytes.NewReader(buff), &codec.CborHandle{})
	var tree jet.Tree
	err = dec.Decode(&tree)
	if err != nil {
		return nil, err
	}

	return &tree, nil
}
