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

package transport

import (
	"context"
	"crypto/rand"
	"encoding/gob"
	"testing"

	"github.com/insolar/insolar/configuration"
	consensus "github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/insolar/insolar/network/transport/relay"
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
	go t.node1.transport.Listen(ctx)
	go t.node2.transport.Listen(ctx)
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
	p := packet.NewBuilder(t.node1.host).Type(types.Ping).Receiver(t.node2.host).Build()
	future, err := t.node1.transport.SendRequest(p)
	t.Assert().NoError(err)

	requestMsg := <-t.node2.transport.Packets()
	t.Assert().Equal(types.Ping, requestMsg.Type)
	t.Assert().Equal(t.node2.host, future.Actor())
	t.Assert().False(requestMsg.IsResponse)

	builder := packet.NewBuilder(t.node2.host).Receiver(requestMsg.Sender).Type(types.Ping)
	err = t.node2.transport.SendResponse(requestMsg.RequestID, builder.Response(nil).Build())
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
	data, _ := generateRandomBytes(1024 * 1024 * 2)
	builder := packet.NewBuilder(t.node1.host).Receiver(t.node2.host).Type(packet.TestPacket)
	requestMsg := builder.Request(&packet.RequestTest{Data: data}).Build()

	_, err := t.node1.transport.SendRequest(requestMsg)
	t.Assert().NoError(err)

	msg := <-t.node2.transport.Packets()
	t.Assert().Equal(packet.TestPacket, msg.Type)
	receivedData := msg.Data.(*packet.RequestTest).Data
	t.Assert().Equal(data, receivedData)
}

func (t *transportSuite) TestSendPacketConsensus() {
	if t.node1.config.Protocol != "PURE_UDP" {
		t.T().Skip("Skipping TestSendPacketConsensus for non-UDP transports")
	}

	builder := packet.NewBuilder(t.node1.host).Receiver(t.node2.host).Type(types.Phase1)
	requestMsg := builder.Request(consensus.NewPhase1Packet()).Build()
	_, err := t.node1.transport.SendRequest(requestMsg)
	t.Assert().NoError(err)

	msg := <-t.node2.transport.Packets()
	t.Assert().Equal(types.Phase1, msg.Type)
}

func TestUTPTransport(t *testing.T) {
	cfg1 := configuration.Transport{Protocol: "UTP", Address: "127.0.0.1:17010", BehindNAT: false}
	cfg2 := configuration.Transport{Protocol: "UTP", Address: "127.0.0.1:17011", BehindNAT: false}

	suite.Run(t, NewSuite(cfg1, cfg2))
}

func TestKCPTransport(t *testing.T) {
	cfg1 := configuration.Transport{Protocol: "KCP", Address: "127.0.0.1:17012", BehindNAT: false}
	cfg2 := configuration.Transport{Protocol: "KCP", Address: "127.0.0.1:17013", BehindNAT: false}

	suite.Run(t, NewSuite(cfg1, cfg2))
}

func TestUDPTransport(t *testing.T) {
	cfg1 := configuration.Transport{Protocol: "PURE_UDP", Address: "127.0.0.1:17014", BehindNAT: false}
	cfg2 := configuration.Transport{Protocol: "PURE_UDP", Address: "127.0.0.1:17015", BehindNAT: false}

	suite.Run(t, NewSuite(cfg1, cfg2))
}

func TestTCPTransport(t *testing.T) {
	cfg1 := configuration.Transport{Protocol: "TCP", Address: "127.0.0.1:17016", BehindNAT: false}
	cfg2 := configuration.Transport{Protocol: "TCP", Address: "127.0.0.1:17017", BehindNAT: false}

	suite.Run(t, NewSuite(cfg1, cfg2))
}

func TestQuicTransport(t *testing.T) {
	cfg1 := configuration.Transport{Protocol: "QUIC", Address: "127.0.0.1:17018", BehindNAT: false}
	cfg2 := configuration.Transport{Protocol: "QUIC", Address: "127.0.0.1:17019", BehindNAT: false}

	suite.Run(t, NewSuite(cfg1, cfg2))
}
