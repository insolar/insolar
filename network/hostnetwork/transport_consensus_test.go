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
	"context"
	"sync"
	"testing"
	"time"

	consensus "github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/utils"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type consensusTransportSuite struct {
	suite.Suite
	crypto core.CryptographyService
}

func createTwoConsensusNetworks(id1, id2 core.ShortNodeID) (t1, t2 network.ConsensusNetwork, err error) {
	m := newMockResolver()

	cn1, err := NewConsensusNetwork("127.0.0.1:0", ID1+DOMAIN, id1, m)
	if err != nil {
		return nil, nil, err
	}
	cn2, err := NewConsensusNetwork("127.0.0.1:0", ID2+DOMAIN, id2, m)
	if err != nil {
		return nil, nil, err
	}

	ref1, err := core.NewRefFromBase58(ID2 + DOMAIN)
	if err != nil {
		return nil, nil, err
	}
	routing1, err := host.NewHostNS(cn1.PublicAddress(), *ref1, id1)
	if err != nil {
		return nil, nil, err
	}
	ref2, err := core.NewRefFromBase58(ID2 + DOMAIN)
	if err != nil {
		return nil, nil, err
	}
	routing2, err := host.NewHostNS(cn2.PublicAddress(), *ref2, id2)
	if err != nil {
		return nil, nil, err
	}
	m.addMappingHost(routing1)
	m.addMappingHost(routing2)

	return cn1, cn2, nil
}

func (t *consensusTransportSuite) sendPacket(packet consensus.ConsensusPacket) (bool, error) {
	cn1, cn2, err := createTwoConsensusNetworks(0, 1)
	if err != nil {
		return false, err
	}
	ctx := context.Background()
	ctx2 := context.Background()

	wg := sync.WaitGroup{}
	wg.Add(1)

	handler := func(incomingPacket consensus.ConsensusPacket, sender core.RecordRef) {
		log.Info("handler triggered")
		wg.Done()
	}
	cn2.RegisterPacketHandler(packet.GetType(), handler)

	cn2.Start(ctx)
	cn1.Start(ctx2)
	defer func() {
		cn1.Stop()
		cn2.Stop()
	}()

	err = cn1.SignAndSendPacket(packet, cn2.GetNodeID(), t.crypto)
	if err != nil {
		return false, err
	}
	return utils.WaitTimeout(&wg, time.Second), nil
}

func (t *consensusTransportSuite) TestSendPacketPhase1() {
	packet := consensus.NewPhase1Packet()
	success, err := t.sendPacket(packet)
	require.NoError(t.T(), err)
	assert.True(t.T(), success)
}

func (t *consensusTransportSuite) TestSendPacketPhase2() {
	var ghs consensus.GlobuleHashSignature
	bitset, err := consensus.NewBitSet(10)
	require.NoError(t.T(), err)
	packet := consensus.NewPhase2Packet(ghs, bitset)
	success, err := t.sendPacket(packet)
	require.NoError(t.T(), err)
	assert.True(t.T(), success)
}

func (t *consensusTransportSuite) TestSendPacketPhase3() {
	var ghs consensus.GlobuleHashSignature
	bitset, err := consensus.NewBitSet(10)
	require.NoError(t.T(), err)
	packet := consensus.NewPhase3Packet(ghs, bitset)
	success, err := t.sendPacket(packet)
	require.NoError(t.T(), err)
	assert.True(t.T(), success)
}

func (t *consensusTransportSuite) TestRegisterPacketHandler() {
	m := newMockResolver()

	cn, err := NewConsensusNetwork("127.0.0.1:0", ID1+DOMAIN, 0, m)
	require.NoError(t.T(), err)
	defer cn.Stop()
	handler := func(incomingPacket consensus.ConsensusPacket, sender core.RecordRef) {
		// do nothing
	}
	f := func() {
		cn.RegisterPacketHandler(consensus.Phase1, handler)
	}
	assert.NotPanics(t.T(), f, "first request handler register should not panic")
	assert.Panics(t.T(), f, "second request handler register should panic because it is already registered")
}

func NewSuite() (*consensusTransportSuite, error) {
	kp := platformpolicy.NewKeyProcessor()
	sk, err := kp.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	cryptoService := cryptography.NewKeyBoundCryptographyService(sk)

	return &consensusTransportSuite{
		Suite:  suite.Suite{},
		crypto: cryptoService,
	}, nil
}

func TestConsensusTransport(t *testing.T) {
	s, err := NewSuite()
	require.NoError(t, err)
	suite.Run(t, s)
}
