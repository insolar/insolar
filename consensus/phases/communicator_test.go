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
	"testing"

	"github.com/stretchr/testify/suite"
)

type communicatorSuite struct {
	suite.Suite
	communicator Communicator
	participants []Communicator
}

func NewSuite() *communicatorSuite {
	return &communicatorSuite{
		Suite:        suite.Suite{},
		communicator: nil,
		participants: nil,
	}
}

func (t *communicatorSuite) SetupTest() {
	//setupNode(t, &t.node1)
	//setupNode(t, &t.node2)
}

func TestNaiveCommunicator(t *testing.T) {
	suite.Run(t, NewSuite())
}
