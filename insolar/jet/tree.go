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
	"errors"
	"fmt"
	"strings"

	"github.com/insolar/insolar/insolar"
)

type jet struct {
	// Active indicates what node is approved as last one, i.e. it's leaf (ignored for non leafs if set)
	Actual bool
	Left   *jet
	Right  *jet
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
	if depth == maxDepth {
		if setActual {
			j.Actual = true
		}
		return
	}

	if j.Right == nil {
		j.Right = &jet{}
	}
	if j.Left == nil {
		j.Left = &jet{}
	}
	if getBit(prefix, depth) {
		j.Right.Update(prefix, setActual, maxDepth, depth+1)
	} else {
		j.Left.Update(prefix, setActual, maxDepth, depth+1)
	}
}

// Clone clones tree either keeping actuality state or resetting it to false.
func (j *jet) Clone(keep bool) *jet {
	res := &jet{
		Actual: keep && j.Actual,
	}
	if j.Left != nil {
		res.Left = j.Left.Clone(keep)
	}
	if j.Right != nil {
		res.Right = j.Right.Clone(keep)
	}
	return res
}

func (j *jet) ExtractLeafIDs(ids *[]insolar.JetID, path []byte, depth uint8) {
	if j == nil {
		return
	}
	if j.Left == nil && j.Right == nil && j.Actual {
		*ids = append(*ids, *insolar.NewJetID(depth, path))
		return
	}

	if j.Left != nil {
		j.Left.ExtractLeafIDs(ids, path, depth+1)
	}
	if j.Right != nil {
		rightPath := make([]byte, len(path))
		copy(rightPath, path)
		setBit(rightPath, depth)
		j.Right.ExtractLeafIDs(ids, rightPath, depth+1)
	}
}

// Tree represents jets tree.
type Tree struct {
	Head *jet
}

// NewTree creates new tree.
func NewTree(isActual bool) *Tree {
	return &Tree{
		Head: &jet{
			Actual: isActual,
		},
	}
}

// Clone clones the tree keeping actuality or setting everything to false
func (t *Tree) Clone(keep bool) *Tree {
	return &Tree{Head: t.Head.Clone(keep)}
}

// Find returns jet for provided record ID.
// If found jet is actual, the second argument will be true.
func (t *Tree) Find(recordID insolar.ID) (insolar.JetID, bool) {
	// if provided record ID is JetID, returns it as actual jet. (kind of hack)
	// TODO: we should remove this and check tests
	if recordID.Pulse() == insolar.PulseNumberJet {
		return insolar.JetID(recordID), true
	}

	hash := recordID.Hash()
	j, depth := t.Head.Find(hash, 0)
	id := *insolar.NewJetID(depth, resetBits(hash, depth))
	return id, j.Actual
}

// Update add missing tree branches for provided prefix.
// If 'setActual' is set, all encountered nodes will be marked as actual.
func (t *Tree) Update(id insolar.JetID, setActual bool) {
	t.Head.Update(id.Prefix(), setActual, id.Depth(), 0)
}

// Split looks for provided jet and creates (and returns) two branches for it.
// If provided jet is not found, an error will be returned.
func (t *Tree) Split(id insolar.JetID) (insolar.JetID, insolar.JetID, error) {
	depth, prefix := id.Depth(), id.Prefix()
	j, foundDepth := t.Head.Find(prefix, 0)
	if depth != foundDepth {
		return insolar.ZeroJetID, insolar.ZeroJetID, errors.New("failed to split: incorrect jet provided")
	}

	left, right := Siblings(id)
	j.Left = &jet{Actual: true}
	j.Right = &jet{Actual: true}
	return left, right, nil
}

func (t *Tree) LeafIDs() []insolar.JetID {
	var ids []insolar.JetID
	t.Head.ExtractLeafIDs(&ids, make([]byte, insolar.RecordHashSize), 0)
	return ids
}

// getBit returns true if bit at index is set to 1 in byte array.
// Panics if index is out of range (value size * 8).
func getBit(value []byte, index uint8) bool {
	if uint(index) >= uint(len(value)*8) {
		panic(fmt.Sprintf("index overflow: value=%08b, index=%v", value, index))
	}
	byteIndex := uint(index / 8)
	bitIndex := uint(7 - index%8)
	mask := byte(1 << bitIndex)
	return value[byteIndex]&mask != 0
}

// setBit sets bit to 1 in byte array at index.
// Panics if index is out of range (value size * 8).
func setBit(value []byte, index uint8) {
	if uint(index) >= uint(len(value)*8) {
		panic("index overflow")
	}
	byteIndex := uint(index / 8)
	bitIndex := uint(7 - index%8)
	mask := byte(1 << bitIndex)
	value[byteIndex] |= mask
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
