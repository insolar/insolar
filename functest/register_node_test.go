// +build functest

/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package functest

import (
	"encoding/json"
	"testing"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/require"
)

var scheme = platformpolicy.NewPlatformCryptographyScheme()
var keyProcessor = platformpolicy.NewKeyProcessor()

const TESTPUBLICKEY = "some_fancy_public_key"

func registerNodeSignedCall(params ...interface{}) (string, error) {
	res, err := signedRequest(&root, "RegisterNode", params...)
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

func TestRegisterNodeVirtual(t *testing.T) {
	const testRole = "virtual"
	ref, err := registerNodeSignedCall(TESTPUBLICKEY, testRole)
	require.NoError(t, err)

	require.NotNil(t, ref)
}

func TestRegisterNodeHeavyMaterial(t *testing.T) {
	const testRole = "heavy_material"
	ref, err := registerNodeSignedCall(TESTPUBLICKEY, testRole)
	require.NoError(t, err)

	require.NotNil(t, ref)
}

func TestRegisterNodeLightMaterial(t *testing.T) {
	const testRole = "light_material"
	ref, err := registerNodeSignedCall(TESTPUBLICKEY, testRole)
	require.NoError(t, err)

	require.NotNil(t, ref)
}

func TestRegisterNodeNotExistRole(t *testing.T) {
	_, err := registerNodeSignedCall(TESTPUBLICKEY, "some_not_fancy_role")
	require.Contains(t, err.Error(),
		"[ RegisterNode ] Can't save as child: [ SaveAsChild ] on calling main API: executer error: "+
			"problem with API call: Can't call constructor NewNodeRecord: Role is not supported: some_not_fancy_role")
}

func TestRegisterNodeByNoRoot(t *testing.T) {
	member := createMember(t, "Member1")
	const testRole = "virtual"
	_, err := signedRequest(member, "RegisterNode", TESTPUBLICKEY, testRole)
	require.Contains(t, err.Error(), "[ RegisterNode ] Only Root member can register node")
}

func TestReceiveNodeCert(t *testing.T) {
	const testRole = "virtual"
	ref, err := registerNodeSignedCall(TESTPUBLICKEY, testRole)
	require.NoError(t, err)

	body := getRPSResponseBody(t, postParams{
		"jsonrpc": "2.0",
		"method":  "cert.Get",
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
			verified := scheme.Verifier(pKey).Verify(core.SignatureFromBytes(discoveryNode.NetworkSign), networkPart)
			require.True(t, verified)
		})
		t.Run("Verify node sign for "+discoveryNode.Host, func(t *testing.T) {
			verified := scheme.Verifier(pKey).Verify(core.SignatureFromBytes(discoveryNode.NodeSign), nodePart)
			require.True(t, verified)
		})
	}
}
