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
	"crypto/rand"
	"testing"

	"github.com/insolar/insolar/configuration"
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
	address   *host.Address
}

type transportSuite struct {
	suite.Suite
	node1 node
	node2 node

	testSendBigPacket bool
}

func NewSuite(cfg1 configuration.Transport, cfg2 configuration.Transport, testSendBigPacket bool) *transportSuite {
	return &transportSuite{
		Suite:             suite.Suite{},
		testSendBigPacket: testSendBigPacket,
		node1:             node{config: cfg1},
		node2:             node{config: cfg2},
	}
}

func setupNode(t *transportSuite, n *node) {
	var err error
	n.address, err = host.NewAddress(n.config.Address)
	t.Assert().NoError(err)

	n.host = host.NewHost(n.address)

	n.transport, err = NewTransport(n.config, relay.NewProxy())
	t.Require().NoError(err)
	t.Require().NotNil(n.transport)
	t.Require().Implements((*Transport)(nil), n.transport)
}

func (t *transportSuite) SetupTest() {
	setupNode(t, &t.node1)
	setupNode(t, &t.node2)
}

func (t *transportSuite) BeforeTest(suiteName, testName string) {
	go t.node1.transport.Start()
	go t.node2.transport.Start()
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
	future, err := t.node1.transport.SendRequest(packet.NewPingPacket(t.node1.host, t.node2.host))
	t.Assert().NoError(err)

	requestMsg := <-t.node2.transport.Packets()
	t.Assert().True(requestMsg.IsValid())
	t.Assert().Equal(types.TypePing, requestMsg.Type)
	t.Assert().Equal(t.node2.host, future.Actor())
	t.Assert().False(requestMsg.IsResponse)

	builder := packet.NewBuilder(t.node2.host).Receiver(requestMsg.Sender).Type(types.TypePing)
	err = t.node2.transport.SendResponse(requestMsg.RequestID, builder.Response(nil).Build())
	t.Assert().NoError(err)

	responseMsg := <-future.Result()
	t.Assert().True(responseMsg.IsValid())
	t.Assert().Equal(types.TypePing, responseMsg.Type)
	t.Assert().True(responseMsg.IsResponse)
}

func (t *transportSuite) TestSendBigPacket() {
	if testing.Short() {
		t.T().Skip("Skipping TestSendBigPacket in short mode")
	}
	if !t.testSendBigPacket {
		t.T().Skip("TestSendBigPacket is skipped because transport does not support transfer " +
			"of messages that are heavier than UDP datagram")
	}
	data, _ := generateRandomBytes(1024 * 1024 * 2)
	builder := packet.NewBuilder(t.node1.host).Receiver(t.node2.host).Type(types.TypeStore)
	requestMsg := builder.Request(&packet.RequestDataStore{data, true}).Build()
	t.Assert().True(requestMsg.IsValid())

	_, err := t.node1.transport.SendRequest(requestMsg)
	t.Assert().NoError(err)

	msg := <-t.node2.transport.Packets()
	t.Assert().True(requestMsg.IsValid())
	t.Assert().Equal(types.TypeStore, requestMsg.Type)
	receivedData := msg.Data.(*packet.RequestDataStore).Data
	t.Assert().Equal(data, receivedData)
}

func TestUTPTransport(t *testing.T) {
	cfg1 := configuration.Transport{Protocol: "UTP", Address: "127.0.0.1:17010", BehindNAT: false}
	cfg2 := configuration.Transport{Protocol: "UTP", Address: "127.0.0.1:17011", BehindNAT: false}

	suite.Run(t, NewSuite(cfg1, cfg2, true))
}

func TestKCPTransport(t *testing.T) {
	cfg1 := configuration.Transport{Protocol: "KCP", Address: "127.0.0.1:17012", BehindNAT: false}
	cfg2 := configuration.Transport{Protocol: "KCP", Address: "127.0.0.1:17013", BehindNAT: false}

	suite.Run(t, NewSuite(cfg1, cfg2, true))
}

func TestUDPTransport(t *testing.T) {
	cfg1 := configuration.Transport{Protocol: "PURE_UDP", Address: "127.0.0.1:17014", BehindNAT: false}
	cfg2 := configuration.Transport{Protocol: "PURE_UDP", Address: "127.0.0.1:17015", BehindNAT: false}

	suite.Run(t, NewSuite(cfg1, cfg2, false))
}
