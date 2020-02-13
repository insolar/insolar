// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package configuration

// Introspection holds introspection configuration.
type Introspection struct {
	// Addr specifies address where introspection server starts
	Addr string
}

// NewIntrospection creates new default configuration for introspection.
func NewIntrospection() Introspection {
	return Introspection{}
}
