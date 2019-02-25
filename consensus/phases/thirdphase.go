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

package phases

import (
	"context"

	"github.com/insolar/insolar/consensus"
	"github.com/insolar/insolar/consensus/claimhandler"
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/merkle"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

type ThirdPhase interface {
	Execute(ctx context.Context, pulse *core.Pulse, state *SecondPhaseState) (*ThirdPhaseState, error)
}

func NewThirdPhase() ThirdPhase {
	return &ThirdPhaseImpl{}
}

type ThirdPhaseImpl struct {
	Cryptography core.CryptographyService `inject:""`
	Communicator Communicator             `inject:""`
	NodeKeeper   network.NodeKeeper       `inject:""`
	Calculator   merkle.Calculator        `inject:""`
}

func (tp *ThirdPhaseImpl) Execute(ctx context.Context, pulse *core.Pulse, state *SecondPhaseState) (*ThirdPhaseState, error) {
	ctx, span := instracer.StartSpan(ctx, "ThirdPhase.Execute")
	span.AddAttributes(trace.Int64Attribute("pulse", int64(state.PulseEntry.Pulse.PulseNumber)))
	defer span.End()
	stats.Record(ctx, consensus.Phase3Exec.M(1))

	logger := inslogger.FromContext(ctx)
	totalCount := state.UnsyncList.Length()

	var gSign [packets.SignatureLength]byte
	copy(gSign[:], state.GlobuleProof.Signature.Bytes()[:packets.SignatureLength])
	packet := packets.NewPhase3Packet(pulse.PulseNumber, gSign, state.BitSet)

	nodes := make([]core.Node, 0)
	for _, node := range state.MatrixState.Active {
		nodes = append(nodes, state.UnsyncList.GetActiveNode(node))
	}
	responses, err := tp.Communicator.ExchangePhase3(ctx, nodes, packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ NET Consensus phase-3 ] Failed to exchange packets")
	}
	logger.Infof("[ NET Consensus phase-3 ] received responses: %d/%d", len(responses), totalCount)
	err = stats.RecordWithTags(ctx, []tag.Mutator{tag.Upsert(consensus.TagPhase, "phase 3")}, consensus.PacketsRecv.M(int64(len(responses))))
	if err != nil {
		logger.Warn("[ NET Consensus phase-3 ] Failed to record received responses metric: " + err.Error())
	}

	handler := claimhandler.NewJoinHandler(totalCount)
	for ref, packet := range responses {
		err = nil
		if !ref.Equal(tp.NodeKeeper.GetOrigin().ID()) {
			err = tp.checkPacketSignature(packet, ref, state.UnsyncList)
		}
		if err != nil {
			logger.Warnf("[ NET Consensus phase-3 ] Failed to check phase3 packet signature from %s: %s", ref, err.Error())
			continue
		}
		// not needed until we implement fraud detection
		// cells, err := packet.GetBitset().GetCells(state.UnsyncList)

		state.UnsyncList.SetGlobuleHashSignature(ref, packet.GetGlobuleHashSignature())
	}

	for _, node := range nodes {
		handler.AddClaims(state.UnsyncList.GetClaims(node.ID()), pulse.Entropy)
	}

	handledJoinClaims := handler.HandleAndReturnClaims()

	for ref := range responses {
		tp.removeExcessJoinClaims(handledJoinClaims, ref, state)
	}

	prevCloudHash := tp.NodeKeeper.GetCloudHash()
	validNodes := 0
	for _, node := range nodes {
		ghs, ok := state.UnsyncList.GetGlobuleHashSignature(node.ID())
		if !ok {
			log.Warnf("[ NET Consensus phase-3 ] No globule hash signature for node %s", node.ID())
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
			validNodes++
		} else {
			logger.Warnf("[ NET Consensus phase-3 ] Failed to validate globule hash from node %s", node.ID())
		}
	}

	if !consensusReachedBFT(validNodes, totalCount) {
		return nil, errors.Errorf("[ NET Consensus phase-3 ] Failed to pass BFT consensus: %d/%d", validNodes, totalCount)
	}

	logger.Infof("[ NET Consensus phase-3 ] BFT consensus passed: %d/%d", validNodes, totalCount)

	return &ThirdPhaseState{
		ActiveNodes:  state.MatrixState.Active,
		UnsyncList:   state.UnsyncList,
		GlobuleProof: state.GlobuleProof,
	}, nil
}

func (tp *ThirdPhaseImpl) removeExcessJoinClaims(joinClaims []*packets.NodeJoinClaim, ref core.RecordRef, state *SecondPhaseState) {
	claims := state.UnsyncList.GetClaims(ref)
	originLen := len(claims)
	if originLen == 0 || len(joinClaims) == 0 {
		return
	}

	updatedClaims := make([]packets.ReferendumClaim, 0)
	for _, join := range joinClaims {
		for i := 0; i < len(claims); i++ {
			claim, ok := claims[i].(*packets.NodeJoinClaim)
			if ok && claim.NodeRef.Equal(join.NodeRef) {
				updatedClaims = append(updatedClaims, join)
				continue
			} else if !ok {
				updatedClaims = append(updatedClaims, claim)
			}
		}
	}

	state.UnsyncList.InsertClaims(ref, updatedClaims)
}

func (tp *ThirdPhaseImpl) checkPacketSignature(packet *packets.Phase3Packet, recordRef core.RecordRef, unsyncList network.UnsyncList) error {
	activeNode := unsyncList.GetActiveNode(recordRef)
	if activeNode == nil {
		return errors.New("failed to get active node")
	}
	key := activeNode.PublicKey()
	return packet.Verify(tp.Cryptography, key)
}
