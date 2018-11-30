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
	Certificate         core.Certificate         `inject:""`
	NetworkSwitcher     core.NetworkSwitcher     `inject:""`
	ContractRequester   core.ContractRequester   `inject:""`
	GenesisDataProvider core.GenesisDataProvider `inject:""`

	realCoordinator core.NetworkCoordinator
	zeroCoordinator core.NetworkCoordinator
}

// New creates new NetworkCoordinator
func New() (*NetworkCoordinator, error) {
	return &NetworkCoordinator{}, nil
}

// Init implements interface of Component
func (nc *NetworkCoordinator) Init(ctx context.Context) error {
	nc.zeroCoordinator = newZeroNetworkCoordinator()
	nc.realCoordinator = newRealNetworkCoordinator()
	return nil
}

func (nc *NetworkCoordinator) getCoordinator() core.NetworkCoordinator {
	if nc.NetworkSwitcher.GetState() == core.CompleteNetworkState {
		return nc.realCoordinator
	}
	return nc.zeroCoordinator
}

// GetCert method returns node certificate
func (nc *NetworkCoordinator) GetCert(ctx context.Context, nodeRef core.RecordRef) (core.Certificate, error) {
	return nc.getCoordinator().GetCert(ctx, nodeRef)
}

// ValidateCert validates node certificate
func (nc *NetworkCoordinator) ValidateCert(ctx context.Context, certificate core.AuthorizationCertificate) (bool, error) {
	return nc.getCoordinator().ValidateCert(ctx, certificate)
}

// WriteActiveNodes writes active nodes to ledger
func (nc *NetworkCoordinator) WriteActiveNodes(ctx context.Context, number core.PulseNumber, activeNodes []core.Node) error {
	return nc.getCoordinator().WriteActiveNodes(ctx, number, activeNodes)
}

// SetPulse writes pulse data on local storage
func (nc *NetworkCoordinator) SetPulse(ctx context.Context, pulse core.Pulse) error {
	return nc.getCoordinator().SetPulse(ctx, pulse)
}
