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

package drop

import (
	"context"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/ledger/storage/record"
)

// Packer is an wrapper interface around process of building jetdrop
// It's considered that implementation of packer uses Bulder under the hood
//go:generate minimock -i github.com/insolar/insolar/ledger/storage/jet/drop.Packer -o ./ -s _mock.go
type Packer interface {
	Pack(ctx context.Context, jetID core.JetID, pulse core.PulseNumber, prevHash []byte) (jet.Drop, error)
}

// NewDbPacker creates db-based impl of packer
func NewDbPacker(hasher core.Hasher, db storage.DBContext) Packer {
	return &packer{
		Builder:   NewBuilder(hasher),
		DBContext: db,
	}
}

type packer struct {
	Builder
	storage.DBContext
}

// Pack creates new Drop through interactions with db and Builder
func (p *packer) Pack(ctx context.Context, jetID core.JetID, pulse core.PulseNumber, prevHash []byte) (jet.Drop, error) {
	p.DBContext.WaitingFlight()
	_, jetPrefix := jetID.Jet()

	var dropSize uint64
	recordPrefix := storage.IDRecordPrefixKey(jetPrefix, pulse)

	err := p.GetBadgerDB().View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(recordPrefix); it.ValidForPrefix(recordPrefix); it.Next() {
			val, err := it.Item().ValueCopy(nil)
			if err != nil {
				return err
			}

			err = p.Append(record.DeserializeRecord(val))
			if err != nil {
				return err
			}
			dropSize += uint64(len(val))
		}
		return nil
	})
	if err != nil {
		return jet.Drop{}, err
	}

	p.Pulse(pulse)
	p.PrevHash(prevHash)
	p.Size(dropSize)

	return p.Build()
}
