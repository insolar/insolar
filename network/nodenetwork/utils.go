/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package nodenetwork

import (
	consensus "github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

type MergedListCopy struct {
	ActiveList                 map[core.RecordRef]core.Node
	NodesJoinedDuringPrevPulse bool
}

func GetMergedCopy(nodes []core.Node, claims []consensus.ReferendumClaim) (*MergedListCopy, error) {
	nodesMap := copyActiveNodes(nodes)

	var nodesJoinedDuringPrevPulse bool
	for _, claim := range claims {
		isJoin, err := mergeClaim(nodesMap, claim)
		if err != nil {
			return nil, errors.Wrap(err, "[ GetMergedCopy ] failed to merge a claim")
		}

		nodesJoinedDuringPrevPulse = nodesJoinedDuringPrevPulse || isJoin
	}

	return &MergedListCopy{
		ActiveList:                 nodesMap,
		NodesJoinedDuringPrevPulse: nodesJoinedDuringPrevPulse,
	}, nil
}

func mergeClaim(nodes map[core.RecordRef]core.Node, claim consensus.ReferendumClaim) (bool, error) {
	isJoinClaim := false
	switch t := claim.(type) {
	case *consensus.NodeJoinClaim:
		isJoinClaim = true
		// TODO: fix version
		node, err := ClaimToNode("", t)
		if err != nil {
			return isJoinClaim, errors.Wrap(err, "[ mergeClaim ] failed to convert Claim -> Node")
		}
		node.(MutableNode).SetState(core.NodeJoining)
		nodes[node.ID()] = node
	case *consensus.NodeLeaveClaim:
		if nodes[t.NodeID] == nil {
			break
		}

		node := nodes[t.NodeID].(MutableNode)
		if t.ETA == 0 || !node.Leaving() {
			node.SetLeavingETA(t.ETA)
		}
	}

	return isJoinClaim, nil
}
