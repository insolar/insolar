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
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/merkle"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

type SecondPhase interface {
	Execute(ctx context.Context, pulse *core.Pulse, state *FirstPhaseState) (*SecondPhaseState, error)
	Execute21(ctx context.Context, pulse *core.Pulse, state *SecondPhaseState) (*SecondPhaseState, error)
}

func NewSecondPhase() SecondPhase {
	return &SecondPhaseImpl{}
}

type SecondPhaseImpl struct {
	NodeKeeper   network.NodeKeeper       `inject:""`
	Calculator   merkle.Calculator        `inject:""`
	Communicator Communicator             `inject:""`
	Cryptography core.CryptographyService `inject:""`
}

func (sp *SecondPhaseImpl) Execute(ctx context.Context, pulse *core.Pulse, state *FirstPhaseState) (*SecondPhaseState, error) {
	logger := inslogger.FromContext(ctx)
	ctx, span := instracer.StartSpan(ctx, "SecondPhase.Execute")
	span.AddAttributes(trace.Int64Attribute("pulse", int64(state.PulseEntry.Pulse.PulseNumber)))
	defer span.End()
	prevCloudHash := sp.NodeKeeper.GetCloudHash()

	state.ValidProofs[sp.NodeKeeper.GetOrigin()] = state.PulseProof

	entry := &merkle.GlobuleEntry{
		PulseEntry:    state.PulseEntry,
		ProofSet:      state.ValidProofs,
		PulseHash:     state.PulseHash,
		PrevCloudHash: prevCloudHash,
		GlobuleID:     sp.NodeKeeper.GetOrigin().GetGlobuleID(),
	}
	globuleHash, globuleProof, err := sp.Calculator.GetGlobuleProof(entry)

	if err != nil {
		return nil, errors.Wrap(err, "[ NET Consensus phase-2.0 ] Failed to calculate globule proof")
	}

	packet := packets.NewPhase2Packet(pulse.PulseNumber)
	err = packet.SetGlobuleHashSignature(globuleProof.Signature.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "[ NET Consensus phase-2.0 ] Failed to set globule proof in Phase2Packet")
	}
	bitset, err := sp.generatePhase2Bitset(state.UnsyncList, state.ValidProofs, pulse.PulseNumber)
	if err != nil {
		return nil, errors.Wrap(err, "[ NET Consensus phase-2.0 ] Failed to generate bitset for Phase2Packet")
	}
	packet.SetBitSet(bitset)
	activeNodes := state.UnsyncList.GetActiveNodes()
	packets, err := sp.Communicator.ExchangePhase2(ctx, state.UnsyncList, state.ClaimHandler, activeNodes, packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ NET Consensus phase-2.0 ] Failed to exchange packets")
	}
	logger.Infof("[ NET Consensus phase-2.0 ] Received responses: %d/%d", len(packets), state.UnsyncList.Length())
	err = stats.RecordWithTags(ctx, []tag.Mutator{tag.Upsert(consensus.TagPhase, "phase 2")}, consensus.PacketsRecv.M(int64(len(packets))))
	if err != nil {
		logger.Warn("[ NET Consensus phase-2.0 ] Failed to record received packets metric: " + err.Error())
	}

	origin := sp.NodeKeeper.GetOrigin().ID()
	stateMatrix := NewStateMatrix(state.UnsyncList)

	for ref, packet := range packets {
		err = nil
		if !ref.Equal(origin) {
			err = sp.checkPacketSignature(packet, ref, state.UnsyncList)
		}
		if err != nil {
			logger.Warnf("[ NET Consensus phase-2.0 ] Failed to check phase2 packet signature from %s: %s", ref, err.Error())
			continue
		}
		state.UnsyncList.SetGlobuleHashSignature(ref, packet.GetGlobuleHashSignature())
		err = stateMatrix.ApplyBitSet(ref, packet.GetBitSet())
		if err != nil {
			logger.Warnf("[ NET Consensus phase-2.0 ] Could not apply bitset from node %s: %s", ref, err.Error())
			continue
		}
	}

	matrixCalculation, err := stateMatrix.CalculatePhase2(origin)
	if err != nil {
		return nil, errors.Wrap(err, "[ NET Consensus phase-2.0 ] Failed to calculate bitset matrix consensus result")
	}

	if len(matrixCalculation.TimedOut) > 0 {
		for _, nodeID := range matrixCalculation.TimedOut {
			state.UnsyncList.RemoveNode(nodeID)
		}

		type none struct{}
		newActive := make(map[core.RecordRef]none)
		for _, active := range matrixCalculation.Active {
			newActive[active] = none{}
		}

		newProofs := make(map[core.Node]*merkle.PulseProof)
		for node, proof := range state.ValidProofs {
			_, ok := newActive[node.ID()]
			if !ok {
				continue
			}
			newProofs[node] = proof
		}
		newProofs[sp.NodeKeeper.GetOrigin()] = state.PulseProof

		state.ValidProofs = newProofs

		entry := &merkle.GlobuleEntry{
			PulseEntry:    state.PulseEntry,
			ProofSet:      state.ValidProofs,
			PulseHash:     state.PulseHash,
			PrevCloudHash: prevCloudHash,
			GlobuleID:     sp.NodeKeeper.GetOrigin().GetGlobuleID(),
		}
		globuleHash, globuleProof, err = sp.Calculator.GetGlobuleProof(entry)

		if err != nil {
			return nil, errors.Wrap(err, "[ NET Consensus phase-2.0 ] Failed to calculate globule proof")
		}
	}

	return &SecondPhaseState{
		FirstPhaseState: state,
		Matrix:          stateMatrix,
		MatrixState:     matrixCalculation,
		BitSet:          bitset,

		GlobuleHash:  globuleHash,
		GlobuleProof: globuleProof,
	}, nil
}

func (sp *SecondPhaseImpl) Execute21(ctx context.Context, pulse *core.Pulse, state *SecondPhaseState) (*SecondPhaseState, error) {
	ctx, span := instracer.StartSpan(ctx, "SecondPhase.Execute21")
	span.AddAttributes(trace.Int64Attribute("pulse", int64(state.PulseEntry.Pulse.PulseNumber)))
	defer span.End()
	stats.Record(ctx, consensus.Phase21Exec.M(1))
	additionalRequests := state.MatrixState.AdditionalRequestsPhase2

	logger := inslogger.FromContext(ctx)
	logger.Infof("[ NET Consensus phase-2.1 ] Additional requests needed: %d", len(additionalRequests))

	results := make(map[uint16]*packets.MissingNodeSupplementaryVote)
	claims := make(map[uint16]*packets.MissingNodeClaim)

	packet := packets.NewPhase2Packet(pulse.PulseNumber)
	err := packet.SetGlobuleHashSignature(state.GlobuleProof.Signature.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "[ NET Consensus phase-2.1 ] Failed to set pulse proof in Phase2Packet")
	}
	packet.SetBitSet(state.BitSet)

	voteAnswers, err := sp.Communicator.ExchangePhase21(ctx, state.UnsyncList, state.ClaimHandler, packet, additionalRequests)
	if err != nil {
		return nil, errors.Wrap(err, "[ NET Consensus phase-2.1 ] Failed to send additional requests")
	}

	if len(additionalRequests) == 0 {
		return state, nil
	}

	for _, vote := range voteAnswers {
		switch v := vote.(type) {
		case *packets.MissingNodeSupplementaryVote:
			results[v.NodeIndex] = v
		case *packets.MissingNodeClaim:
			claims[v.NodeIndex] = v
		}
	}

	err = stats.RecordWithTags(ctx, []tag.Mutator{tag.Upsert(consensus.TagPhase, "phase 21")}, consensus.PacketsRecv.M(int64(len(results))))
	if err != nil {
		logger.Warn("[ NET Consensus phase-2.1 ] Failed to record received results metric: " + err.Error())
	}
	if len(results) != len(additionalRequests) {
		return nil, errors.Errorf("[ NET Consensus phase-2.1 ] Failed to receive enough MissingNodeSupplementaryVote responses,"+
			" received: %d/%d", len(results), len(additionalRequests))
	}

	origin := sp.NodeKeeper.GetOrigin().ID()
	bitsetChanges := make([]packets.BitSetCell, 0)
	for index, result := range results {
		claim := result.NodeClaimUnsigned
		node, err := nodenetwork.ClaimToNode("", &claim)
		if err != nil {
			return nil, errors.Wrapf(err, "[ NET Consensus phase-2.1 ] Failed to convert claim to node, "+
				"ref: %s", claim.NodeRef)
		}

		merkleProof := &merkle.PulseProof{
			BaseProof: merkle.BaseProof{
				Signature: core.SignatureFromBytes(result.NodePulseProof.Signature()),
			},
			StateHash: result.NodePulseProof.StateHash(),
		}

		state.UnsyncList.AddNode(node, index)
		err = sp.NodeKeeper.AddTemporaryMapping(claim.NodeRef, claim.ShortNodeID, claim.NodeAddress.Get())
		if err != nil {
			logger.Warn("Error adding temporary mapping: " + err.Error())
		}
		valid := validateProof(sp.Calculator, state.UnsyncList, state.PulseHash, node.ID(), merkleProof)
		if !valid {
			logger.Warnf("[ NET Consensus phase-2.1 ] Failed to validate proof from %s", node.ID())
			state.UnsyncList.RemoveNode(node.ID())
			continue
		}

		err = state.Matrix.ReceivedProofFromNode(origin, node.ID())
		if err != nil {
			return nil, errors.Wrapf(err, "[ NET Consensus phase-2.1 ] Failed to assign proof from node %s "+
				"to state matrix", claim.NodeRef)
		}
		state.UnsyncList.AddProof(node.ID(), &result.NodePulseProof)
		state.ValidProofs[node] = merkleProof
		bitsetChanges = append(bitsetChanges, packets.BitSetCell{NodeID: node.ID(), State: packets.Legit})
	}

	err = state.BitSet.ApplyChanges(bitsetChanges, state.UnsyncList)
	if err != nil {
		return nil, errors.Wrap(err, "[ NET Consensus phase-2.1 ] Failed to apply changes to current bitset")
	}
	claimMap := make(map[core.RecordRef][]packets.ReferendumClaim)
	for index, claim := range claims {
		ref, err := state.UnsyncList.IndexToRef(int(index))
		if err != nil {
			return nil, errors.Wrapf(err, "[ NET Consensus phase-2.1 ] Failed to map index %d to ref", index)
		}
		list := claimMap[ref]
		if list == nil {
			list = make([]packets.ReferendumClaim, 0)
		}
		list = append(list, claim.Claim)
		claimMap[ref] = list
	}
	for ref, claims := range claimMap {
		state.ClaimHandler.SetClaimsFromNode(ref, claims)
	}
	state.MatrixState, err = state.Matrix.CalculatePhase2(origin)
	if err != nil {
		return nil, errors.Wrap(err, "[ NET Consensus phase-2.1 ] Failed to calculate matrix state")
	}
	addReqCount := len(state.MatrixState.AdditionalRequestsPhase2)
	if addReqCount != 0 {
		return nil, errors.Errorf("[ NET Consensus phase-2.1 ] Failed to get enough data during phase 2.1 "+
			"(still need additional %d requests)", addReqCount)
	}

	prevCloudHash := sp.NodeKeeper.GetCloudHash()
	entry := &merkle.GlobuleEntry{
		PulseEntry:    state.PulseEntry,
		ProofSet:      state.ValidProofs,
		PulseHash:     state.PulseHash,
		PrevCloudHash: prevCloudHash,
		GlobuleID:     sp.NodeKeeper.GetOrigin().GetGlobuleID(),
	}
	state.GlobuleHash, state.GlobuleProof, err = sp.Calculator.GetGlobuleProof(entry)
	if err != nil {
		return nil, errors.Wrap(err, "[ NET Consensus phase-2.1 ] Failed to calculate globule proof")
	}
	var ghs packets.GlobuleHashSignature
	copy(ghs[:], state.GlobuleProof.Signature.Bytes()[:packets.SignatureLength])
	state.UnsyncList.SetGlobuleHashSignature(origin, ghs)

	return state, nil
}

func (sp *SecondPhaseImpl) generatePhase2Bitset(list network.UnsyncList, proofs map[core.Node]*merkle.PulseProof, pulseNumber core.PulseNumber) (packets.BitSet, error) {
	bitset, err := packets.NewBitSet(list.Length())
	if err != nil {
		return nil, err
	}
	cells := make([]packets.BitSetCell, 0)
	for node := range proofs {
		cells = append(cells, packets.BitSetCell{
			NodeID: node.ID(),
			State:  getNodeState(node, pulseNumber),
		})
	}
	cells = append(cells, packets.BitSetCell{
		NodeID: sp.NodeKeeper.GetOrigin().ID(),
		State:  getNodeState(sp.NodeKeeper.GetOrigin(), pulseNumber),
	})
	err = bitset.ApplyChanges(cells, list)
	if err != nil {
		return nil, err
	}
	return bitset, nil
}

func getNodeState(node core.Node, pulseNumber core.PulseNumber) packets.BitSetState {
	state := packets.Legit
	if node.Leaving() && node.LeavingETA() < pulseNumber {
		state = packets.TimedOut
	}

	return state
}

func (sp *SecondPhaseImpl) checkPacketSignature(packet *packets.Phase2Packet, recordRef core.RecordRef, unsyncList network.UnsyncList) error {
	activeNode := unsyncList.GetActiveNode(recordRef)
	if activeNode == nil {
		return errors.New("failed to get active node")
	}
	key := activeNode.PublicKey()
	return packet.Verify(sp.Cryptography, key)
}
