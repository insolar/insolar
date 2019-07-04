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
	"testing"

	"github.com/insolar/insolar/testutils"

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
	
		func (r *One) Hello(str string) (map[string]interface{}, error) {
			return map[string]interface{}{"message": str}, nil
		}`

	// if we running this test with count we need to get unique names
	prototypeRef := uploadContract(t, testutils.RandStringBytes(16), contractCode)
	objectRef := callConstructor(t, prototypeRef)

	testParam := "test"
	methodResult := callMethod(t, objectRef, "Hello", testParam)
	require.Empty(t, methodResult.Error)
	require.Equal(t, testParam, methodResult.ExtractedReply["message"])
}
