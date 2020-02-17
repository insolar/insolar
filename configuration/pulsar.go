// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package configuration

// Pulsar holds configuration for pulsar node.
type Pulsar struct {
	PulseTime int32 // ms

	NumberDelta uint32

	DistributionTransport Transport
	PulseDistributor      PulseDistributor
}

type PulseDistributor struct {
	BootstrapHosts      []string
	PulseRequestTimeout int32 // ms
}

type PulsarNodeAddress struct {
	Address   string
	PublicKey string
}

// NewPulsar creates new default configuration for pulsar node.
func NewPulsar() Pulsar {
	return Pulsar{
		PulseTime: 10000,

		NumberDelta: 10,
		DistributionTransport: Transport{
			Protocol: "TCP",
			Address:  "0.0.0.0:18091",
		},
		PulseDistributor: PulseDistributor{
			BootstrapHosts:      []string{"localhost:53837"},
			PulseRequestTimeout: 1000,
		},
	}
}
