// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
