package inssyslog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseParam(t *testing.T) {
	n, a := toNetworkAndAddress("127.0.0.1")
	require.Equal(t, DefaultSyslogNetwork, n)
	require.Equal(t, "127.0.0.1", a)

	n, a = toNetworkAndAddress("tcp:127.0.0.1")
	require.Equal(t, "tcp", n)
	require.Equal(t, "127.0.0.1", a)

	n, a = toNetworkAndAddress("tcp4:127.0.0.1")
	require.Equal(t, "tcp4", n)
	require.Equal(t, "127.0.0.1", a)

	n, a = toNetworkAndAddress("unix:127.0.0.1")
	require.Equal(t, "unix", n)
	require.Equal(t, "127.0.0.1", a)

	n, a = toNetworkAndAddress("127.0.0.1:555")
	require.Equal(t, DefaultSyslogNetwork, n)
	require.Equal(t, "127.0.0.1:555", a)

	n, a = toNetworkAndAddress("tcp:127.0.0.1:555")
	require.Equal(t, "tcp", n)
	require.Equal(t, "127.0.0.1:555", a)
}
