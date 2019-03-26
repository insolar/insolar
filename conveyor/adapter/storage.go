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
)

type storage struct {
	adapters map[ID]TaskSink
}

func (s *storage) GetAdapterByID(id uint32) TaskSink {
	return s.adapters[ID(id)]
}

func (s *storage) Register(adapter TaskSink) TaskSink {
	id := ID(adapter.GetAdapterID())
	_, ok := s.adapters[id]
	if ok {
		panic(fmt.Sprintf("[ Storage.Register ] adapter ID '%s' already exists", id.String()))
	}

	s.adapters[id] = adapter

	return adapter
}

func (s *storage) GetRegisteredAdapters() []interface{} {
	var result []interface{}

	for _, adapter := range s.adapters {
		result = append(result, adapter)
	}

	return result
}

var Storage storage

func NewStorage() storage {
	return storage{
		adapters: make(map[ID]TaskSink),
	}
}

func init() {
	Storage = NewStorage()
	Storage.Register(NewResponseSendAdapter(idType(SendResponseAdapterID)))
}
