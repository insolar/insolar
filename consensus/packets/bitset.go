//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package packets

import (
	"github.com/insolar/insolar/insolar"
	"github.com/pkg/errors"
)

// BitSetState is state of the communicating node
type BitSetState uint8

const (
	// TimedOut is state indicating that timeout occurred when communicating with node
	TimedOut BitSetState = iota
	// Legit is state indicating OK data from node
	Legit
	// Fraud is state indicating that the node is malicious (fraud)
	Fraud
	// Inconsistent is state indicating that node validation is inconsistent on different nodes
	Inconsistent
)

// BitSetCell is structure that contains the state of the node
type BitSetCell struct {
	NodeID insolar.Reference
	State  BitSetState
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

//go:generate minimock -i github.com/insolar/insolar/consensus/packets.BitSetMapper -o . -s _mock.go

// BitSetMapper contains the mapping from bitset index to node ID (and vice versa)
type BitSetMapper interface {
	// IndexToRef get ID of the node that is stored on the specified internal index
	IndexToRef(index int) (insolar.Reference, error)
	// RefToIndex get bitset internal index where the specified node state is stored
	RefToIndex(nodeID insolar.Reference) (int, error)
	// Length returns required length of the bitset
	Length() int
}

// BitSet is interface
type BitSet interface {
	Serialize() ([]byte, error)
	// GetCells get buckets of bitset
	GetCells(mapper BitSetMapper) ([]BitSetCell, error)
	// GetTristateArray get underlying tristate
	GetTristateArray() ([]BitSetState, error)
	// ApplyChanges returns copy of the current bitset with changes applied
	ApplyChanges(changes []BitSetCell, mapper BitSetMapper) error
}

// NewBitSet creates bitset from a set of buckets and the mapper. Size == cells count.
func NewBitSet(size int) (BitSet, error) {
	return NewBitSetImpl(size, false)
}
