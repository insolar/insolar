package capacity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultPercent(t *testing.T) {
	require.Equal(t, 20, LevelMinimal.DefaultPercent())

	require.Panics(t, func() { LevelCount.DefaultPercent() })
}

func TestChooseInt(t *testing.T) {
	var options [LevelCount]int
	l := LevelMinimal
	options[l] = 5
	require.Equal(t, 5, l.ChooseInt(options))

	require.Panics(t, func() { LevelCount.ChooseInt(options) })
}
