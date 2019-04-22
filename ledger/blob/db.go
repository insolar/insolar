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

package blob

import (
	"bytes"
	"context"

	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/internal/ledger/store"
)

// DB implements persistent blob-storage.
type DB struct {
	db store.DB
}

// NewDB creates a new storage, that holds persistent data.
func NewDB(db store.DB) *DB {
	return &DB{
		db: db,
	}
}

type dbKey struct {
	id insolar.ID
}

func (k *dbKey) Scope() store.Scope {
	return store.ScopeBlob
}

func (k *dbKey) ID() []byte {
	return k.id[:]
}

// ForID returns Blob for provided id.
func (s *DB) ForID(ctx context.Context, id insolar.ID) (Blob, error) {
	b, err := s.db.Get(&dbKey{id: id})
	if err != nil {
		if err == store.ErrNotFound {
			err = ErrNotFound
		}
		return Blob{}, err
	}

	return decode(b)
}

// Set saves new Blob-value in storage.
func (s *DB) Set(ctx context.Context, id insolar.ID, blob Blob) error {
	// Blob override is ok.
	k := &dbKey{id: id}

	_, getErr := s.db.Get(k)
	if getErr != nil && getErr != store.ErrNotFound {
		return errors.Wrapf(getErr, "got db error on key %v get", k)
	} else if getErr == nil {
		return ErrOverride
	}

	b := mustEncode(blob)

	err := s.db.Set(k, b)
	if err != nil {
		return err
	}

	stats.Record(ctx,
		statBlobInStorageSize.M(int64(len(b))),
		statBlobInStorageCount.M(1),
	)
	return nil
}

// mustEncode serializes blob struct.
func mustEncode(blob Blob) []byte {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	err := enc.Encode(blob)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// decode deserializes bytes to blob struct.
func decode(buf []byte) (Blob, error) {
	dec := codec.NewDecoder(bytes.NewReader(buf), &codec.CborHandle{})
	var blob Blob
	err := dec.Decode(&blob)
	return blob, err
}
