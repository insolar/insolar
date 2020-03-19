// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"encoding/json"
	"testing"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/insolar/insolar/applicationbase/testutils/testresponse"

	"github.com/stretchr/testify/require"
)

// MakeRequest - call rpc server and parse results.
func MakeRPCRequest(t testing.TB, m string, params interface{}) (interface{}, error) {
	pp := testresponse.PostParams{
		"jsonrpc": "2.0",
		"method":  m,
		"id":      1,
		"params":  params,
	}
	body := testresponse.GetRPSResponseBody(t, launchnet.TestRPCUrl, pp)
	res := new(interface{})
	err := json.Unmarshal(body, &res)
	return res, err
}

func TestSpecServiceGet(t *testing.T) {
	requester.SetVerbose(true)
	ret, err := MakeRPCRequest(t, "spec.get", map[string]interface{}{})
	require.NoError(t, err)

	r := ret.(*map[string]interface{})
	rr := *r
	require.IsType(t, "string", rr["openapi"], "right openapi")
}
