/*
 *    Copyright 2019 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package core

import (
	"fmt"
)

// JetID should be used, when id is a jetID
type JetID RecordID

// ZeroJetID is value of an empty Jet ID
var ZeroJetID = *NewJetID(0, nil)

// NewJetID creates a new jet with provided ID and index
func NewJetID(depth uint8, prefix []byte) *JetID {
	var id JetID
	copy(id[:PulseNumberSize], PulseNumberJet.Bytes())
	id[PulseNumberSize] = depth
	copy(id[PulseNumberSize+1:], prefix)
	return &id
}

// Jet extracts depth and prefix from jet id.
func (id JetID) Jet() (uint8, []byte) {
	recordID := RecordID(id)
	if recordID.Pulse() != PulseNumberJet {
		panic(fmt.Sprintf("provided id %b is not a jet id", id))
	}
	return id[PulseNumberSize], id[PulseNumberSize+1:]
}

// Parent returns a parent of the jet
func (id JetID) Parent() JetID {
	depth, prefix := id.Jet()
	if depth == 0 {
		return id
	}

	return *NewJetID(depth-1, ResetBits(prefix, depth-1))
}

// ResetBits returns a new byte slice with all bits in 'value' reset, starting from 'start' number of bit. If 'start'
// is bigger than len(value), the original slice will be returned.
func ResetBits(value []byte, start uint8) []byte {
	if int(start) >= len(value)*8 {
		return value
	}

	startByte := start / 8
	startBit := start % 8

	result := make([]byte, len(value))
	copy(result, value[:startByte])

	// Reset bits in starting byte.
	mask := byte(0xFF)
	mask <<= 8 - byte(startBit)
	result[startByte] = value[startByte] & mask

	return result
}
