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

package object

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecordStorage_TruncateHead(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)

	dbMock, err := store.NewBadgerDB(tmpdir)
	defer dbMock.Stop(ctx)
	require.NoError(t, err)

	recordStore := NewRecordDB(dbMock)

	numElements := 100

	// it's used for writing pulses in random order to db
	indexes := make([]int, numElements)
	for i := 0; i < numElements; i++ {
		indexes[i] = i
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(indexes), func(i, j int) { indexes[i], indexes[j] = indexes[j], indexes[i] })

	startPulseNumber := insolar.GenesisPulse.PulseNumber
	ids := make([]insolar.ID, numElements)
	for _, idx := range indexes {
		pulse := startPulseNumber + insolar.PulseNumber(idx)
		ids[idx] = *insolar.NewID(pulse, []byte(testutils.RandomString()))

		err := recordStore.Set(ctx, record.Material{JetID: *insolar.NewJetID(uint8(idx), nil), ID: ids[idx]})
		require.NoError(t, err)
	}

	for i := 0; i < numElements; i++ {
		_, err := recordStore.ForID(ctx, ids[i])
		require.NoError(t, err)
	}

	numLeftElements := numElements / 2
	err = recordStore.TruncateHead(ctx, startPulseNumber+insolar.PulseNumber(numLeftElements))
	require.NoError(t, err)

	for i := 0; i < numLeftElements; i++ {
		_, err := recordStore.ForID(ctx, ids[i])
		require.NoError(t, err)
	}

	for i := numElements - 1; i >= numLeftElements; i-- {
		_, err := recordStore.ForID(ctx, ids[i])
		require.EqualError(t, err, ErrNotFound.Error())
	}
}
