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
	"bytes"
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage/db"
	"github.com/insolar/insolar/ledger/storage/jet"
)

type dropStorageDB struct {
	DB db.DB `inject:""`
}

// NewStorageDB creates a new storage, that holds data in a db.
func NewStorageDB() *dropStorageDB { // nolint: golint
	return &dropStorageDB{}
}

type dropDbKey struct {
	jetPrefix []byte
	pn        core.PulseNumber
}

func (dk *dropDbKey) Scope() db.Scope {
	return db.ScopeJetDrop
}

func (dk *dropDbKey) ID() []byte {
	return bytes.Join([][]byte{dk.jetPrefix, dk.pn.Bytes()}, nil)
}

// ForPulse returns a jet.Drop for a provided pulse, that is stored in a db.
func (ds *dropStorageDB) ForPulse(ctx context.Context, jetID core.JetID, pulse core.PulseNumber) (jet.Drop, error) {
	k := dropDbKey{jetID.Prefix(), pulse}

	// buf, err := db.get(ctx, k)
	buf, err := ds.DB.Get(&k)
	if err != nil {
		return jet.Drop{}, err
	}
	drop, err := jet.Decode(buf)
	if err != nil {
		return jet.Drop{}, err
	}
	return *drop, nil
}

// Set saves a provided jet.Drop to a db.
func (ds *dropStorageDB) Set(ctx context.Context, jetID core.JetID, drop jet.Drop) error {
	k := dropDbKey{jetID.Prefix(), drop.Pulse}

	_, err := ds.DB.Get(&k)
	if err == nil {
		return ErrOverride
	}

	encoded, err := jet.Encode(&drop)
	if err != nil {
		return err
	}
	return ds.DB.Set(&k, encoded)
}

// Delete methods removes a drop from a storage. But the method mustn't be called for a db storage.
// Because db storage must be used only on a heavy-node.
func (ds *dropStorageDB) Delete(pulse core.PulseNumber) {
	panic("mustn't be called. because db storage must work only on a heavy node. heavy mustn't remove any data")
}
