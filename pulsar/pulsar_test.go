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

package pulsar

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"net"
	"net/rpc"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/pulsar/pulsartestutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRpcClientFactoryWrapper struct {
	mock.Mock
}

func (mock *MockRpcClientFactoryWrapper) CreateWrapper() RpcClientWrapper {
	args := mock.Mock.Called()
	return args.Get(0).(RpcClientWrapper)
}

func TestNewPulsar_WithoutNeighbours(t *testing.T) {
	assertObj := assert.New(t)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	expectedPrivateKey, _ := ExportPrivateKey(privateKey)
	config := configuration.Pulsar{
		ConnectionType: "testType",
		ListenAddress:  "listedAddress",
		PrivateKey:     expectedPrivateKey,
	}
	actualConnectionType := ""
	actualAddress := ""

	mockListener := func(connectionType string, address string) (net.Listener, error) {
		actualConnectionType = connectionType
		actualAddress = address
		return &pulsartestutil.MockListener{}, nil
	}
	storage := &pulsartestutil.MockStorage{}
	storage.On("GetLastPulse", mock.Anything).Return(&core.Pulse{PulseNumber: 123}, nil)
	clientFactory := &MockRpcClientFactoryWrapper{}
	clientFactory.On("CreateWrapper").Return(&pulsartestutil.MockRpcClientWrapper{})

	result, err := NewPulsar(config,
		storage,
		clientFactory,
		pulsartestutil.MockEntropyGenerator{},
		mockListener)

	assertObj.NoError(err)
	parsedKey, _ := ImportPrivateKey(expectedPrivateKey)
	assertObj.Equal(parsedKey, result.PrivateKey)
	assertObj.Equal("testType", actualConnectionType)
	assertObj.Equal("listedAddress", actualAddress)
	assertObj.IsType(result.Sock, &pulsartestutil.MockListener{})
	assertObj.NotNil(result.PrivateKey)
	clientFactory.AssertNumberOfCalls(t, "CreateWrapper", 0)
}

func TestNewPulsar_WithNeighbours(t *testing.T) {
	assertObj := assert.New(t)

	firstPrivateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	firstExpectedKey, _ := ExportPublicKey(&firstPrivateKey.PublicKey)

	secondPrivateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	secondExpectedKey, _ := ExportPublicKey(&secondPrivateKey.PublicKey)

	expectedPrivateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	parsedExpectedPrivateKey, _ := ExportPrivateKey(expectedPrivateKey)
	config := configuration.Pulsar{
		ConnectionType: "testType",
		ListenAddress:  "listedAddress",
		PrivateKey:     parsedExpectedPrivateKey,
		ListOfNeighbours: []*configuration.PulsarNodeAddress{
			{ConnectionType: "tcp", Address: "first", PublicKey: firstExpectedKey},
			{ConnectionType: "pct", Address: "second", PublicKey: secondExpectedKey},
		},
	}
	storage := &pulsartestutil.MockStorage{}
	storage.On("GetLastPulse", mock.Anything).Return(&core.Pulse{PulseNumber: 123}, nil)
	clientFactory := &MockRpcClientFactoryWrapper{}
	clientFactory.On("CreateWrapper").Return(&pulsartestutil.MockRpcClientWrapper{})

	result, err := NewPulsar(config, storage, clientFactory,
		pulsartestutil.MockEntropyGenerator{}, func(connectionType string, address string) (net.Listener, error) {
			return &pulsartestutil.MockListener{}, nil
		})

	assertObj.NoError(err)
	assertObj.Equal(2, len(result.Neighbours))
	assertObj.Equal("tcp", result.Neighbours[firstExpectedKey].ConnectionType.String())
	assertObj.Equal("pct", result.Neighbours[secondExpectedKey].ConnectionType.String())
	clientFactory.AssertNumberOfCalls(t, "CreateWrapper", 2)
}

func TestPulsar_EstablishConnection_IsInitialised(t *testing.T) {
	pulsar := &Pulsar{Neighbours: map[string]*Neighbour{}}

	mockClientWrapper := &pulsartestutil.MockRpcClientWrapper{}
	mockClientWrapper.On("IsInitialised").Return(true)

	firstPrivateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	expectedNeighbourKey, _ := ExportPublicKey(&firstPrivateKey.PublicKey)
	expectedNeighbour := &Neighbour{OutgoingClient: mockClientWrapper}
	pulsar.Neighbours[expectedNeighbourKey] = expectedNeighbour

	err := pulsar.EstablishConnection(expectedNeighbourKey)

	assert.NoError(t, err)
	mockClientWrapper.AssertNotCalled(t, "Lock")
}

func TestPulsar_EstablishConnection_IsNotInitialised_ProblemsCreateConnection(t *testing.T) {
	pulsar := &Pulsar{Neighbours: map[string]*Neighbour{}}

	mockClientWrapper := &pulsartestutil.MockRpcClientWrapper{}
	mockClientWrapper.On("IsInitialised").Return(false)
	mockClientWrapper.On("CreateConnection", mock.Anything, mock.Anything).Return(errors.New("test reasons"))
	mockClientWrapper.On("Lock")
	mockClientWrapper.On("Unlock")

	firstPrivateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	expectedNeighbourKey, _ := ExportPublicKey(&firstPrivateKey.PublicKey)
	expectedNeighbour := &Neighbour{OutgoingClient: mockClientWrapper}
	pulsar.Neighbours[expectedNeighbourKey] = expectedNeighbour

	err := pulsar.EstablishConnection(expectedNeighbourKey)

	assert.Error(t, err, "test reasons")
	mockClientWrapper.AssertNumberOfCalls(t, "Lock", 1)
	mockClientWrapper.AssertNumberOfCalls(t, "Unlock", 1)
}

func TestPulsar_EstablishConnection_IsNotInitialised_ProblemsWithRequest(t *testing.T) {
	mainPrivateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	pulsar := &Pulsar{
		Neighbours:       map[string]*Neighbour{},
		PrivateKey:       mainPrivateKey,
		EntropyGenerator: pulsartestutil.MockEntropyGenerator{},
	}

	mockClientWrapper := &pulsartestutil.CustomRpcWrapperMock{}

	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{Error: errors.New("oops, request is broken")}
	replyChan := &rpc.Call{Done: done}
	mockClientWrapper.Done = replyChan

	firstPrivateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	expectedNeighbourKey, _ := ExportPublicKey(&firstPrivateKey.PublicKey)
	expectedNeighbour := &Neighbour{OutgoingClient: mockClientWrapper}
	pulsar.Neighbours[expectedNeighbourKey] = expectedNeighbour

	err = pulsar.EstablishConnection(expectedNeighbourKey)

	assert.Error(t, err, "oops, request is broken")
}

func TestPulsar_EstablishConnection_IsNotInitialised_SignatureFailed(t *testing.T) {
	mainPrivateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	pulsar := &Pulsar{
		Neighbours:       map[string]*Neighbour{},
		PrivateKey:       mainPrivateKey,
		EntropyGenerator: pulsartestutil.MockEntropyGenerator{},
	}

	mockClientWrapper := &pulsartestutil.CustomRpcWrapperMock{}

	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{Reply: &Payload{Body: HandshakePayload{}}}
	replyChan := &rpc.Call{Done: done}
	mockClientWrapper.Done = replyChan

	firstPrivateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	expectedNeighbourKey, _ := ExportPublicKey(&firstPrivateKey.PublicKey)
	expectedNeighbour := &Neighbour{OutgoingClient: mockClientWrapper}
	pulsar.Neighbours[expectedNeighbourKey] = expectedNeighbour

	err = pulsar.EstablishConnection(expectedNeighbourKey)

	assert.Error(t, err, "Signature check failed")
}

func TestPulsar_EstablishConnection_IsNotInitialised_Success(t *testing.T) {
	mainPrivateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	pulsar := &Pulsar{
		Neighbours:       map[string]*Neighbour{},
		PrivateKey:       mainPrivateKey,
		EntropyGenerator: pulsartestutil.MockEntropyGenerator{},
	}
	firstPrivateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	expectedNeighbourKey, _ := ExportPublicKey(&firstPrivateKey.PublicKey)
	payload := Payload{Body: HandshakePayload{Entropy: pulsartestutil.MockEntropy}}
	sign, err := singData(firstPrivateKey, payload.Body)
	payload.Signature = sign
	payload.PublicKey = expectedNeighbourKey

	mockClientWrapper := &pulsartestutil.CustomRpcWrapperMock{}

	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{Reply: &payload}
	replyChan := &rpc.Call{Done: done}
	mockClientWrapper.Done = replyChan

	expectedNeighbour := &Neighbour{OutgoingClient: mockClientWrapper}
	pulsar.Neighbours[expectedNeighbourKey] = expectedNeighbour

	err = pulsar.EstablishConnection(expectedNeighbourKey)

	assert.NoError(t, err)
}

func TestPulsar_stateSwitchedToVerifying_OnePulsar(t *testing.T) {
	mainPrivateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	pulsar := &Pulsar{
		Neighbours:       map[string]*Neighbour{},
		PrivateKey:       mainPrivateKey,
		GeneratedEntropy: pulsartestutil.MockEntropy,
		OwnedBftRow:      map[string]*BftCell{"test": &BftCell{}},
	}

	pulsar.stateSwitchedToVerifying()

}
