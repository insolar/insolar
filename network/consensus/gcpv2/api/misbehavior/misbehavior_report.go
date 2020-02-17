// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package misbehavior

import (
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
)

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior.Report -o . -s _mock.go -g

type Report interface {
	CaptureMark() interface{}
	Details() []interface{}
	ViolatorNode() profiles.BaseNode
	ViolatorHost() endpoints.InboundConnection
	MisbehaviorType() Type
}

func Is(err error) bool {
	_, ok := err.(Report)
	return ok
}

func Of(err error) Report {
	rep, ok := err.(Report)
	if ok {
		return rep
	}
	return nil
}

type ReportFunc func(report Report) interface{}

type Type uint64
type Category int

const (
	_ Category = iota
	Blame
	Fraud
)

func (c Type) Category() Category {
	return Category(c >> 32)
}

func (c Type) Type() int {
	return int(c & (1<<32 - 1))
}

func (c Category) Of(misbehavior int) Type {
	return Type(c<<32) | Type(misbehavior&(1<<32-1))
}
