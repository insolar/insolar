/*
 *    Copyright 2019 Insolar
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

package jet

import (
	"context"
	"io"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/record"
	"github.com/pkg/errors"
)

type Builder interface {
	Append(item Hashable) error
	Size(size uint64)
	PrevHash(prevHash []byte)
	Pulse(pn core.PulseNumber)

	Build() (JetDrop, error)
}

type Hashable interface {
	WriteHashData(w io.Writer) (int, error)
}

type builder struct {
	core.Hasher
	dropSize *uint64
	prevHash []byte
	pn       *core.PulseNumber
}

func NewBuilder(hasher core.Hasher) Builder {
	return &builder{
		Hasher: hasher,
	}
}

func (b *builder) Append(item Hashable) (err error) {
	_, err = item.WriteHashData(b.Hasher)
	return
}

func (b *builder) Size(size uint64) {
	b.dropSize = &size
}

func (b *builder) PrevHash(prevHash []byte) {
	b.prevHash = prevHash
}

func (b *builder) Pulse(pn core.PulseNumber) {
	b.pn = &pn
}

func (b *builder) Build() (JetDrop, error) {
	if b.prevHash == nil {
		return JetDrop{}, errors.New("prevHash is required")
	}
	if b.dropSize == nil {
		return JetDrop{}, errors.New("dropSize is required")
	}
	if b.pn == nil {
		return JetDrop{}, errors.New("pulseNumber is required")
	}

	return JetDrop{
		Pulse:    *b.pn,
		PrevHash: b.prevHash,
		Hash:     b.Hasher.Sum(nil),
		DropSize: *b.dropSize,
	}, nil
}

type Packer interface {
	Pack(ctx context.Context, jetID storage.JetID, pulse core.PulseNumber, prevHash []byte) (JetDrop, error)
}

func NewPacker(hasher core.Hasher, db storage.DBContext) Packer {
	return &packer{
		Builder:   NewBuilder(hasher),
		DBContext: db,
	}
}

type packer struct {
	Builder
	storage.DBContext
}

func (p *packer) Pack(ctx context.Context, jetID storage.JetID, pulse core.PulseNumber, prevHash []byte) (JetDrop, error) {
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
		return JetDrop{}, err
	}

	p.Pulse(pulse)
	p.PrevHash(prevHash)
	p.Size(dropSize)

	return p.Build()
}
