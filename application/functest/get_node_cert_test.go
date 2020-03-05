// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"encoding/json"
	"testing"

	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/insolar/insolar/applicationbase/testutils/testrequest"
	"github.com/insolar/insolar/applicationbase/testutils/testresponse"
	"github.com/insolar/insolar/certificate"

	"github.com/stretchr/testify/require"
)

func TestNodeCert(t *testing.T) {
	publicKey := testrequest.GenerateNodePublicKey(t)
	const testRole = "virtual"
	res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrl, &Root,
		"contract.registerNode", map[string]interface{}{"publicKey": publicKey, "role": testRole})
	require.NoError(t, err)

	body := getRPSResponseBody(t, launchnet.TestRPCUrl, testresponse.PostParams{
		"jsonrpc": "2.0",
		"method":  "cert.get",
		"id":      1,
		"params":  map[string]string{"ref": res.(string)},
	})

	cert := struct {
		Result struct {
			Cert certificate.Certificate
		}
	}{}

	err = json.Unmarshal(body, &cert)
	require.NoError(t, err)
	require.Equal(t, res.(string), cert.Result.Cert.Reference)
}
