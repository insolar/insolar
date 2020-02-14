// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package configuration

// ServiceNetwork is configuration for ServiceNetwork.
type ServiceNetwork struct {
	CacheDirectory string
}

// NewServiceNetwork creates a new ServiceNetwork configuration.
func NewServiceNetwork() ServiceNetwork {
	return ServiceNetwork{
		CacheDirectory: "network_cache",
	}
}
