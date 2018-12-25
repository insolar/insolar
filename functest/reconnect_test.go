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

	err = stopInsgorund()
	// no need to stop test if this fails
	assert.NoError(t, err)

	err = startInsgorund()
	require.NoError(t, err)

	_, err = signedRequest(&root, "DumpAllUsers")
	require.NoError(t, err)
}
