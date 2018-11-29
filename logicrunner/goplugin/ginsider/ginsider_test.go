package ginsider

import (
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils"

	"github.com/insolar/insolar/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	protocol := "unix"
	socket := os.TempDir() + "/" + testutils.RandomString() + ".sock"

	tmpDir, err := ioutil.TempDir("", "contractcache-")
	require.NoError(t, err, "failed to build tmp dir")
	defer os.RemoveAll(tmpDir)

	currentPath, err := os.Getwd()
	require.NoError(t, err)

	insgoccPath := currentPath + "/../../../bin/insgocc"
	contractPath := currentPath + "/healthcheck/healthcheck.go"

	// TODO remove debug
	log.Warnf(currentPath)
	log.Warnf(insgoccPath)
	log.Warnf(contractPath)
	log.Warnf(tmpDir)

	execResult, err := exec.Command(insgoccPath, "compile", "-o", tmpDir, contractPath).CombinedOutput()
	log.Warnf("%s", execResult)
	require.NoError(t, err, "failed to compile contract")

	//start GoInsider
	gi := NewGoInsider(tmpDir, protocol, socket)

	refString := "1111111111111111111111111111111111111111111111111111111111111112"
	ref := core.NewRefFromBase58(refString)
	err = gi.AddPlugin(ref, tmpDir+"/main.so")
	require.NoError(t, err, "failed to add plugin")

	startGoInsider(t, gi, protocol, socket)

	cmd := exec.Command(currentPath+"/../../../bin/healthcheck",
		"-a", socket,
		"-p", protocol,
		"-r", refString)

	output, err := cmd.CombinedOutput()

	log.Warnf("%+v", output)

	assert.NoError(t, err)
}

func startGoInsider(t *testing.T, gi *GoInsider, protocol string, socket string) {
	err := rpc.Register(&RPC{GI: gi})
	require.NoError(t, err, "can't register gi as rpc")
	listener, err := net.Listen(protocol, socket)
	require.NoError(t, err, "can't start listener")
	go rpc.Accept(listener)
}
