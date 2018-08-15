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
	"crypto/rand"
	"log"
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
	node       *node.Node
}

func NewSuite(factory Factory) *transportSuite {
	addr, _ := node.NewAddress("127.0.0.1:3012")
	return &transportSuite{suite.Suite{}, NewUTPTransportFactory(), nil, nil, node.NewNode(addr)}
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

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (t *transportSuite) TestPingPong() {
	future, err := t.transport.SendRequest(message.NewPingMessage(t.node, t.node))
	t.Assert().NoError(err)

	requestMsg := <-t.transport.Messages()
	t.Assert().True(requestMsg.IsValid())
	t.Assert().Equal(message.TypePing, requestMsg.Type)
	t.Assert().Equal(t.node, future.Actor())
	t.Assert().False(requestMsg.IsResponse)

	builder := message.NewBuilder().Sender(t.node).Receiver(requestMsg.Sender).Type(message.TypePing)
	err = t.transport.SendResponse(requestMsg.RequestID, builder.Response(nil).Build())
	t.Assert().NoError(err)

	responseMsg := <-future.Result()
	t.Assert().True(responseMsg.IsValid())
	t.Assert().Equal(message.TypePing, responseMsg.Type)
	t.Assert().True(responseMsg.IsResponse)
}

func (t *transportSuite) TestSendBigMessage() {
	t.T().Skip("fix impl for this test pass")

	data, _ := generateRandomBytes(1024 * 1024)
	builder := message.NewBuilder().Sender(t.node).Receiver(t.node).Type(message.TypeStore)
	requestMsg := builder.Request(&message.RequestDataStore{data, true}).Build()
	t.Assert().True(requestMsg.IsValid())

	_, err := t.transport.SendRequest(requestMsg)
	t.Assert().NoError(err)

	msg := <-t.transport.Messages()
	t.Assert().True(requestMsg.IsValid())
	t.Assert().Equal(message.TypeStore, requestMsg.Type)
	receivedData := msg.Data.(*message.RequestDataStore).Data
	t.Assert().Equal(data, receivedData)
}

func (t *transportSuite) TestSendInvalidMessage() {
	t.T().Skip("fix impl for this test pass")

	builder := message.NewBuilder().Sender(t.node).Receiver(t.node).Type(message.TypeRPC)
	msg := builder.Build()
	t.Assert().False(msg.IsValid())

	future, err := t.transport.SendRequest(msg)
	t.Assert().Error(err)
	log.Println("future: ", future.ID())
}

func TestUTPTransport(t *testing.T) {
	suite.Run(t, NewSuite(NewUTPTransportFactory()))
}

func TestKCPTransport(t *testing.T) {
	suite.Run(t, NewSuite(NewKCPTransportFactory()))
}

/*
func TestKCPSecureTransport(t *testing.T) {
	// secureOptions := {}
	suite.Run(t, NewSuite(NewKCPTransportFactory( /* secureOptions * )))
}
*/
