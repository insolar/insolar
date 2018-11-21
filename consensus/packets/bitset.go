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

package packets

import (
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

// TriState is state of the communicating node
type TriState uint8

const (
	// Legit is bit indicating OK data from node
	Legit TriState = iota + 1
	// TimedOut is bit indicating that timeout occurred when communicating with node
	TimedOut
	// Fraud is bit indicating that the node is malicious (fraud)
	Fraud
)

// BitSetCell is structure that contains the state of the node
type BitSetCell struct {
	NodeID core.RecordRef
	State  TriState
}

// Possible errors in BitSetMapper
var (
	// ErrBitSetOutOfRange is returned when index passed to IndexToRef function is out of range (ERROR)
	ErrBitSetOutOfRange = errors.New("index out of range")
	// ErrBitSetNodeIsMissing is returned in IndexToRef when we have no information about the node on specified index (SPECIAL CASE)
	ErrBitSetNodeIsMissing = errors.New("no information about node on specified index")
	// ErrBitSetIncorrectNode is returned when an incorrect node is passed to RefToIndex (ERROR)
	ErrBitSetIncorrectNode = errors.New("incorrect node ID")
)

// BitSetMapper contains the mapping from bitset index to node ID (and vice versa)
type BitSetMapper interface {
	// IndexToRef get ID of the node that is stored on the specified internal index
	IndexToRef(index int) (core.RecordRef, error)
	// RefToIndex get bitset internal index where the specified node state is stored
	RefToIndex(nodeID core.RecordRef) (int, error)
	// Length returns required length of the bitset
	Length() int
}

// BitSet is interface
type BitSet interface {
	Serializer
	// GetCells get buckets of bitset
	GetCells() []BitSetCell
	// ApplyChanges returns copy of the current bitset with changes applied
	ApplyChanges(changes []BitSetCell)
}

// NewBitSet creates bitset from a set of buckets and the mapper
func NewBitSet(cells []BitSetCell, mapper BitSetMapper) (BitSet, error) {
	return NewTriStateBitSet(cells, mapper)
}
