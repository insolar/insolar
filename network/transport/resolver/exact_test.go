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

package resolver

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewExactResolver(t *testing.T) {
	resolver := NewExactResolver()

	assert.IsType(t, &exactResolver{}, resolver)
}

func TestExactResolver_Resolve(t *testing.T) {
	strAddr := "127.0.0.1:31337"
	resolver := NewExactResolver()

	conn := &MockPacketConn{}
	conn.On("LocalAddr").Return(net.ResolveUDPAddr("udp", strAddr))

	addr, err := resolver.Resolve(conn)

	conn.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, strAddr, addr)
}
