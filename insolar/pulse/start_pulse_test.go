package pulse

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
)

func ReadPulses(t testing.TB, pulser StartPulse) func() {
	return func() {
		pulser.PulseNumber()
	}
}

func TestStartPulseRace(t *testing.T) {
	ctx := inslogger.TestContext(t)
	startPulse := NewStartPulse()
	syncTest := &testutils.SyncT{T: t}
	for i := 0; i < 10; i++ {
		go ReadPulses(syncTest, startPulse)()
	}
	startPulse.SetStartPulse(ctx, insolar.Pulse{PulseNumber: gen.PulseNumber()})
}
