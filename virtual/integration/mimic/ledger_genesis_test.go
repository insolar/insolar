package mimic

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/genesis"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/platformpolicy"
)

const (
	LaunchnetRelativePath = "scripts/insolard/launchnet.sh"
	GenesisRelativePath   = "launchnet/configs/heavy_genesis.json"
)

func GenerateBootstrap(t testing.TB, skipBuild bool) (func(), string, error) {
	artifactsDir, err := ioutil.TempDir("", "mimic")
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to create temporary directory")
	}

	cleanupFunc := func() {
		err := os.RemoveAll(artifactsDir)
		if err != nil {
			t.Logf("[ Error ] Failed to cleanup temporary dir %s: %s", artifactsDir, err.Error())
		}
	}

	cmd := exec.Command(LaunchnetRelativePath, "-b")
	cmd.Dir = insolar.RootModuleDir()
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "INSOLAR_ARTIFACTS_DIR="+artifactsDir)
	if skipBuild {
		cmd.Env = append(cmd.Env, "SKIP_BUILD=1")
	}

	if output, err := cmd.CombinedOutput(); err != nil {
		cleanupFunc()

		t.Logf("[ Error ] Failed to execute bootstrap: %s", err.Error())
		t.Logf("[ Error ] Output of bootstrap is:")

		outputString := string(bytes.TrimSpace(output))
		for _, line := range strings.Split(outputString, "\n") {
			t.Logf("[ Error ] > %s", line)
		}

		return nil, "", errors.Wrapf(err, "Failed to execute bootstrap: %s", err.Error())
	}

	return cleanupFunc, artifactsDir, nil
}

func ReadGenesisContractsConfig(dirPath string) (*insolar.GenesisContractsConfig, error) {
	genesisConfigPath := path.Join(dirPath, GenesisRelativePath)

	fh, err := os.Open(genesisConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open genesis config for reading")
	}

	rv := insolar.GenesisHeavyConfig{}
	if err := json.NewDecoder(fh).Decode(&rv); err != nil {
		return nil, errors.Wrap(err, "failed to decode genesis config")
	}

	return &rv.ContractsConfig, nil
}

func TestMimicLedger_Genesis(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)

	pcs := platformpolicy.NewPlatformCryptographyScheme()
	pulseStorage := pulse.NewStorageMem()
	dmm := drop.NewModifierMock(mc).SetMock.Return(nil)
	imm := object.NewIndexModifierMock(mc).
		SetIndexMock.Return(nil).
		UpdateLastKnownPulseMock.Return(nil)
	rmm := object.NewRecordModifierMock(mc).SetMock.Return(nil)

	mimicLedgerInstance := NewMimicLedger(pcs, pulseStorage)
	mimicStorage := mimicLedgerInstance.(*mimicLedger).storage

	mimicClient := NewClient(mimicStorage)

	cleanup, bootstrapDir, err := GenerateBootstrap(t, true)
	require.NoError(t, err)
	defer cleanup()

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
		DiscoveryNodes:  []insolar.DiscoveryNodeRegister{},
		ContractsConfig: *genesisContractsConfig,
	}

	err = genesisObject.Start(ctx)
	assert.NoError(t, err)
}
