package chaser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewChasingTimer(t *testing.T) {
	chasingDelay := time.Second
	ct := NewChasingTimer(chasingDelay)
	require.Equal(t, chasingDelay, ct.chasingDelay)
}

func TestIsEnabled(t *testing.T) {
	ct := NewChasingTimer(time.Second)
	require.True(t, ct.IsEnabled())

	ct = NewChasingTimer(0)
	require.False(t, ct.IsEnabled())

	ct = NewChasingTimer(-time.Second)
	require.False(t, ct.IsEnabled())
}

func TestWasStarted(t *testing.T) {
	ct := NewChasingTimer(time.Second)
	require.False(t, ct.WasStarted())

	ct.timer = time.NewTimer(time.Second)
	require.True(t, ct.WasStarted())
}

func TestRestartChase(t *testing.T) {
	ct := NewChasingTimer(-time.Second)
	ct.RestartChase()
	require.Nil(t, ct.timer)

	ct = NewChasingTimer(0)
	ct.RestartChase()
	require.Nil(t, ct.timer)

	ct = NewChasingTimer(time.Microsecond)
	ct.RestartChase()
	require.NotNil(t, ct.timer)

	ct.RestartChase()
	require.NotNil(t, ct.timer)
}

func TestChannel(t *testing.T) {
	ct := NewChasingTimer(0)
	require.Nil(t, ct.Channel())

	ct = NewChasingTimer(time.Microsecond)
	ct.RestartChase()
	require.NotNil(t, ct.Channel())
}
