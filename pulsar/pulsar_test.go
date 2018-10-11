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
	"bytes"
	"crypto/ecdsa"
	"net"
	"net/rpc"
	"os"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	ecdsa_helper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/pulsar/pulsartestutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func capture(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

type MockRpcClientFactoryWrapper struct {
	mock.Mock
}

func (mock *MockRpcClientFactoryWrapper) CreateWrapper() RPCClientWrapper {
	args := mock.Mock.Called()
	return args.Get(0).(RPCClientWrapper)
}

func TestNewPulsar_WithoutNeighbours(t *testing.T) {
	assertObj := assert.New(t)
	privateKey, _ := ecdsa_helper.GeneratePrivateKey()
	expectedPrivateKey, _ := ecdsa_helper.ExportPrivateKey(privateKey)
	config := configuration.Pulsar{
		ConnectionType:      "testType",
		MainListenerAddress: "listedAddress",
		PrivateKey:          expectedPrivateKey,
	}
	actualConnectionType := ""
	actualAddress := ""

	mockListener := func(connectionType string, address string) (net.Listener, error) {
		actualConnectionType = connectionType
		actualAddress = address
		return &pulsartestutil.MockListener{}, nil
	}
	storage := &pulsartestutil.MockPulsarStorage{}
	storage.On("GetLastPulse", mock.Anything).Return(&core.Pulse{PulseNumber: 123}, nil)
	clientFactory := &MockRpcClientFactoryWrapper{}
	clientFactory.On("CreateWrapper").Return(&pulsartestutil.MockRPCClientWrapper{})

	result, err := NewPulsar(config,
		storage,
		clientFactory,
		pulsartestutil.MockEntropyGenerator{},
		nil,
		mockListener)

	assertObj.NoError(err)
	parsedKey, _ := ecdsa_helper.ImportPrivateKey(expectedPrivateKey)
	assertObj.Equal(parsedKey, result.PrivateKey)
	assertObj.Equal("testType", actualConnectionType)
	assertObj.Equal("listedAddress", actualAddress)
	assertObj.IsType(result.Sock, &pulsartestutil.MockListener{})
	assertObj.NotNil(result.PrivateKey)
	clientFactory.AssertNumberOfCalls(t, "CreateWrapper", 0)
}

func TestNewPulsar_WithNeighbours(t *testing.T) {
	assertObj := assert.New(t)

	firstPrivateKey, _ := ecdsa_helper.GeneratePrivateKey()
	firstExpectedKey, _ := ecdsa_helper.ExportPublicKey(&firstPrivateKey.PublicKey)

	secondPrivateKey, _ := ecdsa_helper.GeneratePrivateKey()
	secondExpectedKey, _ := ecdsa_helper.ExportPublicKey(&secondPrivateKey.PublicKey)

	expectedPrivateKey, _ := ecdsa_helper.GeneratePrivateKey()
	parsedExpectedPrivateKey, _ := ecdsa_helper.ExportPrivateKey(expectedPrivateKey)
	config := configuration.Pulsar{
		ConnectionType:      "testType",
		MainListenerAddress: "listedAddress",
		PrivateKey:          parsedExpectedPrivateKey,
		Neighbours: []configuration.PulsarNodeAddress{
			{ConnectionType: "tcp", Address: "first", PublicKey: firstExpectedKey},
			{ConnectionType: "pct", Address: "second", PublicKey: secondExpectedKey},
		},
	}
	storage := &pulsartestutil.MockPulsarStorage{}
	storage.On("GetLastPulse", mock.Anything).Return(&core.Pulse{PulseNumber: 123}, nil)
	clientFactory := &MockRpcClientFactoryWrapper{}
	clientFactory.On("CreateWrapper").Return(&pulsartestutil.MockRPCClientWrapper{})

	result, err := NewPulsar(config, storage, clientFactory,
		pulsartestutil.MockEntropyGenerator{}, nil, func(connectionType string, address string) (net.Listener, error) {
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

	mockClientWrapper := &pulsartestutil.MockRPCClientWrapper{}
	mockClientWrapper.On("IsInitialised").Return(true)
	mockClientWrapper.On("Lock")
	mockClientWrapper.On("Unlock")
	mockClientWrapper.On("CreateConnection")

	firstPrivateKey, _ := ecdsa_helper.GeneratePrivateKey()
	expectedNeighbourKey, _ := ecdsa_helper.ExportPublicKey(&firstPrivateKey.PublicKey)
	expectedNeighbour := &Neighbour{OutgoingClient: mockClientWrapper}
	pulsar.Neighbours[expectedNeighbourKey] = expectedNeighbour

	err := pulsar.EstablishConnectionToPulsar(expectedNeighbourKey)

	assert.NoError(t, err)
	mockClientWrapper.AssertNumberOfCalls(t, "Lock", 1)
	mockClientWrapper.AssertNumberOfCalls(t, "Unlock", 1)
	mockClientWrapper.AssertNotCalled(t, "CreateConnection")
}

func TestPulsar_EstablishConnection_IsNotInitialised_ProblemsCreateConnection(t *testing.T) {
	pulsar := &Pulsar{Neighbours: map[string]*Neighbour{}}

	mockClientWrapper := &pulsartestutil.MockRPCClientWrapper{}
	mockClientWrapper.On("IsInitialised").Return(false)
	mockClientWrapper.On("CreateConnection", mock.Anything, mock.Anything).Return(errors.New("test reasons"))
	mockClientWrapper.On("Lock")
	mockClientWrapper.On("Unlock")

	firstPrivateKey, _ := ecdsa_helper.GeneratePrivateKey()
	expectedNeighbourKey, _ := ecdsa_helper.ExportPublicKey(&firstPrivateKey.PublicKey)
	expectedNeighbour := &Neighbour{OutgoingClient: mockClientWrapper}
	pulsar.Neighbours[expectedNeighbourKey] = expectedNeighbour

	err := pulsar.EstablishConnectionToPulsar(expectedNeighbourKey)

	assert.Error(t, err, "test reasons")
	mockClientWrapper.AssertNumberOfCalls(t, "Lock", 1)
	mockClientWrapper.AssertNumberOfCalls(t, "Unlock", 1)
}

func TestPulsar_EstablishConnection_IsNotInitialised_ProblemsWithRequest(t *testing.T) {
	mainPrivateKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)
	pulsar := &Pulsar{
		Neighbours:       map[string]*Neighbour{},
		PrivateKey:       mainPrivateKey,
		EntropyGenerator: pulsartestutil.MockEntropyGenerator{},
	}

	mockClientWrapper := &pulsartestutil.CustomRPCWrapperMock{}

	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{Error: errors.New("oops, request is broken")}
	replyChan := &rpc.Call{Done: done}
	mockClientWrapper.Done = replyChan

	firstPrivateKey, _ := ecdsa_helper.GeneratePrivateKey()
	expectedNeighbourKey, _ := ecdsa_helper.ExportPublicKey(&firstPrivateKey.PublicKey)
	expectedNeighbour := &Neighbour{OutgoingClient: mockClientWrapper}
	pulsar.Neighbours[expectedNeighbourKey] = expectedNeighbour

	err = pulsar.EstablishConnectionToPulsar(expectedNeighbourKey)

	assert.Error(t, err, "oops, request is broken")
}

func TestPulsar_EstablishConnection_IsNotInitialised_SignatureFailed(t *testing.T) {
	mainPrivateKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)
	pulsar := &Pulsar{
		Neighbours:       map[string]*Neighbour{},
		PrivateKey:       mainPrivateKey,
		EntropyGenerator: pulsartestutil.MockEntropyGenerator{},
	}

	mockClientWrapper := &pulsartestutil.CustomRPCWrapperMock{}

	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{Reply: &Payload{Body: HandshakePayload{}}}
	replyChan := &rpc.Call{Done: done}
	mockClientWrapper.Done = replyChan

	firstPrivateKey, _ := ecdsa_helper.GeneratePrivateKey()
	expectedNeighbourKey, _ := ecdsa_helper.ExportPublicKey(&firstPrivateKey.PublicKey)
	expectedNeighbour := &Neighbour{OutgoingClient: mockClientWrapper}
	pulsar.Neighbours[expectedNeighbourKey] = expectedNeighbour

	err = pulsar.EstablishConnectionToPulsar(expectedNeighbourKey)

	assert.Error(t, err, "Signature check failed")
}

func TestPulsar_EstablishConnection_IsNotInitialised_Success(t *testing.T) {
	mainPrivateKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)
	pulsar := &Pulsar{
		Neighbours:       map[string]*Neighbour{},
		PrivateKey:       mainPrivateKey,
		EntropyGenerator: pulsartestutil.MockEntropyGenerator{},
	}
	firstPrivateKey, _ := ecdsa_helper.GeneratePrivateKey()
	expectedNeighbourKey, _ := ecdsa_helper.ExportPublicKey(&firstPrivateKey.PublicKey)
	payload := Payload{Body: HandshakePayload{Entropy: pulsartestutil.MockEntropy}}
	sign, err := signData(firstPrivateKey, payload.Body)
	payload.Signature = sign
	payload.PublicKey = expectedNeighbourKey

	mockClientWrapper := &pulsartestutil.CustomRPCWrapperMock{}

	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{Reply: &payload}
	replyChan := &rpc.Call{Done: done}
	mockClientWrapper.Done = replyChan

	expectedNeighbour := &Neighbour{OutgoingClient: mockClientWrapper}
	pulsar.Neighbours[expectedNeighbourKey] = expectedNeighbour

	err = pulsar.EstablishConnectionToPulsar(expectedNeighbourKey)

	assert.NoError(t, err)
}

func TestPulsar_CheckConnectionsToPulsars_NoProblems(t *testing.T) {
	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{}
	replyChan := &rpc.Call{Done: done}

	mockClientWrapper := &pulsartestutil.MockRPCClientWrapper{}
	mockClientWrapper.On("IsInitialised").Return(true)
	mockClientWrapper.On("Go", HealthCheck.String(), nil, nil, mock.Anything).Return(replyChan)

	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}}
	firstNeighbourPrivateKey, _ := ecdsa_helper.GeneratePrivateKey()
	firstNeighbourExpectedKey, _ := ecdsa_helper.ExportPublicKey(&firstNeighbourPrivateKey.PublicKey)
	pulsar.Neighbours[firstNeighbourExpectedKey] = &Neighbour{
		PublicKey:      &firstNeighbourPrivateKey.PublicKey,
		OutgoingClient: mockClientWrapper,
	}

	pulsar.CheckConnectionsToPulsars()

	mockClientWrapper.AssertNumberOfCalls(t, "IsInitialised", 1)
	mockClientWrapper.AssertNumberOfCalls(t, "Go", 1)
}

func TestPulsar_CheckConnectionsToPulsars_NilClient_FirstConnectionFailed(t *testing.T) {
	mockClientWrapper := &pulsartestutil.MockRPCClientWrapper{}
	mockClientWrapper.On("IsInitialised").Return(false)
	mockClientWrapper.On("Lock")
	mockClientWrapper.On("Unlock")
	mockClientWrapper.On("CreateConnection", mock.Anything, mock.Anything).Return(errors.New("this will have to fall"))

	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}}
	firstNeighbourPrivateKey, _ := ecdsa_helper.GeneratePrivateKey()
	pulsar.Neighbours["thisShouldFailEstablishConnection"] = &Neighbour{
		PublicKey:      &firstNeighbourPrivateKey.PublicKey,
		OutgoingClient: mockClientWrapper,
	}

	resultLog := capture(pulsar.CheckConnectionsToPulsars)

	assert.Contains(t, resultLog, "this will have to fall")
}

func TestPulsar_CheckConnectionsToPulsars_NilClient_SecondConnectionFailed(t *testing.T) {
	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{Error: errors.New("test error")}
	replyChan := &rpc.Call{Done: done}
	mockClientWrapper := &pulsartestutil.CustomRPCWrapperMock{}
	mockClientWrapper.Done = replyChan
	isInitTimeCalled := 0
	mockClientWrapper.IsInitFunc = func() bool {
		if isInitTimeCalled == 0 {
			isInitTimeCalled++
			return true
		} else {
			isInitTimeCalled++
			return false
		}
	}
	mockClientWrapper.CreateConnectionFunc = func() error {
		return errors.New("this should failed")
	}

	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}}
	firstNeighbourPrivateKey, _ := ecdsa_helper.GeneratePrivateKey()
	pulsar.Neighbours["this should fail second connection"] = &Neighbour{
		PublicKey:         &firstNeighbourPrivateKey.PublicKey,
		OutgoingClient:    mockClientWrapper,
		ConnectionAddress: "TestConnectionAddress",
	}

	resultLog := capture(pulsar.CheckConnectionsToPulsars)

	assert.Contains(t, resultLog, "Problems with connection to TestConnectionAddress, with error - test error")
	assert.Contains(t, resultLog, "Attempt of connection to TestConnectionAddress failed with error - this should failed")
	assert.Equal(t, 2, isInitTimeCalled)
}

func TestPulsar_StartConsensusProcess_WithWrongPulseNumber(t *testing.T) {
	pulsar := &Pulsar{stateSwitcher: &MockStateSwitcher{}}
	pulsar.ProcessingPulseNumber = core.PulseNumber(123)
	pulsar.LastPulse = &core.Pulse{PulseNumber: core.PulseNumber(122)}

	err := pulsar.StartConsensusProcess(core.PulseNumber(121))

	assert.Error(t, err, "wrong state status or pulse number, state - waitingForStart, received pulse - 121, last pulse - 122, processing pulse - 123")
}

type MockStateSwitcher struct {
	mock.Mock
	state State
}

func (impl *MockStateSwitcher) getState() State {
	return impl.state
}

func (impl *MockStateSwitcher) setState(state State) {
	impl.state = state
}

func (impl *MockStateSwitcher) SetPulsar(pulsar *Pulsar) {
	impl.Called(pulsar)
}

func (impl *MockStateSwitcher) switchToState(state State, args interface{}) {
	impl.Called(state, args)
}

func TestPulsar_StartConsensusProcess_Success(t *testing.T) {
	mainPrivateKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)
	mockSwitcher := MockStateSwitcher{}
	mockSwitcher.On("switchToState", waitingForEntropySigns, nil)
	pulsar := &Pulsar{
		EntropyGenerator: pulsartestutil.MockEntropyGenerator{},
		PrivateKey:       mainPrivateKey,
		OwnedBftRow:      map[string]*bftCell{},
	}
	pulsar.ProcessingPulseNumber = core.PulseNumber(120)
	pulsar.LastPulse = &core.Pulse{PulseNumber: core.PulseNumber(2)}
	pulsar.stateSwitcher = &mockSwitcher
	expectedPulse := core.PulseNumber(123)

	err = pulsar.StartConsensusProcess(expectedPulse)

	assert.NoError(t, err)
	assert.Equal(t, pulsar.ProcessingPulseNumber, expectedPulse)
	mockSwitcher.AssertNumberOfCalls(t, "switchToState", 1)
	mockSwitcher.AssertCalled(t, "switchToState", waitingForEntropySigns, nil)
}

func TestPulsar_broadcastSignatureOfEntropy_StateFailed(t *testing.T) {
	mockClientWrapper := &pulsartestutil.MockRPCClientWrapper{}
	mockClientWrapper.On("Go", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}, stateSwitcher: &MockStateSwitcher{}}
	pulsar.Neighbours["1"] = &Neighbour{}
	pulsar.ProcessingPulseNumber = core.PulseNumber(123)
	pulsar.stateSwitcher.setState(failed)

	pulsar.broadcastSignatureOfEntropy()

	mockClientWrapper.AssertNotCalled(t, "Go", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestPulsar_broadcastSignatureOfEntropy_SendToNeighbours(t *testing.T) {
	done := make(chan *rpc.Call, 2)
	done <- &rpc.Call{}
	done <- &rpc.Call{Error: errors.New("failed")}
	replyChan := &rpc.Call{Done: done}

	mainPrivateKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)

	mockClientWrapper := &pulsartestutil.MockRPCClientWrapper{}
	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}, PrivateKey: mainPrivateKey, stateSwitcher: &MockStateSwitcher{}}
	pulsar.stateSwitcher.setState(waitingForStart)
	pulsar.ProcessingPulseNumber = 123
	pulsar.GeneratedEntropySign = pulsartestutil.MockEntropy[:]
	pulsar.Neighbours["1"] = &Neighbour{OutgoingClient: mockClientWrapper, ConnectionAddress: "first"}
	pulsar.Neighbours["2"] = &Neighbour{OutgoingClient: mockClientWrapper, ConnectionAddress: "second"}

	mockClientWrapper.On("Go", ReceiveSignatureForEntropy.String(), mock.Anything, nil, (chan *rpc.Call)(nil)).Return(replyChan)

	resultLog := capture(pulsar.broadcastSignatureOfEntropy)

	mockClientWrapper.AssertCalled(t, "Go", ReceiveSignatureForEntropy.String(), mock.Anything, mock.Anything, mock.Anything)
	mockClientWrapper.AssertNumberOfCalls(t, "Go", 2)
	assert.Equal(t, 0, len(done))
	assert.Contains(t, resultLog, "finished with error - failed")
}

func TestPulsar_broadcastVector_StateFailed(t *testing.T) {
	mockClientWrapper := &pulsartestutil.MockRPCClientWrapper{}
	mockClientWrapper.On("Go", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}, OwnedBftRow: map[string]*bftCell{}, stateSwitcher: &MockStateSwitcher{}}
	pulsar.Neighbours["1"] = &Neighbour{}
	pulsar.ProcessingPulseNumber = core.PulseNumber(123)
	pulsar.stateSwitcher.setState(failed)

	pulsar.broadcastVector()

	mockClientWrapper.AssertNotCalled(t, "Go", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestPulsar_broadcastVector_SendToNeighbours(t *testing.T) {
	done := make(chan *rpc.Call, 2)
	done <- &rpc.Call{}
	done <- &rpc.Call{Error: errors.New("failed")}
	replyChan := &rpc.Call{Done: done}

	mainPrivateKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)

	mockClientWrapper := &pulsartestutil.MockRPCClientWrapper{}
	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}, PrivateKey: mainPrivateKey, OwnedBftRow: map[string]*bftCell{}, stateSwitcher: &MockStateSwitcher{}}
	pulsar.stateSwitcher.setState(waitingForStart)
	pulsar.GeneratedEntropySign = pulsartestutil.MockEntropy[:]
	pulsar.GeneratedEntropy = pulsartestutil.MockEntropy
	pulsar.Neighbours["1"] = &Neighbour{OutgoingClient: mockClientWrapper, ConnectionAddress: "first"}
	pulsar.Neighbours["2"] = &Neighbour{OutgoingClient: mockClientWrapper, ConnectionAddress: "second"}

	mockClientWrapper.On("Go", ReceiveVector.String(), mock.Anything, nil, (chan *rpc.Call)(nil)).Return(replyChan)

	resultLog := capture(pulsar.broadcastVector)

	mockClientWrapper.AssertCalled(t, "Go", ReceiveVector.String(), mock.Anything, mock.Anything, mock.Anything)
	mockClientWrapper.AssertNumberOfCalls(t, "Go", 2)
	assert.Equal(t, 0, len(done))
	assert.Contains(t, resultLog, "finished with error - failed")
}

func TestPulsar_broadcastEntropy_StateFailed(t *testing.T) {
	mockClientWrapper := &pulsartestutil.MockRPCClientWrapper{}
	mockClientWrapper.On("Go", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}, OwnedBftRow: map[string]*bftCell{}, stateSwitcher: &MockStateSwitcher{}}
	pulsar.Neighbours["1"] = &Neighbour{}
	pulsar.ProcessingPulseNumber = core.PulseNumber(123)
	pulsar.stateSwitcher.setState(failed)

	pulsar.broadcastEntropy()

	mockClientWrapper.AssertNotCalled(t, "Go", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestPulsar_broadcastEntropy_SendToNeighbours(t *testing.T) {
	done := make(chan *rpc.Call, 2)
	done <- &rpc.Call{}
	done <- &rpc.Call{Error: errors.New("failed")}
	replyChan := &rpc.Call{Done: done}

	mainPrivateKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)

	mockClientWrapper := &pulsartestutil.MockRPCClientWrapper{}
	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}, PrivateKey: mainPrivateKey, stateSwitcher: &MockStateSwitcher{}}
	pulsar.stateSwitcher.setState(waitingForStart)
	pulsar.ProcessingPulseNumber = 123
	pulsar.GeneratedEntropy = pulsartestutil.MockEntropy
	pulsar.Neighbours["1"] = &Neighbour{OutgoingClient: mockClientWrapper, ConnectionAddress: "first"}
	pulsar.Neighbours["2"] = &Neighbour{OutgoingClient: mockClientWrapper, ConnectionAddress: "second"}

	mockClientWrapper.On("Go", ReceiveEntropy.String(), mock.Anything, nil, (chan *rpc.Call)(nil)).Return(replyChan)

	resultLog := capture(pulsar.broadcastEntropy)

	mockClientWrapper.AssertCalled(t, "Go", ReceiveEntropy.String(), mock.Anything, mock.Anything, mock.Anything)
	mockClientWrapper.AssertNumberOfCalls(t, "Go", 2)
	assert.Equal(t, 0, len(done))
	assert.Contains(t, resultLog, "finished with error - failed")
}

func TestPulsar_sendVector_StateFailed(t *testing.T) {
	pulsar := Pulsar{}
	switcher := MockStateSwitcher{}
	switcher.setState(failed)
	pulsar.stateSwitcher = &switcher

	pulsar.sendVector()

	switcher.AssertNotCalled(t, "switchToState")
}

func TestPulsar_sendVector_OnePulsar(t *testing.T) {
	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}}
	switcher := MockStateSwitcher{}
	switcher.setState(waitingForStart)
	switcher.On("switchToState", verifying, nil)
	pulsar.stateSwitcher = &switcher

	pulsar.sendVector()

	switcher.AssertCalled(t, "switchToState", verifying, nil)
	switcher.AssertNumberOfCalls(t, "switchToState", 1)
}

func TestPulsar_sendVector_TwoPulsars(t *testing.T) {
	// Arrange
	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{}
	replyChan := &rpc.Call{Done: done}

	mainPrivateKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)

	mockClientWrapper := &pulsartestutil.MockRPCClientWrapper{}
	pulsar := Pulsar{
		Neighbours:  map[string]*Neighbour{},
		PrivateKey:  mainPrivateKey,
		OwnedBftRow: map[string]*bftCell{},
		BftGrid:     map[string]map[string]*bftCell{},
	}

	pulsar.GeneratedEntropy = pulsartestutil.MockEntropy
	pulsar.Neighbours["1"] = &Neighbour{OutgoingClient: mockClientWrapper, ConnectionAddress: "first"}
	mockClientWrapper.On("Go", ReceiveVector.String(), mock.Anything, nil, (chan *rpc.Call)(nil)).Return(replyChan)

	switcher := MockStateSwitcher{}
	switcher.On("switchToState", verifying, nil)
	switcher.On("switchToState", waitingForVectors, nil)
	switcher.setState(waitingForStart)
	pulsar.stateSwitcher = &switcher

	// Act
	pulsar.sendVector()

	// Assert
	switcher.AssertCalled(t, "switchToState", waitingForVectors, nil)
	switcher.AssertNumberOfCalls(t, "switchToState", 1)
}

func TestPulsar_sendEntropy_OnePulsar(t *testing.T) {
	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}}
	switcher := MockStateSwitcher{}
	switcher.setState(waitingForStart)
	switcher.On("switchToState", verifying, nil)
	pulsar.stateSwitcher = &switcher

	pulsar.sendEntropy()

	switcher.AssertCalled(t, "switchToState", verifying, nil)
	switcher.AssertNumberOfCalls(t, "switchToState", 1)
}

func TestPulsar_sendEntropy_TwoPulsars(t *testing.T) {
	// Arrange
	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{}
	replyChan := &rpc.Call{Done: done}

	mainPrivateKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)

	mockClientWrapper := &pulsartestutil.MockRPCClientWrapper{}
	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}, PrivateKey: mainPrivateKey, OwnedBftRow: map[string]*bftCell{}}
	pulsar.GeneratedEntropy = pulsartestutil.MockEntropy
	pulsar.Neighbours["1"] = &Neighbour{OutgoingClient: mockClientWrapper, ConnectionAddress: "first"}
	mockClientWrapper.On("Go", ReceiveEntropy.String(), mock.Anything, nil, (chan *rpc.Call)(nil)).Return(replyChan)

	switcher := MockStateSwitcher{}
	switcher.setState(waitingForStart)
	switcher.On("switchToState", verifying, nil)
	switcher.On("switchToState", waitingForEntropy, nil)
	pulsar.stateSwitcher = &switcher

	// Act
	pulsar.sendEntropy()

	// Assert
	switcher.AssertCalled(t, "switchToState", waitingForEntropy, nil)
	switcher.AssertNumberOfCalls(t, "switchToState", 1)
}

func TestPulsar_verify_failedState(t *testing.T) {
	pulsar := &Pulsar{stateSwitcher: &MockStateSwitcher{}}
	pulsar.stateSwitcher.setState(failed)
	pulsar.PublicKeyRaw = "testKey"
	pulsar.OwnedBftRow = map[string]*bftCell{}
	pulsar.BftGrid = map[string]map[string]*bftCell{}

	pulsar.verify()
}

func TestPulsar_verify_Standalone_Success(t *testing.T) {
	mockSwitcher := &MockStateSwitcher{}
	mockSwitcher.On("switchToState", sendingPulse, nil)
	pulsar := &Pulsar{stateSwitcher: mockSwitcher}
	pulsar.stateSwitcher.setState(verifying)
	pulsar.PublicKeyRaw = "testKey"
	pulsar.GeneratedEntropy = pulsartestutil.MockEntropy
	pulsar.OwnedBftRow = map[string]*bftCell{}
	pulsar.BftGrid = map[string]map[string]*bftCell{}

	pulsar.verify()

	mockSwitcher.AssertCalled(t, "switchToState", sendingPulse, nil)
	assert.Equal(t, "testKey", pulsar.PublicKeyRaw)
	assert.Equal(t, core.Entropy(pulsartestutil.MockEntropy), pulsar.GeneratedEntropy)
}

func TestPulsar_verify_NotEnoughForConsensus_Success(t *testing.T) {
	mainPrivateKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)

	mockSwitcher := &MockStateSwitcher{}
	mockSwitcher.On("switchToState", failed, mock.Anything)
	pulsar := &Pulsar{stateSwitcher: mockSwitcher}
	pulsar.stateSwitcher.setState(verifying)
	pulsar.PublicKeyRaw = "testKey"
	pulsar.PrivateKey = mainPrivateKey
	pulsar.OwnedBftRow = map[string]*bftCell{}
	pulsar.BftGrid = map[string]map[string]*bftCell{}
	pulsar.Neighbours = map[string]*Neighbour{}
	pulsar.Neighbours["1"] = &Neighbour{}
	pulsar.Neighbours["2"] = &Neighbour{}

	pulsar.verify()

	mockSwitcher.AssertCalled(t, "switchToState", failed, mock.Anything)
}

func generatePrivateAndConvertPublic(t *testing.T) (privateKey *ecdsa.PrivateKey, privateKeyPem string, pubKey string) {
	privateKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)
	pubKey, err = ecdsa_helper.ExportPublicKey(&privateKey.PublicKey)
	assert.NoError(t, err)
	privateKeyPem, err = ecdsa_helper.ExportPrivateKey(privateKey)
	assert.NoError(t, err)

	return
}

func prepareEntropy(t *testing.T, key *ecdsa.PrivateKey) (entropy core.Entropy, sign []byte) {
	entropy = (&StandardEntropyGenerator{}).GenerateEntropy()
	sign, err := signData(key, entropy)
	assert.NoError(t, err)
	return
}

func TestPulsar_verify_Success(t *testing.T) {
	mockSwitcher := &MockStateSwitcher{}
	mockSwitcher.On("switchToState", waitingForPulseSigns, nil)
	mockSwitcher.On("switchToState", sendingPulseSign, nil)

	privateKey, _, currentPulsarPublicKey := generatePrivateAndConvertPublic(t)
	privateKeySecond, _, publicKeySecond := generatePrivateAndConvertPublic(t)
	privateKeyThird, _, publicKeyThird := generatePrivateAndConvertPublic(t)

	clientMock := pulsartestutil.MockRPCClientWrapper{}
	clientMock.On("IsInitialised").Return(true)
	pulsar := &Pulsar{
		stateSwitcher:                  mockSwitcher,
		PrivateKey:                     privateKey,
		PublicKeyRaw:                   currentPulsarPublicKey,
		OwnedBftRow:                    map[string]*bftCell{},
		BftGrid:                        map[string]map[string]*bftCell{},
		CurrentSlotSenderConfirmations: map[string]core.PulseSenderConfirmation{},
		Neighbours: map[string]*Neighbour{
			publicKeySecond: {PublicKey: &privateKeySecond.PublicKey, OutgoingClient: &clientMock},
			publicKeyThird:  {PublicKey: &privateKeyThird.PublicKey, OutgoingClient: &clientMock},
		},
	}
	pulsar.stateSwitcher.setState(verifying)

	firstEntropy, firstSign := prepareEntropy(t, privateKey)
	secondEntropy, secondSign := prepareEntropy(t, privateKeySecond)
	thirdEntropy, thirdSign := prepareEntropy(t, privateKeyThird)

	pulsar.BftGrid[currentPulsarPublicKey] = map[string]*bftCell{
		currentPulsarPublicKey: {Entropy: firstEntropy, Sign: firstSign, IsEntropyReceived: true},
		publicKeySecond:        {Entropy: secondEntropy, Sign: secondSign, IsEntropyReceived: true},
		publicKeyThird:         {Entropy: thirdEntropy, Sign: thirdSign, IsEntropyReceived: true},
	}

	pulsar.BftGrid[publicKeySecond] = map[string]*bftCell{
		currentPulsarPublicKey: {Entropy: firstEntropy, Sign: firstSign, IsEntropyReceived: true},
		publicKeySecond:        {Entropy: secondEntropy, Sign: secondSign, IsEntropyReceived: true},
		publicKeyThird:         {Entropy: thirdEntropy, Sign: thirdSign, IsEntropyReceived: true},
	}
	pulsar.BftGrid[publicKeyThird] = map[string]*bftCell{
		currentPulsarPublicKey: {Entropy: firstEntropy, Sign: firstSign, IsEntropyReceived: true},
		publicKeySecond:        {Entropy: secondEntropy, Sign: secondSign, IsEntropyReceived: true},
		publicKeyThird:         {Entropy: thirdEntropy, Sign: thirdSign, IsEntropyReceived: true},
	}
	var expectedEntropy core.Entropy
	for _, tempEntropy := range []core.Entropy{firstEntropy, secondEntropy, thirdEntropy} {
		for byteIndex := 0; byteIndex < core.EntropySize; byteIndex++ {
			expectedEntropy[byteIndex] ^= tempEntropy[byteIndex]
		}
	}

	pulsar.verify()

	assert.NotNil(t, pulsar.CurrentSlotPulseSender)
	if pulsar.CurrentSlotPulseSender == currentPulsarPublicKey {
		mockSwitcher.AssertCalled(t, "switchToState", waitingForPulseSigns, nil)
	} else {
		mockSwitcher.AssertCalled(t, "switchToState", sendingPulseSign, nil)
	}
	assert.Equal(t, expectedEntropy, pulsar.CurrentSlotEntropy)
}
