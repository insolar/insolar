// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package appfoundation

// VestingType type of vesting process
type VestingType int

//go:generate stringer -type=VestingType
const (
	DefaultVesting VestingType = iota
	Vesting1
	Vesting2
	Vesting3
	Vesting4
)
