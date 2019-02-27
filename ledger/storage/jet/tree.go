/*
 *    Copyright 2019 Insolar Technologies
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

// Clone clones tree either keeping actuality state or resetting it to false
func (j *jet) Clone(keep bool) *jet {
	res := &jet{Actual: keep && j.Actual}
	if j.Left != nil {
		res.Left = j.Left.Clone(keep)
	}
	if j.Right != nil {
		res.Right = j.Right.Clone(keep)
	}
	return res
}

func (j *jet) ExtractLeafIDs(ids *[]core.RecordID, path []byte, depth uint8) {
	if j == nil {
		return
	}
	if j.Left == nil && j.Right == nil {
		*ids = append(*ids, *NewID(depth, path))
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

func (t *Tree) toArray() []*jet {
	queue := []*jet{t.Head}
	queueIndex := 0
	lastElementIndex := 0

	for queueIndex <= lastElementIndex {
		currentNode := queue[queueIndex]
		if currentNode == nil {
			queueIndex++
			continue
		}
		nextLeftIndex := (queueIndex * 2) + 1
		nextRightIndex := (queueIndex * 2) + 2

		if nextRightIndex >= len(queue) {
			for nextRightIndex >= len(queue) {
				queue = append(queue, nil)
			}
		}

		queue[nextLeftIndex] = queue[queueIndex].Left
		queue[nextRightIndex] = queue[queueIndex].Right

		if queue[nextLeftIndex] != nil && nextLeftIndex > lastElementIndex {
			lastElementIndex = nextLeftIndex
		}

		if queue[nextRightIndex] != nil && nextRightIndex > lastElementIndex {
			lastElementIndex = nextRightIndex
		}
		queueIndex++
	}

	return queue[:lastElementIndex+1]
}

// Merge merges two trees to one
func (t *Tree) Merge(newTree *Tree) *Tree {
	maxLengthFunc := func(left int, right int) int {
		if left > right {
			return left
		}

		return right
	}

	savedTree := t.toArray()
	inputTree := newTree.toArray()

	maxLength := maxLengthFunc(len(savedTree), len(inputTree))

	result := make([]*jet, maxLength)

	for index := 0; index < maxLength; index++ {
		var savedItem *jet
		var newItem *jet

		if index < len(savedTree) {
			savedItem = savedTree[index]
		}
		if index < len(inputTree) {
			newItem = inputTree[index]
		}

		if savedItem == nil {
			result[index] = newItem
			continue
		}
		if newItem == nil {
			result[index] = savedItem
			continue
		}

		if !savedItem.Actual && newItem.Actual {
			result[index] = newItem
		} else {
			result[index] = savedItem
		}
	}

	for index := 0; index < (maxLength/2)+1; index++ {
		current := result[index]
		if current != nil {
			leftIndex := index*2 + 1
			if leftIndex < len(result) {
				current.Left = result[index*2+1]
			}

			rightIndex := index*2 + 2
			if rightIndex < len(result) {
				current.Right = result[index*2+2]
			}
		}
	}

	return &Tree{Head: result[0]}
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
func NewTree(isActual bool) *Tree {
	return &Tree{Head: &jet{Actual: isActual}}
}

// Clone clones the tree keeping actuality or setting everything to false
func (t *Tree) Clone(keep bool) *Tree {
	return &Tree{Head: t.Head.Clone(keep)}
}

// Find returns jet for provided reference. If found jet is actual, the second argument will be true.
func (t *Tree) Find(id core.RecordID) (*core.RecordID, bool) {
	if id.Pulse() == core.PulseNumberJet {
		return &id, true
	}
	hash := id.Hash()
	j, depth := t.Head.Find(hash, 0)
	return NewID(uint8(depth), ResetBits(hash, depth)), j.Actual
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
		return nil, nil, errors.New("failed to split: incorrect jet provided")
	}
	j.Right = &jet{}
	j.Left = &jet{}
	leftPrefix := ResetBits(prefix, depth)
	rightPrefix := ResetBits(prefix, depth)
	setBit(rightPrefix, depth)
	return NewID(depth+1, leftPrefix), NewID(depth+1, rightPrefix), nil
}

func (t *Tree) LeafIDs() []core.RecordID {
	var ids []core.RecordID
	t.Head.ExtractLeafIDs(&ids, make([]byte, core.RecordHashSize), 0)
	return ids
}

func getBit(value []byte, index uint8) bool {
	if uint(index) >= uint(len(value)*8) {
		panic(fmt.Sprintf("index overflow: value=%08b, index=%v", value, index))
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
