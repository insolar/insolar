// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
