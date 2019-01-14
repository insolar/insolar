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

package phases

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/merkle"
	"github.com/pkg/errors"
)

type ThirdPhase struct {
	Cryptography core.CryptographyService `inject:""`
	NodeNetwork  core.NodeNetwork         `inject:""`
	Communicator Communicator             `inject:""`
	NodeKeeper   network.NodeKeeper       `inject:""`
	Calculator   merkle.Calculator        `inject:""`
}

func (tp *ThirdPhase) Execute(ctx context.Context, state *SecondPhaseState) (*ThirdPhaseState, error) {
	var gSign [packets.SignatureLength]byte
	copy(gSign[:], state.GlobuleProof.Signature.Bytes()[:packets.SignatureLength])
	packet := packets.NewPhase3Packet(gSign, state.BitSet)

	nodes := make([]core.Node, 0)
	for _, node := range state.MatrixState.Active {
		nodes = append(nodes, state.UnsyncList.GetActiveNode(node))
	}
	responses, err := tp.Communicator.ExchangePhase3(ctx, nodes, packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase 3 ] Failed exchange packets on phase 3")
	}
	inslogger.FromContext(ctx).Infof("[ Phase 3 ] received responses: %d/%d", len(responses), len(nodes))

	for ref, packet := range responses {
		err = nil
		if !ref.Equal(tp.NodeKeeper.GetOrigin().ID()) {
			err = tp.checkPacketSignature(packet, ref, state.UnsyncList)
		}
		if err != nil {
			inslogger.FromContext(ctx).Warnf("Failed to check phase3 packet signature from %s: %s", ref, err.Error())
			continue
		}
		// not needed until we implement fraud detection
		// cells, err := packet.GetBitset().GetCells(state.UnsyncList)

		state.UnsyncList.SetGlobuleHashSignature(ref, packet.GetGlobuleHashSignature())
	}

	totalCount := state.UnsyncList.Length()
	prevCloudHash := tp.NodeKeeper.GetCloudHash()
	validNodes := make([]core.RecordRef, 0)
	for _, node := range nodes {
		ghs, ok := state.UnsyncList.GetGlobuleHashSignature(node.ID())
		if !ok {
			log.Warnf("[ Phase 3 ] No globule hash signature for node %s", node.ID())
			continue
		}
		proof := &merkle.GlobuleProof{
			BaseProof: merkle.BaseProof{
				Signature: core.SignatureFromBytes(ghs[:]),
			},
			PrevCloudHash: prevCloudHash,
			GlobuleID:     state.GlobuleProof.GlobuleID,
			NodeCount:     state.GlobuleProof.NodeCount,
			NodeRoot:      state.GlobuleProof.NodeRoot,
		}
		valid := tp.Calculator.IsValid(proof, state.GlobuleHash, node.PublicKey())
		if valid {
			validNodes = append(validNodes, node.ID())
		} else {
			state.UnsyncList.RemoveNode(node.ID())
		}
	}

	if !consensusReachedBFT(len(validNodes), totalCount) {
		return nil, errors.New(fmt.Sprintf("[ Phase 3 ] Failed to pass BFT consensus: %d/%d", len(validNodes), totalCount))
	}

	inslogger.FromContext(ctx).Infof("Network phase 3 BFT consensus passed: %d/%d", len(validNodes), totalCount)

	return &ThirdPhaseState{
		ActiveNodes:  validNodes,
		UnsyncList:   state.UnsyncList,
		GlobuleProof: state.GlobuleProof,
	}, nil
}

func (tp *ThirdPhase) checkPacketSignature(packet *packets.Phase3Packet, recordRef core.RecordRef, unsyncList network.UnsyncList) error {
	activeNode := unsyncList.GetActiveNode(recordRef)
	if activeNode == nil {
		return errors.New("failed to get active node")
	}
	key := activeNode.PublicKey()
	return packet.Verify(tp.Cryptography, key)
}
