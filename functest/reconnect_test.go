package functest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInsgorundReload(t *testing.T) {
	_, err := signedRequest(&root, "DumpAllUsers")
	require.NoError(t, err)

	stopInsgorund()
	err = startInsgorund()
	require.NoError(t, err)

	_, err = signedRequest(&root, "DumpAllUsers")
	require.NoError(t, err)
}
