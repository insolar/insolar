/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package connection

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConnectionFactory(t *testing.T) {
	pool := NewConnectionFactory()

	require.IsType(t, &udpConnectionFactory{}, pool)
}

func TestUdpConnectionFactory_Create(t *testing.T) {
	addrStr := "127.0.0.1:31337"
	pool := NewConnectionFactory()

	conn, err := pool.Create(addrStr)

	require.NoError(t, err)

	require.IsType(t, &net.UDPConn{}, conn)
	require.Equal(t, addrStr, conn.LocalAddr().String())

}
