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
		// Because first 2 bites of pulse number and first 65536 pulses a are used by system needs and pulse numbers are related to the seconds of Unix time
		// for calculation pulse numbers we use the formula = unix.Now() - firstPulseDate + 65536
		genesisPulse := core.Pulse{PulseNumber: core.FirstPulseNumber}
		predefinedEntropy := []byte{138, 67, 169, 65, 13, 4, 211, 121, 35, 73, 128, 81, 138, 164, 87, 139,
			150, 104, 24, 255, 159, 10, 172, 233, 183, 61, 183, 192, 169, 103, 187, 209,
			181, 235, 43, 188, 164, 151, 138, 213, 231, 222, 27, 244, 42, 194, 55, 133,
			30, 202, 50, 246, 119, 180, 59, 143, 130, 248, 87, 28, 155, 33, 157, 30}
		copy(genesisPulse.Entropy[:], predefinedEntropy[:core.EntropySize])

		err = db.SavePulse(&genesisPulse)
		if err != nil {
			return nil, errors.Wrap(err, "problems with init database")
		}
		err = db.SetLastPulse(&genesisPulse)
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
