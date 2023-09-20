package platformpolicy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPlatformPolicy(t *testing.T) {
	pcs := NewPlatformCryptographyScheme()

	require.NotNil(t, pcs)

	pcsImpl := pcs.(*platformCryptographyScheme)
	require.NotNil(t, pcsImpl.hashProvider)
	require.NotNil(t, pcsImpl.signProvider)
}
