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

// GetRecord returns record from BadgerDB by *record.Reference.
//
// It returns storage.ErrNotFound if the DB does not contain the key.
func (s *Store) GetRecord(ref *record.Reference) (record.Record, error) {
	tx := s.BeginTransaction(false)
	rec, err := tx.GetRecord(ref)
	if err != nil {
		return nil, err
	}
	tx.Discard()
	return rec, nil
}

// SetRecord stores record in BadgerDB and returns *record.Reference of new record.
func (s *Store) SetRecord(rec record.Record) (*record.Reference, error) {
	tx := s.BeginTransaction(true)
	ref, err := tx.SetRecord(rec)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return ref, nil
}

// GetClassIndex fetches class lifeline's index.
func (s *Store) GetClassIndex(ref *record.Reference) (*index.ClassLifeline, error) {
	tx := s.BeginTransaction(false)
	idx, err := tx.GetClassIndex(ref)
	if err != nil {
		return nil, err
	}
	tx.Discard()
	return idx, nil
}

// SetClassIndex stores class lifeline index.
func (s *Store) SetClassIndex(ref *record.Reference, idx *index.ClassLifeline) error {
	tx := s.BeginTransaction(true)
	err := tx.SetClassIndex(ref, idx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// GetObjectIndex fetches object lifeline index.
func (s *Store) GetObjectIndex(ref *record.Reference) (*index.ObjectLifeline, error) {
	tx := s.BeginTransaction(false)
	idx, err := tx.GetObjectIndex(ref)
	if err != nil {
		return nil, err
	}
	tx.Discard()
	return idx, nil
}

// SetObjectIndex stores object lifeline index.
func (s *Store) SetObjectIndex(ref *record.Reference, idx *index.ObjectLifeline) error {
	tx := s.BeginTransaction(true)
	err := tx.SetObjectIndex(ref, idx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// GetDrop returns jet drop for a given pulse number.
func (s *Store) GetDrop(pulse record.PulseNum) (*jetdrop.JetDrop, error) {
	tx := s.BeginTransaction(false)
	drop, err := tx.GetDrop(pulse)
	if err != nil {
		return nil, err
	}
	tx.Discard()
	return drop, nil
}

// SetDrop stores jet drop for given pulse number.
// Previous JetDrop should be provided.
// On success returns saved drop hash.
func (s *Store) SetDrop(pulse record.PulseNum, prevdrop *jetdrop.JetDrop) (*jetdrop.JetDrop, error) {
	tx := s.BeginTransaction(true)
	drop, err := tx.SetDrop(pulse, prevdrop)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return drop, nil
}

// SetEntropy stores given entropy for given pulse in storage.
//
// Entropy is used for calculating node roles.
func (s *Store) SetEntropy(pulse record.PulseNum, entropy []byte) error {
	tx := s.BeginTransaction(true)
	err := tx.SetEntropy(pulse, entropy)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// GetEntropy returns entropy from storage for given pulse.
//
// Entropy is used for calculating node roles.
func (s *Store) GetEntropy(pulse record.PulseNum) ([]byte, error) {
	tx := s.BeginTransaction(false)
	idx, err := tx.GetEntropy(pulse)
	if err != nil {
		return nil, err
	}
	tx.Discard()
	return idx, nil
}

// SetCurrentPulse sets current pulse number.
func (s *Store) SetCurrentPulse(pulse record.PulseNum) {
	s.currentPulse = pulse
}

// GetCurrentPulse returns current pulse number.
func (s *Store) GetCurrentPulse() record.PulseNum {
	return s.currentPulse
}

func (s *Store) BeginTransaction(update bool) TransactionManager {
	return TransactionManager{
		store: s,
		txn:   s.db.NewTransaction(update),
	}
}
