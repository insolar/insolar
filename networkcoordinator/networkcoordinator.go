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
	"github.com/insolar/insolar/core/reply"
	"github.com/pkg/errors"
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
func (nc *NetworkCoordinator) ValidateCert(ctx context.Context, certificate core.Certificate) (bool, error) {
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

func (nc *NetworkCoordinator) CreateNodeCert(ctx context.Context, ref string) (core.Certificate, error) {
	rr := core.NewRefFromBase58(ref)
	res, err := nc.ContractRequester.SendRequest(ctx, &rr, "GetNodeInfo", []interface{}{})
	if err != nil {
		return nil, errors.Wrap(err, "[ CreateNodeCert ] Couldn't call GetNodeInfo")
	}
	z, err := core.UnMarshalResponse(res.(*reply.CallMethod).Result, []interface{}{nil})
	if err != nil {
		return nil, errors.Wrap(err, "[ CreateNodeCert ] Couldn't unmarshall response")
	}
	answer := z[0].(map[interface{}]interface{})

	cert, err := nc.Certificate.NewCertForHost(
		answer["PublicKey"].(string),
		core.NodeRole(answer["Role"].(uint64)).String(),
		ref)
	if err != nil {
		return nil, errors.Wrap(err, "[ CreateNodeCert ] Couldn't create certificate")
	}
	return cert, nil
}
