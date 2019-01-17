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
	"context"
	"crypto"
	"net"
	"net/rpc"
	"os"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/pulsar/pulsartestutils"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func capture(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

func TestNewPulsar_WithoutNeighbours(t *testing.T) {
	actualConnectionType := ""
	actualAddress := ""

	mockListener := func(connectionType string, address string) (net.Listener, error) {
		actualConnectionType = connectionType
		actualAddress = address
		return &pulsartestutils.MockListener{}, nil
	}
	storage := pulsartestutils.NewPulsarStorageMock(t)
	storage.GetLastPulseMock.Return(&core.Pulse{PulseNumber: 123}, nil)

	keyProcessor := platformpolicy.NewKeyProcessor()
	privateKey, err := keyProcessor.GeneratePrivateKey()
	require.NoError(t, err)
	cryptoService := cryptography.NewKeyBoundCryptographyService(privateKey)
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	factoryMock := NewRPCClientWrapperFactoryMock(t)
	clientMock := NewRPCClientWrapperMock(t)
	factoryMock.CreateWrapperMock.Return(clientMock)

	result, err := NewPulsar(configuration.Pulsar{
		ConnectionType:      "testType",
		MainListenerAddress: "listedAddress",
	},
		cryptoService,
		scheme,
		keyProcessor,
		testutils.NewPulseDistributorMock(t),
		storage,
		factoryMock,
		pulsartestutils.MockEntropyGenerator{},
		nil,
		mockListener,
	)

	require.NoError(t, err)
	require.Equal(t, "testType", actualConnectionType)
	require.Equal(t, "listedAddress", actualAddress)
	require.IsType(t, result.Sock, &pulsartestutils.MockListener{})

	require.Equal(t, uint64(0), factoryMock.CreateWrapperCounter)
}

func TestNewPulsar_WithNeighbours(t *testing.T) {
	requireObj := require.New(t)

	keyProcessor := platformpolicy.NewKeyProcessor()

	firstPrivateKey, _ := keyProcessor.GeneratePrivateKey()
	firstExpectedKey, _ := keyProcessor.ExportPublicKeyPEM(keyProcessor.ExtractPublicKey(firstPrivateKey))

	secondPrivateKey, _ := keyProcessor.GeneratePrivateKey()
	secondExpectedKey, _ := keyProcessor.ExportPublicKeyPEM(keyProcessor.ExtractPublicKey(secondPrivateKey))

	pulsarPrivateKey, _ := keyProcessor.GeneratePrivateKey()
	pulsarCryptoService := cryptography.NewKeyBoundCryptographyService(pulsarPrivateKey)
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	storage := pulsartestutils.NewPulsarStorageMock(t)
	storage.GetLastPulseMock.Return(&core.Pulse{PulseNumber: 123}, nil)

	factoryMock := NewRPCClientWrapperFactoryMock(t)
	clientMock := NewRPCClientWrapperMock(t)
	factoryMock.CreateWrapperMock.Return(clientMock)

	result, err := NewPulsar(
		configuration.Pulsar{
			ConnectionType:      "testType",
			MainListenerAddress: "listedAddress",
			Neighbours: []configuration.PulsarNodeAddress{
				{ConnectionType: "tcp", Address: "first", PublicKey: string(firstExpectedKey)},
				{ConnectionType: "pct", Address: "second", PublicKey: string(secondExpectedKey)},
			},
		},
		pulsarCryptoService,
		scheme,
		keyProcessor,
		testutils.NewPulseDistributorMock(t),
		storage,
		factoryMock,
		pulsartestutils.MockEntropyGenerator{},
		nil,
		func(connectionType string, address string) (net.Listener, error) {
			return &pulsartestutils.MockListener{}, nil
		})

	requireObj.NoError(err)
	requireObj.Equal(2, len(result.Neighbours))
	requireObj.Equal("tcp", result.Neighbours[string(firstExpectedKey)].ConnectionType.String())
	requireObj.Equal("pct", result.Neighbours[string(secondExpectedKey)].ConnectionType.String())
	require.Equal(t, uint64(2), factoryMock.CreateWrapperCounter)
}

func TestPulsar_EstablishConnection_IsInitialised(t *testing.T) {
	ctx := inslogger.TestContext(t)
	pulsar := &Pulsar{Neighbours: map[string]*Neighbour{}}

	mockClientWrapper := &pulsartestutils.MockRPCClientWrapper{}
	mockClientWrapper.On("IsInitialised").Return(true)
	mockClientWrapper.On("Lock")
	mockClientWrapper.On("Unlock")
	mockClientWrapper.On("CreateConnection")

	keyProcessor := platformpolicy.NewKeyProcessor()
	firstPrivateKey, _ := keyProcessor.GeneratePrivateKey()
	expectedNeighbourKey, _ := keyProcessor.ExportPublicKeyPEM(keyProcessor.ExtractPublicKey(firstPrivateKey))
	expectedNeighbour := &Neighbour{OutgoingClient: mockClientWrapper}
	pulsar.Neighbours[string(expectedNeighbourKey)] = expectedNeighbour

	err := pulsar.EstablishConnectionToPulsar(ctx, string(expectedNeighbourKey))

	require.NoError(t, err)
	mockClientWrapper.AssertNumberOfCalls(t, "Lock", 1)
	mockClientWrapper.AssertNumberOfCalls(t, "Unlock", 1)
	mockClientWrapper.AssertNotCalled(t, "CreateConnection")
}

func TestPulsar_EstablishConnection_IsNotInitialised_ProblemsCreateConnection(t *testing.T) {
	ctx := inslogger.TestContext(t)
	pulsar := &Pulsar{Neighbours: map[string]*Neighbour{}}

	clientMock := NewRPCClientWrapperMock(t)
	clientMock.IsInitialisedMock.Return(false)
	clientMock.LockMock.Return()
	clientMock.UnlockMock.Return()
	clientMock.CreateConnectionMock.Set(func(p configuration.ConnectionType, p1 string) (r error) {
		return errors.New("test reasons")
	})

	keyProcessor := platformpolicy.NewKeyProcessor()
	firstPrivateKey, _ := keyProcessor.GeneratePrivateKey()
	expectedNeighbourKey, _ := keyProcessor.ExportPublicKeyPEM(keyProcessor.ExtractPublicKey(firstPrivateKey))
	expectedNeighbour := &Neighbour{OutgoingClient: clientMock}
	pulsar.Neighbours[string(expectedNeighbourKey)] = expectedNeighbour

	err := pulsar.EstablishConnectionToPulsar(ctx, string(expectedNeighbourKey))

	require.Error(t, err, "test reasons")
	require.Equal(t, clientMock.LockCounter, uint64(1))
	require.Equal(t, clientMock.UnlockCounter, uint64(1))
}

func TestPulsar_EstablishConnection_IsNotInitialised_ProblemsWithRequest(t *testing.T) {
	ctx := inslogger.TestContext(t)

	keyProcessor := platformpolicy.NewKeyProcessor()
	privateKey, _ := keyProcessor.GeneratePrivateKey()
	cryptoService := cryptography.NewKeyBoundCryptographyService(privateKey)
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	pulsar := &Pulsar{
		Neighbours:                 map[string]*Neighbour{},
		CryptographyService:        cryptoService,
		KeyProcessor:               keyProcessor,
		PlatformCryptographyScheme: scheme,
		EntropyGenerator:           pulsartestutils.MockEntropyGenerator{},
	}

	mockClientWrapper := &pulsartestutils.CustomRPCWrapperMock{}

	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{Error: errors.New("oops, request is broken")}
	replyChan := &rpc.Call{Done: done}
	mockClientWrapper.Done = replyChan

	secondKeyProcessor := platformpolicy.NewKeyProcessor()
	secondPrivateKey, _ := secondKeyProcessor.GeneratePrivateKey()
	expectedNeighbourKey, _ := secondKeyProcessor.ExportPublicKeyPEM(secondKeyProcessor.ExtractPublicKey(secondPrivateKey))
	expectedNeighbour := &Neighbour{OutgoingClient: mockClientWrapper}
	pulsar.Neighbours[string(expectedNeighbourKey)] = expectedNeighbour

	err := pulsar.EstablishConnectionToPulsar(ctx, string(expectedNeighbourKey))

	require.Error(t, err, "oops, request is broken")
}

func TestPulsar_EstablishConnection_IsNotInitialised_SignatureFailed(t *testing.T) {
	ctx := inslogger.TestContext(t)

	keyProcessor := platformpolicy.NewKeyProcessor()
	privateKey, _ := keyProcessor.GeneratePrivateKey()
	cryptoService := cryptography.NewKeyBoundCryptographyService(privateKey)
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	pulsar := &Pulsar{
		Neighbours:                 map[string]*Neighbour{},
		CryptographyService:        cryptoService,
		KeyProcessor:               keyProcessor,
		PlatformCryptographyScheme: scheme,
		EntropyGenerator:           pulsartestutils.MockEntropyGenerator{},
	}

	mockClientWrapper := &pulsartestutils.CustomRPCWrapperMock{}

	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{Reply: &Payload{Body: &HandshakePayload{}}}
	replyChan := &rpc.Call{Done: done}
	mockClientWrapper.Done = replyChan

	secondKeyProcessor := platformpolicy.NewKeyProcessor()
	secondPrivateKey, _ := secondKeyProcessor.GeneratePrivateKey()
	expectedNeighbourKey, err := secondKeyProcessor.ExportPublicKeyPEM(secondKeyProcessor.ExtractPublicKey(secondPrivateKey))
	require.NoError(t, err)
	expectedNeighbour := &Neighbour{OutgoingClient: mockClientWrapper}
	pulsar.Neighbours[string(expectedNeighbourKey)] = expectedNeighbour
	err = pulsar.EstablishConnectionToPulsar(ctx, string(expectedNeighbourKey))

	require.Error(t, err, "Signature check Failed")
}

func TestPulsar_EstablishConnection_IsNotInitialised_Success(t *testing.T) {
	ctx := inslogger.TestContext(t)

	keyProcessor := platformpolicy.NewKeyProcessor()
	privateKey, _ := keyProcessor.GeneratePrivateKey()
	cryptoService := cryptography.NewKeyBoundCryptographyService(privateKey)
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	pulsar := &Pulsar{
		Neighbours:                 map[string]*Neighbour{},
		CryptographyService:        cryptoService,
		KeyProcessor:               keyProcessor,
		PlatformCryptographyScheme: scheme,
		EntropyGenerator:           pulsartestutils.MockEntropyGenerator{},
	}

	secondKeyProcessor := platformpolicy.NewKeyProcessor()
	secondPrivateKey, err := secondKeyProcessor.GeneratePrivateKey()
	require.NoError(t, err)
	secondCryptoService := cryptography.NewKeyBoundCryptographyService(secondPrivateKey)
	expectedNeighbourKey, err := secondKeyProcessor.ExportPublicKeyPEM(secondKeyProcessor.ExtractPublicKey(secondPrivateKey))
	require.NoError(t, err)

	secondScheme := platformpolicy.NewPlatformCryptographyScheme()

	secondPulsar := &Pulsar{
		CryptographyService:        secondCryptoService,
		KeyProcessor:               secondKeyProcessor,
		PlatformCryptographyScheme: secondScheme,
		PublicKeyRaw:               string(expectedNeighbourKey),
	}

	payload, err := secondPulsar.preparePayload(&HandshakePayload{Entropy: pulsartestutils.MockEntropy})
	require.NoError(t, err)

	mockClientWrapper := &pulsartestutils.CustomRPCWrapperMock{}

	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{Reply: payload}
	replyChan := &rpc.Call{Done: done}
	mockClientWrapper.Done = replyChan

	expectedNeighbour := &Neighbour{OutgoingClient: mockClientWrapper}
	pulsar.Neighbours[string(expectedNeighbourKey)] = expectedNeighbour

	err = pulsar.EstablishConnectionToPulsar(ctx, string(expectedNeighbourKey))

	require.NoError(t, err)
}

func TestPulsar_CheckConnectionsToPulsars_NoProblems(t *testing.T) {
	ctx := inslogger.TestContext(t)

	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{}
	replyChan := &rpc.Call{Done: done}

	clientMock := NewRPCClientWrapperMock(t)
	clientMock.IsInitialisedMock.Return(true)
	clientMock.GoMock.Set(func(p string, p1 interface{}, p2 interface{}, p3 chan *rpc.Call) (r *rpc.Call) {
		return replyChan
	})

	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}}
	keyProcessor := platformpolicy.NewKeyProcessor()
	firstNeighbourPrivateKey, _ := keyProcessor.GeneratePrivateKey()
	firstNeighbourPublicKey := keyProcessor.ExtractPublicKey(firstNeighbourPrivateKey)
	firstNeighbourExpectedKey, _ := keyProcessor.ExportPublicKeyPEM(firstNeighbourPublicKey)
	pulsar.Neighbours[string(firstNeighbourExpectedKey)] = &Neighbour{
		PublicKey:      firstNeighbourPublicKey,
		OutgoingClient: clientMock,
	}

	pulsar.CheckConnectionsToPulsars(ctx)

	require.Equal(t, uint64(1), clientMock.IsInitialisedCounter)
	require.Equal(t, uint64(1), clientMock.GoCounter)
}

func TestPulsar_CheckConnectionsToPulsars_NilClient_FirstConnectionFailed(t *testing.T) {
	ctx := inslogger.TestContext(t)

	clientMock := NewRPCClientWrapperMock(t)
	clientMock.IsInitialisedMock.Return(false)
	clientMock.LockMock.Return()
	clientMock.UnlockMock.Return()
	clientMock.CreateConnectionFunc = func(p configuration.ConnectionType, p1 string) (r error) {
		return errors.New("this will have to fall")
	}

	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}}

	keyProcessor := platformpolicy.NewKeyProcessor()

	firstNeighbourPrivateKey, _ := keyProcessor.GeneratePrivateKey()
	pulsar.Neighbours["thisShouldFailEstablishConnection"] = &Neighbour{
		PublicKey:      keyProcessor.ExtractPublicKey(firstNeighbourPrivateKey),
		OutgoingClient: clientMock,
	}

	resultLog := capture(func() {
		pulsar.CheckConnectionsToPulsars(ctx)
	})

	require.Contains(t, resultLog, "this will have to fall")
}

func TestPulsar_CheckConnectionsToPulsars_NilClient_SecondConnectionFailed(t *testing.T) {
	ctx := inslogger.TestContext(t)

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

	keyProcessor := platformpolicy.NewKeyProcessor()
	firstNeighbourPrivateKey, _ := keyProcessor.GeneratePrivateKey()
	pulsar.Neighbours["this should fail second connection"] = &Neighbour{
		PublicKey:         keyProcessor.ExtractPublicKey(firstNeighbourPrivateKey),
		OutgoingClient:    mockClientWrapper,
		ConnectionAddress: "TestConnectionAddress",
	}

	resultLog := capture(func() {
		pulsar.CheckConnectionsToPulsars(ctx)
	})

	require.Contains(t, resultLog, "Problems with connection to TestConnectionAddress, with error - test error")
	require.Contains(t, resultLog, "Attempt of connection to TestConnectionAddress Failed with error - this should Failed")
	require.Equal(t, 2, isInitTimeCalled)
}

func TestPulsar_StartConsensusProcess_WithWrongPulseNumber(t *testing.T) {
	ctx := inslogger.TestContext(t)

	switcherMock := NewStateSwitcherMock(t)
	switcherMock.GetStateMock.Return(WaitingForStart)
	pulsar := &Pulsar{
		StateSwitcher:         switcherMock,
		ProcessingPulseNumber: core.PulseNumber(123),
		lastPulse:             &core.Pulse{PulseNumber: core.PulseNumber(122)},
	}
	err := pulsar.StartConsensusProcess(ctx, core.PulseNumber(121))

	require.Error(t, err, "wrong state status or pulse number, state - WaitingForStart, received pulse - 121, last pulse - 122, processing pulse - 123")
}

func TestPulsar_StartConsensusProcess_Success(t *testing.T) {
	ctx := inslogger.TestContext(t)

	mockSwitcher := NewStateSwitcherMock(t)
	mockSwitcher.GetStateMock.Return(WaitingForStart)
	mockSwitcher.setStateMock.Expect(GenerateEntropy).Return()
	mockSwitcher.SwitchToStateFunc = func(ctx context.Context, p State, p1 interface{}) {
		if p != WaitingForEntropySigns {
			t.Error(t, "Wrong state")
		}
	}

	privateKey, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	cryptoService := cryptography.NewKeyBoundCryptographyService(privateKey)
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	pulsar := &Pulsar{
		EntropyGenerator:           pulsartestutils.MockEntropyGenerator{},
		CryptographyService:        cryptoService,
		PlatformCryptographyScheme: scheme,
		ownedBftRow:                map[string]*BftCell{},
	}
	pulsar.ProcessingPulseNumber = core.PulseNumber(120)
	pulsar.SetLastPulse(&core.Pulse{PulseNumber: core.PulseNumber(2)})
	pulsar.StateSwitcher = mockSwitcher
	expectedPulse := core.PulseNumber(123)

	err := pulsar.StartConsensusProcess(ctx, expectedPulse)

	require.NoError(t, err)
	require.Equal(t, pulsar.ProcessingPulseNumber, expectedPulse)
	require.Equal(t, uint64(1), mockSwitcher.SwitchToStateCounter)
}

func TestPulsar_broadcastSignatureOfEntropy_StateFailed(t *testing.T) {
	ctx := inslogger.TestContext(t)
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

	pulsar.broadcastSignatureOfEntropy(ctx)
}

func TestPulsar_broadcastSignatureOfEntropy_SendToNeighbours(t *testing.T) {
	ctx := inslogger.TestContext(t)
	done := make(chan *rpc.Call, 2)
	done <- &rpc.Call{}
	done <- &rpc.Call{Error: errors.New("Failed")}
	replyChan := &rpc.Call{Done: done}

	mockClientWrapper := &pulsartestutils.MockRPCClientWrapper{}
	switcher := NewStateSwitcherMock(t)
	switcher.GetStateMock.Return(WaitingForStart)

	privateKey, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	cryptoService := cryptography.NewKeyBoundCryptographyService(privateKey)
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	pulsar := Pulsar{
		Neighbours: map[string]*Neighbour{
			"1": {OutgoingClient: mockClientWrapper, ConnectionAddress: "first"},
			"2": {OutgoingClient: mockClientWrapper, ConnectionAddress: "second"},
		},
		CryptographyService:        cryptoService,
		StateSwitcher:              switcher,
		ProcessingPulseNumber:      123,
		GeneratedEntropySign:       pulsartestutils.MockEntropy[:],
		PlatformCryptographyScheme: scheme,
	}

	mockClientWrapper.On("Go", ReceiveSignatureForEntropy.String(), mock.Anything, nil, (chan *rpc.Call)(nil)).Return(replyChan)

	resultLog := capture(func() {
		pulsar.broadcastSignatureOfEntropy(ctx)
	})

	mockClientWrapper.AssertCalled(t, "Go", ReceiveSignatureForEntropy.String(), mock.Anything, mock.Anything, mock.Anything)
	mockClientWrapper.AssertNumberOfCalls(t, "Go", 2)
	require.Equal(t, 0, len(done))
	require.Contains(t, resultLog, "finished with error - Failed")
}

func TestPulsar_broadcastVector_StateFailed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mockClientWrapper := &pulsartestutils.MockRPCClientWrapper{}
	mockClientWrapper.On("Go", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	mockSwitcher := NewStateSwitcherMock(t)
	mockSwitcher.GetStateMock.Return(Failed)
	pulsar := Pulsar{
		Neighbours: map[string]*Neighbour{
			"1": {},
		},
		ownedBftRow:           map[string]*BftCell{},
		StateSwitcher:         mockSwitcher,
		ProcessingPulseNumber: core.PulseNumber(123),
	}

	pulsar.broadcastVector(ctx)

	mockClientWrapper.AssertNotCalled(t, "Go", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestPulsar_broadcastVector_SendToNeighbours(t *testing.T) {
	ctx := inslogger.TestContext(t)
	done := make(chan *rpc.Call, 2)
	done <- &rpc.Call{}
	done <- &rpc.Call{Error: errors.New("Failed")}
	replyChan := &rpc.Call{Done: done}

	mockClientWrapper := &pulsartestutils.MockRPCClientWrapper{}
	mockSwitcher := NewStateSwitcherMock(t)
	mockSwitcher.GetStateMock.Return(WaitingForStart)

	privateKey, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	cryptoService := cryptography.NewKeyBoundCryptographyService(privateKey)
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	generatedEntropy := core.Entropy(pulsartestutils.MockEntropy)
	pulsar := Pulsar{
		Neighbours: map[string]*Neighbour{
			"1": {OutgoingClient: mockClientWrapper, ConnectionAddress: "first"},
			"2": {OutgoingClient: mockClientWrapper, ConnectionAddress: "second"},
		},
		CryptographyService:        cryptoService,
		ownedBftRow:                map[string]*BftCell{},
		StateSwitcher:              mockSwitcher,
		GeneratedEntropySign:       pulsartestutils.MockEntropy[:],
		generatedEntropy:           &generatedEntropy,
		PlatformCryptographyScheme: scheme,
	}
	mockClientWrapper.On("Go", ReceiveVector.String(), mock.Anything, nil, (chan *rpc.Call)(nil)).Return(replyChan)

	resultLog := capture(func() {
		pulsar.broadcastVector(ctx)
	})

	mockClientWrapper.AssertCalled(t, "Go", ReceiveVector.String(), mock.Anything, mock.Anything, mock.Anything)
	mockClientWrapper.AssertNumberOfCalls(t, "Go", 2)
	require.Equal(t, 0, len(done))
	require.Contains(t, resultLog, "finished with error - Failed")
}

func TestPulsar_broadcastEntropy_StateFailed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mockClientWrapper := &pulsartestutils.MockRPCClientWrapper{}
	mockClientWrapper.On("Go", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	mockSwitcher := NewStateSwitcherMock(t)
	mockSwitcher.GetStateMock.Return(Failed)
	pulsar := Pulsar{
		Neighbours: map[string]*Neighbour{
			"1": {},
		},
		ownedBftRow:           map[string]*BftCell{},
		StateSwitcher:         mockSwitcher,
		ProcessingPulseNumber: core.PulseNumber(123),
	}

	pulsar.broadcastEntropy(ctx)

	mockClientWrapper.AssertNotCalled(t, "Go", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestPulsar_broadcastEntropy_SendToNeighbours(t *testing.T) {
	ctx := inslogger.TestContext(t)
	done := make(chan *rpc.Call, 2)
	done <- &rpc.Call{}
	done <- &rpc.Call{Error: errors.New("Failed")}
	replyChan := &rpc.Call{Done: done}

	mockClientWrapper := &pulsartestutils.MockRPCClientWrapper{}
	switcherMock := NewStateSwitcherMock(t)
	switcherMock.GetStateMock.Return(WaitingForStart)

	privateKey, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	cryptoService := cryptography.NewKeyBoundCryptographyService(privateKey)
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	generatedEntropy := core.Entropy(pulsartestutils.MockEntropy)
	pulsar := Pulsar{
		Neighbours: map[string]*Neighbour{
			"1": {OutgoingClient: mockClientWrapper, ConnectionAddress: "first"},
			"2": {OutgoingClient: mockClientWrapper, ConnectionAddress: "second"},
		},
		CryptographyService:        cryptoService,
		StateSwitcher:              switcherMock,
		ProcessingPulseNumber:      123,
		PlatformCryptographyScheme: scheme,
		generatedEntropy:           &generatedEntropy,
	}

	mockClientWrapper.On("Go", ReceiveEntropy.String(), mock.Anything, nil, (chan *rpc.Call)(nil)).Return(replyChan)

	resultLog := capture(func() {
		pulsar.broadcastEntropy(ctx)
	})

	mockClientWrapper.AssertCalled(t, "Go", ReceiveEntropy.String(), mock.Anything, mock.Anything, mock.Anything)
	mockClientWrapper.AssertNumberOfCalls(t, "Go", 2)
	require.Equal(t, 0, len(done))
	require.Contains(t, resultLog, "finished with error - Failed")
}

func TestPulsar_sendVector_StateFailed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	pulsar := Pulsar{}
	switcher := NewStateSwitcherMock(t)
	switcher.GetStateMock.Return(Failed)
	pulsar.StateSwitcher = switcher

	pulsar.sendVector(ctx)

	require.Equal(t, uint64(0), switcher.SwitchToStateCounter)
}

func TestPulsar_sendVector_OnePulsar(t *testing.T) {
	ctx := inslogger.TestContext(t)
	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}}
	switcher := NewStateSwitcherMock(t)
	switcher.GetStateMock.Return(WaitingForStart)
	switcher.SwitchToStateFunc = func(ctx context.Context, p State, p1 interface{}) {
		require.Equal(t, p, Verifying)
	}
	pulsar.StateSwitcher = switcher

	pulsar.sendVector(ctx)

	require.Equal(t, switcher.SwitchToStateCounter, uint64(1))
}

func TestPulsar_sendVector_TwoPulsars(t *testing.T) {
	// Arrange
	ctx := inslogger.TestContext(t)
	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{}
	replyChan := &rpc.Call{Done: done}

	mockClientWrapper := &pulsartestutils.MockRPCClientWrapper{}
	privateKey, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	cryptoService := cryptography.NewKeyBoundCryptographyService(privateKey)
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	pulsar := Pulsar{
		Neighbours:                 map[string]*Neighbour{},
		CryptographyService:        cryptoService,
		PlatformCryptographyScheme: scheme,
		ownedBftRow:                map[string]*BftCell{},
		bftGrid:                    map[string]map[string]*BftCell{},
	}

	generatedEntropy := core.Entropy(pulsartestutils.MockEntropy)
	pulsar.generatedEntropy = &generatedEntropy
	pulsar.Neighbours["1"] = &Neighbour{OutgoingClient: mockClientWrapper, ConnectionAddress: "first"}
	mockClientWrapper.On("Go", ReceiveVector.String(), mock.Anything, nil, (chan *rpc.Call)(nil)).Return(replyChan)

	switcher := NewStateSwitcherMock(t)
	switcher.GetStateMock.Return(WaitingForStart)
	switcher.SwitchToStateMock.Expect(ctx, WaitingForVectors, nil).Return()
	pulsar.StateSwitcher = switcher

	// Act
	pulsar.sendVector(ctx)

	// Assert
	require.Equal(t, uint64(1), switcher.SwitchToStateCounter)
}

func TestPulsar_sendEntropy_OnePulsar(t *testing.T) {
	ctx := inslogger.TestContext(t)
	pulsar := Pulsar{Neighbours: map[string]*Neighbour{}}
	switcher := NewStateSwitcherMock(t)
	switcher.GetStateMock.Return(WaitingForStart)
	switcher.SwitchToStateMock.Expect(ctx, Verifying, nil).Return()
	pulsar.StateSwitcher = switcher

	pulsar.sendEntropy(ctx)

	require.Equal(t, uint64(1), switcher.SwitchToStateCounter)
}

func TestPulsar_sendEntropy_TwoPulsars(t *testing.T) {
	// Arrange
	ctx := inslogger.TestContext(t)
	done := make(chan *rpc.Call, 1)
	done <- &rpc.Call{}
	replyChan := &rpc.Call{Done: done}

	mockClientWrapper := &pulsartestutils.MockRPCClientWrapper{}
	privateKey, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	cryptoService := cryptography.NewKeyBoundCryptographyService(privateKey)
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	pulsar := Pulsar{
		Neighbours:                 map[string]*Neighbour{},
		CryptographyService:        cryptoService,
		ownedBftRow:                map[string]*BftCell{},
		PlatformCryptographyScheme: scheme,
	}
	generatedEntropy := core.Entropy(pulsartestutils.MockEntropy)
	pulsar.generatedEntropy = &generatedEntropy
	pulsar.Neighbours["1"] = &Neighbour{OutgoingClient: mockClientWrapper, ConnectionAddress: "first"}
	mockClientWrapper.On("Go", ReceiveEntropy.String(), mock.Anything, nil, (chan *rpc.Call)(nil)).Return(replyChan)

	switcher := NewStateSwitcherMock(t)
	switcher.GetStateMock.Return(WaitingForStart)
	switcher.SwitchToStateMock.Expect(ctx, WaitingForEntropy, nil).Return()
	pulsar.StateSwitcher = switcher

	// Act
	pulsar.sendEntropy(ctx)

	// Assert
	require.Equal(t, uint64(1), switcher.SwitchToStateCounter)
}

func TestPulsar_verify_failedState(t *testing.T) {
	ctx := inslogger.TestContext(t)
	switcherMock := NewStateSwitcherMock(t)
	switcherMock.GetStateMock.Return(Failed)
	pulsar := &Pulsar{StateSwitcher: switcherMock}
	pulsar.PublicKeyRaw = "testKey"
	pulsar.ownedBftRow = map[string]*BftCell{}
	pulsar.bftGrid = map[string]map[string]*BftCell{}

	pulsar.verify(ctx)
}

func TestPulsar_verify_Standalone_Success(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mockSwitcher := NewStateSwitcherMock(t)
	mockSwitcher.GetStateMock.Return(Verifying)
	mockSwitcher.SwitchToStateMock.Expect(ctx, SendingPulse, nil).Return()
	pulsar := &Pulsar{StateSwitcher: mockSwitcher}
	pulsar.PublicKeyRaw = "testKey"
	generatedEntropy := core.Entropy(pulsartestutils.MockEntropy)
	pulsar.generatedEntropy = &generatedEntropy
	pulsar.ownedBftRow = map[string]*BftCell{}
	pulsar.bftGrid = map[string]map[string]*BftCell{}

	pulsar.verify(ctx)

	require.Equal(t, uint64(1), mockSwitcher.SwitchToStateCounter)
	require.Equal(t, "testKey", pulsar.PublicKeyRaw)
	require.Equal(t, core.Entropy(pulsartestutils.MockEntropy), *pulsar.GetGeneratedEntropy())
}

func TestPulsar_verify_NotEnoughForConsensus_Success(t *testing.T) {
	ctx := inslogger.TestContext(t)

	mockSwitcher := NewStateSwitcherMock(t)
	mockSwitcher.SwitchToStateMock.Set(func(ctx context.Context, p State, p1 interface{}) {
		require.Equal(t, Failed, p)
	})
	mockSwitcher.GetStateMock.Return(Verifying)
	privateKey, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	cryptoService := cryptography.NewKeyBoundCryptographyService(privateKey)

	pulsar := &Pulsar{
		StateSwitcher:       mockSwitcher,
		PublicKeyRaw:        "testKey",
		CryptographyService: cryptoService,
		ownedBftRow:         map[string]*BftCell{},
		bftGrid:             map[string]map[string]*BftCell{},
		Neighbours: map[string]*Neighbour{
			"1": {},
			"2": {},
		},
	}

	pulsar.verify(ctx)

	require.Equal(t, uint64(1), mockSwitcher.SwitchToStateCounter)
}

func prepareEntropy(t *testing.T, service core.CryptographyService) (entropy core.Entropy, sign []byte) {
	entropy = (&entropygenerator.StandardEntropyGenerator{}).GenerateEntropy()
	fetchedSign, err := service.Sign(entropy[:])
	require.NoError(t, err)
	sign = fetchedSign.Bytes()
	return
}

func TestPulsar_verify_Success(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mockSwitcher := NewStateSwitcherMock(t)
	mockSwitcher.SwitchToStateFunc = func(ctx context.Context, p State, p1 interface{}) {
		if p != WaitingForPulseSigns && p != SendingPulseSign {
			t.Error("Wrong state")
		}
	}
	mockSwitcher.GetStateFunc = func() State {
		return Verifying
	}

	keyGenerator := func() (crypto.PrivateKey, crypto.PublicKey, string) {
		kp := platformpolicy.NewKeyProcessor()
		key, _ := kp.GeneratePrivateKey()
		publicKey := kp.ExtractPublicKey(key)
		pubKeyString, _ := kp.ExportPublicKeyPEM(publicKey)

		return key, publicKey, string(pubKeyString)
	}

	pulsarPrivate, _, currentPulsarPublicKey := keyGenerator()

	secondPrivate, pub2, publicKeySecond := keyGenerator()
	thirdPrivate, pub3, publicKeyThird := keyGenerator()

	clientMock := pulsartestutils.MockRPCClientWrapper{}
	clientMock.On("IsInitialised").Return(true)

	pulsarCryptoService := cryptography.NewKeyBoundCryptographyService(pulsarPrivate)
	secondCryptoService := cryptography.NewKeyBoundCryptographyService(secondPrivate)
	thirdCryptoService := cryptography.NewKeyBoundCryptographyService(thirdPrivate)

	pulsar := &Pulsar{
		KeyProcessor:               platformpolicy.NewKeyProcessor(),
		StateSwitcher:              mockSwitcher,
		CryptographyService:        pulsarCryptoService,
		PlatformCryptographyScheme: platformpolicy.NewPlatformCryptographyScheme(),
		PublicKeyRaw:               currentPulsarPublicKey,
		ownedBftRow:                map[string]*BftCell{},
		bftGrid:                    map[string]map[string]*BftCell{},
		CurrentSlotSenderConfirmations: map[string]core.PulseSenderConfirmation{},
		Neighbours: map[string]*Neighbour{
			publicKeySecond: {PublicKey: pub2, OutgoingClient: &clientMock},
			publicKeyThird:  {PublicKey: pub3, OutgoingClient: &clientMock},
		},
	}

	firstEntropy, firstSign := prepareEntropy(t, pulsarCryptoService)
	secondEntropy, secondSign := prepareEntropy(t, secondCryptoService)
	thirdEntropy, thirdSign := prepareEntropy(t, thirdCryptoService)

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

	pulsar.verify(ctx)

	require.NotNil(t, pulsar.CurrentSlotPulseSender)
	require.Equal(t, expectedEntropy, *pulsar.GetCurrentSlotEntropy())
	require.Equal(t, uint64(1), mockSwitcher.SwitchToStateCounter)
}
