// +build functest

package functest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsgorundReload(t *testing.T) {
	_, err := signedRequest(&root, "DumpAllUsers")
	require.NoError(t, err)

	err = stopAllInsgorunds()
	// No need to stop test if this fails. All tests may stack
	assert.NoError(t, err)

	err = startAllInsgorunds()
	require.NoError(t, err)

	_, err = signedRequest(&root, "DumpAllUsers")
	require.NoError(t, err)
}
