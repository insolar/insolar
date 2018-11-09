/*
 *    Copyright 2018 Insolar
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

package localstorage

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage"
)

// LocalStorage allows a node to save local data.
type LocalStorage struct {
	db *storage.DB
}

// NewLocalStorage create new storage instance.
func NewLocalStorage(db *storage.DB) (*LocalStorage, error) {
	return &LocalStorage{db: db}, nil
}

// SetMessage saves data in storage.
func (s *LocalStorage) Set(ctx context.Context, pulse core.PulseNumber, key []byte, data []byte) error {
	return s.db.SetLocalData(ctx, pulse, key, data)
}

// GetMessage retrieves data from storage.
func (s *LocalStorage) Get(ctx context.Context, pulse core.PulseNumber, key []byte) ([]byte, error) {
	buff, err := s.db.GetLocalData(ctx, pulse, key)
	if err == storage.ErrNotFound {
		return nil, ErrNotFound
	}
	return buff, err
}

func (s *LocalStorage) Iterate(ctx context.Context, pulse core.PulseNumber, prefix []byte, handler func(k, v []byte) error) error {
	return s.db.IterateLocalData(ctx, pulse, prefix, handler)
}
