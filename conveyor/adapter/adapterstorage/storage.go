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

package storage

import (
	"fmt"

	"github.com/insolar/insolar/conveyor/adapter"
	"github.com/insolar/insolar/conveyor/adapter/adapterid"
	"github.com/insolar/insolar/ledger/artifactmanager"
)

type Storage struct {
	adapters map[adapterid.ID]adapter.TaskSink
}

func (s *Storage) GetAdapterByID(id adapterid.ID) adapter.TaskSink {
	return s.adapters[id]
}

func (s *Storage) Register(adapter adapter.TaskSink) {
	id := adapter.GetAdapterID()
	_, ok := s.adapters[id]
	if ok {
		panic(fmt.Sprintf("[ Manager.Register ] adapter ID '%s' already exists", id.String()))
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

var Manager Storage

func NewEmptyStorage() Storage {
	return Storage{
		adapters: make(map[adapterid.ID]adapter.TaskSink),
	}
}

var processors []interface{}

type createProcessor func() adapter.Processor

func addAdapter(creator createProcessor, id adapterid.ID) {
	processor := creator()
	Manager.Register(adapter.NewAdapterWithQueue(processor, id))
	processors = append(processors, processor)
}

func init() {
	Manager = NewEmptyStorage()

	addAdapter(adapter.NewSendResponseProcessor, adapterid.SendResponse)
	addAdapter(artifactmanager.NewGetCodeProcessor, adapterid.GetCode)
}

// GetAllProcessors is used for component manager
func GetAllProcessors() []interface{} {
	return processors
}
