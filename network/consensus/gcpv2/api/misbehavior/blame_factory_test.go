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

package misbehavior

import (
	"testing"

	"github.com/insolar/insolar/network/consensus/common/cryptkit"

	"github.com/insolar/insolar/network/consensus/common/endpoints"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"

	"github.com/stretchr/testify/require"
)

func TestBlameIsUnknown(t *testing.T) {
	be := &BlameError{blameType: 0}
	require.False(t, be.IsUnknown())

	be.blameType = 1
	require.True(t, be.IsUnknown())
}

func TestBlameMisbehaviorType(t *testing.T) {
	be := &BlameError{blameType: 0}
	require.Equal(t, Type(1<<32), be.MisbehaviorType())

	be.blameType = 1
	require.Equal(t, Type((1<<32)+1), be.MisbehaviorType())
}

func TestBlameCaptureMark(t *testing.T) {
	cm := interface{}(1)
	be := &BlameError{captureMark: cm}
	require.Equal(t, cm, be.CaptureMark())
}

func TestBlameDetails(t *testing.T) {
	dets := []interface{}{1, 2}
	be := &BlameError{details: dets}
	require.Equal(t, dets, be.Details())
}

func TestBlameViolatorNode(t *testing.T) {
	bn := profiles.NewBaseNodeMock(t)
	be := &BlameError{violatorNode: bn}

	require.Equal(t, bn, be.ViolatorNode())
}

func TestBlameViolatorHost(t *testing.T) {
	inc := endpoints.InboundConnection{}
	be := &BlameError{violatorHost: inc}
	require.Equal(t, inc, be.ViolatorHost())
}

func TestBlameType(t *testing.T) {
	bt := 1
	be := &BlameError{blameType: bt}
	require.Equal(t, bt, be.BlameType())
}

func TestBlameError(t *testing.T) {
	be := &BlameError{}
	require.True(t, be.Error() != "")

	bn := profiles.NewBaseNodeMock(t)
	be.violatorNode = bn
	require.True(t, be.Error() != "")

	be.captureMark = 1
	require.True(t, be.Error() != "")
}

func reportFunc(Report) interface{} {
	return 1
}

func TestNewBlameFactory(t *testing.T) {
	bf := NewBlameFactory(reportFunc)
	require.True(t, bf.capture != nil)
}

func TestNewBlame(t *testing.T) {
	bf := NewBlameFactory(reportFunc)
	fraudType := 1
	msg := "test"
	inc := endpoints.NewInboundMock(t)
	violatorHost := inc
	bn := profiles.NewBaseNodeMock(t)
	violatorNode := bn
	details := []interface{}{1, 2}
	inc.GetNameAddressMock.Set(func() endpoints.Name { return "test" })
	inc.GetTransportKeyMock.Set(func() cryptkit.SignatureKeyHolder { return nil })
	inc.GetTransportCertMock.Set(func() cryptkit.CertificateHolder { return nil })
	be := bf.NewBlame(fraudType, msg, violatorHost, violatorNode, details...)
	require.Equal(t, fraudType, be.blameType)

	require.Equal(t, msg, be.msg)

	require.Equal(t, violatorNode, be.violatorNode)

	require.Equal(t, details[1], be.details[1])

	require.True(t, be.captureMark != nil)

	bf = NewBlameFactory(nil)
	be = bf.NewBlame(fraudType, msg, nil, violatorNode, details...)

	require.True(t, be.captureMark == nil)
}

func TestNewNodeBlame(t *testing.T) {
	bf := NewBlameFactory(reportFunc)
	fraudType := 1
	msg := "test"
	bn := profiles.NewBaseNodeMock(t)
	violatorNode := bn
	details := []interface{}{1, 2}
	be := bf.NewNodeBlame(fraudType, msg, violatorNode, details...)
	require.Equal(t, msg, be.msg)
}

func TestNewHostBlame(t *testing.T) {
	bf := NewBlameFactory(reportFunc)
	fraudType := 1
	msg := "test"
	inc := endpoints.NewInboundMock(t)
	violatorHost := inc
	details := []interface{}{1, 2}
	inc.GetNameAddressMock.Set(func() endpoints.Name { return "test" })
	inc.GetTransportKeyMock.Set(func() cryptkit.SignatureKeyHolder { return nil })
	inc.GetTransportCertMock.Set(func() cryptkit.CertificateHolder { return nil })
	be := bf.NewHostBlame(fraudType, msg, violatorHost, details...)
	require.Equal(t, msg, be.msg)
}

func TestExcessiveIntro(t *testing.T) {
	be := NewBlameFactory(reportFunc).ExcessiveIntro(profiles.NewBaseNodeMock(t), nil, nil)
	require.Equal(t, "excessive intro", be.msg)
}

func TestNewMismatchedPulsarPacket(t *testing.T) {
	be := NewBlameFactory(reportFunc).NewMismatchedPulsarPacket(nil, nil, nil)
	require.Equal(t, "mixed pulsar pulses", be.msg)
}

func TestNewMismatchedPulsePacket(t *testing.T) {
	be := NewBlameFactory(reportFunc).NewMismatchedPulsePacket(nil, nil, nil)
	require.Equal(t, "mixed pulsar pulses", be.msg)
}

func TestNewProtocolViolation(t *testing.T) {
	msg := "test"
	be := NewBlameFactory(reportFunc).NewProtocolViolation(nil, msg)
	require.Equal(t, msg, be.msg)
}
