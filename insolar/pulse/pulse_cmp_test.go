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

package pulse_test

import (
	"crypto/rand"
	"io/ioutil"
	rand2 "math/rand"
	"os"
	"testing"

	"github.com/dgraph-io/badger"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func BadgerDefaultOptions(dir string) badger.Options {
	ops := badger.DefaultOptions(dir)
	ops.CompactL0OnClose = false
	ops.SyncWrites = false

	return ops
}

func TestPulse_Components(t *testing.T) {
	ctx := inslogger.TestContext(t)

	memStorage := pulse.NewStorageMem()

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(ctx)
	dbStorage := pulse.NewDB(db)

	var pulses []insolar.Pulse
	f := fuzz.New().Funcs(func(p *insolar.Pulse, c fuzz.Continue) {
		p.PulseNumber = gen.PulseNumber()
		_, err := rand.Read(p.Entropy[:])
		require.NoError(t, err)
	})
	f.NilChance(0).NumElements(10, 20)
	f.Fuzz(&pulses)

	var appended []insolar.Pulse
	latest := pulses[0]
	for i, p := range pulses {
		// Append appends if Pulse is greater.
		memErr := memStorage.Append(ctx, p)
		dbErr := dbStorage.Append(ctx, p)
		if p.PulseNumber <= latest.PulseNumber && i > 0 {
			assert.Equal(t, pulse.ErrBadPulse, memErr)
			assert.Equal(t, pulse.ErrBadPulse, dbErr)
			continue
		}
		latest = p
		appended = append(appended, p)

		// Latest returns correct Pulse.
		memLatest, memErr := memStorage.Latest(ctx)
		dbLatest, dbErr := dbStorage.Latest(ctx)
		assert.NoError(t, memErr)
		assert.NoError(t, dbErr)
		assert.Equal(t, p, memLatest)
		assert.Equal(t, p, dbLatest)

		// ForPulse returns correct value
		memForPulse, memErr := memStorage.ForPulseNumber(ctx, p.PulseNumber)
		dbForPulse, dbErr := dbStorage.ForPulseNumber(ctx, p.PulseNumber)
		assert.NoError(t, memErr)
		assert.NoError(t, dbErr)
		assert.Equal(t, p, memForPulse)
		assert.Equal(t, p, dbForPulse)
	}

	// Forwards returns correct value.
	{
		steps := rand2.Intn(len(appended))
		memPulse, memErr := memStorage.Forwards(ctx, appended[0].PulseNumber, steps)
		dbPulse, dbErr := dbStorage.Forwards(ctx, appended[0].PulseNumber, steps)
		assert.NoError(t, memErr)
		assert.NoError(t, dbErr)
		assert.Equal(t, appended[steps], memPulse)
		assert.Equal(t, appended[steps], dbPulse)
	}
	// Backwards returns correct value.
	{
		steps := rand2.Intn(len(appended))
		memPulse, memErr := memStorage.Backwards(ctx, appended[len(appended)-1].PulseNumber, steps)
		dbPulse, dbErr := dbStorage.Backwards(ctx, appended[len(appended)-1].PulseNumber, steps)
		assert.NoError(t, memErr)
		assert.NoError(t, dbErr)
		assert.Equal(t, appended[len(appended)-steps-1], memPulse)
		assert.Equal(t, appended[len(appended)-steps-1], dbPulse)
	}
}
