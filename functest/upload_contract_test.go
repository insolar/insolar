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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCallUploadedContract(t *testing.T) {
	t.Skip("this test fixing right now")
	contractCode := `
		package main
		import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
		type One struct {
			foundation.BaseContract
		}
		func New() (*One, error){
			return &One{}, nil}
	
		func (r *One) Hello(str string) (string, error) {
			return str, nil
		}`

	prototypeRef := uploadContract(t, contractCode)
	objectRef := callConstructor(t, prototypeRef)

	testParam := "test"
	args := append(make([]string, 0), testParam)
	argsSerialized, err := insolar.Serialize(args)
	require.NoError(t, err)

	methodResult := callMethod(t, objectRef, "Hello", argsSerialized)

	require.Equal(t, testParam, methodResult.(string))
}
