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

package adapter

import (
	"fmt"

	"github.com/insolar/insolar/conveyor/adapter/adapterid"
)

type Storage struct {
	adapters map[adapterid.ID]TaskSink
}

func (s *Storage) GetAdapterByID(id adapterid.ID) TaskSink {
	return s.adapters[id]
}

func (s *Storage) Register(adapter TaskSink) {
	id := adapter.GetAdapterID()
	_, ok := s.adapters[id]
	if ok {
		panic(fmt.Sprintf("[ StorageManager.Register ] adapter ID '%s' already exists", id.String()))
	}

	s.adapters[id] = adapter
}

func (s *Storage) GetRegisteredAdapters() []interface{} {
	var result []interface{}

	for _, adapter := range s.adapters {
		result = append(result, adapter)
	}

	return result
}

var StorageManager Storage

func NewStorage() Storage {
	return Storage{
		adapters: make(map[adapterid.ID]TaskSink),
	}
}

func init() {
	StorageManager = NewStorage()
	StorageManager.Register(NewResponseSendAdapter(adapterid.SendResponseAdapterID))
}
