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

package internal

import (
	"bytes"
	"crypto"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/serialization"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/utils"
)

type Identities []Identity

func (is Identities) CreateNodes(discoveries []Identity) (Nodes, error) {
	nodes := make([]Node, len(is))

	for i, identity := range is {
		n, err := identity.CreateNode(discoveries)
		if err != nil {
			return nil, err
		}

		nodes[i] = *n
	}

	return nodes, nil
}

type Identity struct {
	addr       string
	id         insolar.ShortNodeID
	ref        insolar.Reference
	role       insolar.StaticRole
	privateKey crypto.PrivateKey
	publicKey  crypto.PublicKey
}

func (i Identity) createAnnounce(cert insolar.Certificate) ([]byte, *insolar.Signature, error) {
	brief := serialization.NodeBriefIntro{}
	brief.ShortID = i.id
	brief.SetPrimaryRole(adapters.StaticRoleToPrimaryRole(i.role))
	if utils.IsDiscovery(i.ref, cert) {
		brief.SpecialRoles = member.SpecialRoleDiscovery
	}
	brief.StartPower = 10

	addr, err := endpoints.NewIPAddress(i.addr)
	if err != nil {
		return nil, nil, err
	}
	copy(brief.Endpoint[:], addr[:])

	pk, err := keyProcessor.ExportPublicKeyBinary(i.publicKey)
	if err != nil {
		return nil, nil, err
	}

	copy(brief.NodePK[:], pk)

	buf := &bytes.Buffer{}
	err = brief.SerializeTo(nil, buf)
	if err != nil {
		return nil, nil, err
	}

	data := buf.Bytes()
	data = data[:len(data)-64]

	digest := scheme.IntegrityHasher().Hash(data)
	sign, err := scheme.DigestSigner(i.privateKey).Sign(digest)
	if err != nil {
		return nil, nil, err
	}

	return digest, sign, nil
}

func (i Identity) createNetworkNode(cert insolar.Certificate) (insolar.NetworkNode, error) {
	n := node.NewNode(
		i.ref,
		i.role,
		i.publicKey,
		i.addr,
		"",
	)
	mn := n.(node.MutableNode)
	mn.SetShortID(i.id)

	digest, signature, err := i.createAnnounce(cert)
	if err != nil {
		return nil, err
	}

	mn.SetSignature(digest, *signature)

	return mn, err
}

func (i Identity) createCertificate(discoveries []Identity) (insolar.Certificate, error) {
	publicKey, _ := keyProcessor.ExportPublicKeyPEM(i.publicKey)
	bootstrapNodes := make([]certificate.BootstrapNode, len(discoveries))

	for i, discovery := range discoveries {
		publicKeyBytes, err := keyProcessor.ExportPublicKeyPEM(discovery.publicKey)
		if err != nil {
			return nil, err
		}

		bootstrapNodes[i] = *certificate.NewBootstrapNode(
			publicKey,
			string(publicKeyBytes[:]),
			discovery.addr,
			discovery.ref.String(),
		)
	}

	return &certificate.Certificate{
		AuthorizationCertificate: certificate.AuthorizationCertificate{
			PublicKey: string(publicKey[:]),
			Reference: i.ref.String(),
			Role:      i.role.String(),
		},
		BootstrapNodes: bootstrapNodes,
	}, nil
}

func (i Identity) CreateNode(discoveries []Identity) (*Node, error) {
	cert, err := i.createCertificate(discoveries)
	if err != nil {
		return nil, err
	}

	networkNode, err := i.createNetworkNode(cert)
	if err != nil {
		return nil, err
	}

	n := Node{
		networkNode: networkNode,
		profile:     adapters.NewStaticProfile(networkNode, cert, keyProcessor),
		identity:    i,
		certificate: cert,
	}
	return &n, nil
}
