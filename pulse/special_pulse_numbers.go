// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulse

// Number is a type for pulse numbers.
//
// Special values:
// 0 					Unknown
// 1 .. 256				Reserved for package internal usage
// 257 .. 65535			Reserved for platform wide usage
// 65536				Local relative pulse number
// 65537 .. 1<<30 - 1	Regular time based pulse numbers
//
// NB! Range 0..256 IS RESERVED for internal operations
// There MUST BE NO references with PN < 256 ever visible to contracts / users.
const (
	_ Number = 256 + iota

	// Jet is a special pulse number value that signifies jet ID.
	Jet

	// BuiltinContract declares special pulse number that creates namespace for builtin contracts
	BuiltinContract
)

func (n Number) IsJet() bool {
	return n == Jet
}
