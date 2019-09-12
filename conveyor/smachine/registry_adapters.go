///
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
///

package smachine

import "sync"

func NewSharedRegistry() *SharedRegistry {
	return &SharedRegistry{}
}

type SharedRegistry struct {
	mutex sync.RWMutex

	adapters map[AdapterID]*adapterExecHelper
}

func (m *SharedRegistry) RegisterAdapter(adapterID AdapterID, adapterExecutor AdapterExecutor) ExecutionAdapter {
	if adapterID.IsEmpty() {
		panic("illegal value")
	}
	if adapterExecutor == nil {
		panic("illegal value")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.adapters == nil {
		m.adapters = make(map[AdapterID]*adapterExecHelper)
	}
	if m.adapters[adapterID] != nil {
		panic("duplicate adapter id: " + adapterID)
	}
	//adapterExecutor.RegisterOn(m.containerState)
	r := &adapterExecHelper{adapterID, adapterExecutor}
	m.adapters[adapterID] = r

	return r
}

func (m *SharedRegistry) GetAdapter(adapterID AdapterID) ExecutionAdapter {
	m.mutex.RLock()
	a := m.adapters[adapterID]
	m.mutex.RUnlock()
	return a
}

func (m *SharedRegistry) migrate(state SlotMachineState, migrationCount uint16) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for _, adapter := range m.adapters {
		adapter.executor.Migrate(state, migrationCount)
	}
}
