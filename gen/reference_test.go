package gen

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/core"
)

func TestGen_JetID(t *testing.T) {
	for i := 0; i < 10000; i++ {
		jetID := JetID()
		recID := (*core.RecordID)(&jetID)
		require.Equalf(t,
			core.PulseNumberJet, recID.Pulse(),
			"pulse number should be core.PulseNumberJet. jet: %v", recID.DebugString())
		require.GreaterOrEqualf(t,
			uint8(core.JetMaximumDepth), jetID.Depth(),
			"jet depth %v should be less than maximum value %v. jet: %v",
			jetID.Depth(), core.JetMaximumDepth, jetID.DebugString(),
		)
	}
}
