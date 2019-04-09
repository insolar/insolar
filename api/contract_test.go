///
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
///

// +build functest

package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/testutils"
)

func TestUpload(t *testing.T) {
	cfg := configuration.NewAPIRunner()
	ar, _ := NewRunner(&cfg)

	amMock := artifacts.NewClientMock(t)
	genesisRef := testutils.RandomRef()
	amMock.GenesisRefMock.Return(&genesisRef)
	amMock.RegisterRequestFunc = func(p context.Context, p1 insolar.Reference, p2 insolar.Parcel) (r *insolar.ID, r1 error) {
		ID := testutils.RandomID()
		return &ID, nil
	}
	amMock.DeployCodeFunc = func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte, p4 insolar.MachineType) (r *insolar.ID, r1 error) {
		ID := testutils.RandomID()
		return &ID, nil
	}
	amMock.ActivatePrototypeMock.Return(nil, nil)

	ar.ArtifactManager = amMock

	service := NewContractService(ar)

	request := &http.Request{}

	contractCode := `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type One struct {
	foundation.BaseContract
}

func (r *One) Hello(str string) (string, error) {
	return r.GetPrototype().String() + str, nil
}
`
	params := &UploadArgs{
		Name: "test",
		Code: contractCode,
	}

	reply := &UploadReply{}

	err := service.Upload(request, params, reply)
	require.NoError(t, err)

}
