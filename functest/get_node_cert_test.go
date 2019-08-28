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

	"github.com/insolar/insolar/testutils/launchnet"

	"github.com/stretchr/testify/require"
)

func TestNodeCert(t *testing.T) {
	const TESTPUBLICKEY = "some_fancy_public_key"
	const testRole = "virtual"
	res, err := signedRequest(t, launchnet.TestRPCUrl, &launchnet.Root,
		"contract.registerNode", map[string]interface{}{"publicKey": TESTPUBLICKEY, "role": testRole})
	require.NoError(t, err)

	body := getRPSResponseBody(t, launchnet.TestRPCUrl, postParams{
		"jsonrpc": "2.0",
		"method":  "cert.get",
		"id":      1,
		"params":  map[string]string{"ref": res.(string)},
	})

	require.NotEqual(t, "", string(body))
}
