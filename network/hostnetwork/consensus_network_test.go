//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package hostnetwork

import (
	"context"
	"sync"
	"testing"

	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensusv1/packets"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/platformpolicy"
)

type consensusNetworkSuite struct {
	suite.Suite
	crypto     insolar.CryptographyService
	id1, id2   string
	sid1, sid2 insolar.ShortNodeID
	ref1, ref2 insolar.Reference
}

type consensusTestCase struct {
	parent   *consensusNetworkSuite
	ctx      context.Context
	cn1, cn2 network.ConsensusNetwork
	resolver *MockResolver
}

func (s *consensusNetworkSuite) newTestCase() *consensusTestCase {
	resolver := newMockResolver()

	cm1 := component.NewManager(nil)
	f1 := transport.NewFactory(configuration.NewHostNetwork().Transport)
	cn1, err := NewConsensusNetwork(s.id1, s.sid1)
	require.NoError(s.T(), err)
	cm1.Inject(f1, cn1, resolver)

	cm2 := component.NewManager(nil)
	f2 := transport.NewFactory(configuration.NewHostNetwork().Transport)
	cn2, err := NewConsensusNetwork(s.id2, s.sid2)
	require.NoError(s.T(), err)
	cm2.Inject(f2, cn2, resolver)

	ctx := context.Background()

	err = cn1.Init(ctx)
	require.NoError(s.T(), err)
	err = cn2.Init(ctx)
	require.NoError(s.T(), err)

	return &consensusTestCase{
		parent:   s,
		ctx:      ctx,
		cn1:      cn1,
		cn2:      cn2,
		resolver: resolver,
	}
}

func (ctc *consensusTestCase) Start() {
	// start the second consensusNetwork before the first because test cases perform sending packets first -> second,
	// so the second consensusNetwork should be ready to receive packets when the first starts to send
	err := ctc.cn2.Start(ctc.ctx)
	require.NoError(ctc.parent.T(), err)
	err = ctc.cn1.Start(ctc.ctx)
	require.NoError(ctc.parent.T(), err)

	routing1, err := host.NewHostNS(ctc.cn1.PublicAddress(), ctc.parent.ref1, ctc.parent.sid1)
	require.NoError(ctc.parent.T(), err)
	routing2, err := host.NewHostNS(ctc.cn2.PublicAddress(), ctc.parent.ref2, ctc.parent.sid2)
	require.NoError(ctc.parent.T(), err)
	ctc.resolver.addMappingHost(routing1)
	ctc.resolver.addMappingHost(routing2)
}

func (ctc *consensusTestCase) Stop() {
	// stop consensusNetworks in the reverse order of their start
	_ = ctc.cn1.Stop(ctc.ctx)
	_ = ctc.cn2.Stop(ctc.ctx)
}

func (s *consensusNetworkSuite) sendPacket(packet packets.ConsensusPacket) {
	ctc := s.newTestCase()
	defer ctc.Stop()

	wg := sync.WaitGroup{}
	wg.Add(1)

	handler := func(incomingPacket packets.ConsensusPacket, sender insolar.Reference) {
		log.Info("handler triggered")
		wg.Done()
	}
	ctc.cn2.RegisterPacketHandler(packet.GetType(), handler)

	ctc.Start()

	err := ctc.cn1.SignAndSendPacket(packet, s.ref2, s.crypto)
	s.Require().NoError(err)
	wg.Wait()
}

func newPhase1Packet() *packets.Phase1Packet {
	return packets.NewPhase1Packet(insolar.Pulse{})
}

func newPhase2Packet() (*packets.Phase2Packet, error) {
	bitset, err := packets.NewBitSet(10)
	if err != nil {
		return nil, err
	}
	result := packets.NewPhase2Packet(insolar.PulseNumber(0))
	result.SetBitSet(bitset)
	return result, nil
}

func newPhase3Packet() (*packets.Phase3Packet, error) {
	var ghs packets.GlobuleHashSignature
	bitset, err := packets.NewBitSet(10)
	if err != nil {
		return nil, err
	}
	return packets.NewPhase3Packet(insolar.PulseNumber(0), ghs, bitset), nil
}

func (s *consensusNetworkSuite) TestSendPacketPhase1() {
	packet := newPhase1Packet()
	s.sendPacket(packet)
}

func (s *consensusNetworkSuite) TestSendPacketPhase2() {
	packet, err := newPhase2Packet()
	require.NoError(s.T(), err)
	s.sendPacket(packet)
}

func (s *consensusNetworkSuite) TestSendPacketPhase3() {
	packet, err := newPhase3Packet()
	require.NoError(s.T(), err)
	s.sendPacket(packet)
}

func (s *consensusNetworkSuite) sendPacketAndVerify(packet packets.ConsensusPacket) {
	ctc := s.newTestCase()
	defer ctc.Stop()

	result := make(chan bool, 1)

	handler := func(incomingPacket packets.ConsensusPacket, sender insolar.Reference) {
		log.Info("handler triggered")
		pk, err := s.crypto.GetPublicKey()
		if err != nil {
			log.Error("handler get public key error: " + err.Error())
			result <- false
			return
		}
		err = incomingPacket.Verify(s.crypto, pk)
		if err != nil {
			log.Error("verify signature error: " + err.Error())
			result <- false
			return
		}
		result <- true
	}
	ctc.cn2.RegisterPacketHandler(packet.GetType(), handler)

	ctc.Start()

	err := ctc.cn1.SignAndSendPacket(packet, s.ref2, s.crypto)
	s.Require().NoError(err)
	s.True(<-result)
}

func (s *consensusNetworkSuite) TestVerifySignPhase1() {
	packet := newPhase1Packet()
	s.sendPacketAndVerify(packet)
}

func (s *consensusNetworkSuite) TestVerifySignPhase2() {
	packet, err := newPhase2Packet()
	require.NoError(s.T(), err)
	s.sendPacketAndVerify(packet)
}

func (s *consensusNetworkSuite) TestVerifySignPhase3() {
	packet, err := newPhase3Packet()
	require.NoError(s.T(), err)
	s.sendPacketAndVerify(packet)
}

func NewSuite(t *testing.T) *consensusNetworkSuite {
	kp := platformpolicy.NewKeyProcessor()
	sk, err := kp.GeneratePrivateKey()
	require.NoError(t, err)
	cryptoService := cryptography.NewKeyBoundCryptographyService(sk)

	id1 := ID1 + DOMAIN
	id2 := ID2 + DOMAIN
	sid1 := insolar.ShortNodeID(0)
	sid2 := insolar.ShortNodeID(1)
	ref1, err := insolar.NewReferenceFromBase58(id1)
	require.NoError(t, err)
	ref2, err := insolar.NewReferenceFromBase58(id2)
	require.NoError(t, err)

	return &consensusNetworkSuite{
		Suite:  suite.Suite{},
		crypto: cryptoService,
		id1:    id1, id2: id2, sid1: sid1, sid2: sid2, ref1: *ref1, ref2: *ref2,
	}
}

func TestConsensusNetwork(t *testing.T) {
	suite.Run(t, NewSuite(t))
}

func TestNetworkConsensus_SignAndSendPacket_NotStarted(t *testing.T) {
	cn, err := NewConsensusNetwork(ID1+DOMAIN, 1)
	require.NoError(t, err)

	err = cn.SignAndSendPacket(nil, testutils.RandomRef(), nil)
	require.EqualError(t, err, "consensus network is not started")
}
