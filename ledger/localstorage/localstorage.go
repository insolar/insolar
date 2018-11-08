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
	"bytes"
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
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

// SetMessage saves message in storage.
func (s *LocalStorage) SetMessage(ctx context.Context, msg core.SignedMessage) (*core.RecordID, error) {
	buff, err := message.SignedToBytes(msg)
	if err != nil {
		return nil, err
	}

	return s.db.SetBlob(ctx, msg.Pulse(), buff)
}

// GetMessage retrieves message from storage.
func (s *LocalStorage) GetMessage(ctx context.Context, id core.RecordID) (core.SignedMessage, error) {
	buff, err := s.db.GetBlob(ctx, &id)
	if err != nil {
		return nil, err
	}

	return message.DeserializeSigned(bytes.NewBuffer(buff))
}
