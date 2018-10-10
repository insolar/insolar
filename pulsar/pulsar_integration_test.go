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
	"os"
	"testing"
	"time"

	ecdsa_helper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/ledger/ledgertestutil"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/pulsar/pulsartestutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTwoPulsars_Handshake(t *testing.T) {
	firstKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	firstPublic, err := ecdsa_helper.ExportPublicKey(&firstKey.PublicKey)
	assert.NoError(t, err)
	firstPublicExported, err := ecdsa_helper.ExportPrivateKey(firstKey)
	assert.NoError(t, err)

	secondKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	secondPublic, err := ecdsa_helper.ExportPublicKey(&secondKey.PublicKey)
	assert.NoError(t, err)
	secondPublicExported, err := ecdsa_helper.ExportPrivateKey(secondKey)
	assert.NoError(t, err)

	storage := &pulsartestutil.MockStorage{}
	storage.On("GetLastPulse", mock.Anything).Return(&core.Pulse{PulseNumber: 123}, nil)
	firstPulsar, err := NewPulsar(configuration.Pulsar{
		ConnectionType:      "tcp",
		MainListenerAddress: ":1639",
		PrivateKey:          firstPublicExported,
		Neighbours: []configuration.PulsarNodeAddress{
			{ConnectionType: "tcp", Address: "127.0.0.1:1639", PublicKey: firstPublic},
			{ConnectionType: "tcp", Address: "127.0.0.1:1640", PublicKey: secondPublic},
			{ConnectionType: "tcp", Address: "127.0.0.1:1641"},
		}},
		storage,
		&RPCClientWrapperFactoryImpl{},
		pulsartestutil.MockEntropyGenerator{},
		nil,
		net.Listen,
	)
	assert.NoError(t, err)

	secondPulsar, err := NewPulsar(configuration.Pulsar{
		ConnectionType:      "tcp",
		MainListenerAddress: ":1640",
		PrivateKey:          secondPublicExported,
		Neighbours: []configuration.PulsarNodeAddress{
			{ConnectionType: "tcp", Address: "127.0.0.1:1639", PublicKey: firstPublic},
			{ConnectionType: "tcp", Address: "127.0.0.1:1640", PublicKey: secondPublic},
			{ConnectionType: "tcp", Address: "127.0.0.1:1641"},
		}},
		storage,
		&RPCClientWrapperFactoryImpl{},
		pulsartestutil.MockEntropyGenerator{},
		nil,
		net.Listen,
	)
	assert.NoError(t, err)

	go firstPulsar.StartServer()
	go secondPulsar.StartServer()
	err = secondPulsar.EstablishConnectionToPulsar(firstPublic)

	assert.NoError(t, err)
	assert.NotNil(t, firstPulsar.Neighbours[secondPublic].OutgoingClient)
	assert.NotNil(t, secondPulsar.Neighbours[firstPublic].OutgoingClient)

	defer func() {
		firstPulsar.StopServer()
		secondPulsar.StopServer()
	}()
}

func TestOnePulsar_FullStatesTransition(t *testing.T) {
	t.Skip("should be re-written after refactoring the body of pulsar")
	firstKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	firstPublicExported, err := ecdsa_helper.ExportPrivateKey(firstKey)
	assert.NoError(t, err)

	storage := &pulsartestutil.MockStorage{}
	firstPulse := 123
	storage.On("GetLastPulse", mock.Anything).Return(&core.Pulse{PulseNumber: core.PulseNumber(firstPulse)}, nil)
	pulsar, err := NewPulsar(configuration.Pulsar{
		ConnectionType:         "tcp",
		MainListenerAddress:    ":1639",
		PrivateKey:             firstPublicExported,
		Neighbours:             []configuration.PulsarNodeAddress{},
		PulseTime:              10000,
		ReceivingSignTimeout:   1000,
		ReceivingNumberTimeout: 1000,
		ReceivingVectorTimeout: 1000},
		storage,

		&RPCClientWrapperFactoryImpl{},
		pulsartestutil.MockEntropyGenerator{},
		nil,
		net.Listen,
	)
	assert.NoError(t, err)

	pulsar.StartConsensusProcess(core.PulseNumber(firstPulse + 1))

	for pulsar.stateSwitcher.getState() != sendingPulse {
		time.Sleep(1 * time.Millisecond)
	}

	assert.NoError(t, err)

	defer pulsar.StopServer()
}

func TestTwoPulsars_Full_Consensus(t *testing.T) {
	t.Skip("should be re-written after refactoring the body of pulsar")
	firstKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	firstPublic, err := ecdsa_helper.ExportPublicKey(&firstKey.PublicKey)
	assert.NoError(t, err)
	firstPublicExported, err := ecdsa_helper.ExportPrivateKey(firstKey)
	assert.NoError(t, err)

	secondKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	secondPublic, err := ecdsa_helper.ExportPublicKey(&secondKey.PublicKey)
	assert.NoError(t, err)
	secondPublicExported, err := ecdsa_helper.ExportPrivateKey(secondKey)
	assert.NoError(t, err)

	storage := &pulsartestutil.MockStorage{}
	firstPulse := 123
	storage.On("GetLastPulse", mock.Anything).Return(&core.Pulse{PulseNumber: core.PulseNumber(firstPulse)}, nil)
	firstPulsar, err := NewPulsar(configuration.Pulsar{
		ConnectionType:      "tcp",
		MainListenerAddress: ":1639",
		PrivateKey:          firstPublicExported,
		Neighbours: []configuration.PulsarNodeAddress{
			{ConnectionType: "tcp", Address: "127.0.0.1:1640", PublicKey: secondPublic},
		}},
		storage,
		&RPCClientWrapperFactoryImpl{},
		pulsartestutil.MockEntropyGenerator{},
		nil,
		net.Listen,
	)
	assert.NoError(t, err)

	secondPulsar, err := NewPulsar(configuration.Pulsar{
		ConnectionType:      "tcp",
		MainListenerAddress: ":1640",
		PrivateKey:          secondPublicExported,
		Neighbours: []configuration.PulsarNodeAddress{
			{ConnectionType: "tcp", Address: "127.0.0.1:1639", PublicKey: firstPublic},
		}},
		storage,
		&RPCClientWrapperFactoryImpl{},
		pulsartestutil.MockEntropyGenerator{},
		nil,
		net.Listen,
	)
	assert.NoError(t, err)

	go firstPulsar.StartServer()
	go secondPulsar.StartServer()
	err = secondPulsar.EstablishConnectionToPulsar(firstPublic)
	assert.NoError(t, err)
	assert.NotNil(t, firstPulsar.Neighbours[secondPublic].OutgoingClient)
	assert.NotNil(t, secondPulsar.Neighbours[firstPublic].OutgoingClient)

	firstPulsar.StartConsensusProcess(core.PulseNumber(firstPulse + 1))

	for len(secondPulsar.OwnedBftRow) != 1 {
		time.Sleep(1 * time.Millisecond)
	}

	defer func() {
		firstPulsar.StopServer()
		secondPulsar.StopServer()
	}()
}

func TestPulsar_ConnectToNode(t *testing.T) {
	lr, err := logicrunner.NewLogicRunner(&configuration.LogicRunner{
		BuiltIn: &configuration.BuiltIn{},
	})
	assert.NoError(t, err)

	os.MkdirAll("bootstrapLedger", os.ModePerm)
	bootstrapLedger, bootstrapLedgerCleaner := ledgertestutil.TmpLedger(t, lr, "bootstrapLedger")
	bootstrapNodeConfig := configuration.NewConfiguration()
	bootstrapNodeNetwork, err := servicenetwork.NewServiceNetwork(bootstrapNodeConfig.Host, bootstrapNodeConfig.Node)
	assert.NoError(t, err)
	err = bootstrapNodeNetwork.Start(core.Components{Ledger: bootstrapLedger})
	assert.NoError(t, err)
	bootstrapAddress := bootstrapNodeNetwork.GetAddress()

	os.MkdirAll("usualLedger", os.ModePerm)
	usualLedger, usualLedgerCleaner := ledgertestutil.TmpLedger(t, lr, "usualLedger")
	usualNodeConfig := configuration.NewConfiguration()
	usualNodeConfig.Host.BootstrapHosts = []string{bootstrapAddress}
	usualNodeNetwork, err := servicenetwork.NewServiceNetwork(usualNodeConfig.Host, usualNodeConfig.Node)
	assert.NoError(t, err)
	err = usualNodeNetwork.Start(core.Components{Ledger: usualLedger})
	assert.NoError(t, err)

	pulsarPrivateKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)
	firstPublicExported, err := ecdsa_helper.ExportPrivateKey(pulsarPrivateKey)
	assert.NoError(t, err)
	storage := &pulsartestutil.MockStorage{}
	storage.On("GetLastPulse").Return(core.GenesisPulse, nil)

	stateSwitcher := &StateSwitcherImpl{}
	newPulsar, err := NewPulsar(configuration.Pulsar{
		ConnectionType:      "tcp",
		MainListenerAddress: ":1640",
		PrivateKey:          firstPublicExported,
		BootstrapNodes:      []string{bootstrapAddress},
		BootstrapListener:   configuration.Transport{Protocol: "UTP", Address: "127.0.0.1:18091", BehindNAT: false},
		Neighbours:          []configuration.PulsarNodeAddress{}},
		storage,
		&RPCClientWrapperFactoryImpl{},
		pulsartestutil.MockEntropyGenerator{},
		stateSwitcher,
		net.Listen,
	)
	stateSwitcher.SetPulsar(newPulsar)
	newPulsar.StartConsensusProcess(core.GenesisPulse.PulseNumber + 1)

	time.Sleep(100 * time.Millisecond)
	usualNodeNetwork.Stop()
	bootstrapNodeNetwork.Stop()
	newPulsar.StopServer()
	bootstrapLedgerCleaner()

	currentPulse, err := usualLedger.GetPulseManager().Current()
	assert.NoError(t, err)
	assert.Equal(t, currentPulse.PulseNumber, core.GenesisPulse.PulseNumber+1)

	defer func() {
		usualLedgerCleaner()
		err = os.RemoveAll("bootstrapLedger")
		if err != nil {
			assert.NoError(t, err)
		}
		err = os.RemoveAll("usualLedger")
		if err != nil {
			assert.NoError(t, err)
		}
	}()
}
