package ginsider

import (
	"go/build"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/insolar/insolar/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	protocol := "unix"
	socket := os.TempDir() + "/" + testutils.RandomString() + ".sock"

	// start GoInsider
	gi := NewGoInsider("", protocol, socket)
	ref := addContractCode(t, gi)
	startGoInsider(t, gi, protocol, socket)

	caller := testutils.RandomRef()
	res := rpctypes.DownCallMethodResp{}
	req := rpctypes.DownCallMethodReq{
		Context:   &core.LogicCallContext{Caller: &caller},
		Code:      ref,
		Data:      goplugintestutils.CBORMarshal(t, []interface{}{}),
		Method:    "Check",
		Arguments: goplugintestutils.CBORMarshal(t, []interface{}{}),
	}

	client, err := rpc.Dial("unix", socket)
	require.NoError(t, err)

	err = client.Call("RPC.CallMethod", req, &res)
	require.NoError(t, err)

	unMarshaledResponse := goplugintestutils.CBORUnMarshal(t, res.Ret)

	assert.Equal(t, unMarshaledResponse, []interface{}{true, interface{}(nil)})
}

func addContractCode(t *testing.T, gi *GoInsider) core.RecordRef {
	ref := testutils.RandomRef()
	dir, err := build.Default.Import("github.com/insolar/insolar", "", build.FindOnly)
	require.NoError(t, err)

	pluginPath := filepath.Join(dir.Dir, "logicrunner", "goplugin", "ginsider", "healthcheck", "healthcheck.so")

	err = gi.registerCustomPlugin(ref, pluginPath)
	require.NoError(t, err)

	return ref
}

func startGoInsider(t *testing.T, gi *GoInsider, protocol string, socket string) {
	err := rpc.Register(&RPC{GI: gi})
	require.NoError(t, err)
	listener, err := net.Listen(protocol, socket)
	require.NoError(t, err)
	go rpc.Accept(listener)
}
