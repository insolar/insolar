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

package host

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAddress(t *testing.T) {
	addrStr := "127.0.0.1:31337"
	udpAddr, _ := net.ResolveUDPAddr("udp", addrStr)
	expectedAddr := &Address{*udpAddr}
	actualAddr, err := NewAddress(addrStr)

	assert.NoError(t, err)
	assert.Equal(t, expectedAddr, actualAddr)
	assert.True(t, actualAddr.Equal(*expectedAddr))
}

func TestAddress_Equal(t *testing.T) {
	addr1, _ := NewAddress("127.0.0.1:31337")
	addr2, _ := NewAddress("127.0.0.1:31337")
	addr3, _ := NewAddress("10.10.11.11:12345")

	assert.True(t, addr1.Equal(*addr2))
	assert.True(t, addr2.Equal(*addr1))
	assert.False(t, addr1.Equal(*addr3))
	assert.False(t, addr3.Equal(*addr1))
}
