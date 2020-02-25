// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulsar

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/pulse"
	"github.com/insolar/insolar/testutils"
)

func TestPulsar_Send(t *testing.T) {
	distMock := testutils.NewPulseDistributorMock(t)
	var pn pulse.Number = pulse.MinTimePulse

	distMock.DistributeMock.Set(func(ctx context.Context, p1 insolar.Pulse) {
		require.Equal(t, pn, p1.PulseNumber)
		require.NotNil(t, p1.Entropy)
	})

	pcs := platformpolicy.NewPlatformCryptographyScheme()
	crypto := testutils.NewCryptographyServiceMock(t)
	crypto.SignMock.Return(&insolar.Signature{}, nil)
	proc := platformpolicy.NewKeyProcessor()
	key, err := proc.GeneratePrivateKey()
	require.NoError(t, err)
	crypto.GetPublicKeyMock.Return(proc.ExtractPublicKey(key), nil)

	p := NewPulsar(
		configuration.NewPulsar(),
		crypto,
		pcs,
		platformpolicy.NewKeyProcessor(),
		distMock,
		&entropygenerator.StandardEntropyGenerator{},
	)

	err = p.Send(context.TODO(), pn)

	require.NoError(t, err)
	require.Equal(t, pn, p.LastPN())

	distMock.MinimockWait(1 * time.Minute)
}
