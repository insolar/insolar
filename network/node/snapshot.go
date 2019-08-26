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
	"reflect"

	"github.com/insolar/insolar/insolar"
	protonode "github.com/insolar/insolar/network/node/internal/node"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
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

func (s *Snapshot) Equal(s2 *Snapshot) bool {
	if s.pulse != s2.pulse || s.state != s2.state {
		return false
	}

	for t, list := range s.nodeList {
		if len(list) != len(s2.nodeList[t]) {
			return false
		}
		for i, n := range list {
			n2 := s2.nodeList[t][i]
			if !reflect.DeepEqual(n, n2) {
				return false
			}
		}
	}
	return true
}

// NewSnapshot create new snapshot for pulse.
func NewSnapshot(number insolar.PulseNumber, nodes []insolar.NetworkNode) *Snapshot {
	return &Snapshot{
		pulse: number,
		// TODO: pass actual state
		state:    insolar.NoNetworkState,
		nodeList: splitNodes(nodes),
	}
}

// splitNodes temporary method to create snapshot lists. Will be replaced by special function that will take in count
// previous snapshot and approved claims.
func splitNodes(nodes []insolar.NetworkNode) [ListLength][]insolar.NetworkNode {
	var result [ListLength][]insolar.NetworkNode
	for i := 0; i < int(ListLength); i++ {
		result[i] = make([]insolar.NetworkNode, 0)
	}
	for _, node := range nodes {
		listType := nodeStateToListType(node)
		if listType == ListLength {
			continue
		}
		result[listType] = append(result[listType], node)
	}
	return result
}

func nodeStateToListType(node insolar.NetworkNode) ListType {
	switch node.GetState() {
	case insolar.NodeReady:
		if node.GetPower() > 0 {
			return ListWorking
		}
		return ListIdle
	case insolar.NodeJoining:
		return ListJoiner
	case insolar.NodeUndefined, insolar.NodeLeaving:
		return ListLeaving
	}
	// special case for no match
	return ListLength
}

func (s *Snapshot) Encode() ([]byte, error) {
	ss := protonode.Snapshot{}
	ss.PulseNumber = uint32(s.pulse)
	ss.State = uint32(s.state)
	keyProc := platformpolicy.NewKeyProcessor()

	ss.Nodes = make(map[uint32]*protonode.NodeList)
	for t, list := range s.nodeList {
		protoNodeList := make([]*protonode.Node, len(list))
		for i, n := range list {

			exportedKey, err := keyProc.ExportPublicKeyBinary(n.PublicKey())
			if err != nil {
				return nil, errors.Wrap(err, "Failed to export a public key")
			}

			protoNode := &protonode.Node{
				NodeID:         n.ID().Bytes(),
				NodeShortID:    uint32(n.ShortID()),
				NodeRole:       uint32(n.Role()),
				NodePublicKey:  exportedKey,
				NodeAddress:    n.Address(),
				NodeVersion:    n.Version(),
				NodeLeavingETA: uint32(n.LeavingETA()),
				State:          uint32(n.GetState()),
			}

			protoNodeList[i] = protoNode
		}

		l := &protonode.NodeList{}
		l.List = protoNodeList
		ss.Nodes[uint32(t)] = l
	}

	return ss.Marshal()
}

func (s *Snapshot) Decode(buff []byte) error {
	ss := protonode.Snapshot{}
	err := ss.Unmarshal(buff)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal node")
	}

	keyProc := platformpolicy.NewKeyProcessor()
	s.pulse = insolar.PulseNumber(ss.PulseNumber)
	s.state = insolar.NetworkState(ss.State)

	for t, nodes := range ss.Nodes {
		nodeList := make([]insolar.NetworkNode, len(nodes.List))
		for i, n := range nodes.List {

			pk, err := keyProc.ImportPublicKeyBinary(n.NodePublicKey)
			if err != nil {
				return errors.Wrap(err, "Failed to ImportPublicKeyBinary")
			}

			ref := insolar.NewReferenceFromBytes(n.NodeID)
			nodeList[i] = newMutableNode(*ref, insolar.StaticRole(n.NodeRole), pk, insolar.NodeState(n.State), n.NodeAddress, n.NodeVersion)
		}
		s.nodeList[t] = nodeList
	}

	return nil
}

func Select(nodes []insolar.NetworkNode, typ ListType) []insolar.NetworkNode {
	lists := splitNodes(nodes)
	return lists[typ]
}
