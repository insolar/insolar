package mimic

import (
	"testing"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/genesis"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/platformpolicy"
)

func TestMimicLedger_Genesis(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)

	pcs := platformpolicy.NewPlatformCryptographyScheme()
	pulseStorage := pulse.NewStorageMem()
	dmm := drop.NewModifierMock(mc).SetMock.Return(nil)
	imm := object.NewIndexModifierMock(mc).SetIndexMock.Return(nil)
	rmm := object.NewRecordModifierMock(mc).SetMock.Return(nil)

	mimicLedgerInstance := NewMimicLedger(pcs, pulseStorage)
	mimicStorage := mimicLedgerInstance.(*mimicLedger).storage

	mimicClient := NewClient(mimicStorage)

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
		DiscoveryNodes: []insolar.DiscoveryNodeRegister{},
		ContractsConfig: insolar.GenesisContractsConfig{
			PKShardCount: 10,
			MAShardCount: 10,
		},
	}

	err := genesisObject.Start(ctx)
	assert.NoError(t, err)
}
