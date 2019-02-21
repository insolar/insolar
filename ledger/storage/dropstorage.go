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

package storage

import (
	"context"
	"sync"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/pkg/errors"
)

// DropStorage jet-drops
//go:generate minimock -i github.com/insolar/insolar/ledger/storage.DropStorage -o ./ -s _mock.go
type DropStorage interface {
	CreateDrop(ctx context.Context, jetID core.RecordID, pulse core.PulseNumber, prevHash []byte) (
		*jet.JetDrop,
		[][]byte,
		uint64,
		error,
	)
	SetDrop(ctx context.Context, jetID core.RecordID, drop *jet.JetDrop) error
	GetDrop(ctx context.Context, jetID core.RecordID, pulse core.PulseNumber) (*jet.JetDrop, error)

	AddDropSize(ctx context.Context, dropSize *jet.DropSize) error
	SetDropSizeHistory(ctx context.Context, jetID core.RecordID, dropSizeHistory jet.DropSizeHistory) error
	GetDropSizeHistory(ctx context.Context, jetID core.RecordID) (jet.DropSizeHistory, error)

	GetJetSizesHistoryDepth() int
}

type dropStorage struct {
	DB                         DBContext                       `inject:""`
	PlatformCryptographyScheme core.PlatformCryptographyScheme `inject:""`

	addBlockSizeLock     sync.RWMutex
	jetSizesHistoryDepth int
}

func NewDropStorage(jetSizesHistoryDepth int) DropStorage {
	return &dropStorage{jetSizesHistoryDepth: jetSizesHistoryDepth}
}

// CreateDrop creates and stores jet drop for given pulse number.
//
// On success returns saved drop object, slot records, drop size.
func (ds *dropStorage) CreateDrop(ctx context.Context, jetID core.RecordID, pulse core.PulseNumber, prevHash []byte) (
	*jet.JetDrop,
	[][]byte,
	uint64,
	error,
) {
	var err error
	ds.DB.waitingFlight()

	hw := ds.PlatformCryptographyScheme.ReferenceHasher()
	_, err = hw.Write(prevHash)
	if err != nil {
		return nil, nil, 0, err
	}

	var messages [][]byte
	_, jetPrefix := jet.Jet(jetID)
	// messagesPrefix := prefixkey(scopeIDMessage, jetPrefix, pulse.Bytes())

	// err = db.db.View(func(txn *badger.Txn) error {
	// 	it := txn.NewIterator(badger.DefaultIteratorOptions)
	// 	defer it.Close()
	//
	// 	for it.Seek(messagesPrefix); it.ValidForPrefix(messagesPrefix); it.Next() {
	// 		val, err := it.Item().ValueCopy(nil)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		messages = append(messages, val)
	// 	}
	// 	return nil
	// })
	// if err != nil {
	// 	return nil, nil, 0, err
	// }

	var dropSize uint64
	recordPrefix := prefixkey(scopeIDRecord, jetPrefix, pulse.Bytes())

	err = ds.DB.GetBadgerDB().View(func(txn *badger.Txn) error {
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
			dropSize += uint64(len(val))
		}
		return nil
	})
	if err != nil {
		return nil, nil, 0, err
	}

	drop := jet.JetDrop{
		Pulse:    pulse,
		PrevHash: prevHash,
		Hash:     hw.Sum(nil),
	}
	return &drop, messages, dropSize, nil
}

// SetDrop saves provided JetDrop in db.
func (ds *dropStorage) SetDrop(ctx context.Context, jetID core.RecordID, drop *jet.JetDrop) error {
	_, prefix := jet.Jet(jetID)
	k := prefixkey(scopeIDJetDrop, prefix, drop.Pulse.Bytes())
	_, err := ds.DB.get(ctx, k)
	if err == nil {
		return ErrOverride
	}

	encoded, err := jet.Encode(drop)
	if err != nil {
		return err
	}
	return ds.DB.set(ctx, k, encoded)
}

// GetDrop returns jet drop for a given pulse number and jet id.
func (ds *dropStorage) GetDrop(ctx context.Context, jetID core.RecordID, pulse core.PulseNumber) (*jet.JetDrop, error) {
	_, prefix := jet.Jet(jetID)
	k := prefixkey(scopeIDJetDrop, prefix, pulse.Bytes())

	// buf, err := db.get(ctx, k)
	buf, err := ds.DB.get(ctx, k)
	if err != nil {
		return nil, err
	}
	drop, err := jet.Decode(buf)
	if err != nil {
		return nil, err
	}
	return drop, nil
}

// AddDropSize adds Jet drop size stats (required for split decision).
func (ds *dropStorage) AddDropSize(ctx context.Context, dropSize *jet.DropSize) error {
	ds.addBlockSizeLock.Lock()
	defer ds.addBlockSizeLock.Unlock()

	k := dropSizesPrefixKey(dropSize.JetID)
	buff, err := ds.DB.get(ctx, k)
	if err != nil && err != core.ErrNotFound {
		return errors.Wrapf(err, "[ AddDropSize ] Can't get object: %s", string(k))
	}

	var dropSizes = jet.DropSizeHistory{}
	if err != core.ErrNotFound {
		dropSizes, err = jet.DeserializeJetDropSizeHistory(ctx, buff)
		if err != nil {
			return errors.Wrapf(err, "[ AddDropSize ] Can't decode dropSizes")
		}

		if len([]jet.DropSize(dropSizes)) >= ds.jetSizesHistoryDepth {
			dropSizes = dropSizes[1:]
		}
	}

	dropSizes = append(dropSizes, *dropSize)

	return ds.DB.set(ctx, k, dropSizes.Bytes())
}

// SetDropSizeHistory saves drop sizes history.
func (ds *dropStorage) SetDropSizeHistory(ctx context.Context, jetID core.RecordID, dropSizeHistory jet.DropSizeHistory) error {
	ds.addBlockSizeLock.Lock()
	defer ds.addBlockSizeLock.Unlock()

	k := dropSizesPrefixKey(jetID)
	err := ds.DB.set(ctx, k, dropSizeHistory.Bytes())
	return errors.Wrap(err, "[ ResetDropSizeHistory ] Can't db.set")
}

// GetDropSizeHistory returns last drops sizes.
func (ds *dropStorage) GetDropSizeHistory(ctx context.Context, jetID core.RecordID) (jet.DropSizeHistory, error) {
	ds.addBlockSizeLock.RLock()
	defer ds.addBlockSizeLock.RUnlock()

	k := dropSizesPrefixKey(jetID)
	buff, err := ds.DB.get(ctx, k)
	if err != nil && err != core.ErrNotFound {
		return nil, errors.Wrap(err, "[ GetDropSizeHistory ] Can't db.set")
	}

	if err == core.ErrNotFound {
		return jet.DropSizeHistory{}, nil
	}

	dropSizes, err := jet.DeserializeJetDropSizeHistory(ctx, buff)
	if err != nil {
		return nil, errors.Wrapf(err, "[ GetDropSizeHistory ] Can't decode dropSizes")
	}

	return dropSizes, nil
}

// GetJetSizesHistoryDepth returns max amount of drop sizes
func (ds *dropStorage) GetJetSizesHistoryDepth() int {
	return ds.jetSizesHistoryDepth
}

func dropSizesPrefixKey(jetID core.RecordID) []byte {
	return prefixkey(scopeIDSystem, []byte{sysDropSizeHistory}, jetID.Bytes())
}
