package ginsider

import (
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"testing"

	"github.com/insolar/insolar/log"

	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	//protocol := "unix"
	//socket := os.TempDir() + "/" + testutils.RandomString() + ".sock"

	tmpDir, err := ioutil.TempDir("", "contractcache-")
	require.NoError(t, err, "failed to build tmp dir")
	//defer os.RemoveAll(tmpDir)

	currentPath, err := os.Getwd()
	require.NoError(t, err, "FUCK")
	//currentPath := filepath.Dir(currentFile)
	insgoccPath := currentPath + "/../../../bin/insgocc"
	contractPath := currentPath + "/healthcheck/healthcheck.go"

	log.Warnf(currentPath)
	log.Warnf(insgoccPath)
	log.Warnf(contractPath)
	log.Warnf(tmpDir)

	execResult, err := exec.Command(insgoccPath, "compile", "-o", tmpDir, contractPath).CombinedOutput()
	log.Warnf("%s", execResult)
	require.NoError(t, err, "failed to compile contract "+err.Error())

	// TODO check file exists
	//_, err = os.Stat(path)

	// start GoInsider
	//gi := NewGoInsider(tmpDir, protocol, socket)
	//
	//ref := core.RecordRef{}.FromSlice(append(make([]byte, 63), 1))
	//err = gi.AddPlugin(ref, tmpDir+"/healthcheck/main.so")
	//require.NoError(t, err, "failed to add plugin"+err.Error())
	//
	//startGoInsider(t, gi, protocol, socket)
	//
	//cmd := exec.Command(currentPath+"/../../../bin/healthcheck",
	//	"-a", socket,
	//	"-p", protocol)
	//
	//_, err = cmd.CombinedOutput()
	//
	//assert.Equal(t, "exit status 0", err.Error())
}

func startGoInsider(t *testing.T, gi *GoInsider, protocol string, socket string) {
	err := rpc.Register(&RPC{GI: gi})
	require.NoError(t, err, "can't register gi as rpc"+err.Error())
	listener, err := net.Listen(protocol, socket)
	require.NoError(t, err, "can't start listener"+err.Error())
	go rpc.Accept(listener)
}
