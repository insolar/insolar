//
// Copyright 2019 Insolar Technologies GmbH
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
//
// +build slowtest

package mimic

import (
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/application/genesis"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/platformpolicy"
)

func TestMimicLedger_Genesis(t *testing.T) {
	cleanup, bootstrapDir, err := GenerateBootstrap(true)
	require.NoError(t, err)
	defer cleanup()

	senderMock := bus.NewSenderMock(t)

	t.Run("WithMocks", func(t *testing.T) {
		// t.Parallel()
		ctx := inslogger.TestContext(t)

		mc := minimock.NewController(t)

		pcs := platformpolicy.NewPlatformCryptographyScheme()
		pulseStorage := pulse.NewStorageMem()
		dmm := drop.NewModifierMock(mc).SetMock.Return(nil)
		imm := object.NewIndexModifierMock(mc).
			SetIndexMock.Return(nil).
			UpdateLastKnownPulseMock.Return(nil)
		rmm := object.NewRecordModifierMock(mc).SetMock.Return(nil)

		mimicLedgerInstance := NewMimicLedger(ctx, pcs, pulseStorage, pulseStorage, senderMock)
		mimicStorage := mimicLedgerInstance.(*mimicLedger).storage

		mimicClient := NewClient(mimicStorage)

		genesisContractsConfig, err := ReadGenesisContractsConfig(bootstrapDir)
		require.NoError(t, err)

		genesisObject := genesis.Genesis{
			ArtifactManager: mimicClient,
			BaseRecord: &genesis.BaseRecord{
				DB:             mimicStorage,
				DropModifier:   dmm,
				PulseAppender:  pulseStorage,
				PulseAccessor:  pulseStorage,
				RecordModifier: rmm,
				IndexModifier:  imm,
			},
			DiscoveryNodes:  []application.DiscoveryNodeRegister{},
			ContractsConfig: *genesisContractsConfig,
		}

		err = genesisObject.Start(ctx)
		assert.NoError(t, err)
	})

	t.Run("WithoutMocks", func(t *testing.T) {
		// t.Parallel()
		ctx := inslogger.TestContext(t)

		pcs := platformpolicy.NewPlatformCryptographyScheme()
		pulseStorage := pulse.NewStorageMem()

		mimicLedgerInstance := NewMimicLedger(ctx, pcs, pulseStorage, pulseStorage, senderMock)

		err := mimicLedgerInstance.LoadGenesis(ctx, bootstrapDir)
		assert.NoError(t, err)
	})
}
