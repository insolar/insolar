package pulse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEpoch_IsUnknown(t *testing.T) {
	require.True(t, InvalidPulseEpoch.IsUnknown())
	require.False(t, EphemeralPulseEpoch.IsUnknown())
	require.False(t, Epoch(MaxSpecialEpoch+1).IsUnknown())
	require.False(t, Epoch(MinTimePulse).IsUnknown())
	require.False(t, Epoch(MaxTimePulse).IsUnknown())
	require.False(t, Epoch(MaxTimePulse+1).IsUnknown())
}

func TestEpoch_IsEphemeral(t *testing.T) {
	require.False(t, InvalidPulseEpoch.IsEphemeral())
	require.True(t, EphemeralPulseEpoch.IsEphemeral())
	require.False(t, Epoch(MaxSpecialEpoch+1).IsEphemeral())
	require.False(t, Epoch(MinTimePulse).IsEphemeral())
	require.False(t, Epoch(MaxTimePulse).IsEphemeral())
	require.False(t, Epoch(MaxTimePulse+1).IsEphemeral())
}

func TestEpoch_IsArticulation(t *testing.T) {
	require.False(t, InvalidPulseEpoch.IsArticulation())
	require.False(t, EphemeralPulseEpoch.IsArticulation())
	require.True(t, ArticulationPulseEpoch.IsArticulation())
	require.False(t, Epoch(MaxSpecialEpoch+1).IsArticulation())
	require.False(t, Epoch(MinTimePulse).IsArticulation())
	require.False(t, Epoch(MaxTimePulse).IsArticulation())
	require.False(t, Epoch(MaxTimePulse+1).IsArticulation())
}

func TestEpoch_IsTimeEpoch(t *testing.T) {
	require.False(t, InvalidPulseEpoch.IsTimeEpoch())
	require.False(t, EphemeralPulseEpoch.IsTimeEpoch())
	require.False(t, ArticulationPulseEpoch.IsTimeEpoch())
	require.False(t, Epoch(MaxSpecialEpoch+1).IsTimeEpoch())
	require.True(t, Epoch(MinTimePulse).IsTimeEpoch())
	require.True(t, Epoch(MaxTimePulse).IsTimeEpoch())
	require.False(t, Epoch(MaxTimePulse+1).IsTimeEpoch())
}

func TestEpoch_IsValidEpoch(t *testing.T) {
	require.False(t, InvalidPulseEpoch.IsValidEpoch())
	require.True(t, EphemeralPulseEpoch.IsValidEpoch())
	require.True(t, ArticulationPulseEpoch.IsValidEpoch())
	require.True(t, Epoch(MaxSpecialEpoch-1).IsValidEpoch())
	require.True(t, Epoch(MaxSpecialEpoch).IsValidEpoch())
	require.False(t, Epoch(MaxSpecialEpoch+1).IsValidEpoch())
	require.False(t, Epoch(MinTimePulse-1).IsValidEpoch())
	require.True(t, Epoch(MinTimePulse).IsValidEpoch())
	require.True(t, Epoch(MaxTimePulse).IsValidEpoch())
	require.False(t, Epoch(MaxTimePulse+1).IsValidEpoch())
}

func TestEpoch_IsCompatible_Invalid(t *testing.T) {
	epoch := InvalidPulseEpoch
	require.False(t, epoch.IsCompatible(InvalidPulseEpoch))
	require.False(t, epoch.IsCompatible(EphemeralPulseEpoch))
	require.False(t, epoch.IsCompatible(ArticulationPulseEpoch))
	require.False(t, epoch.IsCompatible(MaxSpecialEpoch))
	require.False(t, epoch.IsCompatible(MaxSpecialEpoch+1))
	require.False(t, epoch.IsCompatible(MinTimePulse-1))
	require.False(t, epoch.IsCompatible(MinTimePulse))
	require.False(t, epoch.IsCompatible(MaxTimePulse))
	require.False(t, epoch.IsCompatible(MaxTimePulse+1))
}

func TestEpoch_IsCompatible_Ephemeral(t *testing.T) {
	epoch := EphemeralPulseEpoch
	require.False(t, epoch.IsCompatible(InvalidPulseEpoch))
	require.True(t, epoch.IsCompatible(EphemeralPulseEpoch))
	require.False(t, epoch.IsCompatible(ArticulationPulseEpoch))
	require.False(t, epoch.IsCompatible(MaxSpecialEpoch+1))
	require.False(t, epoch.IsCompatible(MinTimePulse-1))
	require.False(t, epoch.IsCompatible(MinTimePulse))
	require.False(t, epoch.IsCompatible(MaxTimePulse))
	require.False(t, epoch.IsCompatible(MaxTimePulse+1))
}

func TestEpoch_IsCompatible_Articulation(t *testing.T) {
	epoch := ArticulationPulseEpoch
	require.False(t, epoch.IsCompatible(InvalidPulseEpoch))
	require.False(t, epoch.IsCompatible(EphemeralPulseEpoch))
	require.True(t, epoch.IsCompatible(ArticulationPulseEpoch))
	require.False(t, epoch.IsCompatible(MaxSpecialEpoch+1))
	require.False(t, epoch.IsCompatible(MinTimePulse-1))
	require.True(t, epoch.IsCompatible(MinTimePulse))
	require.True(t, epoch.IsCompatible(MaxTimePulse))
	require.False(t, epoch.IsCompatible(MaxTimePulse+1))
}

func TestEpoch_IsCompatible_TimePulse(t *testing.T) {
	for _, epoch := range []Epoch{MinTimePulse, MinTimePulse + 1, MaxTimePulse - 1, MaxTimePulse} {
		require.False(t, epoch.IsCompatible(InvalidPulseEpoch), epoch)
		require.False(t, epoch.IsCompatible(EphemeralPulseEpoch), epoch)
		require.True(t, epoch.IsCompatible(ArticulationPulseEpoch), epoch)
		require.False(t, epoch.IsCompatible(MaxSpecialEpoch+1), epoch)
		require.False(t, epoch.IsCompatible(MinTimePulse-1), epoch)
		require.True(t, epoch.IsCompatible(MinTimePulse), epoch)
		require.True(t, epoch.IsCompatible(MaxTimePulse), epoch)
		require.False(t, epoch.IsCompatible(MaxTimePulse+1), epoch)
	}
}
