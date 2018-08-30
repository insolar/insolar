/*
 *    Copyright 2018 INS Ecosystem
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

package badgerdb

import (
	"path/filepath"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/pkg/errors"
)

const (
	scopeIDLifeline byte = 1
	scopeIDRecord   byte = 2
	scopeIDJetDrop  byte = 3
	scopeIDEntropy  byte = 4
)

// Store represents BadgerDB storage implementation.
type Store struct {
	db           *badger.DB
	currentPulse record.PulseNum
	zeroRef      record.Reference
}

func setOptions(o *badger.Options) *badger.Options {
	newo := &badger.Options{}
	if o != nil {
		*newo = *o
	} else {
		*newo = badger.DefaultOptions
	}
	return newo
}

// NewStore returns badgerdb.Store with BadgerDB instance initialized by opts.
// Creates database in provided dir or in current directory if dir parameter is empty.
func NewStore(dir string, opts *badger.Options) (*Store, error) {
	opts = setOptions(opts)
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	opts.Dir = dir
	opts.ValueDir = dir

	db, err := badger.Open(*opts)
	if err != nil {
		return nil, errors.Wrap(err, "local database open failed")
	}

	bl := &Store{
		db: db,
	}

	return bl, nil
}

// Close wraps BadgerDB Close method.
//
// From https://godoc.org/github.com/dgraph-io/badger#DB.Close:
// «It's crucial to call it to ensure all the pending updates make their way to disk.
// Calling DB.Close() multiple times is not safe and wouldcause panic.»
func (s *Store) Close() error {
	// TODO: add close flag and mutex guard on Close method
	return s.db.Close()
}

// Get gets value by key in BadgerDB storage.
func (s *Store) Get(key []byte) ([]byte, error) {
	var buf []byte
	// TODO: handle transaction conflicts.
	txerr := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return storage.ErrNotFound
			}
			return err
		}
		buf, err = item.Value()
		if err != nil {
			return err
		}
		return err
	})
	if txerr != nil {
		buf = nil
	}
	return buf, txerr
}

// Set stores value by key in BadgerDB.
func (s *Store) Set(key, value []byte) error {
	// TODO: handle transaction conflicts.
	txerr := s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
	return txerr
}

// GetRecord returns record from BadgerDB by *record.Reference.
//
// It returns storage.ErrNotFound if the DB does not contain the key.
func (s *Store) GetRecord(ref *record.Reference) (record.Record, error) {
	var raw *record.Raw

	k := prefixkey(scopeIDRecord, ref.Bytes())
	txerr := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(k)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return storage.ErrNotFound
			}
			return err
		}
		buf, err := item.Value()
		if err != nil {
			return err
		}
		raw, err = record.DecodeToRaw(buf)
		return err
	})
	// TODO: check transaction conflict
	if txerr != nil {
		return nil, txerr
	}
	return raw.ToRecord(), nil
}

// SetRecord stores record in BadgerDB and returns *record.Reference of new record.
func (s *Store) SetRecord(rec record.Record) (*record.Reference, error) {
	raw, err := record.EncodeToRaw(rec)
	if err != nil {
		return nil, err
	}
	ref := &record.Reference{
		Domain: rec.Domain().Record,
		Record: record.ID{
			Pulse: s.GetCurrentPulse(),
			Hash:  raw.Hash(),
		},
	}
	k := prefixkey(scopeIDRecord, ref.Bytes())
	val := record.MustEncodeRaw(raw)
	txerr := s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(k, val)
	})
	// TODO: check transaction conflict
	if txerr != nil {
		return nil, txerr
	}
	return ref, nil
}

func (s *Store) GetClassIndex(*record.Reference) (*index.ClassLifeline, error) {
	panic("not implemented")
}

func (s *Store) SetClassIndex(*record.Reference, *index.ClassLifeline) error {
	panic("not implemented")
}

func (s *Store) GetObjectIndex(*record.Reference) (*index.ObjectLifeline, error) {
	panic("not implemented")
}

func (s *Store) SetObjectIndex(*record.Reference, *index.ObjectLifeline) error {
	panic("not implemented")
}

func (s *Store) GetDrop(record.PulseNum) (*jetdrop.JetDrop, error) {
	panic("not implemented")
}

func (s *Store) SetDrop(record.PulseNum, *jetdrop.JetDrop) (*jetdrop.JetDrop, error) {
	panic("not implemented")
}

// SetEntropy stores given entropy for given pulse in storage.
//
// Entropy is used for calculating node roles.
func (s *Store) SetEntropy(pulse record.PulseNum, entropy []byte) error {
	k := prefixkey(scopeIDEntropy, record.EncodePulseNum(pulse))
	return s.Set(k, entropy)
}

// GetEntropy returns entropy from storage for given pulse.
//
// Entropy is used for calculating node roles.
func (s *Store) GetEntropy(pulse record.PulseNum) ([]byte, error) {
	k := prefixkey(scopeIDEntropy, record.EncodePulseNum(pulse))
	return s.Get(k)
}

// SetCurrentPulse sets current pulse number.
func (s *Store) SetCurrentPulse(pulse record.PulseNum) {
	s.currentPulse = pulse
}

// GetCurrentPulse returns current pulse number.
func (s *Store) GetCurrentPulse() record.PulseNum {
	return s.currentPulse
}
