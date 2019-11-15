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

package goplugintestutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

var contractOneCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	recursive "github.com/insolar/insolar/application/proxy/recursive_call_one"
)
type One struct {
	foundation.BaseContract
}

func New() (*One, error) {
	return &One{}, nil
}

var INSATTR_Recursive_API = true
func (r *One) Recursive() (error) {
	remoteSelf := recursive.GetObject(r.GetReference())
	err := remoteSelf.Recursive()
	return err
}

`

func TestContractsBuilder_Build(t *testing.T) {

	insgocc, err := BuildPreprocessor()
	assert.NoError(t, err)

	am := artifacts.NewClientMock(t)
	am.RegisterIncomingRequestMock.Set(func(ctx context.Context, request *record.IncomingRequest) (rp1 *payload.RequestInfo, err error) {
		rp1 = &payload.RequestInfo{RequestID: gen.ID()}
		return
	})
	am.DeployCodeMock.Set(func(ctx context.Context, code []byte, machineType insolar.MachineType) (ip1 *insolar.ID, err error) {
		assert.Equal(t, insolar.MachineTypeGoPlugin, machineType)
		id := gen.ID()
		return &id, nil
	})
	am.ActivatePrototypeMock.Set(func(ctx context.Context, request insolar.Reference, parent insolar.Reference, code insolar.Reference, memory []byte) (err error) {
		return nil
	})

	pa := pulse.NewAccessorMock(t)
	pa.LatestMock.Set(func(ctx context.Context) (p1 insolar.Pulse, err error) {
		return *insolar.GenesisPulse, nil
	})

	j := jet.NewCoordinatorMock(t)
	j.MeMock.Set(func() (r1 insolar.Reference) {
		return gen.Reference()
	})

	cb := NewContractBuilder(insgocc, am, pa, j)
	defer cb.Clean()

	contractMap := make(map[string]string)
	contractMap["recursive_call_one"] = contractOneCode

	err = cb.Build(context.Background(), contractMap)
	assert.NoError(t, err)

	reference := cb.Prototypes["recursive_call_one"]
	PrototypeRef := reference.String()
	assert.NotEmpty(t, PrototypeRef)
}
