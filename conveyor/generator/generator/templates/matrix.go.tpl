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

package matrix

import (
    "github.com/insolar/insolar/conveyor/generator/common"
    {{range .Imports}}"{{.}}"
    {{end}}
)

type matrix struct {
    matrix  []*common.StateMachine
}

type MachineType int
var Matrix matrix

const (
    {{range $i, $m := .Machines}}{{$m.Name}}{{if (isNull $i)}} MachineType = iota + 1{{end}}{{end}}
)

func init() {
    Matrix := matrix{}
    Matrix.matrix = append(Matrix.matrix, nil,
    {{range .Machines}}{{.Package}}.SMRH{{.Name}}Factory(),
{{end}})
}

func (m *matrix) GetStateMachineByType(mType MachineType) *common.StateMachine {
    return m.matrix[int(mType)]
}
