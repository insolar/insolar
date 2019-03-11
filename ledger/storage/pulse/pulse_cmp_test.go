package pulse_test

import (
	"crypto/rand"
	rand2 "math/rand"
	"testing"

	"github.com/google/gofuzz"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/db"
	"github.com/insolar/insolar/ledger/storage/pulse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPulse_Components(t *testing.T) {
	ctx := inslogger.TestContext(t)
	memStorage := pulse.NewStorageMem()
	dbStorage := pulse.NewStorageDB()
	dbStorage.DB = db.NewMockDB()

	var pulses []core.Pulse
	f := fuzz.New().Funcs(func(p *core.Pulse, c fuzz.Continue) {
		p.PulseNumber = gen.PulseNumber()
		_, err := rand.Read(p.Entropy[:])
		require.NoError(t, err)
	})
	f.NilChance(0).NumElements(10, 20)
	f.Fuzz(&pulses)

	var appended []core.Pulse
	latest := pulses[0]
	for i, p := range pulses {
		// Append appends if pulse is greater.
		memErr := memStorage.Append(ctx, p)
		dbErr := dbStorage.Append(ctx, p)
		if p.PulseNumber <= latest.PulseNumber && i > 0 {
			assert.Equal(t, pulse.ErrBadPulse, memErr)
			assert.Equal(t, pulse.ErrBadPulse, dbErr)
			continue
		}
		latest = p
		appended = append(appended, p)

		// Latest returns correct pulse.
		memLatest, memErr := memStorage.Latest(ctx)
		dbLatest, dbErr := memStorage.Latest(ctx)
		assert.NoError(t, memErr)
		assert.NoError(t, dbErr)
		assert.Equal(t, p, memLatest)
		assert.Equal(t, p, dbLatest)

		// ForPulse returns correct value
		memForPulse, memErr := memStorage.ForPulseNumber(ctx, p.PulseNumber)
		dbForPulse, dbErr := memStorage.ForPulseNumber(ctx, p.PulseNumber)
		assert.NoError(t, memErr)
		assert.NoError(t, dbErr)
		assert.Equal(t, p, memForPulse)
		assert.Equal(t, p, dbForPulse)
	}

	// Forwards returns correct value.
	{
		steps := rand2.Intn(len(appended))
		memPulse, memErr := memStorage.Forwards(ctx, appended[0].PulseNumber, steps)
		dbPulse, dbErr := memStorage.Forwards(ctx, appended[0].PulseNumber, steps)
		assert.NoError(t, memErr)
		assert.NoError(t, dbErr)
		assert.Equal(t, appended[steps], memPulse)
		assert.Equal(t, appended[steps], dbPulse)
	}
	// Backwards returns correct value.
	{
		steps := rand2.Intn(len(appended))
		memPulse, memErr := memStorage.Backwards(ctx, appended[len(appended)-1].PulseNumber, steps)
		dbPulse, dbErr := memStorage.Backwards(ctx, appended[len(appended)-1].PulseNumber, steps)
		assert.NoError(t, memErr)
		assert.NoError(t, dbErr)
		assert.Equal(t, appended[len(appended)-steps-1], memPulse)
		assert.Equal(t, appended[len(appended)-steps-1], dbPulse)
	}
}
