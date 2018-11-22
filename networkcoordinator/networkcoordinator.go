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
	"github.com/pkg/errors"
)

// NetworkCoordinator encapsulates logic of network configuration
type NetworkCoordinator struct {
	Certificate         core.Certificate         `inject:""`
	KeyProcessor        core.KeyProcessor        `inject:""`
	ContractRequester   core.ContractRequester   `inject:""`
	GenesisDataProvider core.GenesisDataProvider `inject:""`
}

// New creates new NetworkCoordinator
func New() (*NetworkCoordinator, error) {
	return &NetworkCoordinator{}, nil
}

// WriteActiveNodes writes active nodes to ledger
func (nc *NetworkCoordinator) WriteActiveNodes(ctx context.Context, number core.PulseNumber, activeNodes []core.Node) error {
	return errors.New("not implemented")
}

// Authorize authorizes node by verifying it's signature
/*func (nc *NetworkCoordinator) Authorize(ctx context.Context, nodeRef core.RecordRef, seed []byte, signatureRaw []byte) (string, core.NodeRole, error) {
	nodeDomainRef, err := nc.GenesisDataProvider.GetNodeDomain(ctx)
	if err != nil {
		return "", core.RoleUnknown, errors.Wrap(err, "[ Authorize ] Can't get nodeDomainRef")
	}

	routResult, err := nc.ContractRequester.SendRequest(ctx, nodeDomainRef, "Authorize", []interface{}{nodeRef, seed, signatureRaw})

	if err != nil {
		return "", core.RoleUnknown, errors.Wrap(err, "[ Authorize ] Can't send request")
	}

	pubKey, role, err := extractor.ExtractAuthorizeResponse(routResult.(*reply.CallMethod).Result)
	if err != nil {
		return "", core.RoleUnknown, errors.Wrap(err, "[ Authorize ] Can't extract response")
	}

	return pubKey, role, nil
}*/

// RegisterNode registers node in nodedomain
/*func (nc *NetworkCoordinator) RegisterNode(ctx context.Context, publicKey crypto.PublicKey, numberOfBootstrapNodes int, majorityRule int, role string, ip string) ([]byte, error) {
	nodeDomainRef, err := nc.GenesisDataProvider.GetNodeDomain(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "[ RegisterNode ] Can't get nodeDomainRef")
	}
	publicKeyStr, err := nc.KeyProcessor.ExportPublicKey(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ RegisterNode ] Can't import public key")
	}

	routResult, err := nc.ContractRequester.SendRequest(ctx, nodeDomainRef, "RegisterNode", []interface{}{publicKeyStr, numberOfBootstrapNodes, majorityRule, role})
	if err != nil {
		return nil, errors.Wrap(err, "[ RegisterNode ] Can't send request")
	}

	rawCertificate, err := extractor.ExtractRegisterNodeResponse(routResult.(*reply.CallMethod).Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ RegisterNode ] Can't extract response")
	}

	return rawCertificate, nil
}
*/
