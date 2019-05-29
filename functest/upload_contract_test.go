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
	"encoding/json"
	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/insolar"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCallUploadedContract(t *testing.T) {
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

	objectBody := getRPSResponseBody(t, postParams{
		"jsonrpc": "2.0",
		"method":  "contract.CallConstructor",
		"id":      "",
		"params": map[string]string{
			"PrototypeRefString": prototypeRef.String(),
		},
	})
	require.NotEmpty(t, objectBody)

	callConstructorRes := struct {
		Version string                   `json:"jsonrpc"`
		ID      string                   `json:"id"`
		Result  api.CallConstructorReply `json:"result"`
	}{}

	err := json.Unmarshal(objectBody, &callConstructorRes)
	require.NoError(t, err)

	objectRef, err := insolar.NewReferenceFromBase58(callConstructorRes.Result.ObjectRef)
	require.NoError(t, err)

	require.NotEqual(t, insolar.Reference{}.FromSlice(make([]byte, insolar.RecordRefSize)), objectRef)

	testParam := "test"
	args := make([]string, 0)
	args = append(args, testParam)
	argsSerialized, err := insolar.Serialize(args)
	require.NoError(t, err)

	callMethodBody := getRPSResponseBody(t, postParams{
		"jsonrpc": "2.0",
		"method":  "contract.CallMethod",
		"id":      "",
		"params": map[string]interface{}{
			"ObjectRefString": callConstructorRes.Result.ObjectRef,
			"Method":          "Hello",
			"MethodArgs":      argsSerialized,
		},
	})
	require.NotEmpty(t, callMethodBody)

	callRes := struct {
		Version string              `json:"jsonrpc"`
		ID      string              `json:"id"`
		Result  api.CallMethodReply `json:"result"`
	}{}

	err = json.Unmarshal(callMethodBody, &callRes)
	require.NoError(t, err)

	methodResult := callRes.Result.ExtractedReply.(string)
	require.Equal(t, testParam, methodResult)
}
