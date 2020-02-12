// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
