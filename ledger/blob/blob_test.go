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

package blob

import (
	"testing"

	"github.com/google/gofuzz"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClone(t *testing.T) {
	t.Parallel()

	jetID := gen.JetID()

	cases := []struct {
		name  string
		value []byte
	}{
		{
			name:  "rand value",
			value: slice(),
		},
		{
			name:  "nil value",
			value: nil,
		},
		{
			name:  "empty value",
			value: []byte{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			blob := Blob{
				JetID: jetID,
				Value: c.value,
			}

			clonedBlob := Clone(blob)

			assert.Equal(t, blob, clonedBlob)
			assert.False(t, &blob == &clonedBlob)
		})
	}

}

func TestStorageMemory_ForPN(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	ms := NewStorageMemory()
	pcs := testutils.NewPlatformCryptographyScheme()

	searchJetID := gen.JetID()
	searchPN := gen.PulseNumber()

	searchBlobs := map[insolar.ID]struct{}{}
	for i := 0; i < 5; i++ {
		b := Blob{}
		fuzz.New().NilChance(0).Fuzz(&b)
		b.JetID = searchJetID

		bID := object.CalculateIDForBlob(pcs, searchPN, MustEncode(&b))
		searchBlobs[*bID] = struct{}{}
		_ = ms.Set(ctx, *bID, b)
	}

	for i := 0; i < 500; i++ {
		b := Blob{}
		fuzz.New().NilChance(0).Fuzz(&b)
		bID := object.CalculateIDForBlob(pcs, gen.PulseNumber(), MustEncode(&b))
		_ = ms.Set(ctx, *bID, b)
	}

	res := ms.ForPulse(ctx, searchJetID, searchPN)

	require.Equal(t, len(searchBlobs), len(res))
	for _, b := range res {
		bID := object.CalculateIDForBlob(pcs, searchPN, MustEncode(&b))
		_, ok := searchBlobs[*bID]
		require.Equal(t, true, ok)
	}
}
