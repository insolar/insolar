package member

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsEvicted(t *testing.T) {
	require.False(t, ModeRestrictedAnnouncement.IsEvicted())

	require.True(t, ModeEvictedGracefully.IsEvicted())

	require.True(t, ModeEvictedAsSuspected.IsEvicted())
}

func TestIsEvictedForcefully(t *testing.T) {
	require.False(t, ModeRestrictedAnnouncement.IsEvictedForcefully())

	require.False(t, ModeEvictedGracefully.IsEvictedForcefully())

	require.True(t, ModeEvictedAsSuspected.IsEvictedForcefully())
}

func TestIsRestricted(t *testing.T) {
	require.False(t, ModeSuspected.IsRestricted())

	require.True(t, ModeEvictedGracefully.IsRestricted())
}

func TestCanIntroduceJoiner(t *testing.T) {
	require.False(t, ModeRestrictedAnnouncement.CanIntroduceJoiner(false))

	require.False(t, ModePossibleFraudAndSuspected.CanIntroduceJoiner(false))

	require.False(t, ModePossibleFraud.CanIntroduceJoiner(true))

	require.True(t, ModePossibleFraud.CanIntroduceJoiner(false))
}

func TestIsMistrustful(t *testing.T) {
	require.False(t, ModeSuspected.IsMistrustful())

	require.True(t, ModePossibleFraudAndSuspected.IsMistrustful())
}

func TestIsSuspended(t *testing.T) {
	require.False(t, ModePossibleFraud.IsSuspended())

	require.True(t, ModePossibleFraudAndSuspected.IsSuspended())
}

func TestIsPowerless(t *testing.T) {
	require.False(t, ModePossibleFraud.IsPowerless())

	require.True(t, ModePossibleFraudAndSuspected.IsPowerless())

	require.True(t, ModeEvictedAsFraud.IsPowerless())
}

func TestAsUnit32(t *testing.T) {
	require.Equal(t, uint32(ModePossibleFraud), ModePossibleFraud.AsUnit32())

	require.Panics(t, func() { OpMode(1 << ModeBits).AsUnit32() })
}

func TestOpModeString(t *testing.T) {
	require.NotEmpty(t, ModeNormal.String())

	require.NotEmpty(t, ModeSuspected.String())

	require.NotEmpty(t, ModePossibleFraud.String())

	require.NotEmpty(t, ModePossibleFraudAndSuspected.String())

	require.NotEmpty(t, ModeRestrictedAnnouncement.String())

	require.NotEmpty(t, ModeEvictedGracefully.String())

	require.NotEmpty(t, ModeEvictedAsFraud.String())

	require.NotEmpty(t, ModeEvictedAsSuspected.String())

	require.NotEmpty(t, OpMode(1<<ModeBits).String())
}
