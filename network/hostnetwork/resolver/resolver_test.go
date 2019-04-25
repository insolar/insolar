/*
 *    Copyright 2019 Insolar Technologies
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

package resolver

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ResolverSuite struct {
	suite.Suite
}

func (s *ResolverSuite) TestSuccessPublic() {
	localAddress := "127.0.0.1:12345"
	externalAddress := "192.168.0.1"

	realAddress, err := Resolve(externalAddress, localAddress)
	s.NoError(err)
	s.Equal("192.168.0.1:12345", realAddress)
}

func (s *ResolverSuite) TestSuccessExact() {
	localAddress := "127.0.0.1:12345"

	realAddress, err := Resolve("", localAddress)
	s.NoError(err)
	s.Equal(localAddress, realAddress)
}

func (s *ResolverSuite) TestFailure() {
	localAddress := "empty_port"
	externalAddress := "192.168.0.1"

	_, err := Resolve(externalAddress, localAddress)
	s.Error(err)
}

func TestResolver(t *testing.T) {
	suite.Run(t, new(ResolverSuite))
}
