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
	ecdsahelper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/pulsar/pulsartestutils"
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

func TestNewPulsar_WithoutNeighbours(t *testing.T) {

	privateKey, privateKeyExported, _ := generatePrivateAndConvertPublic(t)
	actualConnectionType := ""
	actualAddress := ""

	mockListener := func(connectionType string, address string) (net.Listener, error) {
		actualConnectionType = connectionType
		actualAddress = address
		return &pulsartestutils.MockListener{}, nil
	}
	storage := pulsartestutils.NewPulsarStorageMock(t)
	storage.GetLastPulseMock.Return(&core.Pulse{PulseNumber: 123}, nil)

	factoryMock := NewRPCClientWrapperFactoryMock(t)
	clientMock := NewRPCClientWrapperMock(t)
	factoryMock.CreateWrapperMock.Return(clientMock)

	result, err := NewPulsar(configuration.Configuration{
		Pulsar: configuration.Pulsar{
			ConnectionType:      "testType",
			MainListenerAddress: "listedAddress",
		},
		PrivateKey: privateKeyExported,
	},
		storage,
		factoryMock,
		pulsartestutils.MockEntropyGenerator{},
		nil,
		mockListener,
	)

	assert.NoError(t, err)
	assert.Equal(t, privateKey, result.PrivateKey)
	assert.Equal(t, "testType", actualConnectionType)
	assert.Equal(t, "listedAddress", actualAddress)
	assert.IsType(t, result.Sock, &pulsartestutils.MockListener{})
	assert.NotNil(t, result.PrivateKey)

	assert.Equal(t, uint64(0), factoryMock.CreateWrapperCounter)
}

func TestNewPulsar_WithNeighbours(t *testing.T) {
	assertObj := assert.New(t)

	firstPrivateKey, _ := ecdsahelper.GeneratePrivateKey()
	firstExpectedKey, _ := ecdsahelper.ExportPublicKey(&firstPrivateKey.PublicKey)

	secondPrivateKey, _ := ecdsahelper.GeneratePrivateKey()
	secondExpectedKey, _ := ecdsahelper.ExportPublicKey(&secondPrivateKey.PublicKey)

	expectedPrivateKey, _ := ecdsahelper.GeneratePrivateKey()
	parsedExpectedPrivateKey, _ := ecdsahelper.ExportPrivateKey(expectedPrivateKey)

	storage := pulsartestutils.NewPulsarStorageMock(t)
	storage.GetLastPulseMock.Return(&core.Pulse{PulseNumber: 123}, nil)

	factoryMock := NewRPCClientWrapperFactoryMock(t)
	clientMock := NewRPCClientWrapperMock(t)
	factoryMock.CreateWrapperMock.Return(clientMock)

	result, err := NewPulsar(
		configuration.Configuration{
			Pulsar: configuration.Pulsar{
				ConnectionType:      "testType",
				MainListenerAddress: "listedAddress",
				Neighbours: []configuration.PulsarNodeAddress{
					{ConnectionType: "tcp", Address: "first", PublicKey: firstExpectedKey},
					{ConnectionType: "pct", Address: "second", PublicKey: secondExpectedKey},
				},
			},
			PrivateKey: parsedExpectedPrivateKey,
		},
		storage,
		factoryMock,
		pulsartestutils.MockEntropyGenerator{}, nil, func(connectionType string, address string) (net.Listener, error) {
			return &pulsartestutils.MockListener{}, nil
		})

	assertObj.NoError(err)
	assertObj.Equal(2, len(result.Neighbours))
	assertObj.Equal("tcp", result.Neighbours[firstExpectedKey].ConnectionType.String())
	assertObj.Equal("pct", result.Neighbours[secondExpectedKey].ConnectionType.String())
	assert.Equal(t, uint64(2), factoryMock.CreateWrapperCounter)
}

func TestPulsar_EstablishConnection_IsInitialised(t *testing.T) {
	pulsar := &Pulsar{Neighbours: map[string]*Neighbour{}}

	mockClientWrapper := &pulsartestutils.MockRPCClientWrapper{}
	mockClientWrapper.On("IsInitialised").Return(true)
	mockClientWrapper.On("Lock")
	mockClientWrapper.On("Unlock")
	mockClientWrapper.On("CreateConnection")

	firstPrivateKey, _ := ecdsahelper.GeneratePrivateKey()
	expectedNeighbourKey, _ := ecdsahelper.ExportPublicKey(&firstPrivateKey.PublicKey)
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

	clientMock := NewRPCClientWrapperMock(t)
	clientMock.IsInitialisedMock.Return(false)
	clientMock.LockMock.Return()
	clientMock.UnlockMock.Return()
	clientMock.CreateConnectionMock.Set(func(p configuration.ConnectionType, p1 string) (r error) {
		return errors.New("test reasons")
	})

	firstPrivateKey, _ := ecdsahelper.GeneratePrivateKey()
	expectedNeighbourKey, _ := ecdsahelper.ExportPublicKey(&firstPrivateKey.PublicKey)
	expectedNeighbour := &Neighbour{OutgoingClient: clientMock}
	pulsar.Neighbours[expectedNeighbourKey] = expectedNeighbour

	err := pulsar.EstablishConnectionToPulsar(expectedNeighbourKey)

	assert.Error(t, err, "test reasons")
	assert.Equal(t, clientMock.LockCounter, uint64(1))
	assert.Equal(t, clientMock.UnlockCounter, uint64(1))
}

func TestPulsar_EstablishConnection_IsNotInitialised_ProblemsWithRequest(t *testing.T) {
	mainPrivateKey, err := ecdsahelper.GeneratePrivateKey()
	assert.NoError(t, err)
	pulsar := &Pulsar{
		Neighbours:       map[string]*Neighbour{},
		PrivateKey:       mainPrivateKey,
		EntropyGenerator: pulsartestutils.MockEntropyGenerator{},
	}

	mockClientWrapper := &pulsartestutils.CustomRPCWrapperMock{}

	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{Error: errors.New("oops, request is broken")}
	replyChan := &rpc.Call{Done: done}
	mockClientWrapper.Done = replyChan

	firstPrivateKey, _ := ecdsahelper.GeneratePrivateKey()
	expectedNeighbourKey, _ := ecdsahelper.ExportPublicKey(&firstPrivateKey.PublicKey)
	expectedNeighbour := &Neighbour{OutgoingClient: mockClientWrapper}
	pulsar.Neighbours[expectedNeighbourKey] = expectedNeighbour

	err = pulsar.EstablishConnectionToPulsar(expectedNeighbourKey)

	assert.Error(t, err, "oops, request is broken")
}

func TestPulsar_EstablishConnection_IsNotInitialised_SignatureFailed(t *testing.T) {
	mainPrivateKey, err := ecdsahelper.GeneratePrivateKey()
	assert.NoError(t, err)
	pulsar := &Pulsar{
		Neighbours:       map[string]*Neighbour{},
		PrivateKey:       mainPrivateKey,
		EntropyGenerator: pulsartestutils.MockEntropyGenerator{},
	}

	mockClientWrapper := &pulsartestutils.CustomRPCWrapperMock{}

	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{Reply: &Payload{Body: HandshakePayload{}}}
	replyChan := &rpc.Call{Done: done}
	mockClientWrapper.Done = replyChan

	firstPrivateKey, _ := ecdsahelper.GeneratePrivateKey()
	expectedNeighbourKey, _ := ecdsahelper.ExportPublicKey(&firstPrivateKey.PublicKey)
	expectedNeighbour := &Neighbour{OutgoingClient: mockClientWrapper}
	pulsar.Neighbours[expectedNeighbourKey] = expectedNeighbour

	err = pulsar.EstablishConnectionToPulsar(expectedNeighbourKey)

	assert.Error(t, err, "Signature check Failed")
}

func TestPulsar_EstablishConnection_IsNotInitialised_Success(t *testing.T) {
	mainPrivateKey, err := ecdsahelper.GeneratePrivateKey()
	assert.NoError(t, err)
	pulsar := &Pulsar{
		Neighbours:       map[string]*Neighbour{},
		PrivateKey:       mainPrivateKey,
		EntropyGenerator: pulsartestutils.MockEntropyGenerator{},
	}
	firstPrivateKey, _ := ecdsahelper.GeneratePrivateKey()
	expectedNeighbourKey, _ := ecdsahelper.ExportPublicKey(&firstPrivateKey.PublicKey)
	payload := Payload{Body: HandshakePayload{Entropy: pulsartestutils.MockEntropy}}
	sign, err := signData(firstPrivateKey, payload.Body)
	payload.Signature = sign
	payload.PublicKey = expectedNeighbourKey

	mockClientWrapper := &pulsartestutils.CustomRPCWrapperMock{}

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

	clientMock := NewRPCClientWrapperMock(t)
	clientMock.IsInitialisedMock.Return(true)
	clientMock.GoMock.Set(func(p string, p1 interface{}, p2 interface{}, p3 chan *rpc.Call) (r *rpc.Call) {
		return replyChan
	})

	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}}
	firstNeighbourPrivateKey, _ := ecdsahelper.GeneratePrivateKey()
	firstNeighbourExpectedKey, _ := ecdsahelper.ExportPublicKey(&firstNeighbourPrivateKey.PublicKey)
	pulsar.Neighbours[firstNeighbourExpectedKey] = &Neighbour{
		PublicKey:      &firstNeighbourPrivateKey.PublicKey,
		OutgoingClient: clientMock,
	}

	pulsar.CheckConnectionsToPulsars()

	assert.Equal(t, uint64(1), clientMock.IsInitialisedCounter)
	assert.Equal(t, uint64(1), clientMock.GoCounter)
}

func TestPulsar_CheckConnectionsToPulsars_NilClient_FirstConnectionFailed(t *testing.T) {
	clientMock := NewRPCClientWrapperMock(t)
	clientMock.IsInitialisedMock.Return(false)
	clientMock.LockMock.Return()
	clientMock.UnlockMock.Return()
	clientMock.CreateConnectionFunc = func(p configuration.ConnectionType, p1 string) (r error) {
		return errors.New("this will have to fall")
	}

	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}}

	firstNeighbourPrivateKey, _ := ecdsahelper.GeneratePrivateKey()
	pulsar.Neighbours["thisShouldFailEstablishConnection"] = &Neighbour{
		PublicKey:      &firstNeighbourPrivateKey.PublicKey,
		OutgoingClient: clientMock,
	}

	resultLog := capture(pulsar.CheckConnectionsToPulsars)

	assert.Contains(t, resultLog, "this will have to fall")
}

func TestPulsar_CheckConnectionsToPulsars_NilClient_SecondConnectionFailed(t *testing.T) {
	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{Error: errors.New("test error")}
	replyChan := &rpc.Call{Done: done}
	mockClientWrapper := &pulsartestutils.CustomRPCWrapperMock{}
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
		return errors.New("this should Failed")
	}

	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}}

	firstNeighbourPrivateKey, _ := ecdsahelper.GeneratePrivateKey()
	pulsar.Neighbours["this should fail second connection"] = &Neighbour{
		PublicKey:         &firstNeighbourPrivateKey.PublicKey,
		OutgoingClient:    mockClientWrapper,
		ConnectionAddress: "TestConnectionAddress",
	}

	resultLog := capture(pulsar.CheckConnectionsToPulsars)

	assert.Contains(t, resultLog, "Problems with connection to TestConnectionAddress, with error - test error")
	assert.Contains(t, resultLog, "Attempt of connection to TestConnectionAddress Failed with error - this should Failed")
	assert.Equal(t, 2, isInitTimeCalled)
}

func TestPulsar_StartConsensusProcess_WithWrongPulseNumber(t *testing.T) {
	switcherMock := NewStateSwitcherMock(t)
	switcherMock.GetStateMock.Return(WaitingForStart)
	pulsar := &Pulsar{
		StateSwitcher:         switcherMock,
		ProcessingPulseNumber: core.PulseNumber(123),
		lastPulse:             &core.Pulse{PulseNumber: core.PulseNumber(122)},
	}
	err := pulsar.StartConsensusProcess(core.PulseNumber(121))

	assert.Error(t, err, "wrong state status or pulse number, state - WaitingForStart, received pulse - 121, last pulse - 122, processing pulse - 123")
}

func TestPulsar_StartConsensusProcess_Success(t *testing.T) {
	mainPrivateKey, err := ecdsahelper.GeneratePrivateKey()
	assert.NoError(t, err)

	mockSwitcher := NewStateSwitcherMock(t)
	mockSwitcher.GetStateMock.Return(WaitingForStart)
	mockSwitcher.setStateMock.Expect(GenerateEntropy).Return()
	mockSwitcher.SwitchToStateFunc = func(p State, p1 interface{}) {
		if p != WaitingForEntropySigns {
			t.Error(t, "Wrong state")
		}
	}

	pulsar := &Pulsar{
		EntropyGenerator: pulsartestutils.MockEntropyGenerator{},
		PrivateKey:       mainPrivateKey,
		OwnedBftRow:      map[string]*BftCell{},
	}
	pulsar.ProcessingPulseNumber = core.PulseNumber(120)
	pulsar.SetLastPulse(&core.Pulse{PulseNumber: core.PulseNumber(2)})
	pulsar.StateSwitcher = mockSwitcher
	expectedPulse := core.PulseNumber(123)

	err = pulsar.StartConsensusProcess(expectedPulse)

	assert.NoError(t, err)
	assert.Equal(t, pulsar.ProcessingPulseNumber, expectedPulse)
	assert.Equal(t, uint64(1), mockSwitcher.SwitchToStateCounter)
}

func TestPulsar_broadcastSignatureOfEntropy_StateFailed(t *testing.T) {
	mockClientWrapper := &pulsartestutils.MockRPCClientWrapper{}
	mockClientWrapper.On("Go", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	switcher := NewStateSwitcherMock(t)
	switcher.GetStateMock.Return(Failed)
	pulsar := Pulsar{
		Neighbours: map[string]*Neighbour{
			"1": {},
		},
		StateSwitcher:         switcher,
		ProcessingPulseNumber: core.PulseNumber(123),
	}

	pulsar.broadcastSignatureOfEntropy()
}

func TestPulsar_broadcastSignatureOfEntropy_SendToNeighbours(t *testing.T) {
	done := make(chan *rpc.Call, 2)
	done <- &rpc.Call{}
	done <- &rpc.Call{Error: errors.New("Failed")}
	replyChan := &rpc.Call{Done: done}

	mainPrivateKey, err := ecdsahelper.GeneratePrivateKey()
	assert.NoError(t, err)

	mockClientWrapper := &pulsartestutils.MockRPCClientWrapper{}
	switcher := NewStateSwitcherMock(t)
	switcher.GetStateMock.Return(WaitingForStart)
	pulsar := Pulsar{
		Neighbours: map[string]*Neighbour{
			"1": {OutgoingClient: mockClientWrapper, ConnectionAddress: "first"},
			"2": {OutgoingClient: mockClientWrapper, ConnectionAddress: "second"},
		},
		PrivateKey:            mainPrivateKey,
		StateSwitcher:         switcher,
		ProcessingPulseNumber: 123,
		GeneratedEntropySign:  pulsartestutils.MockEntropy[:],
	}

	mockClientWrapper.On("Go", ReceiveSignatureForEntropy.String(), mock.Anything, nil, (chan *rpc.Call)(nil)).Return(replyChan)

	resultLog := capture(pulsar.broadcastSignatureOfEntropy)

	mockClientWrapper.AssertCalled(t, "Go", ReceiveSignatureForEntropy.String(), mock.Anything, mock.Anything, mock.Anything)
	mockClientWrapper.AssertNumberOfCalls(t, "Go", 2)
	assert.Equal(t, 0, len(done))
	assert.Contains(t, resultLog, "finished with error - Failed")
}

func TestPulsar_broadcastVector_StateFailed(t *testing.T) {
	mockClientWrapper := &pulsartestutils.MockRPCClientWrapper{}
	mockClientWrapper.On("Go", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	mockSwitcher := NewStateSwitcherMock(t)
	mockSwitcher.GetStateMock.Return(Failed)
	pulsar := Pulsar{
		Neighbours: map[string]*Neighbour{
			"1": {},
		},
		OwnedBftRow:           map[string]*BftCell{},
		StateSwitcher:         mockSwitcher,
		ProcessingPulseNumber: core.PulseNumber(123),
	}

	pulsar.broadcastVector()

	mockClientWrapper.AssertNotCalled(t, "Go", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestPulsar_broadcastVector_SendToNeighbours(t *testing.T) {
	done := make(chan *rpc.Call, 2)
	done <- &rpc.Call{}
	done <- &rpc.Call{Error: errors.New("Failed")}
	replyChan := &rpc.Call{Done: done}

	mainPrivateKey, err := ecdsahelper.GeneratePrivateKey()
	assert.NoError(t, err)

	mockClientWrapper := &pulsartestutils.MockRPCClientWrapper{}
	mockSwitcher := NewStateSwitcherMock(t)
	mockSwitcher.GetStateMock.Return(WaitingForStart)
	pulsar := Pulsar{
		Neighbours: map[string]*Neighbour{
			"1": {OutgoingClient: mockClientWrapper, ConnectionAddress: "first"},
			"2": {OutgoingClient: mockClientWrapper, ConnectionAddress: "second"},
		},
		PrivateKey:           mainPrivateKey,
		OwnedBftRow:          map[string]*BftCell{},
		StateSwitcher:        mockSwitcher,
		GeneratedEntropySign: pulsartestutils.MockEntropy[:],
		GeneratedEntropy:     pulsartestutils.MockEntropy,
	}
	mockClientWrapper.On("Go", ReceiveVector.String(), mock.Anything, nil, (chan *rpc.Call)(nil)).Return(replyChan)

	resultLog := capture(pulsar.broadcastVector)

	mockClientWrapper.AssertCalled(t, "Go", ReceiveVector.String(), mock.Anything, mock.Anything, mock.Anything)
	mockClientWrapper.AssertNumberOfCalls(t, "Go", 2)
	assert.Equal(t, 0, len(done))
	assert.Contains(t, resultLog, "finished with error - Failed")
}

func TestPulsar_broadcastEntropy_StateFailed(t *testing.T) {
	mockClientWrapper := &pulsartestutils.MockRPCClientWrapper{}
	mockClientWrapper.On("Go", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	mockSwitcher := NewStateSwitcherMock(t)
	mockSwitcher.GetStateMock.Return(Failed)
	pulsar := Pulsar{
		Neighbours: map[string]*Neighbour{
			"1": {},
		},
		OwnedBftRow:           map[string]*BftCell{},
		StateSwitcher:         mockSwitcher,
		ProcessingPulseNumber: core.PulseNumber(123),
	}

	pulsar.broadcastEntropy()

	mockClientWrapper.AssertNotCalled(t, "Go", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestPulsar_broadcastEntropy_SendToNeighbours(t *testing.T) {
	done := make(chan *rpc.Call, 2)
	done <- &rpc.Call{}
	done <- &rpc.Call{Error: errors.New("Failed")}
	replyChan := &rpc.Call{Done: done}

	mainPrivateKey, err := ecdsahelper.GeneratePrivateKey()
	assert.NoError(t, err)

	mockClientWrapper := &pulsartestutils.MockRPCClientWrapper{}
	switcherMock := NewStateSwitcherMock(t)
	switcherMock.GetStateMock.Return(WaitingForStart)
	pulsar := Pulsar{
		Neighbours: map[string]*Neighbour{
			"1": {OutgoingClient: mockClientWrapper, ConnectionAddress: "first"},
			"2": {OutgoingClient: mockClientWrapper, ConnectionAddress: "second"},
		},
		PrivateKey:            mainPrivateKey,
		StateSwitcher:         switcherMock,
		ProcessingPulseNumber: 123,
		GeneratedEntropy:      pulsartestutils.MockEntropy,
	}

	mockClientWrapper.On("Go", ReceiveEntropy.String(), mock.Anything, nil, (chan *rpc.Call)(nil)).Return(replyChan)

	resultLog := capture(pulsar.broadcastEntropy)

	mockClientWrapper.AssertCalled(t, "Go", ReceiveEntropy.String(), mock.Anything, mock.Anything, mock.Anything)
	mockClientWrapper.AssertNumberOfCalls(t, "Go", 2)
	assert.Equal(t, 0, len(done))
	assert.Contains(t, resultLog, "finished with error - Failed")
}

func TestPulsar_sendVector_StateFailed(t *testing.T) {
	pulsar := Pulsar{}
	switcher := NewStateSwitcherMock(t)
	switcher.GetStateMock.Return(Failed)
	pulsar.StateSwitcher = switcher

	pulsar.sendVector()

	assert.Equal(t, uint64(0), switcher.SwitchToStateCounter)
}

func TestPulsar_sendVector_OnePulsar(t *testing.T) {
	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}}
	switcher := NewStateSwitcherMock(t)
	switcher.GetStateMock.Return(WaitingForStart)
	switcher.SwitchToStateFunc = func(p State, p1 interface{}) {
		assert.Equal(t, p, Verifying)
	}
	pulsar.StateSwitcher = switcher

	pulsar.sendVector()

	assert.Equal(t, switcher.SwitchToStateCounter, uint64(1))
}

func TestPulsar_sendVector_TwoPulsars(t *testing.T) {
	// Arrange
	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{}
	replyChan := &rpc.Call{Done: done}

	mainPrivateKey, err := ecdsahelper.GeneratePrivateKey()
	assert.NoError(t, err)

	mockClientWrapper := &pulsartestutils.MockRPCClientWrapper{}
	pulsar := Pulsar{
		Neighbours:  map[string]*Neighbour{},
		PrivateKey:  mainPrivateKey,
		OwnedBftRow: map[string]*BftCell{},
		bftGrid:     map[string]map[string]*BftCell{},
	}

	pulsar.GeneratedEntropy = pulsartestutils.MockEntropy
	pulsar.Neighbours["1"] = &Neighbour{OutgoingClient: mockClientWrapper, ConnectionAddress: "first"}
	mockClientWrapper.On("Go", ReceiveVector.String(), mock.Anything, nil, (chan *rpc.Call)(nil)).Return(replyChan)

	switcher := NewStateSwitcherMock(t)
	switcher.GetStateMock.Return(WaitingForStart)
	switcher.SwitchToStateMock.Expect(WaitingForVectors, nil).Return()
	pulsar.StateSwitcher = switcher

	// Act
	pulsar.sendVector()

	// Assert
	assert.Equal(t, uint64(1), switcher.SwitchToStateCounter)
	assert.Equal(t, WaitingForVectors, switcher.SwitchToStateMock.mockExpectations.p)
}

func TestPulsar_sendEntropy_OnePulsar(t *testing.T) {
	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}}
	switcher := NewStateSwitcherMock(t)
	switcher.GetStateMock.Return(WaitingForStart)
	switcher.SwitchToStateMock.Expect(Verifying, nil).Return()
	pulsar.StateSwitcher = switcher

	pulsar.sendEntropy()

	assert.Equal(t, uint64(1), switcher.SwitchToStateCounter)
	assert.Equal(t, Verifying, switcher.SwitchToStateMock.mockExpectations.p)
}

func TestPulsar_sendEntropy_TwoPulsars(t *testing.T) {
	// Arrange
	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{}
	replyChan := &rpc.Call{Done: done}

	mainPrivateKey, err := ecdsahelper.GeneratePrivateKey()
	assert.NoError(t, err)

	mockClientWrapper := &pulsartestutils.MockRPCClientWrapper{}
	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}, PrivateKey: mainPrivateKey, OwnedBftRow: map[string]*BftCell{}}
	pulsar.GeneratedEntropy = pulsartestutils.MockEntropy
	pulsar.Neighbours["1"] = &Neighbour{OutgoingClient: mockClientWrapper, ConnectionAddress: "first"}
	mockClientWrapper.On("Go", ReceiveEntropy.String(), mock.Anything, nil, (chan *rpc.Call)(nil)).Return(replyChan)

	switcher := NewStateSwitcherMock(t)
	switcher.GetStateMock.Return(WaitingForStart)
	switcher.SwitchToStateMock.Expect(WaitingForEntropy, nil).Return()
	pulsar.StateSwitcher = switcher

	// Act
	pulsar.sendEntropy()

	// Assert
	assert.Equal(t, uint64(1), switcher.SwitchToStateCounter)
}

func TestPulsar_verify_failedState(t *testing.T) {
	switcherMock := NewStateSwitcherMock(t)
	switcherMock.GetStateMock.Return(Failed)
	pulsar := &Pulsar{StateSwitcher: switcherMock}
	pulsar.PublicKeyRaw = "testKey"
	pulsar.OwnedBftRow = map[string]*BftCell{}
	pulsar.bftGrid = map[string]map[string]*BftCell{}

	pulsar.verify()
}

func TestPulsar_verify_Standalone_Success(t *testing.T) {
	mockSwitcher := NewStateSwitcherMock(t)
	mockSwitcher.GetStateMock.Return(Verifying)
	mockSwitcher.SwitchToStateMock.Expect(SendingPulse, nil).Return()
	pulsar := &Pulsar{StateSwitcher: mockSwitcher}
	pulsar.PublicKeyRaw = "testKey"
	pulsar.GeneratedEntropy = pulsartestutils.MockEntropy
	pulsar.OwnedBftRow = map[string]*BftCell{}
	pulsar.bftGrid = map[string]map[string]*BftCell{}

	pulsar.verify()

	assert.Equal(t, uint64(1), mockSwitcher.SwitchToStateCounter)
	assert.Equal(t, "testKey", pulsar.PublicKeyRaw)
	assert.Equal(t, core.Entropy(pulsartestutils.MockEntropy), pulsar.GeneratedEntropy)
}

func TestPulsar_verify_NotEnoughForConsensus_Success(t *testing.T) {
	mainPrivateKey, err := ecdsahelper.GeneratePrivateKey()
	assert.NoError(t, err)

	mockSwitcher := NewStateSwitcherMock(t)
	mockSwitcher.SwitchToStateMock.Set(func(p State, p1 interface{}) {
		assert.Equal(t, Failed, p)
	})
	mockSwitcher.GetStateMock.Return(Verifying)

	pulsar := &Pulsar{StateSwitcher: mockSwitcher}
	pulsar.PublicKeyRaw = "testKey"
	pulsar.PrivateKey = mainPrivateKey
	pulsar.OwnedBftRow = map[string]*BftCell{}
	pulsar.bftGrid = map[string]map[string]*BftCell{}
	pulsar.Neighbours = map[string]*Neighbour{}
	pulsar.Neighbours["1"] = &Neighbour{}
	pulsar.Neighbours["2"] = &Neighbour{}

	pulsar.verify()

	assert.Equal(t, uint64(1), mockSwitcher.SwitchToStateCounter)
}

func generatePrivateAndConvertPublic(t *testing.T) (privateKey *ecdsa.PrivateKey, privateKeyPem string, pubKey string) {
	privateKey, err := ecdsahelper.GeneratePrivateKey()
	assert.NoError(t, err)
	pubKey, err = ecdsahelper.ExportPublicKey(&privateKey.PublicKey)
	assert.NoError(t, err)
	privateKeyPem, err = ecdsahelper.ExportPrivateKey(privateKey)
	assert.NoError(t, err)

	return
}

func prepareEntropy(t *testing.T, key *ecdsa.PrivateKey) (entropy core.Entropy, sign []byte) {
	entropy = (&entropygenerator.StandardEntropyGenerator{}).GenerateEntropy()
	sign, err := signData(key, entropy)
	assert.NoError(t, err)
	return
}

func TestPulsar_verify_Success(t *testing.T) {
	mockSwitcher := NewStateSwitcherMock(t)
	mockSwitcher.SwitchToStateFunc = func(p State, p1 interface{}) {
		if p != WaitingForPulseSigns && p != SendingPulseSign {
			t.Error("Wrong state")
		}
	}
	mockSwitcher.GetStateFunc = func() State {
		return Verifying
	}

	privateKey, _, currentPulsarPublicKey := generatePrivateAndConvertPublic(t)
	privateKeySecond, _, publicKeySecond := generatePrivateAndConvertPublic(t)
	privateKeyThird, _, publicKeyThird := generatePrivateAndConvertPublic(t)

	clientMock := pulsartestutils.MockRPCClientWrapper{}
	clientMock.On("IsInitialised").Return(true)
	pulsar := &Pulsar{
		StateSwitcher:                  mockSwitcher,
		PrivateKey:                     privateKey,
		PublicKeyRaw:                   currentPulsarPublicKey,
		OwnedBftRow:                    map[string]*BftCell{},
		bftGrid:                        map[string]map[string]*BftCell{},
		CurrentSlotSenderConfirmations: map[string]core.PulseSenderConfirmation{},
		Neighbours: map[string]*Neighbour{
			publicKeySecond: {PublicKey: &privateKeySecond.PublicKey, OutgoingClient: &clientMock},
			publicKeyThird:  {PublicKey: &privateKeyThird.PublicKey, OutgoingClient: &clientMock},
		},
	}

	firstEntropy, firstSign := prepareEntropy(t, privateKey)
	secondEntropy, secondSign := prepareEntropy(t, privateKeySecond)
	thirdEntropy, thirdSign := prepareEntropy(t, privateKeyThird)

	pulsar.bftGrid[currentPulsarPublicKey] = map[string]*BftCell{
		currentPulsarPublicKey: {Entropy: firstEntropy, Sign: firstSign, IsEntropyReceived: true},
		publicKeySecond:        {Entropy: secondEntropy, Sign: secondSign, IsEntropyReceived: true},
		publicKeyThird:         {Entropy: thirdEntropy, Sign: thirdSign, IsEntropyReceived: true},
	}

	pulsar.bftGrid[publicKeySecond] = map[string]*BftCell{
		currentPulsarPublicKey: {Entropy: firstEntropy, Sign: firstSign, IsEntropyReceived: true},
		publicKeySecond:        {Entropy: secondEntropy, Sign: secondSign, IsEntropyReceived: true},
		publicKeyThird:         {Entropy: thirdEntropy, Sign: thirdSign, IsEntropyReceived: true},
	}
	pulsar.bftGrid[publicKeyThird] = map[string]*BftCell{
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
	assert.Equal(t, expectedEntropy, pulsar.CurrentSlotEntropy)
	assert.Equal(t, uint64(1), mockSwitcher.SwitchToStateCounter)
}
