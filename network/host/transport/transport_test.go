/*
 *    Copyright 2018 INS Ecosystem
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
	"net"
	"testing"

	"github.com/insolar/insolar/network/host/connection"
	"github.com/insolar/insolar/network/host/message"
	"github.com/insolar/insolar/network/host/node"
	"github.com/insolar/insolar/network/host/relay"
	"github.com/stretchr/testify/suite"
)

type transportSuite struct {
	suite.Suite
	factory    Factory
	connection net.PacketConn
	transport  Transport
}

func NewSuite(factory Factory) *transportSuite {
	return &transportSuite{suite.Suite{}, NewUTPTransportFactory(), nil, nil}
}

func (t *transportSuite) SetupTest() {
	t.connection, _ = connection.NewConnectionFactory().Create("127.0.0.1:3012")
	var err error
	t.transport, err = t.factory.Create(t.connection, relay.NewProxy())
	t.Assert().NoError(err)
	t.Assert().Implements((*Transport)(nil), t.transport)
}

func (t *transportSuite) BeforeTest(suiteName, testName string) {
	go t.transport.Start()
}

func (t *transportSuite) AfterTest(suiteName, testName string) {
	go t.transport.Stop()
	<-t.transport.Stopped()
	t.transport.Close()
}

func (t *transportSuite) TestPingPong() {

	addr, _ := node.NewAddress("127.0.0.1:3012")
	nodeOne := node.NewNode(addr)

	future, err := t.transport.SendRequest(message.NewPingMessage(nodeOne, nodeOne))
	t.Assert().NoError(err)

	requestMsg := <-t.transport.Messages()
	t.Assert().Equal(true, requestMsg.IsValid())
	t.Assert().Equal(message.TypePing, requestMsg.Type)
	t.Assert().Equal(nodeOne, future.Actor())
	t.Assert().Equal(false, requestMsg.IsResponse)

	builder := message.NewBuilder().Sender(nodeOne).Receiver(requestMsg.Sender).Type(message.TypePing)
	err = t.transport.SendResponse(requestMsg.RequestID, builder.Response(nil).Build())
	t.Assert().NoError(err)

	responseMsg := <-future.Result()
	t.Assert().Equal(true, responseMsg.IsValid())
	t.Assert().Equal(message.TypePing, responseMsg.Type)
	t.Assert().Equal(true, responseMsg.IsResponse)
}

func TestUTPTransport(t *testing.T) {
	suite.Run(t, NewSuite(NewUTPTransportFactory()))
}

func TestKCPTransport(t *testing.T) {
	suite.Run(t, NewSuite(NewKCPTransportFactory()))
}

func TestKCPSecureTransport(t *testing.T) {
	// secureOptions := {}
	suite.Run(t, NewSuite(NewKCPTransportFactory( /* secureOptions */ )))
}
