package ginsider

import (
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils"

	"github.com/insolar/insolar/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var binaryPath string

func TestHealthCheck(t *testing.T) {
	protocol := "unix"
	socket := os.TempDir() + "/" + testutils.RandomString() + ".sock"

	tmpDir, err := ioutil.TempDir("", "contractcache-")
	require.NoError(t, err, "failed to build tmp dir")
	defer os.RemoveAll(tmpDir)

	currentPath, err := os.Getwd()
	require.NoError(t, err)

	insgoccPath := binaryPath + "/insgocc"
	healthcheckPath := binaryPath + "/healthcheck"
	contractPath := currentPath + "/healthcheck/healthcheck.go"
	if _, err = os.Stat(healthcheckPath); err != nil {
		t.Fatalf("Binary file %s is not found, please run make build", healthcheckPath)
	}

	pathToTmp, err := filepath.Rel(currentPath, tmpDir)

	execResult, err := exec.Command(insgoccPath, "compile", "-o", pathToTmp, contractPath).CombinedOutput()
	log.Warnf("%s", execResult)
	require.NoError(t, err, "failed to compile contract")

	// start GoInsider
	gi := NewGoInsider(tmpDir, protocol, socket)

	refString := "1111111111111111111111111111111111111111111111111111111111111112"
	ref := core.NewRefFromBase58(refString)
	err = gi.AddPlugin(ref, tmpDir+"/main.so")
	require.NoError(t, err, "failed to add plugin")

	startGoInsider(t, gi, protocol, socket)

	cmd := exec.Command(healthcheckPath,
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

func init() {
	var ok bool

	binaryPath, ok = os.LookupEnv("BIN_DIR")
	if !ok {
		wd, err := os.Getwd()
		binaryPath = filepath.Join(wd, "..", "..", "..", "bin")

		if err != nil {
			panic(err.Error())
		}
	}
}
