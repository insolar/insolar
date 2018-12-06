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
	"github.com/insolar/insolar/ledger/jet"
	"github.com/ugorji/go/codec"
)

// Jet contain jet record.
type Jet struct {
	ID core.RecordID
}

// JetTree stores jet in a binary tree.
type JetTree struct {
	// TODO: implement tree.
	Jets []Jet
}

// Bytes serializes pulse.
func (t *JetTree) Bytes() []byte {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	enc.MustEncode(t)
	return buf.Bytes()
}

// GetDrop returns jet drop for a given pulse number.
func (db *DB) GetDrop(ctx context.Context, pulse core.PulseNumber) (*jet.JetDrop, error) {
	k := prefixkey(scopeIDJetDrop, pulse.Bytes())
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
func (db *DB) CreateDrop(ctx context.Context, pulse core.PulseNumber, prevHash []byte) (
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

	drop := jet.JetDrop{
		Pulse:    pulse,
		PrevHash: prevHash,
		Hash:     hw.Sum(nil),
	}
	return &drop, messages, nil
}

// SetDrop saves provided JetDrop in db.
func (db *DB) SetDrop(ctx context.Context, drop *jet.JetDrop) error {
	k := prefixkey(scopeIDJetDrop, drop.Pulse.Bytes())
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
func (db *DB) SetJetTree(ctx context.Context, pulse core.PulseNumber, tree *JetTree) error {
	k := prefixkey(scopeIDSystem, append([]byte{sysJetTree}, pulse.Bytes()...))
	_, err := db.get(ctx, k)
	if err == nil {
		return ErrOverride
	}

	return db.set(ctx, k, tree.Bytes())
}

// GetJetTree fetches tree for specified pulse.
func (db *DB) GetJetTree(ctx context.Context, pulse core.PulseNumber) (*JetTree, error) {
	k := prefixkey(scopeIDSystem, append([]byte{sysJetTree}, pulse.Bytes()...))
	buff, err := db.get(ctx, k)
	if err == nil {
		return nil, err
	}

	dec := codec.NewDecoder(bytes.NewReader(buff), &codec.CborHandle{})
	var tree JetTree
	err = dec.Decode(&tree)
	if err != nil {
		return nil, err
	}

	return &tree, nil
}
