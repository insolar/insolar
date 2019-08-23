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

package adapters

import (
	"crypto/ecdsa"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/adapters/candidate"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
)

type Candidate candidate.Profile

func (c Candidate) StaticProfile(keyProcessor insolar.KeyProcessor) *StaticProfile {
	publicKey, err := keyProcessor.ImportPublicKeyBinary(c.PublicKey)
	if err != nil {
		panic("Failed to import public key")
	}

	signHolder := cryptkit.NewSignature(
		longbits.NewBits512FromBytes(c.Signature),
		SHA3512Digest.SignedBy(SECP256r1Sign),
	).AsSignatureHolder()

	extension := newStaticProfileExtension(
		c.ShortID,
		c.Ref,
		signHolder,
	)

	return newStaticProfile(
		c.ShortID,
		c.PrimaryRole,
		c.SpecialRole,
		extension,
		NewOutbound(c.Address),
		NewECDSAPublicKeyStore(publicKey.(*ecdsa.PublicKey)),
		NewECDSASignatureKeyHolder(publicKey.(*ecdsa.PublicKey), keyProcessor),
		cryptkit.NewSignedDigest(
			cryptkit.NewDigest(longbits.NewBits512FromBytes(c.Digest), SHA3512Digest),
			cryptkit.NewSignature(longbits.NewBits512FromBytes(c.Signature), SHA3512Digest.SignedBy(SECP256r1Sign)),
		).AsSignedDigestHolder(),
	)
}

func (c Candidate) Profile() candidate.Profile {
	return candidate.Profile(c)
}

func NewCandidate(staticProfile *StaticProfile, keyProcessor insolar.KeyProcessor) *Candidate {
	pubKey, err := keyProcessor.ExportPublicKeyBinary(staticProfile.store.(*ECDSAPublicKeyStore).publicKey)
	if err != nil {
		panic("failed to export public key")
	}

	signedDigest := staticProfile.GetBriefIntroSignedDigest()

	return &Candidate{
		Address:     staticProfile.GetDefaultEndpoint().GetIPAddress().String(),
		Ref:         staticProfile.GetExtension().GetReference(),
		ShortID:     staticProfile.GetStaticNodeID(),
		PrimaryRole: staticProfile.GetPrimaryRole(),
		SpecialRole: staticProfile.GetSpecialRoles(),
		Digest:      signedDigest.GetDigestHolder().AsBytes(),
		Signature:   signedDigest.GetSignatureHolder().AsBytes(),
		PublicKey:   pubKey,
	}
}
