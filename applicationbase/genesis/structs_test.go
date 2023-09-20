package genesis

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	genesisIDHex  = "00010001ac000000000000000000000000000000000000000000000000000000"
	genesisRefHex = genesisIDHex + genesisIDHex
)

func TestGenesisRecordID(t *testing.T) {
	require.Equal(t, genesisIDHex, hex.EncodeToString(Record.ID().Bytes()), "genesis ID should always be the same")
}

func TestReference(t *testing.T) {
	require.Equal(t, genesisRefHex, hex.EncodeToString(Record.Ref().Bytes()), "genesisRef should always be the same")
}
