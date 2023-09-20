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

	r := *(ret.(*interface{}))
	rr := r.(map[string]interface{})
	require.IsType(t, map[string]interface{}{}, rr["result"], "spec.get returns resilt")
}
