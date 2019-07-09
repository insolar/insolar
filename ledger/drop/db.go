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

package drop

import (
	"bytes"
	"context"
	"math"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/jbenet/go-base58"
	"github.com/pkg/errors"
)

type DB struct {
	db store.DB
}

// NewDB creates a new storage, that holds data in a db.
func NewDB(db store.DB) *DB {
	return &DB{db: db}
}

type dropDbKey struct {
	jetPrefix []byte
	pn        insolar.PulseNumber
}

func (dk *dropDbKey) Scope() store.Scope {
	return store.ScopeJetDrop
}

func (dk *dropDbKey) ID() []byte {
	// order ( pn + jetPrefix ) is important: we use this logic for removing not finalized drops
	return bytes.Join([][]byte{dk.pn.Bytes(), dk.jetPrefix}, nil)
}

func NewDropDbKey(raw []byte) dropDbKey {
	dk := dropDbKey{}
	dk.pn = insolar.NewPulseNumber(raw)
	dk.jetPrefix = raw[dk.pn.Size():]

	return dk
}

// ForPulse returns a Drop for a provided pulse, that is stored in a db.
func (ds *DB) ForPulse(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber) (Drop, error) {
	k := dropDbKey{jetID.Prefix(), pulse}

	buf, err := ds.db.Get(&k)
	if err != nil {
		return Drop{}, err
	}
	drop, err := Decode(buf)
	if err != nil {
		return Drop{}, err
	}
	return *drop, nil
}

// Set saves a provided Drop to a db.
func (ds *DB) Set(ctx context.Context, drop Drop) error {
	k := dropDbKey{drop.JetID.Prefix(), drop.Pulse}

	_, err := ds.db.Get(&k)
	if err == nil {
		return ErrOverride
	}

	encoded := MustEncode(&drop)
	return ds.db.Set(&k, encoded)
}

func (ds *DB) TruncateHead(ctx context.Context, lastPulse insolar.PulseNumber) error {
	it := ds.db.NewIterator(&dropDbKey{jetPrefix: []byte{}, pn: math.MaxUint32}, true)
	defer it.Close()

	for it.Next() {
		key := NewDropDbKey(it.Key())
		if key.pn.Equal(lastPulse) {
			break
		}
		err := ds.db.Delete(&key)
		if err != nil {
			return errors.Wrapf(err, "[ DB.TruncateHead ] Can't Delete key: %+v", key)
		}

		inslogger.FromContext(ctx).Infof("[ DB.TruncateHead ] erased key. Pulse number: %s. Jet prefix: %s", key.pn.String(), base58.Encode(key.jetPrefix))
	}
	return nil
}
