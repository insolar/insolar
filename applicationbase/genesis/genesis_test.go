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
