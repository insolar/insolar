/*
 *    Copyright 2019 Insolar Technologies
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

package pulsarstorage

import (
	"bytes"
	"encoding/gob"
	"path/filepath"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

type RecordID string

const (
	LastPulseRecordID RecordID = "lastPulse"
	PulseRecordID     RecordID = "pulse"
)

// NewDB returns pulsar.storage.db with BadgerDB instance initialized by opts.
// Creates database in provided dir or in current directory if dir parameter is empty.
func NewStorageBadger(conf configuration.Pulsar, opts *badger.Options) (PulsarStorage, error) {
	gob.Register(core.Pulse{})
	opts = setOptions(opts)
	dir, err := filepath.Abs(conf.Storage.DataDirectory)
	if err != nil {
		return nil, err
	}

	opts.Dir = dir
	opts.ValueDir = dir

	bdb, err := badger.Open(*opts)
	if err != nil {
		return nil, errors.Wrap(err, "local database open failed")
	}

	db := &BadgerStorageImpl{
		db: bdb,
	}

	pulse, err := db.GetLastPulse()
	if pulse.PulseNumber == 0 || err != nil {
		err = db.SavePulse(core.GenesisPulse)
		if err != nil {
			return nil, errors.Wrap(err, "problems with init database")
		}
		err = db.SetLastPulse(core.GenesisPulse)
		if err != nil {
			return nil, errors.Wrap(err, "problems with init database")
		}
	}

	return db, nil
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

type BadgerStorageImpl struct {
	db *badger.DB
}

func (storage *BadgerStorageImpl) GetLastPulse() (*core.Pulse, error) {
	var pulseNumber core.Pulse

	err := storage.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(LastPulseRecordID))
		if err != nil {
			return err
		}
		val, err := item.Value()
		if err != nil {
			return err
		}

		r := bytes.NewBuffer(val)
		decoder := gob.NewDecoder(r)
		err = decoder.Decode(&pulseNumber)
		if err != nil {
			return err
		}

		return nil
	})
	return &pulseNumber, err
}

func (storage *BadgerStorageImpl) SetLastPulse(pulse *core.Pulse) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(pulse)
	if err != nil {
		return err
	}
	return storage.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(LastPulseRecordID), buffer.Bytes())
		return err
	})
}

func (storage *BadgerStorageImpl) SavePulse(pulse *core.Pulse) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(pulse)
	if err != nil {
		return err
	}
	pulseNumber := pulse.PulseNumber.Bytes()
	key := []byte(PulseRecordID)
	key = append(key, pulseNumber...)

	return storage.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, buffer.Bytes())
		return err
	})
}

func (storage *BadgerStorageImpl) Close() error {
	return storage.db.Close()
}
