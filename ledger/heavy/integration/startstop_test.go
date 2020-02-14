// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build slowtest

package integration_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/application"
)

func TestStartStop(t *testing.T) {
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
	s, err := NewBadgerServer(context.Background(), cfg, heavyConfig, nil)
	assert.NoError(t, err)
	s.Stop()
}
