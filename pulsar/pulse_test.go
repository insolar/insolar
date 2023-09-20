package pulsar

import (
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/pulsar/pulsartestutils"
	"github.com/stretchr/testify/require"
)

func TestNewPulse(t *testing.T) {
	generator := &pulsartestutils.MockEntropyGenerator{}
	previousPulse := insolar.PulseNumber(876)
	expectedPulse := previousPulse + insolar.PulseNumber(configuration.NewPulsar().NumberDelta)

	result := NewPulse(configuration.NewPulsar().NumberDelta, previousPulse, generator)

	require.Equal(t, result.Entropy[:], pulsartestutils.MockEntropy[:])
	require.Equal(t, result.PulseNumber, insolar.PulseNumber(expectedPulse))
}
