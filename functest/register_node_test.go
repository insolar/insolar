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
	"encoding/json"
	"testing"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils/launchnet"
	"github.com/stretchr/testify/require"
)

var scheme = platformpolicy.NewPlatformCryptographyScheme()
var keyProcessor = platformpolicy.NewKeyProcessor()

const TESTPUBLICKEY = "some_fancy_public_key"

func registerNodeSignedCall(t *testing.T, params map[string]interface{}) (string, error) {
	res, err := signedRequest(t, launchnet.TestRPCUrl, &launchnet.Root, "contract.registerNode", params)
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

func TestRegisterNodeVirtual(t *testing.T) {
	const testRole = "virtual"
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": TESTPUBLICKEY, "role": testRole})
	require.NoError(t, err)

	require.NotNil(t, ref)
}

func TestRegisterNodeHeavyMaterial(t *testing.T) {
	const testRole = "heavy_material"
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": TESTPUBLICKEY, "role": testRole})
	require.NoError(t, err)

	require.NotNil(t, ref)
}

func TestRegisterNodeLightMaterial(t *testing.T) {
	const testRole = "light_material"
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": TESTPUBLICKEY, "role": testRole})
	require.NoError(t, err)

	require.NotNil(t, ref)
}

func TestRegisterNodeNotExistRole(t *testing.T) {
	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl, &launchnet.Root,
		"contract.registerNode", map[string]interface{}{"publicKey": TESTPUBLICKEY, "role": "some_not_fancy_role"})
	require.Error(t, err)
	require.Contains(t, err.Error(),
		"role is not supported: some_not_fancy_role")
}

func TestRegisterNodeByNoRoot(t *testing.T) {
	member := createMember(t)
	const testRole = "virtual"
	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl, member, "contract.registerNode",
		map[string]interface{}{"publicKey": TESTPUBLICKEY, "role": testRole})
	require.Error(t, err)
	require.Contains(t, err.Error(), "only root member can register node")
}

func TestReceiveNodeCert(t *testing.T) {
	const testRole = "virtual"
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": TESTPUBLICKEY, "role": testRole})
	require.NoError(t, err)

	body := getRPSResponseBody(t, launchnet.TestRPCUrl, postParams{
		"jsonrpc": "2.0",
		"method":  "cert.get",
		"id":      "",
		"params":  map[string]string{"ref": ref},
	})

	res := struct {
		Result struct {
			Cert certificate.Certificate
		}
	}{}

	err = json.Unmarshal(body, &res)
	require.NoError(t, err)

	networkPart := res.Result.Cert.SerializeNetworkPart()
	nodePart := res.Result.Cert.SerializeNodePart()

	for _, discoveryNode := range res.Result.Cert.BootstrapNodes {
		pKey, err := keyProcessor.ImportPublicKeyPEM([]byte(discoveryNode.PublicKey))
		require.NoError(t, err)

		t.Run("Verify network sign for "+discoveryNode.Host, func(t *testing.T) {
			verified := scheme.DataVerifier(pKey, scheme.IntegrityHasher()).Verify(insolar.SignatureFromBytes(discoveryNode.NetworkSign), networkPart)
			require.True(t, verified)
		})
		t.Run("Verify node sign for "+discoveryNode.Host, func(t *testing.T) {
			verified := scheme.DataVerifier(pKey, scheme.IntegrityHasher()).Verify(insolar.SignatureFromBytes(discoveryNode.NodeSign), nodePart)
			require.True(t, verified)
		})
	}
}
