// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package coreapi

type PacketVerifyFlags uint32

const DefaultVerify PacketVerifyFlags = 0

const (
	SkipVerify PacketVerifyFlags = 1 << iota
	RequireStrictVerify
	AllowUnverified
	SuccessfullyVerified
)
