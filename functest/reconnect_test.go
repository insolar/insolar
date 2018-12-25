// +build functest

package functest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInsgorundReload(t *testing.T) {
	_, err := signedRequest(&root, "DumpAllUsers")
	require.NoError(t, err)

	stopAllInsgorunds()
	err = startAllInsgorunds()
	require.NoError(t, err)

	_, err = signedRequest(&root, "DumpAllUsers")
	require.NoError(t, err)
}
