//
// Copyright 2019 Insolar Technologies GmbH
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
//

package insolar

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	// JetSize is a Jet's size (depth+prefix).
	JetSize = RecordIDSize - PulseNumberSize
	// JetPrefixSize is a Jet's prefix size.
	JetPrefixSize = JetSize - 1
	// JetMaximumDepth is a Jet's maximum depth (maximum offset in bits).
	JetMaximumDepth = JetPrefixSize*8 - 1
	// JetPrefixOffset is an offset where prefix starts in jet id.
	JetPrefixOffset = PulseNumberSize + 1
)

// JetID should be used, when id is a jetID
type JetID ID

// Size is a protobuf required method. It returns size of JetID
func (id *JetID) Size() int { return RecordIDSize }

// MarshalTo is a protobuf required method. It marshals data
func (id *JetID) MarshalTo(data []byte) (n int, err error) {
	copy(data, id[:])
	return RecordIDSize, nil
}

// Unmarshal is a protobuf required method. It unmarshals data
func (id *JetID) Unmarshal(data []byte) error {
	if len(data) != RecordIDSize {
		return errors.New("Not enough bytes to unpack JetID")
	}
	copy(id[:], data)
	return nil
}

// ZeroJetID is value of an empty Jet ID
var ZeroJetID = *NewJetID(0, nil)

// NewJetID creates a new jet with provided ID and index
func NewJetID(depth uint8, prefix []byte) *JetID {
	var id JetID
	copy(id[:PulseNumberSize], PulseNumberJet.Bytes())
	id[PulseNumberSize] = depth
	copy(id[JetPrefixOffset:], prefix)
	return &id
}

// Depth extracts depth from a jet id.
func (id JetID) Depth() uint8 {
	recordID := ID(id)
	if recordID.Pulse() != PulseNumberJet {
		panic(fmt.Sprintf("provided id %b is not a jet id", id))
	}
	return id[PulseNumberSize]
}

// Prefix extracts prefix from a jet id.
func (id JetID) Prefix() []byte {
	recordID := ID(id)
	if recordID.Pulse() != PulseNumberJet {
		panic(fmt.Sprintf("provided id %b is not a jet id", id))
	}
	return id[JetPrefixOffset:]
}

// DebugString prints JetID in human readable form.
func (id JetID) DebugString() string {
	depth := int(id[PulseNumberSize])
	if depth == 0 {
		return "[JET 0 -]"
	}

	prefix := id[PulseNumberSize+1:]
	var res strings.Builder
	res.WriteString("[JET ")
	res.WriteString(strconv.Itoa(depth))
	res.WriteString(" ")
	if len(prefix)*8 < depth {
		return fmt.Sprintf("[JET: <wrong format> %d %b]", depth, prefix)
	}

ScanPrefix:
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
				break ScanPrefix
			}
		}
	}

	return res.String()
}

func (id JetID) Marshal() ([]byte, error) {
	return id[:], nil
}

func (id JetID) Equal(other JetID) bool {
	return id == other
}
