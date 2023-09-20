package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dgraph-io/badger"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/store"
)

func openDB(dataDir string) (*store.BadgerDB, func()) {
	if err := checkDirectory(dataDir); err != nil {
		fatalf("Database directory '%v' open failed. Error: \"%v\"", dataDir, err)
	}

	ops := badger.DefaultOptions(dataDir)
	dbWrapped, err := store.NewBadgerDB(ops)
	if err != nil {
		fatalf("failed open database directory %v: %v", dataDir, err)
	}
	close := func() {
		err := dbWrapped.Backend().Close()
		if err != nil {
			fatalf("failed close database directory %v: %v", dataDir, err)
		}
	}
	return dbWrapped, close
}

func checkDirectory(dir string) error {
	_, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	manifest := filepath.Join(dir, "MANIFEST")
	if _, err := os.Stat(manifest); err != nil {
		return fmt.Errorf("failed %v file check", manifest)
	}
	return nil
}

type pulseKey struct {
	scope store.Scope
	pn    insolar.PulseNumber
}

func (k pulseKey) Scope() store.Scope {
	return k.scope
}

func (k pulseKey) ID() []byte {
	return k.pn.Bytes()
}

type iteration func(k, v []byte) error

type iterOptions struct {
	keysOnly     bool
	prefetchSize int
	counter      bool
}

func iterate(
	db *store.BadgerDB,
	start store.Key,
	opts *iterOptions,
	fns ...iteration,
) {
	txn := db.Backend().NewTransaction(false)
	defer txn.Discard()

	badgerOpts := badger.DefaultIteratorOptions
	if opts == nil {
		opts = &iterOptions{
			prefetchSize: 100,
		}
	}
	if opts.keysOnly {
		badgerOpts.PrefetchValues = false
	} else if opts.prefetchSize > 0 {
		badgerOpts.PrefetchSize = opts.prefetchSize
	}
	it := txn.NewIterator(badgerOpts)
	defer it.Close()
	prefix := append(start.Scope().Bytes(), start.ID()...)

	i := 0
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		i++
		if opts.counter {
			_, _ = fmt.Printf("\r%20s scanned", formatInt(i, " "))
		}
		var err error
		var b []byte

		item := it.Item()
		k := item.Key()
		if !opts.keysOnly {
			b = make([]byte, 1024)
			b, err = item.ValueCopy(b)
			if err != nil {
				panic(err)
			}
		}
		for _, fn := range fns {
			if fn == nil {
				continue
			}
			err = fn(k, b)
			if err != nil {
				panic(err)
			}
		}
	}
	if opts.counter {
		fmt.Println()
	}
}

type scopeKey struct {
	scope store.Scope
}

func (k scopeKey) Scope() store.Scope {
	return k.scope
}

func (k scopeKey) ID() []byte {
	return nil
}

func readValueByKey(db *store.BadgerDB, key []byte) ([]byte, error) {
	var value []byte
	err := db.Backend().View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		value, err = item.ValueCopy(value)
		return err
	})
	return value, err
}
