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
	"strconv"
	"strings"
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

// Depth extracts depth from a jet id.
func (id JetID) Depth() uint8 {
	recordID := RecordID(id)
	if recordID.Pulse() != PulseNumberJet {
		panic(fmt.Sprintf("provided id %b is not a jet id", id))
	}
	return id[PulseNumberSize]
}

// Prefix extracts prefix from a jet id.
func (id JetID) Prefix() []byte {
	recordID := RecordID(id)
	if recordID.Pulse() != PulseNumberJet {
		panic(fmt.Sprintf("provided id %b is not a jet id", id))
	}
	return id[PulseNumberSize+1:]
}

// DebugString returns a debug representation of a jet
func (id JetID) DebugString() string {
	pulse := NewPulseNumber(id[:PulseNumberSize])
	if pulse != PulseNumberJet {
		return fmt.Sprintf("[JET: <wrong pulse number>]")
	}

	depth := int(id[PulseNumberSize])
	if depth == 0 {
		return "[JET 0 -]"
	}

	prefix := id[PulseNumberSize+1:]
	var res strings.Builder
	res.WriteString("[JET ")
	res.WriteString(strconv.Itoa(depth))
	res.WriteString(" ")

	for _, b := range prefix {
		for j := 7; j >= 0; j-- {
			if 0 == (b >> uint(j) & 0x01) {
				res.WriteString("0")
			} else {
				res.WriteString("1")
			}

			depth--
			if depth == 0 {
				res.WriteString("]")
				return res.String()
			}
		}
	}

	return fmt.Sprintf("[JET: <wrong format> %d %b]", depth, prefix)
}
