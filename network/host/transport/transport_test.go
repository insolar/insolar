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
	"github.com/insolar/insolar/network/host/relay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TransportSuite struct {
	suite.Suite
	factory    Factory
	connection net.PacketConn
	transport  Transport
}

func (t *TransportSuite) SetupTest() {
	t.connection, _ = connection.NewConnectionFactory().Create("127.0.0.1:3012")
	var err error
	t.transport, err = t.factory.Create(t.connection, relay.NewProxy())
	assert.NoError(t.T(), err)
	assert.Implements(t.T(), (*Transport)(nil), t.transport)
}

func (t *TransportSuite) TestStartStop() {

	go t.transport.Start()
	go t.transport.Stop()
	<-t.transport.Stopped()
	t.transport.Close()
}

func TestUTPTransport(t *testing.T) {
	suite.Run(t, &TransportSuite{suite.Suite{}, NewUTPTransportFactory(), nil, nil})
}

func TestKCPTransport(t *testing.T) {
	suite.Run(t, &TransportSuite{suite.Suite{}, NewKCPTransportFactory(), nil, nil})
}
