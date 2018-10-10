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
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/pkg/errors"
)

// NetworkCoordinator encapsulates logic of network configuration
type NetworkCoordinator struct {
	logicRunner   core.LogicRunner
	messageBus    core.MessageBus
	nodeDomainRef core.RecordRef
}

// New creates new NetworkCoordinator
func New() (*NetworkCoordinator, error) {
	return &NetworkCoordinator{}, nil
}

// Start implements interface of Component
func (nc *NetworkCoordinator) Start(c core.Components) error {
	nc.logicRunner = c.LogicRunner
	nc.messageBus = c.MessageBus
	nc.nodeDomainRef = *c.Bootstrapper.GetNodeDomainRef()

	return nil
}

// Stop implements interface of Component
func (nc *NetworkCoordinator) Stop() error {
	return nil
}

func (nc *NetworkCoordinator) routeCall(ref core.RecordRef, method string, args core.Arguments) (core.Reply, error) {
	if nc.messageBus == nil {
		return nil, errors.New("[ NetworkCoordinator::routeCall ] message bus was not set during initialization")
	}

	e := &message.CallMethod{
		ObjectRef: ref,
		Method:    method,
		Arguments: args,
	}

	res, err := nc.messageBus.Send(e)
	if err != nil {
		return nil, errors.Wrap(err, "[ NetworkCoordinator::routeCall ] couldn't send message: "+ref.String())
	}

	return res, nil
}

func (nc *NetworkCoordinator) sendRequest(ref core.RecordRef, method string, argsIn []interface{}) (core.Reply, error) {
	args, err := core.MarshalArgs(argsIn...)
	if err != nil {
		return nil, errors.Wrap(err, "[ NetworkCoordinator::sendRequest ]")
	}

	routResult, err := nc.routeCall(ref, method, args)
	if err != nil {
		return nil, errors.Wrap(err, "[ NetworkCoordinator::sendRequest ]")
	}

	return routResult, nil
}

// Authorize authorizes node by verifying it's signature
func (nc *NetworkCoordinator) Authorize(nodeRef core.RecordRef, seed []byte, signatureRaw []byte) (string, core.NodeRole, error) {
	routResult, err := nc.sendRequest(nc.nodeDomainRef, "Authorize", []interface{}{nodeRef, seed, signatureRaw})
	if err != nil {
		return "", core.RoleUnknown, errors.Wrap(err, "[ Authorize ] Can't send request")
	}

	pubKey, role, err := extractAuthorizeResponse(routResult.(*reply.CallMethod).Result)
	if err != nil {
		return "", core.RoleUnknown, errors.Wrap(err, "[ Authorize ] Can't extract response")
	}

	return pubKey, role, nil
}

// RegisterNode registers node in nodedomain
func (nc *NetworkCoordinator) RegisterNode(publicKey string, role string) (*core.RecordRef, error) {
	routResult, err := nc.sendRequest(nc.nodeDomainRef, "RegisterNode", []interface{}{publicKey, role})
	if err != nil {
		return nil, errors.Wrap(err, "[ RegisterNode ] Can't send request")
	}

	nodeRef, err := extractRegisterNodeResponse(routResult.(*reply.CallMethod).Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ RegisterNode ] Can't extract response")
	}

	return nodeRef, nil
}
