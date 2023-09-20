package genesis

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/record"
)

func TestGenesisRecordMarshalUnmarshal(t *testing.T) {
	genIn := record.Genesis{
		Hash: Record,
	}

	virtGenIn := record.Wrap(&genIn)

	data, err := virtGenIn.Marshal()
	require.NoError(t, err)

	require.Equal(t, "aa0604a20101ac", hex.EncodeToString(data),
		"genesis binary representation always the same")

	virtGenOut := record.Virtual{}
	err = virtGenOut.Unmarshal(data)
	require.NoError(t, err, "genesis record unmarshal w/o error")

	genOut := record.Unwrap(&virtGenOut)

	require.Equal(t, &genIn, genOut, "marshal-unmarshal-marshal gives the same struct")

	data2, err := virtGenOut.Marshal()
	require.NoError(t, err)
	require.Equal(t, data, data2, "marshal-unmarshal-marshal gives the same binary result")
}
