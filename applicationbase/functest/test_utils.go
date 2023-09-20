// +build functest

package functest

import (
	"testing"

	"github.com/insolar/insolar/applicationbase/testutils/testrequest"
	"github.com/insolar/insolar/applicationbase/testutils/testresponse"
	"github.com/insolar/insolar/insolar/secrets"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
)

type rpcStatusResponse struct {
	testresponse.RPCResponse
	Result testresponse.StatusResponse `json:"result"`
}

func getStatus(t testing.TB) testresponse.StatusResponse {
	body := testresponse.GetRPSResponseBody(t, launchnet.TestRPCUrl, testresponse.PostParams{
		"jsonrpc": "2.0",
		"method":  "node.getStatus",
		"id":      "1",
	})
	rpcStatusResponse := &rpcStatusResponse{}
	testresponse.UnmarshalRPCResponse(t, body, rpcStatusResponse)
	require.NotNil(t, rpcStatusResponse.Result)
	return rpcStatusResponse.Result
}

func newUserWithKeys() (*CommonUser, error) {
	privateKey, err := secrets.GeneratePrivateKeyEthereum()
	if err != nil {
		return nil, err
	}

	privKeyStr, err := secrets.ExportPrivateKeyPEM(privateKey)
	if err != nil {
		return nil, err
	}
	publicKey := secrets.ExtractPublicKey(privateKey)
	pubKeyStr, err := secrets.ExportPublicKeyPEM(publicKey)
	if err != nil {
		return nil, err
	}
	return &CommonUser{
		PrivKey: string(privKeyStr),
		PubKey:  string(pubKeyStr),
	}, nil
}

func createMember(t *testing.T) *CommonUser {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	member.Ref = Root.GetReference()

	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, member, "member.create", nil)
	require.NoError(t, err)
	ref, ok := result.(map[string]interface{})["reference"].(string)
	require.True(t, ok)
	member.Ref = ref
	return member
}
