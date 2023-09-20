package resolver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExactResolver(t *testing.T) {
	localAddress := "127.0.0.1:12345"

	r := NewExactResolver()
	require.IsType(t, &exactResolver{}, r)
	realAddress, err := r.Resolve(localAddress)
	require.NoError(t, err)
	require.Equal(t, localAddress, realAddress)
}
