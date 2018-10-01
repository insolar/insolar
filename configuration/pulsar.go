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

package configuration

type ConnectionType string

const (
	TCP ConnectionType = "tcp"
)

func (ct ConnectionType) String() string {
	return string(ct)
}

// Pulsar holds configuration for pulsar node.
type Pulsar struct {
	ConnectionType      ConnectionType
	MainListenerAddress string
	PrivateKey          string
	Storage             Storage

	PulseTime                      int32 // ms
	ReceivingSignTimeout           int32 // ms
	ReceivingNumberTimeout         int32 // ms
	ReceivingVectorTimeout         int32 // ms
	ReceivingSignsForChosenTimeout int32 // ms

	Neighbours []PulsarNodeAddress

	NumberOfRandomHosts int
	BootstrapListener   Transport
	BootstrapNodes      []string
}

type PulsarNodeAddress struct {
	Address        string
	ConnectionType ConnectionType
	PublicKey      string
}

// NewPulsar creates new default configuration for pulsar node.
func NewPulsar() Pulsar {
	return Pulsar{
		MainListenerAddress: "0.0.0.0:18090",

		ConnectionType: TCP,
		PrivateKey: `-----BEGIN PRIVATE KEY-----
MHcCAQEEID6XJHMb2aiaK1bp2GHHw0r4LrzZZ4exlcmx8GrjGsMFoAoGCCqGSM49
AwEHoUQDQgAE7DE4ArqxIYbY/UAyLLFBGuFu2gROPaqp4vxbEeie7mnZeqsYexmN
BkrXBEFO5LF4diHC7OJ3xsfebvI0moQRLw==
-----END PRIVATE KEY-----`,

		PulseTime:              10000,
		ReceivingSignTimeout:   1000,
		ReceivingNumberTimeout: 1000,
		ReceivingVectorTimeout: 1000,

		Neighbours: []PulsarNodeAddress{},
		Storage:    Storage{DataDirectory: "./data/pulsar"},

		NumberOfRandomHosts: 1,
		BootstrapListener:   Transport{Protocol: "UTP", Address: "0.0.0.0:18091", BehindNAT: false},
		BootstrapNodes:      []string{"127.0.0.1:64278"},
	}
}
