package merkle

import (
	"testing"

	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/require"
)

func TestFromList(t *testing.T) {
	cs := platformpolicy.NewPlatformCryptographyScheme()

	mt, err := treeFromHashList([][]byte{
		cs.IntegrityHasher().Hash([]byte("123")),
		cs.IntegrityHasher().Hash([]byte("456")),
	}, cs.IntegrityHasher())
	require.NoError(t, err)

	root := mt.Root()

	require.NotNil(t, root)
}
