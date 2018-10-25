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
	"crypto/rand"
	"encoding/binary"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/inscontext"
	"github.com/pkg/errors"
)

// NetworkCoordinator encapsulates logic of network configuration
type NetworkCoordinator struct {
	logicRunner   core.LogicRunner
	messageBus    core.MessageBus
	nodeDomainRef *core.RecordRef
	rootDomainRef *core.RecordRef
}

// New creates new NetworkCoordinator
func New() (*NetworkCoordinator, error) {
	return &NetworkCoordinator{}, nil
}

// Start implements interface of Component
func (nc *NetworkCoordinator) Start(c core.Components) error {
	nc.logicRunner = c.LogicRunner
	nc.messageBus = c.MessageBus
	nc.rootDomainRef = c.Bootstrapper.GetRootDomainRef()

	return nil
}

// Stop implements interface of Component
func (nc *NetworkCoordinator) Stop() error {
	return nil
}

func RandomUint64() uint64 {
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return binary.LittleEndian.Uint64(buf)
}

func (nc *NetworkCoordinator) routeCall(ref core.RecordRef, method string, args core.Arguments) (core.Reply, error) {
	if nc.messageBus == nil {
		return nil, errors.New("[ NetworkCoordinator::routeCall ] message bus was not set during initialization")
	}

	e := &message.CallMethod{
		BaseLogicMessage: message.BaseLogicMessage{Nonce: RandomUint64()},
		ObjectRef:        ref,
		Method:           method,
		Arguments:        args,
	}

	res, err := nc.messageBus.Send(inscontext.TODO(), e)
	if err != nil {
		return nil, errors.Wrap(err, "[ NetworkCoordinator::routeCall ] couldn't send message: "+ref.String())
	}

	return res, nil
}

func (nc *NetworkCoordinator) sendRequest(ref *core.RecordRef, method string, argsIn []interface{}) (core.Reply, error) {
	args, err := core.MarshalArgs(argsIn...)
	if err != nil {
		return nil, errors.Wrap(err, "[ NetworkCoordinator::sendRequest ]")
	}

	routResult, err := nc.routeCall(*ref, method, args)
	if err != nil {
		return nil, errors.Wrap(err, "[ NetworkCoordinator::sendRequest ]")
	}

	return routResult, nil
}

func (nc *NetworkCoordinator) getNodeDomainRef() (*core.RecordRef, error) {
	if nc.nodeDomainRef == nil {
		nodeDomainRef, err := nc.fetchNodeDomainRef()
		if err != nil {
			return nil, errors.Wrap(err, "[ getNodeDomainRef ] can't fetch nodeDomainRef")
		}
		nc.nodeDomainRef = nodeDomainRef
	}
	return nc.nodeDomainRef, nil
}

func (nc *NetworkCoordinator) fetchNodeDomainRef() (*core.RecordRef, error) {
	routResult, err := nc.sendRequest(nc.rootDomainRef, "GetNodeDomainRef", []interface{}{})
	if err != nil {
		return nil, errors.Wrap(err, "[ fetchNodeDomainRef ] Can't send request")
	}

	nodeDomainRef, err := extractReferenceResponse(routResult.(*reply.CallMethod).Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ fetchNodeDomainRef ] Can't extract response")
	}

	return nodeDomainRef, nil
}

// WriteActiveNodes writes active nodes to ledger
func (nc *NetworkCoordinator) WriteActiveNodes(number core.PulseNumber, activeNodes []*core.ActiveNode) error {
	return errors.New("not implemented")
}

// Authorize authorizes node by verifying it's signature
func (nc *NetworkCoordinator) Authorize(nodeRef core.RecordRef, seed []byte, signatureRaw []byte) (string, []core.NodeRole, error) {
	nodeDomainRef, err := nc.getNodeDomainRef()
	if err != nil {
		return "", nil, errors.Wrap(err, "[ Authorize ] Can't get nodeDomainRef")
	}

	routResult, err := nc.sendRequest(nodeDomainRef, "Authorize", []interface{}{nodeRef, seed, signatureRaw})

	if err != nil {
		return "", nil, errors.Wrap(err, "[ Authorize ] Can't send request")
	}

	pubKey, role, err := extractAuthorizeResponse(routResult.(*reply.CallMethod).Result)
	if err != nil {
		return "", nil, errors.Wrap(err, "[ Authorize ] Can't extract response")
	}

	return pubKey, role, nil
}

// RegisterNode registers node in nodedomain
func (nc *NetworkCoordinator) RegisterNode(publicKey string, numberOfBootstrapNodes int, majorityRule int, roles []string, ip string) ([]byte, error) {
	nodeDomainRef, err := nc.getNodeDomainRef()
	if err != nil {
		return nil, errors.Wrap(err, "[ RegisterNode ] Can't get nodeDomainRef")
	}
	routResult, err := nc.sendRequest(nodeDomainRef, "RegisterNode", []interface{}{publicKey, numberOfBootstrapNodes, majorityRule, roles, ip})

	if err != nil {
		return nil, errors.Wrap(err, "[ RegisterNode ] Can't send request")
	}

	rawCertificate, err := ExtractRegisterNodeResponse(routResult.(*reply.CallMethod).Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ RegisterNode ] Can't extract response")
	}

	return rawCertificate, nil
}
