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

type RecordId string

const (
	LastPulseRecordId RecordId = "lastPulse"
)

// NewDB returns pulsar.storage.db with BadgerDB instance initialized by opts.
// Creates database in provided dir or in current directory if dir parameter is empty.
func NewStorageBadger(conf configuration.Pulsar, opts *badger.Options) (*BadgerStorageImpl, error) {
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
		item, err := txn.Get([]byte(LastPulseRecordId))
		if err != nil {
			return err
		}
		val, err := item.Value()
		if err != nil {
			return err
		}

		r := bytes.NewBuffer(val)
		decoder := gob.NewDecoder(r)
		err = decoder.Decode(pulseNumber)
		if err != nil {
			return err
		}

		return nil
	})
	return &pulseNumber, err
}

func (storage *BadgerStorageImpl) UpdatePulse(pulse *core.Pulse) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(pulse)
	if err != nil {
		return err
	}
	return storage.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(LastPulseRecordId), buffer.Bytes())
		return err
	})
}
