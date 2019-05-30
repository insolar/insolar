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

package functest

import (
	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSingleContractError(t *testing.T) {
	var contractCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type One struct {
	foundation.BaseContract
	Number int
}

func (c *One) Inc() (int, error) {
	c.Number++
	return c.Number, nil
}

func (c *One) Get() (int, error) {
	return c.Number, nil
}

func (c *One) Dec() (int, error) {
	c.Number--
	return c.Number, nil
}
`
	objectRef := callConstructor(t, uploadContract(t, contractCode))
	emptyArgs, _ := insolar.Serialize(make([]interface{}, 0))

	result := callMethod(t, objectRef, "Get", emptyArgs)
	require.Equal(t, float64(0), result)

	result = callMethod(t, objectRef, "Inc", emptyArgs)
	require.Equal(t, float64(1), result)

	result = callMethod(t, objectRef, "Get", emptyArgs)
	require.Equal(t, float64(1), result)

	result = callMethod(t, objectRef, "Dec", emptyArgs)
	require.Equal(t, float64(0), result)

	result = callMethod(t, objectRef, "Get", emptyArgs)
	require.Equal(t, float64(0), result)
}
