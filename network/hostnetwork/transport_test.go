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

package hostnetwork

import (
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/stretchr/testify/assert"
)

const (
	TestType types.PacketType = 1024
)

type MockResolver struct {
	mapping map[core.RecordRef]string
}

func (m *MockResolver) Resolve(nodeID core.RecordRef) (string, error) {
	return m.mapping[nodeID], nil
}

func (m *MockResolver) AddToKnownHosts(h *host.Host) {
}

func mockConfiguration(nodeID string, address string) configuration.Configuration {
	result := configuration.Configuration{}
	result.Host.Transport = configuration.Transport{Protocol: "UTP", Address: address, BehindNAT: false}
	result.Node.Node = &configuration.Node{nodeID}
	return result
}

func TestNewInternalTransport(t *testing.T) {
	// broken address
	_, err := NewInternalTransport(mockConfiguration("123", "123"))
	assert.Error(t, err)
	address := "127.0.0.1:0"
	tp, err := NewInternalTransport(mockConfiguration("123", address))
	assert.NoError(t, err)
	defer tp.Disconnect()
	// assert that new address with correct port has been assigned
	assert.NotEqual(t, address, tp.PublicAddress())
}

func TestNewInternalTransport2(t *testing.T) {
	tp, err := NewInternalTransport(mockConfiguration("123", "127.0.0.1:0"))
	assert.NoError(t, err)
	go tp.Listen()
	// no assertion, check that Disconnect does not block
	defer func(t *testing.T) {
		tp.Disconnect()
		assert.True(t, true)
	}(t)
}
