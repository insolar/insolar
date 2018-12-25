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
	"fmt"
	"strings"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

type jet struct {
	Left   *jet
	Right  *jet
	Actual bool
}

// Find returns jet for provided reference.
func (j *jet) Find(val []byte, depth uint8) (*jet, uint8) {
	if j == nil || val == nil {
		return nil, 0
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
	return j, depth
}

// Update add missing tree branches for provided prefix.
func (j *jet) Update(prefix []byte, setActual bool, maxDepth, depth uint8) {
	if setActual {
		j.Actual = true
	}

	if depth >= maxDepth {
		return
	}

	if getBit(prefix, depth) {
		if j.Right == nil {
			j.Right = &jet{}
		}
		j.Right.Update(prefix, setActual, maxDepth, depth+1)
	} else {
		if j.Left == nil {
			j.Left = &jet{}
		}
		j.Left.Update(prefix, setActual, maxDepth, depth+1)
	}
}

// ResetActual resets actual mark, which will signify uncertain state on nodes and require actualization.
func (j *jet) ResetActual() {
	if j.Left != nil {
		j.Left.ResetActual()
	}
	j.Actual = false
	if j.Right != nil {
		j.Right.ResetActual()
	}
}

// Tree stores jet in a binary tree.
type Tree struct {
	Head *jet
}

// String visualizes Jet's tree.
func (t Tree) String() string {
	if t.Head == nil {
		return "<nil>"
	}
	return nodeDeepFmt(0, "", t.Head)
}

func nodeDeepFmt(deep int, binPrefix string, node *jet) string {
	prefix := strings.Repeat(" ", deep)
	if deep == 0 {
		prefix = "root"
	}
	s := fmt.Sprintf("%s%v (level=%v actual=%v)\n", prefix, binPrefix, deep, node.Actual)

	if node.Left != nil {
		s += nodeDeepFmt(deep+1, binPrefix+"0", node.Left)
	}
	if node.Right != nil {
		s += nodeDeepFmt(deep+1, binPrefix+"1", node.Right)
	}
	return s
}

// NewTree creates new tree.
func NewTree() *Tree {
	return &Tree{Head: &jet{Actual: true}}
}

// Find returns jet for provided reference. If found jet is actual, the second argument will be true.
func (t *Tree) Find(id core.RecordID) (*core.RecordID, bool) {
	if id.Pulse() == core.PulseNumberJet {
		return &id, true
	}
	j, depth := t.Head.Find(id.Hash(), 0)
	return NewID(uint8(depth), resetBits(id.Hash(), depth)), j.Actual
}

// Update add missing tree branches for provided prefix. If 'setActual' is set, all encountered nodes will be marked as
// actual.
func (t *Tree) Update(id core.RecordID, setActual bool) {
	maxDepth, prefix := Jet(id)
	t.Head.Update(prefix, setActual, maxDepth, 0)
}

// Bytes serializes pulse.
func (t *Tree) Bytes() []byte {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	enc.MustEncode(t)
	return buf.Bytes()
}

// Split looks for provided jet and creates (and returns) two branches for it. If provided jet is not found, an error
// will be returned.
func (t *Tree) Split(jetID core.RecordID) (*core.RecordID, *core.RecordID, error) {
	depth, prefix := Jet(jetID)
	j, foundDepth := t.Head.Find(prefix, 0)
	if depth != foundDepth {
		return nil, nil, errors.New("failed to split: jet is not present in the tree")
	}
	j.Right = &jet{}
	j.Left = &jet{}
	leftPrefix := resetBits(prefix, depth)
	rightPrefix := resetBits(prefix, depth)
	setBit(rightPrefix, depth)
	return NewID(depth+1, leftPrefix), NewID(depth+1, rightPrefix), nil
}

// ResetActual resets actual mark, which will signify uncertain state on nodes and require actualization.
func (t *Tree) ResetActual() {
	t.Head.ResetActual()
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

func setBit(value []byte, index uint8) {
	if uint(index) >= uint(len(value)*8) {
		panic("index overflow")
	}
	byteIndex := uint(index / 8)
	bitIndex := uint(7 - index%8)
	mask := byte(1 << bitIndex)
	value[byteIndex] |= mask
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
