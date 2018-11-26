package functest

import (
	"net/rpc"
	"testing"

	"github.com/insolar/insolar/testutils"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	client, err := rpc.Dial("tcp", "127.0.0.1:18181")
	require.NoError(t, err)
	caller := testutils.RandomRef()
	res := rpctypes.DownCallMethodResp{}
	req := rpctypes.DownCallMethodReq{
		Context:   &core.LogicCallContext{Caller: &caller},
		Code:      core.RecordRef{}.FromSlice(append(make([]byte, 63), 1)),
		Data:      goplugintestutils.CBORMarshal(t, []interface{}{}),
		Method:    "Check",
		Arguments: goplugintestutils.CBORMarshal(t, []interface{}{}),
	}

	err = client.Call("RPC.CallMethod", req, &res)
	require.NoError(t, err)

	unMarshaledResponse := goplugintestutils.CBORUnMarshal(t, res.Ret)

	assert.Equal(t, unMarshaledResponse, []interface{}{true, interface{}(nil)})
}
