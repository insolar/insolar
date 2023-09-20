package payload_test

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/stretchr/testify/require"
)

func TestPolymorphProducesExpectedBinary(t *testing.T) {
	morph := rand.Uint32()
	pl := payload.Error{
		Polymorph: morph,
	}
	data, err := pl.Marshal()
	require.NoError(t, err)
	buf := proto.NewBuffer(data)

	_, err = buf.DecodeVarint()
	require.NoError(t, err)
	morph64, err := buf.DecodeVarint()
	require.NoError(t, err)

	require.Equal(t, morph, uint32(morph64))
}

func TestMarshalUnmarshalType(t *testing.T) {
	for _, expectedType := range payload.TypesMap {
		buf, err := payload.MarshalType(expectedType)
		require.NoError(t, err)

		tp, err := payload.UnmarshalType(buf)
		require.NoError(t, err)
		require.Equal(t, expectedType, tp)
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	for _, expectedType := range payload.TypesMap {
		if expectedType == payload.TypeUnknown {
			continue
		}

		typeBuf, err := payload.MarshalType(expectedType)
		require.NoError(t, err)
		pl, err := payload.Unmarshal(typeBuf)
		require.NoError(t, err, "Unmarshal() unknown type %s", expectedType.String())
		r := reflect.ValueOf(pl)
		f := reflect.Indirect(r).FieldByName("Polymorph")
		require.Equal(t, uint32(expectedType), uint32(f.Uint()), "Unmarshal() failed on type %s", expectedType.String())

		buf, err := payload.Marshal(pl)
		require.NoError(t, err, "Marshal() unknown type %s", expectedType.String())
		tp, err := payload.UnmarshalType(buf)
		require.NoError(t, err)
		require.Equal(t, expectedType, tp, "type mismatch between Marshal() and Unmarshal()")
	}
}
