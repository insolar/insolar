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
	"crypto/rand"
	"encoding/binary"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/network"
	"github.com/pkg/errors"
)

// NetworkCoordinator encapsulates logic of network configuration
type NetworkCoordinator struct {
	MessageBus      core.MessageBus   `inject:""`
	Certificate     core.Certificate  `inject:""`
	KeyProcessor    core.KeyProcessor `inject:""`
	NetworkSwitcher network.Switcher  `inject:""`
	nodeDomainRef   *core.RecordRef
	rootDomainRef   *core.RecordRef

	realCoordinator realNetworkCoordinator
	zeroCoordinator zeroNetworkCoordinator
}

// New creates new NetworkCoordinator
func New() (*NetworkCoordinator, error) {
	return &NetworkCoordinator{}, nil
}

func (nc *NetworkCoordinator) Init(ctx context.Context) error {
	return nil
}

// Start implements interface of Component
func (nc *NetworkCoordinator) Start(ctx context.Context) error {
	nc.rootDomainRef = nc.Certificate.GetRootDomainReference()

	return nil
}

func (nc *NetworkCoordinator) getCoordinator() core.NetworkCoordinator {
	return nil
}

// GetCert method returns node certificate
func (nc *NetworkCoordinator) GetCert(ctx context.Context, nodeRef core.RecordRef) (core.Certificate, error) {
	return nil, errors.New("not implemented")
}

// WriteActiveNodes writes active nodes to ledger
func (nc *NetworkCoordinator) ValidateCert(ctx context.Context, number core.PulseNumber, activeNodes []core.Node) (bool, error) {
	return false, errors.New("not implemented")
}

// WriteActiveNodes writes active nodes to ledger
func (nc *NetworkCoordinator) WriteActiveNodes(ctx context.Context, number core.PulseNumber, activeNodes []core.Node) error {
	return errors.New("not implemented")
}

// RandomUint64 generates random uint64
func RandomUint64() uint64 {
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return binary.LittleEndian.Uint64(buf)
}

func (nc *NetworkCoordinator) routeCall(ctx context.Context, ref core.RecordRef, method string, args core.Arguments) (core.Reply, error) {
	if nc.MessageBus == nil {
		return nil, errors.New("[ NetworkCoordinator::routeCall ] message bus was not set during initialization")
	}

	e := &message.CallMethod{
		BaseLogicMessage: message.BaseLogicMessage{Nonce: RandomUint64()},
		ObjectRef:        ref,
		Method:           method,
		Arguments:        args,
	}

	res, err := nc.MessageBus.Send(ctx, e)
	if err != nil {
		return nil, errors.Wrap(err, "[ NetworkCoordinator::routeCall ] couldn't send message: "+ref.String())
	}

	return res, nil
}

func (nc *NetworkCoordinator) sendRequest(ctx context.Context, ref *core.RecordRef, method string, argsIn []interface{}) (core.Reply, error) {
	args, err := core.MarshalArgs(argsIn...)
	if err != nil {
		return nil, errors.Wrap(err, "[ NetworkCoordinator::sendRequest ]")
	}

	routResult, err := nc.routeCall(ctx, *ref, method, args)
	if err != nil {
		return nil, errors.Wrap(err, "[ NetworkCoordinator::sendRequest ]")
	}

	return routResult, nil
}

func (nc *NetworkCoordinator) getNodeDomainRef(ctx context.Context) (*core.RecordRef, error) {
	if nc.nodeDomainRef == nil {
		nodeDomainRef, err := nc.fetchNodeDomainRef(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "[ getNodeDomainRef ] can't fetch nodeDomainRef")
		}
		nc.nodeDomainRef = nodeDomainRef
	}
	return nc.nodeDomainRef, nil
}

func (nc *NetworkCoordinator) fetchNodeDomainRef(ctx context.Context) (*core.RecordRef, error) {
	routResult, err := nc.sendRequest(ctx, nc.rootDomainRef, "GetNodeDomainRef", []interface{}{})
	if err != nil {
		return nil, errors.Wrap(err, "[ fetchNodeDomainRef ] Can't send request")
	}

	nodeDomainRef, err := extractReferenceResponse(routResult.(*reply.CallMethod).Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ fetchNodeDomainRef ] Can't extract response")
	}

	return nodeDomainRef, nil
}
