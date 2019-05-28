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
	//contractCode := `package main\r\n\r\nimport \"github.com\/insolar\/insolar\/logicrunner\/goplugin\/foundation\"\r\n\r\ntype One struct {\r\nfoundation.BaseContract\r\n}\r\n\r\nfunc New() (*One, error){\r\nreturn &One{}, nil\r\n}\r\n\r\n\r\nfunc (r *One) Hello(str string) (string, error) {\r\nreturn r.GetPrototype().String() + str, nil\r\n}`

	contractCode := `
		package main
		import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
		type One struct {
			foundation.BaseContract
		}
		func New() (*One, error){
			return &One{}, nil}
	
		func (r *One) Hello(str string) (string, error) {
			return r.GetPrototype().String() + str, nil
		}`

	body := getRPSResponseBody(t, postParams{
		"jsonrpc": "2.0",
		"method":  "contract.Upload",
		"id":      "",
		"params": map[string]string{
			"name": "test",
			"code": contractCode,
		},
	})
	require.NotEmpty(t, body)

	res := struct {
		Version string          `json:"jsonrpc"`
		ID      string          `json:"id"`
		Result  api.UploadReply `json:"result"`
	}{}



	err := json.Unmarshal(body, &res)
	require.NoError(t, err)

	contractRef := insolar.Reference{}.FromSlice([]byte(res.Result.PrototypeRef))
	emptyRef := make([]byte, insolar.RecordRefSize)
	require.NotEqual(t, insolar.Reference{}.FromSlice(emptyRef), contractRef)
}
