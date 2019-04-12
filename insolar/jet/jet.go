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

package jet

import (
	"context"
	"strconv"
	"strings"

	"github.com/insolar/insolar/insolar"
)

//go:generate minimock -i github.com/insolar/insolar/insolar/jet.Accessor -o ./ -s _mock.go

// Accessor provides an interface for accessing jet IDs.
type Accessor interface {
	All(ctx context.Context, pulse insolar.PulseNumber) []insolar.JetID
	ForID(ctx context.Context, pulse insolar.PulseNumber, recordID insolar.ID) (insolar.JetID, bool)
}

//go:generate minimock -i github.com/insolar/insolar/insolar/jet.Modifier -o ./ -s _mock.go

// Modifier provides an interface for modifying jet IDs.
type Modifier interface {
	Update(ctx context.Context, pulse insolar.PulseNumber, actual bool, ids ...insolar.JetID)
	Split(ctx context.Context, pulse insolar.PulseNumber, id insolar.JetID) (insolar.JetID, insolar.JetID, error)
	Clone(ctx context.Context, from, to insolar.PulseNumber)
	Delete(ctx context.Context, pulse insolar.PulseNumber)
}

type Calculator interface {
	MineForPulse(ctx context.Context, pn insolar.PulseNumber) []insolar.JetID
}

//go:generate minimock -i github.com/insolar/insolar/insolar/jet.Storage -o ./ -s _mock.go

// Storage composes Accessor and Modifier interfaces.
type Storage interface {
	Accessor
	Modifier
}

// Parent returns a parent of the jet or jet itself if depth of provided JetID is zero.
func Parent(id insolar.JetID) insolar.JetID {
	depth, prefix := id.Depth(), id.Prefix()
	if depth == 0 {
		return id
	}

	return *insolar.NewJetID(depth-1, resetBits(prefix, depth-1))
}

// resetBits returns a new byte slice with all bits in 'value' reset,
// starting from 'start' number of bit.
//
// If 'start' is bigger than len(value), the original slice will be returned.
func resetBits(value []byte, start uint8) []byte {
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

// NewIDFromString creates new JetID from string represents binary prefix.
//
// "0"     -> prefix=[0..0], depth=1
// "1"     -> prefix=[1..0], depth=1
// "1010"  -> prefix=[1010..0], depth=4
func NewIDFromString(s string) insolar.JetID {
	id := insolar.NewJetID(uint8(len(s)), parsePrefix(s))
	return *id
}

func parsePrefix(s string) []byte {
	var prefix []byte
	tail := s[:]
	for len(tail) > 0 {
		offset := 8
		if len(tail) < 8 {
			tail += strings.Repeat("0", 8-len(tail))
		}
		parsed, err := strconv.ParseUint(tail[:offset], 2, 8)
		if err != nil {
			panic(err)
		}
		prefix = append(prefix, byte(parsed))
		tail = tail[offset:]
	}
	return prefix
}
