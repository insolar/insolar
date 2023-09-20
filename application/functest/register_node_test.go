// +build functest

package functest

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/insolar/insolar/applicationbase/testutils/testrequest"
	"github.com/insolar/insolar/applicationbase/testutils/testresponse"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy"
)

var scheme = platformpolicy.NewPlatformCryptographyScheme()
var keyProcessor = platformpolicy.NewKeyProcessor()

func registerNodeSignedCall(t *testing.T, params map[string]interface{}) (string, error) {
	res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrl, &Root, "contract.registerNode", params)
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

func TestRegisterNodeVirtual(t *testing.T) {
	const testRole = "virtual"
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": testrequest.GenerateNodePublicKey(t), "role": testRole})
	require.NoError(t, err)

	require.NotNil(t, ref)
}

func TestRegisterNodeHeavyMaterial(t *testing.T) {
	const testRole = "heavy_material"
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": testrequest.GenerateNodePublicKey(t), "role": testRole})
	require.NoError(t, err)

	require.NotNil(t, ref)
}

func TestRegisterNodeLightMaterial(t *testing.T) {
	const testRole = "light_material"
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": testrequest.GenerateNodePublicKey(t), "role": testRole})
	require.NoError(t, err)

	require.NotNil(t, ref)
}

func TestRegisterNodeWithSamePK(t *testing.T) {
	const testRole = "light_material"
	testPublicKey := testrequest.GenerateNodePublicKey(t)
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": testPublicKey, "role": testRole})
	require.NoError(t, err)
	require.NotNil(t, ref)

	_, err = testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl, &Root,
		"contract.registerNode", map[string]interface{}{"publicKey": testPublicKey, "role": testRole})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "node already exist with this public key")
}

func TestRegisterNodeNotExistRole(t *testing.T) {
	_, err := testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl, &Root,
		"contract.registerNode", map[string]interface{}{"publicKey": testrequest.GenerateNodePublicKey(t), "role": "some_not_fancy_role"})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "role is not supported")
}

func TestRegisterNodeByNoRoot(t *testing.T) {
	member := createMember(t)
	const testRole = "virtual"
	_, err := testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl, member, "contract.registerNode",
		map[string]interface{}{"publicKey": testrequest.GenerateNodePublicKey(t), "role": testRole})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "only root member can register node")
}

func TestReceiveNodeCert(t *testing.T) {
	const testRole = "virtual"
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": testrequest.GenerateNodePublicKey(t), "role": testRole})
	require.NoError(t, err)

	body := getRPSResponseBody(t, launchnet.TestRPCUrl, testresponse.PostParams{
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
	require.NotEmpty(t, res.Result.Cert.BootstrapNodes)

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
