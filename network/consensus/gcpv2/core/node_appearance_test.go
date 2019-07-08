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

package core

import (
	"fmt"
	"testing"

	"github.com/insolar/insolar/network/consensus/gcpv2/errors"

	"github.com/insolar/insolar/network/consensus/common"
	gcommon "github.com/insolar/insolar/network/consensus/gcpv2/common"

	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
	"github.com/stretchr/testify/require"
)

func TestNewNodeAppearanceAsSelf(t *testing.T) {
	lp := gcommon.NewLocalNodeProfileMock(t)
	lp.LocalNodeProfileMock.Set(func() {})
	callback := &nodeContext{}
	r := NewNodeAppearanceAsSelf(lp, callback)
	require.Equal(t, r.state, packets.NodeStateLocalActive)

	require.Equal(t, r.trust, packets.SelfTrust)

	require.Equal(t, r.profile, lp)

	require.Equal(t, r.callback, callback)

	require.NotEqual(t, r.announceHandler, nil)
}

func TestInit(t *testing.T) {
	lp := gcommon.NewLocalNodeProfileMock(t)
	lp.LocalNodeProfileMock.Set(func() {})
	callback := &nodeContext{}
	r := NewNodeAppearanceAsSelf(lp, callback)
	require.Panics(t, func() { r.init(nil, callback, 0) })

	r.init(lp, callback, 0)
	require.Equal(t, r.state, packets.NodeStateLocalActive)

	require.Equal(t, r.trust, packets.SelfTrust)

	require.Equal(t, r.profile, lp)

	require.Equal(t, r.callback, callback)
}

func TestString(t *testing.T) {
	lp := gcommon.NewLocalNodeProfileMock(t)
	lp.LocalNodeProfileMock.Set(func() {})
	callback := &nodeContext{}
	r := NewNodeAppearanceAsSelf(lp, callback)
	require.Equal(t, r.String(), fmt.Sprintf("node:{%v}", lp))
}

func TestLessByNeighbourWeightForNodeAppearance(t *testing.T) {
	lp := gcommon.NewLocalNodeProfileMock(t)
	lp.LocalNodeProfileMock.Set(func() {})
	callback := &nodeContext{}
	r1 := NewNodeAppearanceAsSelf(lp, callback)
	r2 := NewNodeAppearanceAsSelf(lp, callback)
	r1.neighbourWeight = 0
	r2.neighbourWeight = 1
	require.True(t, LessByNeighbourWeightForNodeAppearance(r1, r2))

	require.False(t, LessByNeighbourWeightForNodeAppearance(r2, r1))

	r2.neighbourWeight = 0
	require.False(t, LessByNeighbourWeightForNodeAppearance(r2, r1))
}

func TestCopySelfTo(t *testing.T) {
	lp := gcommon.NewLocalNodeProfileMock(t)
	lp.LocalNodeProfileMock.Set(func() {})
	callback := &nodeContext{}

	source := NewNodeAppearanceAsSelf(lp, callback)
	source.stateEvidence = gcommon.NewNodeStateHashEvidenceMock(t)
	source.announceSignature = gcommon.NewMemberAnnouncementSignatureMock(t)
	source.requestedPower = 1
	source.state = packets.NodeStateLocalActive
	source.trust = packets.TrustBySome

	target := NewNodeAppearanceAsSelf(lp, callback)
	target.stateEvidence = gcommon.NewNodeStateHashEvidenceMock(t)
	target.announceSignature = gcommon.NewMemberAnnouncementSignatureMock(t)
	target.requestedPower = 2
	target.state = packets.NodeStateReceivedPhases
	target.trust = packets.TrustByNeighbors

	target.copySelfTo(source)

	require.Equal(t, source.stateEvidence, target.stateEvidence)

	require.Equal(t, source.announceSignature, target.announceSignature)

	require.Equal(t, source.requestedPower, target.requestedPower)

	require.Equal(t, source.state, target.state)

	require.Equal(t, source.trust, target.trust)
}

func TestIsJoiner(t *testing.T) {
	lp := gcommon.NewLocalNodeProfileMock(t)
	lp.LocalNodeProfileMock.Set(func() {})
	lp.GetStateMock.Set(func() gcommon.MembershipState { return gcommon.Undefined })
	callback := &nodeContext{}

	r := NewNodeAppearanceAsSelf(lp, callback)
	require.False(t, r.IsJoiner())
}

func TestGetIndex(t *testing.T) {
	lp := gcommon.NewLocalNodeProfileMock(t)
	lp.LocalNodeProfileMock.Set(func() {})
	index := 1
	lp.GetIndexMock.Set(func() int { return index })
	callback := &nodeContext{}
	r := NewNodeAppearanceAsSelf(lp, callback)
	require.Equal(t, r.GetIndex(), index)
}

func TestGetShortNodeID(t *testing.T) {
	lp := gcommon.NewLocalNodeProfileMock(t)
	lp.LocalNodeProfileMock.Set(func() {})
	lp.GetShortNodeIDMock.Set(func() common.ShortNodeID { return common.AbsentShortNodeID })
	callback := &nodeContext{}
	r := NewNodeAppearanceAsSelf(lp, callback)
	require.Equal(t, r.GetShortNodeID(), common.AbsentShortNodeID)
}

func TestGetTrustLevel(t *testing.T) {
	lp := gcommon.NewLocalNodeProfileMock(t)
	lp.LocalNodeProfileMock.Set(func() {})
	callback := &nodeContext{}
	r := NewNodeAppearanceAsSelf(lp, callback)
	r.trust = packets.TrustBySome
	require.Equal(t, r.GetTrustLevel(), packets.TrustBySome)
}

func TestGetProfile(t *testing.T) {
	lp := gcommon.NewLocalNodeProfileMock(t)
	lp.LocalNodeProfileMock.Set(func() {})
	callback := &nodeContext{}
	r := NewNodeAppearanceAsSelf(lp, callback)
	require.Equal(t, r.GetProfile(), lp)
}

func TestVerifyPacketAuthenticity(t *testing.T) {
	lp := gcommon.NewLocalNodeProfileMock(t)
	lp.LocalNodeProfileMock.Set(func() {})
	var isAcceptable bool
	lp.IsAcceptableHostMock.Set(func(common.HostIdentityHolder) bool { return *(&isAcceptable) })
	sv := common.NewSignatureVerifierMock(t)
	var isSignOfSignatureMethodSupported bool
	sv.IsSignOfSignatureMethodSupportedMock.Set(func(common.SignatureMethod) bool { return *(&isSignOfSignatureMethodSupported) })
	var isValidDigestSignature bool
	sv.IsValidDigestSignatureMock.Set(func(common.DigestHolder, common.SignatureHolder) bool { return *(&isValidDigestSignature) })
	lp.GetSignatureVerifierMock.Set(func() common.SignatureVerifier { return sv })
	callback := &nodeContext{}
	r := NewNodeAppearanceAsSelf(lp, callback)
	packet := packets.NewPacketParserMock(t)
	packet.GetPacketSignatureMock.Set(func() common.SignedDigest { return common.SignedDigest{} })
	from := common.NewHostIdentityHolderMock(t)
	strictFrom := true
	isAcceptable = false
	require.NotEqual(t, r.VerifyPacketAuthenticity(packet, from, strictFrom), nil)

	strictFrom = false
	isSignOfSignatureMethodSupported = false
	require.NotEqual(t, r.VerifyPacketAuthenticity(packet, from, strictFrom), nil)

	isSignOfSignatureMethodSupported = true
	isValidDigestSignature = false
	require.NotEqual(t, r.VerifyPacketAuthenticity(packet, from, strictFrom), nil)

	isValidDigestSignature = true
	require.Equal(t, r.VerifyPacketAuthenticity(packet, from, strictFrom), nil)
}

func TestSetReceivedPhase(t *testing.T) {
	lp := gcommon.NewLocalNodeProfileMock(t)
	lp.LocalNodeProfileMock.Set(func() {})
	callback := &nodeContext{}
	r := NewNodeAppearanceAsSelf(lp, callback)
	require.True(t, r.SetReceivedPhase(packets.Phase1))

	require.False(t, r.SetReceivedPhase(packets.Phase1))
}

func TestSetReceivedByPacketType(t *testing.T) {
	lp := gcommon.NewLocalNodeProfileMock(t)
	lp.LocalNodeProfileMock.Set(func() {})
	callback := &nodeContext{}
	r := NewNodeAppearanceAsSelf(lp, callback)
	require.True(t, r.SetReceivedByPacketType(packets.PacketPhase1))

	require.False(t, r.SetReceivedByPacketType(packets.PacketPhase1))

	require.False(t, r.SetReceivedByPacketType(packets.MaxPacketType))
}

func TestSetSentPhase(t *testing.T) {
	lp := gcommon.NewLocalNodeProfileMock(t)
	lp.LocalNodeProfileMock.Set(func() {})
	callback := &nodeContext{}
	r := NewNodeAppearanceAsSelf(lp, callback)
	require.True(t, r.SetSentPhase(packets.Phase1))

	require.False(t, r.SetSentPhase(packets.Phase1))
}

func TestSetSentByPacketType(t *testing.T) {
	lp := gcommon.NewLocalNodeProfileMock(t)
	lp.LocalNodeProfileMock.Set(func() {})
	callback := &nodeContext{}
	r := NewNodeAppearanceAsSelf(lp, callback)
	require.True(t, r.SetSentByPacketType(packets.PacketPhase1))

	require.True(t, r.SetSentByPacketType(packets.PacketPhase1))

	require.False(t, r.SetSentByPacketType(packets.MaxPacketType))
}

func TestSetReceivedWithDupCheck(t *testing.T) {
	lp := gcommon.NewLocalNodeProfileMock(t)
	lp.LocalNodeProfileMock.Set(func() {})
	callback := &nodeContext{}
	r := NewNodeAppearanceAsSelf(lp, callback)
	require.Equal(t, r.SetReceivedWithDupCheck(packets.PacketPhase1), nil)

	require.Equal(t, r.SetReceivedWithDupCheck(packets.PacketPhase1), errors.ErrRepeatedPhasePacket)

	require.Equal(t, r.SetReceivedWithDupCheck(packets.MaxPacketType), errors.ErrRepeatedPhasePacket)
}

func TestGetSignatureVerifier(t *testing.T) {
	lp := gcommon.NewLocalNodeProfileMock(t)
	lp.LocalNodeProfileMock.Set(func() {})
	sv1 := common.NewSignatureVerifierMock(t)
	lp.GetSignatureVerifierMock.Set(func() common.SignatureVerifier { return sv1 })
	lp.GetNodePublicKeyStoreMock.Set(func() common.PublicKeyStore { return nil })
	callback := &nodeContext{}
	r := NewNodeAppearanceAsSelf(lp, callback)
	svf := common.NewSignatureVerifierFactoryMock(t)
	sv2 := common.NewSignatureVerifierMock(t)
	svf.GetSignatureVerifierWithPKSMock.Set(func(common.PublicKeyStore) common.SignatureVerifier { return sv2 })
	require.Equal(t, r.GetSignatureVerifier(svf), sv1)

	lp.GetSignatureVerifierMock.Set(func() common.SignatureVerifier { return nil })
	require.Equal(t, r.GetSignatureVerifier(svf), sv2)
}

func TestCreateSignatureVerifier(t *testing.T) {
	lp := gcommon.NewLocalNodeProfileMock(t)
	lp.LocalNodeProfileMock.Set(func() {})
	lp.GetNodePublicKeyStoreMock.Set(func() common.PublicKeyStore { return nil })
	callback := &nodeContext{}
	r := NewNodeAppearanceAsSelf(lp, callback)
	svf := common.NewSignatureVerifierFactoryMock(t)
	sv := common.NewSignatureVerifierMock(t)
	svf.GetSignatureVerifierWithPKSMock.Set(func(common.PublicKeyStore) common.SignatureVerifier { return sv })
	require.Equal(t, r.CreateSignatureVerifier(svf), sv)
}
