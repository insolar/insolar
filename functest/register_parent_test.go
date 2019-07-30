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

// +build functest

package functest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContractWithEmbeddedConstructor(t *testing.T) {
	var contractOneCode = `
package main

import ("github.com/insolar/insolar/logicrunner/builtin/foundation")

type One struct {
	foundation.BaseContract
	Number int
}

func New() (*One, error) {
	return &One{Number: 0}, nil
}

func NewWithNumber(num int) (*One, error) {
	return &One{Number: num}, nil
}

var INSATTR_Get_API = true

func (c *One) Get() (int, error) {
	return c.Number, nil
}

var INSATTR_DoNothing_API = true

func (r *One) DoNothing() (error) {
	return nil
}
`

	var contractTwoCode = `
package main

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	one "github.com/insolar/insolar/application/proxy/parent_one"
)

type Two struct {
	foundation.BaseContract
	Number int
	OneRef insolar.Reference
}


func New() (*Two, error) {
	return &Two{Number: 10, OneRef: insolar.Reference{}}, nil
}

func NewWithOne(oneNumber int) (*Two, error) {
	holder := one.NewWithNumber(oneNumber)

	objOne, err := holder.AsChild(foundation.GetRequestReference())

	if err != nil {
		return nil, err
	}

	return &Two{Number: oneNumber, OneRef: objOne.GetReference() }, nil
}

var INSATTR_DoNothing_API = true

func (r *Two) DoNothing() (error) {
	return nil
}

var INSATTR_Get_API = true

func (c * Two) Get() (int, error) {
	return c.Number, nil
}
`
	codeOneRef := uploadContractOnce(t, "parent_one", contractOneCode)
	codeTwoRef := uploadContractOnce(t, "parent_two", contractTwoCode)

	objectOneRef := callConstructor(t, codeOneRef, "New")
	objectTwoRef := callConstructor(t, codeTwoRef, "NewWithOne", 10)

	resp := callMethod(t, objectOneRef, "Get")
	require.Empty(t, resp.Error)
	require.Equal(t, float64(0), resp.ExtractedReply)

	resp = callMethod(t, objectTwoRef, "Get")
	require.Empty(t, resp.Error)
	require.Equal(t, float64(10), resp.ExtractedReply)
}
