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

package node

import (
	"bytes"

	"github.com/insolar/insolar/insolar"
	"github.com/ugorji/go/codec"
)

type ListType int

const (
	ListWorking ListType = iota
	ListIdle
	ListLeaving
	ListSuspected
	ListJoiner

	ListLength
)

type Snapshot struct {
	pulse insolar.PulseNumber
	state insolar.NetworkState

	nodeList [ListLength][]insolar.NetworkNode
}

func (s *Snapshot) GetPulse() insolar.PulseNumber {
	return s.pulse
}

func (s *Snapshot) Copy() *Snapshot {
	result := &Snapshot{
		pulse: s.pulse,
		state: s.state,
	}
	for i := 0; i < int(ListLength); i++ {
		result.nodeList[i] = make([]insolar.NetworkNode, len(s.nodeList[i]))
		copy(result.nodeList[i], s.nodeList[i])
	}
	return result
}

// NewSnapshot create new snapshot for pulse.
func NewSnapshot(number insolar.PulseNumber, nodes map[insolar.Reference]insolar.NetworkNode) *Snapshot {
	return &Snapshot{
		pulse: number,
		// TODO: pass actual state
		state:    insolar.NoNetworkState,
		nodeList: splitNodes(nodes),
	}
}

// splitNodes temporary method to create snapshot lists. Will be replaced by special function that will take in count
// previous snapshot and approved claims.
func splitNodes(nodes map[insolar.Reference]insolar.NetworkNode) [ListLength][]insolar.NetworkNode {
	var result [ListLength][]insolar.NetworkNode
	for i := 0; i < int(ListLength); i++ {
		result[i] = make([]insolar.NetworkNode, 0)
	}
	for _, node := range nodes {
		listType := nodeStateToListType(node.GetState())
		if listType == ListLength {
			continue
		}
		result[listType] = append(result[listType], node)
	}
	return result
}

func nodeStateToListType(state insolar.NodeState) ListType {
	switch state {
	case insolar.NodeReady:
		return ListWorking
	case insolar.NodePending:
		return ListJoiner
	case insolar.NodeUndefined, insolar.NodeLeaving:
		return ListLeaving
	}
	// special case for no match
	return ListLength
}

func EncodeSnapshot(s *Snapshot) []byte {
	var buff bytes.Buffer
	enc := codec.NewEncoder(&buff, &codec.JsonHandle{})
	enc.MustEncode(s)
	return buff.Bytes()
}

func DecodeSnapshot(buf []byte) (s Snapshot) {
	dec := codec.NewDecoderBytes(buf, &codec.JsonHandle{})
	dec.MustDecode(&s)
	return s
}
