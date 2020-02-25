// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
