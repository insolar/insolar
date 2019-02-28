package db

import (
	"context"
	"path/filepath"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

type badgerDB struct {
	backend *badger.DB
}

func NewBadgerDB(conf configuration.Ledger) (*badgerDB, error) {
	dir, err := filepath.Abs(conf.Storage.DataDirectory)
	if err != nil {
		return nil, err
	}

	bdb, err := badger.Open(badger.Options{
		Dir:      dir,
		ValueDir: dir,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to open badger")
	}

	db := &badgerDB{
		backend: bdb,
	}
	return db, nil
}

func (b *badgerDB) Get(key Key) (value []byte, err error) {
	fullKey := append(key.Scope().Bytes(), key.Key()...)

	err = b.backend.View(func(txn *badger.Txn) error {
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
			return nil, core.ErrNotFound
		}
		return nil, err
	}

	return
}

func (b *badgerDB) Set(key Key, value []byte) error {
	fullKey := append(key.Scope().Bytes(), key.Key()...)

	err := b.backend.Update(func(txn *badger.Txn) error {
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

// Stop stops DB component.
func (b *badgerDB) Stop(ctx context.Context) error {
	return b.backend.Close()
}
