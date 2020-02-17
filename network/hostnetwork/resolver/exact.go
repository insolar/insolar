// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package resolver

type exactResolver struct {
}

// NewExactResolver returns new no-op resolver.
func NewExactResolver() PublicAddressResolver {
	return newExactResolver()
}

func newExactResolver() *exactResolver {
	return &exactResolver{}
}

// Resolve returns host's current network address.
func (er *exactResolver) Resolve(address string) (string, error) {
	return address, nil
}
