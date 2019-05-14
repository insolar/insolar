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
	"github.com/insolar/insolar/network/consensus/packets"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/platformpolicy"
)

type consensusNetworkSuite struct {
	suite.Suite
	crypto insolar.CryptographyService
}

func createTwoConsensusNetworks(id1, id2 insolar.ShortNodeID) (t1, t2 network.ConsensusNetwork, err error) {
	m := newMockResolver()

	cm1 := component.NewManager(nil)
	f1 := transport.NewFactory(configuration.NewHostNetwork().Transport)
	cn1, err := NewConsensusNetwork(ID1+DOMAIN, id1)
	if err != nil {
		return nil, nil, err
	}
	cm1.Inject(f1, cn1, m)

	cm2 := component.NewManager(nil)
	f2 := transport.NewFactory(configuration.NewHostNetwork().Transport)
	cn2, err := NewConsensusNetwork(ID2+DOMAIN, id2)
	if err != nil {
		return nil, nil, err
	}
	cm2.Inject(f2, cn2, m)

	ctx := context.Background()

	err = cn1.Init(ctx)
	if err != nil {
		return nil, nil, err
	}
	err = cn2.Init(ctx)
	if err != nil {
		return nil, nil, err
	}

	err = cn1.Start(ctx)
	if err != nil {
		return nil, nil, err
	}
	err = cn2.Start(ctx)
	if err != nil {
		return nil, nil, err
	}

	ref1, err := insolar.NewReferenceFromBase58(ID2 + DOMAIN)
	if err != nil {
		return nil, nil, err
	}
	routing1, err := host.NewHostNS(cn1.PublicAddress(), *ref1, id1)
	if err != nil {
		return nil, nil, err
	}
	ref2, err := insolar.NewReferenceFromBase58(ID2 + DOMAIN)
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

func (t *consensusNetworkSuite) sendPacket(packet packets.ConsensusPacket) {
	cn1, cn2, err := createTwoConsensusNetworks(0, 1)
	t.Require().NoError(err)
	ctx := context.Background()
	ctx2 := context.Background()

	wg := sync.WaitGroup{}
	wg.Add(1)

	handler := func(incomingPacket packets.ConsensusPacket, sender insolar.Reference) {
		log.Info("handler triggered")
		wg.Done()
	}
	cn2.RegisterPacketHandler(packet.GetType(), handler)

	err = cn2.Start(ctx2)
	t.Require().NoError(err)
	err = cn1.Start(ctx)
	t.Require().NoError(err)
	defer func() {
		cn1.Stop(ctx)
		cn2.Stop(ctx2)
	}()

	ref2, err := insolar.NewReferenceFromBase58(ID2 + DOMAIN)
	t.Require().NoError(err)

	err = cn1.SignAndSendPacket(packet, *ref2, t.crypto)
	t.Require().NoError(err)
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

func (t *consensusNetworkSuite) TestSendPacketPhase1() {
	packet := newPhase1Packet()
	t.sendPacket(packet)
}

func (t *consensusNetworkSuite) TestSendPacketPhase2() {
	packet, err := newPhase2Packet()
	require.NoError(t.T(), err)
	t.sendPacket(packet)
}

func (t *consensusNetworkSuite) TestSendPacketPhase3() {
	packet, err := newPhase3Packet()
	require.NoError(t.T(), err)
	t.sendPacket(packet)
}

func (t *consensusNetworkSuite) sendPacketAndVerify(packet packets.ConsensusPacket) {
	cn1, cn2, err := createTwoConsensusNetworks(0, 1)
	t.Require().NoError(err)
	ctx := context.Background()
	ctx2 := context.Background()

	result := make(chan bool, 1)

	handler := func(incomingPacket packets.ConsensusPacket, sender insolar.Reference) {
		log.Info("handler triggered")
		pk, err := t.crypto.GetPublicKey()
		if err != nil {
			log.Error("handler get public key error: " + err.Error())
			result <- false
			return
		}
		err = incomingPacket.Verify(t.crypto, pk)
		if err != nil {
			log.Error("verify signature error: " + err.Error())
			result <- false
			return
		}
		result <- true
	}
	cn2.RegisterPacketHandler(packet.GetType(), handler)

	err = cn2.Start(ctx2)
	t.Require().NoError(err)
	err = cn1.Start(ctx)
	t.Require().NoError(err)
	defer func() {
		cn1.Stop(ctx)
		cn2.Stop(ctx2)
	}()

	ref2, err := insolar.NewReferenceFromBase58(ID2 + DOMAIN)
	t.Require().NoError(err)

	err = cn1.SignAndSendPacket(packet, *ref2, t.crypto)
	t.Require().NoError(err)
	t.True(<-result)
}

func (t *consensusNetworkSuite) TestVerifySignPhase1() {
	packet := newPhase1Packet()
	t.sendPacketAndVerify(packet)
}

func (t *consensusNetworkSuite) TestVerifySignPhase2() {
	packet, err := newPhase2Packet()
	require.NoError(t.T(), err)
	t.sendPacketAndVerify(packet)
}

func (t *consensusNetworkSuite) TestVerifySignPhase3() {
	packet, err := newPhase3Packet()
	require.NoError(t.T(), err)
	t.sendPacketAndVerify(packet)
}

func NewSuite() (*consensusNetworkSuite, error) {
	kp := platformpolicy.NewKeyProcessor()
	sk, err := kp.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	cryptoService := cryptography.NewKeyBoundCryptographyService(sk)

	return &consensusNetworkSuite{
		Suite:  suite.Suite{},
		crypto: cryptoService,
	}, nil
}

func TestConsensusNetwork(t *testing.T) {
	s, err := NewSuite()
	require.NoError(t, err)
	suite.Run(t, s)
}

func TestNetworkConsensus_SignAndSendPacket_NotStarted(t *testing.T) {
	cn, err := NewConsensusNetwork(ID1+DOMAIN, 1)
	require.NoError(t, err)

	err = cn.SignAndSendPacket(nil, testutils.RandomRef(), nil)
	require.EqualError(t, err, "consensus network is not started")
}
