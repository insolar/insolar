// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build slowtest

package integration_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
)

func Test_FinalizePulse(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultHeavyConfig()
	defer os.RemoveAll(cfg.Ledger.Storage.DataDirectory)
	testPk := "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEwDcgWZ1SbG+nbiXZkmYUZEfk2nkk\n1PEmEWoj4g6DLEkdaQVorOkqlloEz1zXclQaAE1S8i3F7OFNrNxLkm34ow==\n-----END PUBLIC KEY-----\n"
	heavyConfig := application.GenesisHeavyConfig{
		ContractsConfig: application.GenesisContractsConfig{
			PKShardCount:            10,
			MAShardCount:            10,
			MigrationAddresses:      make([][]string, 10),
			RootPublicKey:           testPk,
			FeePublicKey:            testPk,
			MigrationAdminPublicKey: testPk,
		},
	}
	s, err := NewBadgerServer(ctx, cfg, heavyConfig, nil)
	assert.NoError(t, err)
	defer s.Stop()

	s.SetPulse(ctx)
	s.SetPulse(ctx)

	targetPulse := s.Pulse() - PulseStep

	_, done := s.Send(ctx, &payload.GotHotConfirmation{
		JetID: insolar.ZeroJetID,
		Pulse: targetPulse,
		Split: false,
	})
	done()

	require.Equal(t, insolar.GenesisPulse.PulseNumber, s.JetKeeper.TopSyncPulse())

	d := drop.Drop{
		Pulse: targetPulse,
		JetID: insolar.ZeroJetID,
		Split: false,
	}

	_, done = s.Send(ctx, &payload.Replication{
		JetID: insolar.ZeroJetID,
		Pulse: targetPulse,
		Drop:  d,
	})
	done()

	numIterations := 20
	for s.JetKeeper.TopSyncPulse() == insolar.GenesisPulse.PulseNumber && numIterations > 0 {
		time.Sleep(500 * time.Millisecond)
		numIterations--
	}
	require.Equal(t, targetPulse, s.JetKeeper.TopSyncPulse())
}
