// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
	require.NotEmpty(t, be.Error())

	bn := profiles.NewBaseNodeMock(t)
	be.violatorNode = bn
	require.NotEmpty(t, be.Error())

	be.captureMark = 1
	require.NotEmpty(t, be.Error())
}

func reportFunc(_ Report) interface{} {
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
