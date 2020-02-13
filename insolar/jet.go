// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package insolar

import (
	"fmt"
	"strings"

	"github.com/insolar/insolar/insolar/bits"
	"github.com/insolar/insolar/pulse"
	"github.com/insolar/insolar/reference"

	"github.com/pkg/errors"
)

const (
	// JetSize is a Jet's size (depth+prefix).
	JetSize = RecordIDSize - PulseNumberSize
	// JetPrefixSize is a Jet's prefix size.
	JetPrefixSize = JetSize - 1
	// JetMaximumDepth is a Jet's maximum depth (maximum offset in bits).
	JetMaximumDepth = JetPrefixSize*8 - 1
	// JetDepthPosition is an position where depth of jet id is located
	JetDepthPosition = 0
	// JetPrefixOffset is an offset where prefix starts in jet id.
	JetPrefixOffset = JetDepthPosition + 1
)

// JetID should be used, when id is a jetID
type JetID ID

// Size is a protobuf required method. It returns size of JetID
func (id *JetID) Size() int { return reference.LocalBinarySize }

// MarshalTo is a protobuf required method. It marshals data
func (id *JetID) MarshalTo(data []byte) (n int, err error) {
	return (*ID)(id).MarshalTo(data)
}

// Unmarshal is a protobuf required method. It unmarshals data
func (id *JetID) Unmarshal(data []byte) error {
	if err := (*ID)(id).Unmarshal(data); err != nil {
		return errors.New("Not enough bytes to unpack JetID")
	}
	return nil
}

// IsValid returns true is JetID has a predefined reserved pulse number.
func (id *JetID) IsValid() bool {
	return (*ID)(id).Pulse().IsJet()
}

// IsEmpty - check for void
func (id JetID) IsEmpty() bool {
	return id.Equal(JetID{})
}

// ZeroJetID is value of an empty Jet ID
var ZeroJetID = *NewJetID(0, nil)

// NewJetID creates a new jet with provided ID and index
func NewJetID(depth uint8, prefix []byte) *JetID {
	hash := [reference.LocalBinaryHashSize]byte{depth}
	copy(hash[JetPrefixOffset:], bits.ResetBits(prefix, depth))

	return (*JetID)(NewID(pulse.Jet, hash[:]))
}

// Depth extracts depth from a jet id.
func (id JetID) Depth() uint8 {
	if !id.IsValid() {
		panic(fmt.Sprintf("provided id %b is not a jet id", id))
	}
	return ID(id).GetHash()[JetDepthPosition]
}

// Prefix extracts prefix from a jet id.
func (id JetID) Prefix() []byte {
	if !id.IsValid() {
		panic(fmt.Sprintf("provided id %b is not a jet id", id))
	}
	return ID(id).Hash()[JetPrefixOffset:]
}

// DebugString prints JetID in human readable form.
func (id *JetID) DebugString() string {
	return (*ID)(id).DebugString()
}

type JetIDCollection []JetID

func (ids JetIDCollection) DebugString() string {
	builder := strings.Builder{}
	builder.WriteRune('[')
	for i, id := range ids {
		builder.WriteString(id.DebugString())
		if i < len(ids)-1 {
			builder.WriteRune(',')
		}
	}
	builder.WriteRune(']')
	return builder.String()
}

func (id JetID) Marshal() ([]byte, error) {
	return ID(id).Marshal()
}

func (id *JetID) Equal(other JetID) bool {
	return (*ID)(id).Equal(ID(other))
}
