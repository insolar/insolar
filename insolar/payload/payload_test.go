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

package payload_test

import (
	"math/rand"
	"testing"

	"github.com/gogo/protobuf/proto"
	fuzz "github.com/google/gofuzz"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/stretchr/testify/assert"
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

func TestMarshalUnmarshal(t *testing.T) {
	type data struct {
		tp payload.Type
		pl payload.Payload
	}
	table := []data{
		{tp: payload.TypeError, pl: &payload.Error{}},
		{tp: payload.TypeID, pl: &payload.ID{}},
		{tp: payload.TypeState, pl: &payload.State{}},
		{tp: payload.TypeGetObject, pl: &payload.GetObject{}},
		{tp: payload.TypePassState, pl: &payload.PassState{}},
		{tp: payload.TypeIndex, pl: &payload.Index{}},
		{tp: payload.TypePass, pl: &payload.Pass{}},
		{tp: payload.TypeCode, pl: &payload.Code{}},
		{tp: payload.TypeGetCode, pl: &payload.GetCode{}},
		{tp: payload.TypeSetCode, pl: &payload.SetCode{}},
	}

	for _, d := range table {
		t.Run(d.tp.String(), func(t *testing.T) {
			fuzz.New().Fuzz(d.pl)
			encoded, err := payload.Marshal(d.pl)
			assert.NoError(t, err)
			decoded, err := payload.Unmarshal(encoded)
			assert.NoError(t, err)
			assert.Equal(t, d.pl, decoded)
		})
	}
}
