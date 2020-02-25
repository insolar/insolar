// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package bootstrap

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
)

func createCryptographyService(t *testing.T) insolar.CryptographyService {
	keyProcessor := platformpolicy.NewKeyProcessor()
	privateKey, err := keyProcessor.GeneratePrivateKey()
	require.NoError(t, err)
	return cryptography.NewKeyBoundCryptographyService(privateKey)
}

func TestCreateAndVerifyPermit(t *testing.T) {
	origin, err := host.NewHostN("127.0.0.1:123", gen.Reference())
	assert.NoError(t, err)
	redirect, err := host.NewHostN("127.0.0.1:321", gen.Reference())
	assert.NoError(t, err)

	cryptographyService := createCryptographyService(t)

	permit, err := CreatePermit(origin.NodeID, redirect, []byte{}, cryptographyService)
	assert.NoError(t, err)
	assert.NotNil(t, permit)

	cert := testutils.NewCertificateMock(t)
	cert.GetDiscoveryNodesMock.Set(func() (r []insolar.DiscoveryNode) {
		pk, _ := cryptographyService.GetPublicKey()
		node := certificate.NewBootstrapNode(pk, "", origin.Address.String(), origin.NodeID.String(), insolar.StaticRoleVirtual.String())
		return []insolar.DiscoveryNode{node}
	})

	// validate
	err = ValidatePermit(permit, cert, createCryptographyService(t))
	assert.NoError(t, err)

	assert.Less(t, time.Now().Unix(), permit.Payload.ExpireTimestamp)
}
