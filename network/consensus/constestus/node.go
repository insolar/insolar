//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package constestus

import (
	"context"
	"crypto"
	"fmt"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/testutils"
)

type nodeIdentity struct {
	addr       string
	id         insolar.ShortNodeID
	ref        insolar.Reference
	role       insolar.StaticRole
	privateKey crypto.PrivateKey
	publicKey  crypto.PublicKey
}

func (i nodeIdentity) createNode() insolar.NetworkNode {
	n := node.NewNode(
		i.ref,
		i.role,
		i.publicKey,
		i.addr,
		"",
	)
	mn := n.(node.MutableNode)
	mn.SetShortID(i.id)

	return mn
}

type identityGenerator struct {
	keyProcessor insolar.KeyProcessor

	baseAddr string

	mu         *sync.Mutex
	portOffset uint16
	idOffset   uint32
}

func (g *identityGenerator) generateShared() (insolar.ShortNodeID, uint16) {
	g.mu.Lock()
	defer g.mu.Unlock()

	id := g.idOffset
	g.idOffset++

	port := g.portOffset
	g.portOffset++

	return insolar.ShortNodeID(id), port
}

func (g *identityGenerator) generateIdentity(role insolar.StaticRole) (*nodeIdentity, error) {
	privateKey, err := g.keyProcessor.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	id, port := g.generateShared()

	identity := &nodeIdentity{
		addr:       fmt.Sprintf("%s:%d", g.baseAddr, port),
		id:         id,
		ref:        testutils.RandomRef(),
		role:       role,
		privateKey: privateKey,
		publicKey:  g.keyProcessor.ExtractPublicKey(privateKey),
	}

	return identity, nil
}

type nodeComponents struct {
	controller   consensus.Controller
	nodeKeeper   network.NodeKeeper
	transport    transport.DatagramTransport
	pulseHandler network.PulseHandler
}

type networkNode struct {
	identity   nodeIdentity
	components nodeComponents
	ctx        context.Context
}

func (n networkNode) Connect() error {
	return n.components.transport.Start(n.ctx)
}

func (n networkNode) Join(cloud C) error {
	panic("not implemented")
}

func (n networkNode) Disconnect() error {
	if err := n.components.transport.Stop(n.ctx); err != nil {
		return err
	}

	return nil
}

func (n networkNode) Leave(reason uint32) error {
	<-n.components.controller.Leave(reason)
	if err := n.Disconnect(); err != nil {
		return err
	}

	n.components.controller.Abort()
	return nil
}

type cloudNode struct {
	C
	networkNode
}

func (n cloudNode) Connect() {
	n.C.Require().NoError(n.networkNode.Connect())
}

func (n cloudNode) Join(cloud C) {
	n.C.Require().NoError(n.networkNode.Join(cloud))
}

func (n cloudNode) Disconnect() {
	n.C.Require().NoError(n.networkNode.Disconnect())
}

func (n cloudNode) Leave(reason uint32) {
	n.C.Require().NoError(n.networkNode.Leave(reason))
}

func (n cloudNode) Intercept(nodes ...N) TypedInterceptor {
	panic("not implemented")
}
