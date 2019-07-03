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

package bootstrap

import (
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
)

const PERMIT_TTL = 300

// CreatePermit creates permit as signed protobuf for joiner node to
func CreatePermit(authorityNodeRef insolar.Reference, reconnectHost *host.Host, joinerPublicKey []byte, signer insolar.Signer) (*packet.Permit, error) {
	payload := packet.PermitPayload{
		AuthorityNodeRef: authorityNodeRef,
		ExpireTimestamp:  time.Now().Unix() + PERMIT_TTL,
		ReconnectTo:      reconnectHost,
		JoinerPublicKey:  joinerPublicKey,
	}

	data, err := payload.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal bootstrap permit")
	}
	signature, err := signer.Sign(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign bootstrap permit")
	}
	return &packet.Permit{Payload: payload, Signature: signature.Bytes()}, nil
}

// ValidatePermit validate granted permit and verifies signature of Authority Node
func ValidatePermit(permit *packet.Permit, cert insolar.Certificate, verifier insolar.CryptographyService) error {
	discovery := network.FindDiscoveryByRef(cert, permit.Payload.AuthorityNodeRef)
	if discovery == nil {
		return errors.New("failed to find a discovery node from reference in permit")
	}

	payload, err := permit.Payload.Marshal()
	if err != nil || payload == nil {
		return errors.New("failed to marshal bootstrap permission payload part")
	}

	verified := verifier.Verify(discovery.GetPublicKey(), insolar.SignatureFromBytes(permit.Signature), payload)

	if !verified {
		return errors.New("bootstrap permission payload verification failed")
	}
	return nil
}
