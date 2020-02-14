// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package misbehavior

import (
	"fmt"

	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
)

func IsFraud(err error) bool {
	_, ok := err.(*FraudError)
	return ok
}

func FraudOf(err error) *FraudError {
	rep, ok := err.(*FraudError)
	if ok {
		return rep
	}
	return nil
}

var _ Report = &FraudError{}

type FraudError struct {
	fraudType    int
	msg          string
	violatorHost endpoints.InboundConnection
	violatorNode profiles.BaseNode
	details      []interface{}
	captureMark  interface{}
}

func (p *FraudError) IsUnknown() bool {
	return p.fraudType == 0
}

func (p *FraudError) MisbehaviorType() Type {
	return Fraud.Of(p.fraudType)
}

func (p *FraudError) CaptureMark() interface{} {
	return p.captureMark
}

func (p *FraudError) Details() []interface{} {
	return p.details
}

func (p *FraudError) ViolatorNode() profiles.BaseNode {
	return p.violatorNode
}

func (p *FraudError) ViolatorHost() endpoints.InboundConnection {
	return p.violatorHost
}

func (p FraudError) FraudType() int {
	return p.fraudType
}

func (p FraudError) Error() string {
	var v interface{} = p.violatorNode
	if v == nil {
		v = p.violatorHost
	}
	c := p.captureMark
	if c == nil {
		c = ""
	}
	return fmt.Sprintf("fraud: type=%v(%v)%v, violator=%v, details=%+v", p.fraudType, p.msg, c, v, p.details)
}

const (
	_ = iota
	FraudMultipleNsh
	MismatchedRank
	MismatchedNeighbour
	WrongPower
)

func NewFraudFactory(capture ReportFunc) FraudFactory {
	return FraudFactory{capture: capture}
}

type FraudFactory struct {
	capture ReportFunc
}

func (p FraudFactory) NewFraud(fraudType int, msg string, violatorHost endpoints.Inbound, violatorNode profiles.BaseNode, details ...interface{}) FraudError {
	err := FraudError{
		fraudType:    fraudType,
		msg:          msg,
		violatorNode: violatorNode,
		details:      details,
	}
	if violatorHost != nil {
		err.violatorHost = endpoints.NewHostIdentityFromHolder(violatorHost)
	}
	if p.capture != nil {
		err.captureMark = p.capture(&err)
	}
	return err
}

func (p FraudFactory) NewNodeFraud(fraudType int, msg string, violatorNode profiles.BaseNode, details ...interface{}) FraudError {
	return p.NewFraud(fraudType, msg, nil, violatorNode, details...)
}

func (p FraudFactory) NewHostFraud(fraudType int, msg string, violatorHost endpoints.Inbound, details ...interface{}) FraudError {
	return p.NewFraud(fraudType, msg, violatorHost, nil, details...)
}

func (p FraudFactory) NewInconsistentMembershipAnnouncement(violator profiles.ActiveNode,
	evidence1 profiles.MembershipAnnouncement, evidence2 profiles.MembershipAnnouncement) FraudError {
	return p.NewNodeFraud(FraudMultipleNsh, "multiple membership profile", violator, evidence1, evidence2)
}

func (p FraudFactory) NewMismatchedMembershipRank(violator profiles.ActiveNode, mp profiles.MembershipProfile) FraudError {
	return p.NewNodeFraud(MismatchedRank, "mismatched membership profile rank", violator, mp)
}

func (p FraudFactory) NewMismatchedMembershipRankOrNodeCount(violator profiles.ActiveNode, mp profiles.MembershipProfile, nodeCount int) FraudError {
	return p.NewNodeFraud(MismatchedRank, "mismatched membership profile node count", violator, mp, nodeCount)
}

func (p FraudFactory) NewUnknownNeighbour(violator profiles.BaseNode) error {
	return p.NewNodeFraud(MismatchedNeighbour, "unknown neighbour", violator)
}

func (p FraudFactory) NewMismatchedNeighbourRank(violator profiles.BaseNode) error {
	return p.NewNodeFraud(MismatchedNeighbour, "mismatched neighbour rank", violator)
}

func (p FraudFactory) NewNeighbourMissingTarget(violator profiles.BaseNode) error {
	return p.NewNodeFraud(MismatchedNeighbour, "neighbour must include target node", violator)
}

func (p FraudFactory) NewNeighbourContainsSource(violator profiles.BaseNode) error {
	return p.NewNodeFraud(MismatchedNeighbour, "neighbour must NOT include source node", violator)
}

func (p FraudFactory) NewInconsistentNeighbourAnnouncement(violator profiles.BaseNode) FraudError {
	return p.NewNodeFraud(MismatchedNeighbour, "multiple neighbour profile", violator)
}

func (p FraudFactory) NewInvalidPowerLevel(violator profiles.BaseNode) FraudError {
	return p.NewNodeFraud(WrongPower, "power level is incorrect", violator)
}
