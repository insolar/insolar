// +build functest

package functest

import (
	"testing"

	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var insgorundCleaner func()

func startInsgorund() (err error) {
	// It starts on ports of "virtual" node
	insgorundCleaner, err = goplugintestutils.StartInsgorund(insgorundPath, "tcp", "127.0.0.1:38181", "tcp", "127.0.0.1:38182")
	if err != nil {
		return errors.Wrap(err, "[ startInsgorund ] could't wait for insolard to start completely: ")
	}
	return nil
}

func stopInsgorund() error {
	if insgorundCleaner != nil {
		insgorundCleaner()
	}
	return nil
}

func TestInsgorundReload(t *testing.T) {
	_, err := signedRequest(&root, "DumpAllUsers")
	require.NoError(t, err)

	stopInsgorund()
	err = startInsgorund()
	require.NoError(t, err)

	_, err = signedRequest(&root, "DumpAllUsers")
	require.NoError(t, err)
}
