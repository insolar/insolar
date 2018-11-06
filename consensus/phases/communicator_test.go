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

package phases

import (
	"context"
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/suite"
)

type communicatorSuite struct {
	suite.Suite
	componentManager component.Manager
	communicator     Communicator
	participants     []core.Node
	hostNetwork      *network.HostNetworkMock
}

func NewSuite() *communicatorSuite {
	return &communicatorSuite{
		Suite:        suite.Suite{},
		communicator: NewNaiveCommunicator(),
		participants: nil,
	}
}

func (s *communicatorSuite) SetupTest() {
	s.hostNetwork = network.NewHostNetworkMock(s.T())
	s.componentManager.Register(s.communicator, s.hostNetwork)
	err := s.componentManager.Start(nil)
	s.NoError(err)
}

func (s *communicatorSuite) TestExchangeData() {
	s.Assert().NotNil(s.communicator)
	_, err := s.communicator.ExchangeData(context.Background(), s.participants, Phase1Packet{})
	s.Assert().NoError(err)
}

func TestNaiveCommunicator(t *testing.T) {
	suite.Run(t, NewSuite())
}
