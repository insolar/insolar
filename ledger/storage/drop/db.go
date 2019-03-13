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

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage"
)

type dropStorageDB struct {
	DB storage.DBContext `inject:""`
}

// NewStorageDB creates a new storage, that holds data in a db
func NewStorageDB() *dropStorageDB { // nolint: golint
	return &dropStorageDB{}
}

// ForPulse returns a jet.Drop for a provided pulse, that is stored in a db
func (ds *dropStorageDB) ForPulse(ctx context.Context, jetID core.JetID, pulse core.PulseNumber) (Drop, error) {
	prefix := jetID.Prefix()
	k := storage.JetDropPrefixKey(prefix, pulse)

	// buf, err := db.get(ctx, k)
	buf, err := ds.DB.Get(ctx, k)
	if err != nil {
		return Drop{}, err
	}
	drop, err := Decode(buf)
	if err != nil {
		return Drop{}, err
	}
	return *drop, nil
}

// Set saves a provided jet.Drop to a db
func (ds *dropStorageDB) Set(ctx context.Context, jetID core.JetID, drop Drop) error {
	prefix := jetID.Prefix()
	k := storage.JetDropPrefixKey(prefix, drop.Pulse)
	_, err := ds.DB.Get(ctx, k)
	if err == nil {
		return storage.ErrOverride
	}

	encoded, err := Encode(&drop)
	if err != nil {
		return err
	}
	return ds.DB.Set(ctx, k, encoded)
}
