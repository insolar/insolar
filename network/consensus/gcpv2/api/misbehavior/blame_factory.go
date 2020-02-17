// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package misbehavior

import (
	"fmt"

	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
)

var _ Report = &BlameError{}

type BlameError struct {
	blameType    int
	msg          string
	violatorHost endpoints.InboundConnection
	violatorNode profiles.BaseNode
	details      []interface{}
	captureMark  interface{}
}

func (p *BlameError) IsUnknown() bool {
	return p.blameType != 0
}

func (p *BlameError) MisbehaviorType() Type {
	return Blame.Of(p.blameType)
}

func (p *BlameError) CaptureMark() interface{} {
	return p.captureMark
}

func (p *BlameError) Details() []interface{} {
	return p.details
}

func (p *BlameError) ViolatorNode() profiles.BaseNode {
	return p.violatorNode
}

func (p *BlameError) ViolatorHost() endpoints.InboundConnection {
	return p.violatorHost
}

func (p BlameError) BlameType() int {
	return p.blameType
}

func (p BlameError) Error() string {
	var v interface{} = p.violatorNode
	if v == nil {
		v = p.violatorHost
	}
	c := p.captureMark
	if c == nil {
		c = ""
	}
	return fmt.Sprintf("blame: type=%v(%v)%v, violator=%v, details=%+v", p.blameType, p.msg, c, v, p.details)
}

const (
	_ = iota
	BlameExcessiveIntro
	MismatchedPulsarPacket
	ProtocolViolation
)

func NewBlameFactory(capture ReportFunc) BlameFactory {
	return BlameFactory{capture: capture}
}

type BlameFactory struct {
	capture ReportFunc
}

func (p BlameFactory) NewBlame(fraudType int, msg string, violatorHost endpoints.Inbound, violatorNode profiles.BaseNode, details ...interface{}) BlameError {
	err := BlameError{
		blameType:    fraudType,
		msg:          msg,
		violatorNode: violatorNode,
		details:      details}
	if violatorHost != nil {
		err.violatorHost = endpoints.NewHostIdentityFromHolder(violatorHost)
	}
	if p.capture != nil {
		err.captureMark = p.capture(&err)
	}
	return err
}

func (p BlameFactory) NewNodeBlame(fraudType int, msg string, violatorNode profiles.BaseNode, details ...interface{}) BlameError {
	return p.NewBlame(fraudType, msg, nil, violatorNode, details...)
}

func (p BlameFactory) NewHostBlame(fraudType int, msg string, violatorHost endpoints.Inbound, details ...interface{}) BlameError {
	return p.NewBlame(fraudType, msg, violatorHost, nil, details...)
}

func (p BlameFactory) ExcessiveIntro(violator profiles.BaseNode, evidence1 cryptkit.SignedEvidenceHolder, evidence2 cryptkit.SignedEvidenceHolder) BlameError {
	return p.NewNodeBlame(BlameExcessiveIntro, "excessive intro", violator, evidence1, evidence2)
}

func (p BlameFactory) NewMismatchedPulsarPacket(from endpoints.Inbound, expected proofs.OriginalPulsarPacket, received proofs.OriginalPulsarPacket) BlameError {
	return p.NewHostBlame(MismatchedPulsarPacket, "mixed pulsar pulses", from, expected, received)
}

func (p BlameFactory) NewMismatchedPulsePacket(from profiles.BaseNode, expected proofs.OriginalPulsarPacket, received proofs.OriginalPulsarPacket) BlameError {
	return p.NewNodeBlame(MismatchedPulsarPacket, "mixed pulsar pulses", from, expected, received)
}

func (p BlameFactory) NewProtocolViolation(violator profiles.BaseNode, msg string) BlameError {
	return p.NewNodeBlame(ProtocolViolation, msg, violator)
}
