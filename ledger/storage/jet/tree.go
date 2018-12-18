/*
 *    Copyright 2018 Insolar
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

package jet

import (
	"bytes"

	"github.com/insolar/insolar/core"
	"github.com/ugorji/go/codec"
)

type jet struct {
	Left  *jet
	Right *jet
}

// Find returns jet for provided reference.
func (j *jet) Find(val []byte, depth uint8) uint8 {
	if j == nil || val == nil {
		return 0
	}

	if getBit(val, depth) {
		if j.Right != nil {
			return j.Right.Find(val, depth+1)
		}
	} else {
		if j.Left != nil {
			return j.Left.Find(val, depth+1)
		}
	}
	return depth
}

// Update add missing tree branches for provided prefix.
func (j *jet) Update(prefix []byte, maxDepth, depth uint8) {
	if depth >= maxDepth {
		return
	}

	if getBit(prefix, depth) {
		if j.Right == nil {
			j.Right = &jet{}
		}
		j.Right.Update(prefix, maxDepth, depth+1)
	} else {
		if j.Left == nil {
			j.Left = &jet{}
		}
		j.Left.Update(prefix, maxDepth, depth+1)
	}
}

// Tree stores jet in a binary tree.
type Tree struct {
	Head *jet
}

func NewTree() *Tree {
	return &Tree{Head: &jet{}}
}

// Find returns jet for provided reference.
func (t *Tree) Find(id core.RecordID) *core.RecordID {
	if id.Pulse() == core.PulseNumberJet {
		return &id
	}
	depth := t.Head.Find(id.Hash(), 0)
	return NewID(uint8(depth), resetBits(id.Hash(), depth))
}

// Update add missing tree branches for provided prefix.
func (t *Tree) Update(id core.RecordID) {
	maxDepth, prefix := Jet(id)
	t.Head.Update(prefix, maxDepth, 0)
}

// Bytes serializes pulse.
func (t *Tree) Bytes() []byte {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	enc.MustEncode(t)
	return buf.Bytes()
}

func getBit(value []byte, index uint8) bool {
	if uint(index) >= uint(len(value)*8) {
		panic("index overflow")
	}
	byteIndex := uint(index / 8)
	bitIndex := uint(7 - index%8)
	mask := byte(1 << bitIndex)
	return value[byteIndex]&mask != 0
}

// ResetBits returns a new byte slice with all bits in 'value' reset, starting from 'start' number of bit. If 'start'
// is bigger than len(value), the original slice will be returned.
func resetBits(value []byte, start uint8) []byte {
	if int(start) > len(value)*8 {
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

// PredefinedTree is used to test multi-jet functionality.
// TODO: remove me after functional split.
var PredefinedTree = Tree{
	Head: &jet{
		Left: &jet{
			Left:  &jet{},
			Right: &jet{},
		},
		Right: &jet{
			Left:  &jet{},
			Right: &jet{},
		},
	},
}
