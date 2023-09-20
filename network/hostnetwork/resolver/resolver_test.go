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
