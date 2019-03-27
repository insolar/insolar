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

package adapterstorage

import (
	"fmt"

	"github.com/insolar/insolar/conveyor/adapter"
	"github.com/insolar/insolar/conveyor/adapter/adapterid"
	"github.com/insolar/insolar/ledger/artifactmanager"
)

// Storage contains all adapters
type Storage struct {
	adapters map[adapterid.ID]adapter.TaskSink
}

// GetAdapterByID returns adapter by id
func (s *Storage) GetAdapterByID(id adapterid.ID) adapter.TaskSink {
	fmt.Println("GetAdapterByID", id)
	return s.adapters[id]
}

// GetAdapterByID registers adapters
func (s *Storage) Register(adapter adapter.TaskSink) {
	id := adapter.GetAdapterID()
	_, ok := s.adapters[id]
	if ok {
		panic(fmt.Sprintf("[ Manager.Register ] adapter ID '%s' already exists", id.String()))
	}

	s.adapters[id] = adapter
}

// Manager is global instance of Storage
var Manager Storage

// NewEmptyStorage creates new storage without any adapters
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
	addAdapter(adapter.NewNodeStateProcessor, adapterid.NodeState)
}

// GetAllProcessors is used for component manager
func GetAllProcessors() []interface{} {
	return processors
}
