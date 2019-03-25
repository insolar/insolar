//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package storage

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/ledger/storage/object"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/storage.ObjectStorage -o ./ -s _mock.go

// ObjectStorage returns objects and their meta
type ObjectStorage interface {
	GetBlob(ctx context.Context, jetID insolar.ID, id *insolar.ID) ([]byte, error)
	SetBlob(ctx context.Context, jetID insolar.ID, pulseNumber insolar.PulseNumber, blob []byte) (*insolar.ID, error)

	GetRecord(ctx context.Context, jetID insolar.ID, id *insolar.ID) (object.VirtualRecord, error)
	SetRecord(ctx context.Context, jetID insolar.ID, pulseNumber insolar.PulseNumber, rec object.VirtualRecord) (*insolar.ID, error)

	IterateIndexIDs(
		ctx context.Context,
		jetID insolar.ID,
		handler func(id insolar.ID) error,
	) error

	GetObjectIndex(
		ctx context.Context,
		jetID insolar.ID,
		id *insolar.ID,
	) (*object.Lifeline, error)

	SetObjectIndex(
		ctx context.Context,
		jetID insolar.ID,
		id *insolar.ID,
		idx *object.Lifeline,
	) error
}

type objectStorage struct {
	DB                         DBContext                          `inject:""`
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`
}

func NewObjectStorage() ObjectStorage {
	return new(objectStorage)
}

// GetBlob returns binary value stored by record ID.
// TODO: switch from reference to passing blob id for consistency - @nordicdyno 6.Dec.2018
func (os *objectStorage) GetBlob(ctx context.Context, jetID insolar.ID, id *insolar.ID) ([]byte, error) {
	jetPrefix := insolar.JetID(jetID).Prefix()
	k := prefixkey(scopeIDBlob, jetPrefix, id[:])
	return os.DB.Get(ctx, k)
}

// SetBlob saves binary value for provided pulse.
func (os *objectStorage) SetBlob(ctx context.Context, jetID insolar.ID, pulseNumber insolar.PulseNumber, blob []byte) (*insolar.ID, error) {
	id := object.CalculateIDForBlob(os.PlatformCryptographyScheme, pulseNumber, blob)
	jetPrefix := insolar.JetID(jetID).Prefix()
	k := prefixkey(scopeIDBlob, jetPrefix, id[:])

	// TODO: @andreyromancev. 16.01.19. Blob override is ok.
	// geterr := muxs.db.db.View(func(tx *badger.Txn) error {
	// 	_, err := tx.Get(k)
	// 	return err
	// })
	// if geterr == nil {
	// 	return id, ErrOverride
	// }
	// if geterr != badger.ErrKeyNotFound {
	// 	return nil, ErrNotFound
	// }

	err := os.DB.Set(ctx, k, blob)
	if err != nil {
		return nil, err
	}
	return id, nil
}

// GetRecord wraps matching transaction manager method.
func (os *objectStorage) GetRecord(ctx context.Context, jetID insolar.ID, id *insolar.ID) (object.VirtualRecord, error) {
	jetPrefix := insolar.JetID(jetID).Prefix()
	k := prefixkey(scopeIDRecord, jetPrefix, id[:])
	buf, err := os.DB.Get(ctx, k)
	if err != nil {
		return nil, err
	}
	return object.DeserializeRecord(buf), nil
}

// SetRecord wraps matching transaction manager method.
func (os *objectStorage) SetRecord(ctx context.Context, jetID insolar.ID, pulseNumber insolar.PulseNumber, rec object.VirtualRecord) (*insolar.ID, error) {
	var (
		id  *insolar.ID
		err error
	)
	err = os.DB.Update(ctx, func(tx *TransactionManager) error {
		id, err = tx.SetRecord(ctx, jetID, pulseNumber, rec)
		return err
	})
	if err != nil {
		return nil, err
	}
	return id, nil
}

// IterateIndexIDs iterates over index IDs on provided Jet ID.
func (os *objectStorage) IterateIndexIDs(
	ctx context.Context,
	jetID insolar.ID,
	handler func(id insolar.ID) error,
) error {
	jetPrefix := insolar.JetID(jetID).Prefix()
	prefix := prefixkey(scopeIDLifeline, jetPrefix)

	return os.DB.iterate(ctx, prefix, func(k, v []byte) error {
		pn := pulseNumFromKey(0, k)
		id := insolar.NewID(pn, k[insolar.PulseNumberSize:])
		err := handler(*id)
		if err != nil {
			return err
		}
		return nil
	})
}

// GetObjectIndex wraps matching transaction manager method.
func (os *objectStorage) GetObjectIndex(
	ctx context.Context,
	jetID insolar.ID,
	id *insolar.ID,
) (*object.Lifeline, error) {
	tx, err := os.DB.BeginTransaction(false)
	if err != nil {
		return nil, err
	}
	defer tx.Discard()

	idx, err := tx.GetObjectIndex(ctx, jetID, id)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

// SetObjectIndex wraps matching transaction manager method.
func (os *objectStorage) SetObjectIndex(
	ctx context.Context,
	jetID insolar.ID,
	id *insolar.ID,
	idx *object.Lifeline,
) error {
	return os.DB.Update(ctx, func(tx *TransactionManager) error {
		return tx.SetObjectIndex(ctx, jetID, id, idx)
	})
}
