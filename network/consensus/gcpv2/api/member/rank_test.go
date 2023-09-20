package member

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPower(t *testing.T) {
	require.Equal(t, Rank(1).GetPower(), Power(1))
}

func TestGetIndex(t *testing.T) {
	require.Equal(t, JoinerIndex, JoinerRank.GetIndex())

	require.Zero(t, Rank((1<<8)-1).GetIndex())

	require.Equal(t, Index(1), Rank(1<<8).GetIndex())
}

func TestGetTotalCount(t *testing.T) {
	require.Zero(t, Rank((1<<18)-1).GetTotalCount())

	require.Equal(t, uint16(1), Rank(1<<18).GetTotalCount())
}

func TestIsJoiner(t *testing.T) {
	require.False(t, Rank(1).IsJoiner())

	require.True(t, JoinerRank.IsJoiner())
}

func TestString(t *testing.T) {
	joiner := "{joiner}"

	require.Equal(t, JoinerRank.String(), joiner)
	require.NotEqual(t, Rank(1).String(), joiner)
	require.Equal(t, joiner, JoinerRank.String())
	require.NotEqual(t, joiner, Rank(1).String())
}

func TestNewMembershipRank(t *testing.T) {
	require.Panics(t, func() { NewMembershipRank(ModeNormal, Power(1), 1, 1) })

	require.Panics(t, func() { NewMembershipRank(ModeNormal, Power(1), 0x03FF+1, 1) })

	require.Panics(t, func() { NewMembershipRank(ModeNormal, Power(1), 1, 0x03FF+1) })

	require.Panics(t, func() { NewMembershipRank(ModeNormal, Power(1), 1, 0x03FF+1) })

	require.Equal(t, Rank(0x80101), NewMembershipRank(ModeNormal, Power(1), 1, 2))
}

func TestAsMembershipRank(t *testing.T) {
	fr := FullRank{}
	fr.OpMode = ModeNormal
	fr.Power = 1
	fr.TotalIndex = 1
	require.Equal(t, Rank(0x80101), fr.AsMembershipRank(2))

	require.Panics(t, func() { fr.AsMembershipRank(0) })
}
