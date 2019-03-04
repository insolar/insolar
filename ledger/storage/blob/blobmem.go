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

package blob

import (
	"context"
	"sync"

	"github.com/insolar/insolar/core"
)

// simple aliases for key/value in memory map
type objectID = core.RecordID
type blobValue = []byte

type Storage struct {
	lock  sync.RWMutex
	blobs map[objectID]blobValue
}

func NewStorage() *Storage {
	return &Storage{blobs: map[objectID]blobValue{}}
}

func (s *Storage) Set(ctx context.Context, jetID core.JetID, pulseNumber core.PulseNumber, blob []byte) (*core.RecordID, error) {
	panic("implement me")
}

func (s *Storage) Get(ctx context.Context, jetID core.JetID, id *core.RecordID) ([]byte, error) {
	panic("implement me")
}
