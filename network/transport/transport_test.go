/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package transport

import (
	"context"
	"crypto/rand"
	"encoding/gob"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/insolar/insolar/network/transport/relay"
	"github.com/insolar/insolar/network/transport/resolver"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type node struct {
	config    configuration.Transport
	transport Transport
	host      *host.Host
}

type transportSuite struct {
	suite.Suite
	node1 node
	node2 node
}

func NewSuite(cfg1 configuration.Transport, cfg2 configuration.Transport) *transportSuite {
	return &transportSuite{
		Suite: suite.Suite{},
		node1: node{config: cfg1},
		node2: node{config: cfg2},
	}
}

func setupNode(t *transportSuite, n *node) {
	var err error
	n.host, err = host.NewHost(n.config.Address)
	t.Assert().NoError(err)

	n.transport, err = NewTransport(n.config, relay.NewProxy())
	t.Require().NoError(err)
	t.Require().NotNil(n.transport)
	t.Require().Implements((*Transport)(nil), n.transport)
}

func (t *transportSuite) SetupTest() {
	gob.Register(&packet.RequestTest{})
	setupNode(t, &t.node1)
	setupNode(t, &t.node2)
}

func (t *transportSuite) BeforeTest(suiteName, testName string) {
	ctx := context.Background()
	started1 := make(chan struct{}, 1)
	started2 := make(chan struct{}, 1)

	go t.node1.transport.Listen(ctx, started1)
	go t.node2.transport.Listen(ctx, started2)

	<-started1
	<-started2
}

func (t *transportSuite) AfterTest(suiteName, testName string) {
	go t.node1.transport.Stop()
	<-t.node1.transport.Stopped()
	t.node1.transport.Close()

	go t.node2.transport.Stop()
	<-t.node2.transport.Stopped()
	t.node2.transport.Close()
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (t *transportSuite) TestPingPong() {
	if t.node1.config.Protocol == "PURE_UDP" {
		t.T().Skip("Skipping TestPingPong for PURE_UDP")
	}
	ctx := context.Background()
	p := packet.NewBuilder(t.node1.host).Type(types.Ping).Receiver(t.node2.host).Build()
	future, err := t.node1.transport.SendRequest(ctx, p)
	t.Assert().NoError(err)

	requestMsg := <-t.node2.transport.Packets()
	t.Assert().Equal(types.Ping, requestMsg.Type)
	t.Assert().Equal(t.node2.host, future.Actor())
	t.Assert().False(requestMsg.IsResponse)

	builder := packet.NewBuilder(t.node2.host).Receiver(requestMsg.Sender).Type(types.Ping)
	err = t.node2.transport.SendResponse(ctx, requestMsg.RequestID, builder.Response(nil).Build())
	t.Assert().NoError(err)

	responseMsg := <-future.Result()
	t.Assert().Equal(types.Ping, responseMsg.Type)
	t.Assert().True(responseMsg.IsResponse)
}

func (t *transportSuite) TestSendBigPacket() {
	if testing.Short() {
		t.T().Skip("Skipping TestSendBigPacket in short mode")
	}
	if t.node1.config.Protocol == "PURE_UDP" {
		t.T().Skip("Skipping TestSendBigPacket for PURE_UDP")
	}
	ctx := context.Background()
	data, _ := generateRandomBytes(1024 * 1024 * 2)
	builder := packet.NewBuilder(t.node1.host).Receiver(t.node2.host).Type(packet.TestPacket)
	requestMsg := builder.Request(&packet.RequestTest{Data: data}).Build()

	_, err := t.node1.transport.SendRequest(ctx, requestMsg)
	t.Assert().NoError(err)

	msg := <-t.node2.transport.Packets()
	t.Assert().Equal(packet.TestPacket, msg.Type)
	receivedData := msg.Data.(*packet.RequestTest).Data
	t.Assert().Equal(data, receivedData)
}

// func (t *consensusSuite) TestSendPacketConsensus() {
// 	t.T().Skip("fix tests for consensus udp transport")
// 	ctx := context.Background()
// 	builder := packet.NewBuilder(t.node1.host).Receiver(t.node2.host).Type(types.Phase1)
// 	requestMsg := builder.Request(consensus.NewPhase1Packet()).Build()
// 	_, err := t.node1.transport.SendRequest(ctx, requestMsg)
// 	t.Assert().NoError(err)
//
// 	<-t.node2.transport.Packets()
// }

// func TestUDPTransport(t *testing.T) {
// 	cfg1 := configuration.Transport{Protocol: "PURE_UDP", Address: "127.0.0.1:17014", BehindNAT: false}
// 	cfg2 := configuration.Transport{Protocol: "PURE_UDP", Address: "127.0.0.1:17015", BehindNAT: false}
//
// 	suite.Run(t, NewConsensusSuite(cfg1, cfg2))
// }

func TestTCPTransport(t *testing.T) {
	cfg1 := configuration.Transport{Protocol: "TCP", Address: "127.0.0.1:17016", BehindNAT: false}
	cfg2 := configuration.Transport{Protocol: "TCP", Address: "127.0.0.1:17017", BehindNAT: false}

	suite.Run(t, NewSuite(cfg1, cfg2))
}

func TestQuicTransport(t *testing.T) {
	t.Skip("QUIC internals racing atm. Skip until we want to use it in production")

	cfg1 := configuration.Transport{Protocol: "QUIC", Address: "127.0.0.1:17018", BehindNAT: false}
	cfg2 := configuration.Transport{Protocol: "QUIC", Address: "127.0.0.1:17019", BehindNAT: false}

	suite.Run(t, NewSuite(cfg1, cfg2))
}

func Test_createResolver(t *testing.T) {
	a := assert.New(t)

	cfg1 := configuration.Transport{Protocol: "TCP", Address: "127.0.0.1:17018", BehindNAT: false, FixedPublicAddress: "192.168.0.1"}
	r, err := createResolver(cfg1)
	a.NoError(err)
	a.IsType(resolver.NewFixedAddressResolver(""), r)

	cfg2 := configuration.Transport{Protocol: "TCP", Address: "127.0.0.1:17018", BehindNAT: true, FixedPublicAddress: ""}
	r, err = createResolver(cfg2)
	a.NoError(err)
	a.IsType(resolver.NewStunResolver(""), r)

	cfg3 := configuration.Transport{Protocol: "TCP", Address: "127.0.0.1:17018", BehindNAT: false, FixedPublicAddress: ""}
	r, err = createResolver(cfg3)
	a.NoError(err)
	a.IsType(resolver.NewExactResolver(), r)

	cfg4 := configuration.Transport{Protocol: "TCP", Address: "127.0.0.1:17018", BehindNAT: true, FixedPublicAddress: "192.168.0.1"}
	r, err = createResolver(cfg4)
	a.Error(err)
	a.Nil(r)
}
