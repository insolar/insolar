//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package machinesmanager

import (
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner/machinesmanager.MachinesManager -o ./ -s _mock.go -g

type MachinesManager interface {
	RegisterExecutor(t insolar.MachineType, e insolar.MachineLogicExecutor) error
	GetExecutor(t insolar.MachineType) (insolar.MachineLogicExecutor, error)
}

type mmanager struct {
	Executors [insolar.MachineTypesLastID]insolar.MachineLogicExecutor
}

func NewMachinesManager() MachinesManager {
	return &mmanager{}
}

// RegisterExecutor registers an executor for particular `MachineType`
func (m *mmanager) RegisterExecutor(t insolar.MachineType, e insolar.MachineLogicExecutor) error {
	m.Executors[int(t)] = e
	return nil
}

// GetExecutor returns an executor for the `MachineType` if it was registered (`RegisterExecutor`),
// returns error otherwise
func (m *mmanager) GetExecutor(t insolar.MachineType) (insolar.MachineLogicExecutor, error) {
	if res := m.Executors[int(t)]; res != nil {
		return res, nil
	}

	return nil, errors.Errorf("No executor registered for machine %d", int(t))
}
