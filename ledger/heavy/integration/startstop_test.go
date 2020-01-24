// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	s, err := NewServer(context.Background(), cfg, heavyConfig, nil)
	assert.NoError(t, err)
	s.Stop()
}
