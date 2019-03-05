package matrix

import (
    "github.com/insolar/insolar/conveyor/generator/common"
    "github.com/insolar/insolar/conveyor/generator/state_machines/sample"
    
)

type matrix struct {
    matrix  []*common.StateMachine
}

type MachineType int
var Matrix matrix

const (
    TestStateMachine MachineType = iota + 1
)

func init() {
    Matrix := matrix{}
    Matrix.matrix = append(Matrix.matrix, nil,
        sample.SMRHTestStateMachineFactory(),
        )
}

func (m *matrix) GetStateMachineByType(mType MachineType) *common.StateMachine {
    return m.matrix[int(mType)]
}
