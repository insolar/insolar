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

package networkcoordinator

import (
	"context"

	"github.com/insolar/insolar/core"
)

// NetworkCoordinator encapsulates logic of network configuration
type NetworkCoordinator struct {
	CertificateManager core.CertificateManager  `inject:""`
	NetworkSwitcher    core.NetworkSwitcher     `inject:""`
	ContractRequester  core.ContractRequester   `inject:""`
	MessageBus         core.MessageBus          `inject:""`
	CS                 core.CryptographyService `inject:""`
	PS                 core.PulseStorage        `inject:""`

	realCoordinator Coordinator
	zeroCoordinator Coordinator
	isStarted       bool
}

// New creates new NetworkCoordinator
func New() (*NetworkCoordinator, error) {
	return &NetworkCoordinator{}, nil
}

// Start implements interface of Component
func (nc *NetworkCoordinator) Start(ctx context.Context) error {
	nc.MessageBus.MustRegister(core.TypeNodeSignRequest, nc.signCertHandler)

	nc.zeroCoordinator = newZeroNetworkCoordinator()
	nc.realCoordinator = newRealNetworkCoordinator(
		nc.CertificateManager,
		nc.ContractRequester,
		nc.MessageBus,
		nc.CS,
	)
	nc.isStarted = true
	return nil
}

func (nc *NetworkCoordinator) getCoordinator() Coordinator {
	if nc.NetworkSwitcher.GetState() == core.CompleteNetworkState {
		return nc.realCoordinator
	}
	return nc.zeroCoordinator
}

// IsStarted returns true if component was started and false in other way
func (nc *NetworkCoordinator) IsStarted() bool {
	return nc.isStarted
}

// GetCert method returns node certificate by requesting sign from discovery nodes
func (nc *NetworkCoordinator) GetCert(ctx context.Context, registeredNodeRef *core.RecordRef) (core.Certificate, error) {
	return nc.getCoordinator().GetCert(ctx, registeredNodeRef)
}

// ValidateCert validates node certificate
func (nc *NetworkCoordinator) ValidateCert(ctx context.Context, certificate core.AuthorizationCertificate) (bool, error) {
	return nc.CertificateManager.VerifyAuthorizationCertificate(certificate)
}

// signCertHandler is MsgBus handler that signs certificate for some node with node own key
func (nc *NetworkCoordinator) signCertHandler(ctx context.Context, p core.Parcel) (core.Reply, error) {
	return nc.getCoordinator().signCertHandler(ctx, p)
}

// WriteActiveNodes writes active nodes to ledger
func (nc *NetworkCoordinator) WriteActiveNodes(ctx context.Context, number core.PulseNumber, activeNodes []core.Node) error {
	return nc.getCoordinator().WriteActiveNodes(ctx, number, activeNodes)
}

// SetPulse writes pulse data on local storage
func (nc *NetworkCoordinator) SetPulse(ctx context.Context, pulse core.Pulse) error {
	return nc.getCoordinator().SetPulse(ctx, pulse)
}
