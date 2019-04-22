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

package genesis

import (
	"encoding/hex"
	"testing"

	"github.com/insolar/insolar/ledger/object"
	"github.com/stretchr/testify/require"
)

func TestGenesisRecordEncodeDecode(t *testing.T) {
	recIn := &object.GenesisRecord{}
	recIn = recIn.Init()

	b := object.EncodeVirtual(recIn)
	require.Equal(t, "00000064a16d5669727475616c5265636f726441ac", hex.EncodeToString(b),
		"genesis binary representation always the same")

	recOut, err := object.DecodeVirtual(b)
	require.NoError(t, err, "genesis record decode w/o error")
	require.Equal(t, recIn, recOut, "encode-decode-encode gives the same struct")

	b2 := object.EncodeVirtual(recOut)
	require.Equal(t, b, b2, "encode-decode-encode gives the same binary result")
}
