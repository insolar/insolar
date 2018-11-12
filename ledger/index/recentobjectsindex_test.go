package index

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRecentObjectsIndexInMemoryProvider(t *testing.T) {
	provider := NewRecentObjectsIndexInMemoryProvider()
	require.NotNil(t, provider.cache)
	require.NotNil(t, provider.cache.Fetched)
	require.NotNil(t, provider.cache.Updated)
}