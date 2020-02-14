// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package resolver

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type FixedAddressResolverSuite struct {
	suite.Suite
}

func (s *FixedAddressResolverSuite) TestSuccess() {
	localAddress := "127.0.0.1:12345"
	externalAddress := "192.168.0.1"

	r := NewFixedAddressResolver(externalAddress)
	s.Require().IsType(&fixedAddressResolver{}, r)
	realAddress, err := r.Resolve(localAddress)
	s.NoError(err)
	s.Equal("192.168.0.1:12345", realAddress)
}

func (s *FixedAddressResolverSuite) TestFailure_EmptyPort() {
	localAddress := "empty_port"
	externalAddress := "192.168.0.1"

	r := NewFixedAddressResolver(externalAddress)
	s.Require().IsType(&fixedAddressResolver{}, r)
	_, err := r.Resolve(localAddress)
	s.Error(err)
}

func TestFixedAddressResolver(t *testing.T) {
	suite.Run(t, new(FixedAddressResolverSuite))
}
