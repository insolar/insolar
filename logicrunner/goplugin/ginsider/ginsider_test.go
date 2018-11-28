package ginsider

import (
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"testing"

	"github.com/insolar/insolar/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	protocol := "unix"
	socket := os.TempDir() + "/" + testutils.RandomString() + ".sock"

	tmpDir, err := ioutil.TempDir("", "contractcache-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// start GoInsider
	gi := NewGoInsider(tmpDir, protocol, socket)
	startGoInsider(t, gi, protocol, socket)

	currentPath, err := os.Getwd()
	require.NoError(t, err)

	cmd := exec.Command(currentPath+"/../../../bin/healthcheck",
		"-c", currentPath+"/healthcheck/healthcheck.go",
		"-d", tmpDir,
		"-a", socket,
		"-p", protocol)

	_, err = cmd.CombinedOutput()

	assert.Equal(t, "exit status 0", err.Error())
}

func startGoInsider(t *testing.T, gi *GoInsider, protocol string, socket string) {
	err := rpc.Register(&RPC{GI: gi})
	require.NoError(t, err)
	listener, err := net.Listen(protocol, socket)
	require.NoError(t, err)
	go rpc.Accept(listener)
}
