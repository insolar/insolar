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

	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/id"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticationRequest(t *testing.T) {
	senderAddress, _ := host.NewAddress("0.0.0.0:0")
	sender := host.NewHost(senderAddress)
	sender.ID, _ = id.NewID()
	receiverAddress, _ := host.NewAddress("0.0.0.0:0")
	receiver := host.NewHost(receiverAddress)
	receiver.ID, _ = id.NewID()
	hh := newMockHostHandler()

	err := AuthenticationRequest(hh, "begin", receiver.ID.String())
	assert.Error(t, err, "AuthenticationRequest: target for auth request not found")

	hh.FoundHost = receiver
	err = AuthenticationRequest(hh, "begin", receiver.ID.String())
	assert.NoError(t, err)
	err = AuthenticationRequest(hh, "revoke", receiver.ID.String())
	assert.NoError(t, err)
	err = AuthenticationRequest(hh, "unknown", receiver.ID.String())
	assert.Error(t, err, "AuthenticationRequest: unknown command")
}

func TestCheckOriginRequest(t *testing.T) {
	senderAddress, _ := host.NewAddress("0.0.0.0:0")
	sender := host.NewHost(senderAddress)
	sender.ID, _ = id.NewID()
	receiverAddress, _ := host.NewAddress("0.0.0.0:0")
	receiver := host.NewHost(receiverAddress)
	receiver.ID, _ = id.NewID()
	hh := newMockHostHandler()

	err := CheckOriginRequest(hh, receiver.ID.String())
	assert.Error(t, err, "CheckOriginRequest: target for relay request not found")

	hh.FoundHost = receiver
	err = CheckOriginRequest(hh, receiver.ID.String())
	assert.NoError(t, err)
}

func TestObtainIPRequest(t *testing.T) {
	senderAddress, _ := host.NewAddress("0.0.0.0:0")
	sender := host.NewHost(senderAddress)
	sender.ID, _ = id.NewID()
	receiverAddress, _ := host.NewAddress("0.0.0.0:0")
	receiver := host.NewHost(receiverAddress)
	receiver.ID, _ = id.NewID()
	hh := newMockHostHandler()

	err := ObtainIPRequest(hh, receiver.ID.String())
	assert.Error(t, err, "ObtainIPRequest: target for relay request not found")

	hh.FoundHost = receiver
	ObtainIPRequest(hh, receiver.ID.String())
}

func TestRelayRequest(t *testing.T) {
	senderAddress, _ := host.NewAddress("0.0.0.0:0")
	sender := host.NewHost(senderAddress)
	sender.ID, _ = id.NewID()
	receiverAddress, _ := host.NewAddress("0.0.0.0:0")
	receiver := host.NewHost(receiverAddress)
	receiver.ID, _ = id.NewID()
	hh := newMockHostHandler()

	err := RelayRequest(hh, "begin auth", receiver.ID.String())
	assert.Error(t, err, "RelayRequest: target for relay request not found")

	hh.FoundHost = receiver
	err = RelayRequest(hh, "begin auth", receiver.ID.String())
	assert.Error(t, err, "unknown command")
	err = RelayRequest(hh, "start", receiver.ID.String())
	assert.NoError(t, err)
	err = RelayRequest(hh, "stop", receiver.ID.String())
	assert.NoError(t, err)
}

func TestRelayOwnershipRequest(t *testing.T) {
	senderAddress, _ := host.NewAddress("0.0.0.0:0")
	sender := host.NewHost(senderAddress)
	sender.ID, _ = id.NewID()
	receiverAddress, _ := host.NewAddress("0.0.0.0:0")
	receiver := host.NewHost(receiverAddress)
	receiver.ID, _ = id.NewID()
	hh := newMockHostHandler()

	err := RelayOwnershipRequest(hh, receiver.ID.String())
	assert.Error(t, err, "RelayRequest: target for relay request not found")

	hh.FoundHost = receiver
	err = RelayOwnershipRequest(hh, receiver.ID.String())
	assert.NoError(t, err)
}
