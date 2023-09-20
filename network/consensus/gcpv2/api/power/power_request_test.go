package power

import (
	"testing"

	"github.com/insolar/insolar/network/consensus/common/capacity"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/stretchr/testify/require"
)

func TestNewRequestByLevel(t *testing.T) {
	require.Equal(t, -Request(capacity.LevelMinimal)-1, NewRequestByLevel(capacity.LevelMinimal))
}

func TestNewRequest(t *testing.T) {
	require.Equal(t, Request(1)+1, NewRequest(member.Power(1)))
}

func TestAsCapacityLevel(t *testing.T) {
	b, l := Request(-1).AsCapacityLevel()
	require.True(t, b)
	require.Equal(t, capacity.Level(0), l)

	b, l = Request(1).AsCapacityLevel()
	require.False(t, b)

	r := Request(-2)
	require.Equal(t, capacity.Level(r), l)

	b, l = Request(0).AsCapacityLevel()
	require.False(t, b)

	r = Request(-1)
	require.Equal(t, capacity.Level(r), l)
}

func TestAsMemberPower(t *testing.T) {
	b, l := Request(1).AsMemberPower()
	require.True(t, b)
	require.Zero(t, l)

	b, l = Request(-1).AsMemberPower()
	require.False(t, b)

	r := Request(-2)
	require.Equal(t, member.Power(r), l)

	b, l = Request(0).AsMemberPower()
	require.False(t, b)

	r = Request(-1)
	require.Equal(t, member.Power(r), l)
}

func TestIsEmpty(t *testing.T) {
	require.True(t, EmptyRequest.IsEmpty())

	require.False(t, Request(1).IsEmpty())
}

func TestUpdate(t *testing.T) {
	pws := member.PowerSet([...]member.Power{10, 20, 30, 40})
	pwBase := member.Power(1)
	pw := pwBase

	require.True(t, Request(-1).Update(&pw, pws))

	require.Zero(t, pw)

	pw = pwBase
	require.True(t, Request(-2).Update(&pw, pws))

	require.Equal(t, member.Power(10), pw)

	pw = pwBase
	require.True(t, Request(10).Update(&pw, pws))

	require.Equal(t, member.Power(10), pw)

	pw = pwBase
	require.True(t, Request(100).Update(&pw, pws))

	require.Equal(t, member.Power(40), pw)

	pw = pwBase
	require.False(t, Request(0).Update(&pw, pws))

	require.Equal(t, member.Power(1), pw)
}
