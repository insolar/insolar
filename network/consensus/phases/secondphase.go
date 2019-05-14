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

package phases

import (
	"context"

	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/node"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus/packets"
	"github.com/insolar/insolar/network/merkle"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

type SecondPhase interface {
	Execute(ctx context.Context, pulse *insolar.Pulse, state *FirstPhaseState) (*SecondPhaseState, error)
	Execute21(ctx context.Context, pulse *insolar.Pulse, state *SecondPhaseState) (*SecondPhaseState, error)
}

func NewSecondPhase() SecondPhase {
	return &SecondPhaseImpl{}
}

type SecondPhaseImpl struct {
	NodeKeeper   network.NodeKeeper          `inject:""`
	Calculator   merkle.Calculator           `inject:""`
	Communicator Communicator                `inject:""`
	Cryptography insolar.CryptographyService `inject:""`
}

func (sp *SecondPhaseImpl) Execute(ctx context.Context, pulse *insolar.Pulse, state *FirstPhaseState) (*SecondPhaseState, error) {
	logger := inslogger.FromContext(ctx)
	ctx, span := instracer.StartSpan(ctx, "SecondPhase.Execute")
	span.AddAttributes(trace.Int64Attribute("pulse", int64(state.PulseEntry.Pulse.PulseNumber)))
	defer span.End()
	prevCloudHash := sp.NodeKeeper.GetCloudHash()

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
	packet.SetBitSet(state.BitSet)
	participants := state.NodesMutator.GetActiveNodes()

	packets, err := sp.Communicator.ExchangePhase2(ctx, state.ConsensusState, participants, packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ NET Consensus phase-2.0 ] Failed to exchange packets")
	}
	logger.Infof("[ NET Consensus phase-2.0 ] Received responses: %d/%d",
		len(packets), state.BitsetMapper.Length())
	err = stats.RecordWithTags(ctx, []tag.Mutator{tag.Upsert(consensus.TagPhase, "phase 2")}, consensus.PacketsRecv.M(int64(len(packets))))
	if err != nil {
		logger.Warn("[ NET Consensus phase-2.0 ] Failed to record received packets metric: " + err.Error())
	}

	origin := sp.NodeKeeper.GetOrigin().ID()
	stateMatrix := NewStateMatrix(state.BitsetMapper)

	for ref, packet := range packets {
		err = nil
		if !ref.Equal(origin) {
			err = sp.checkPacketSignature(packet, ref, state.NodesMutator)
		}
		if err != nil {
			logger.Warnf("[ NET Consensus phase-2.0 ] Failed to check phase2 packet signature from %s: %s", ref, err.Error())
			continue
		}
		state.HashStorage.SetGlobuleHashSignature(ref, packet.GetGlobuleHashSignature())
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
		type none struct{}
		newActive := make(map[insolar.Reference]none)
		for _, active := range matrixCalculation.Active {
			newActive[active] = none{}
		}

		newProofs := make(map[insolar.NetworkNode]*merkle.PulseProof)
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

		GlobuleHash:  globuleHash,
		GlobuleProof: globuleProof,
	}, nil
}

func (sp *SecondPhaseImpl) Execute21(ctx context.Context, pulse *insolar.Pulse, state *SecondPhaseState) (*SecondPhaseState, error) {
	ctx, span := instracer.StartSpan(ctx, "SecondPhase.Execute21")
	span.AddAttributes(trace.Int64Attribute("pulse", int64(state.PulseEntry.Pulse.PulseNumber)))
	defer span.End()
	stats.Record(ctx, consensus.Phase21Exec.M(1))
	additionalRequests := state.MatrixState.AdditionalRequestsPhase2

	logger := inslogger.FromContext(ctx)
	logger.Infof("[ NET Consensus phase-2.1 ] Additional requests needed: %d", len(additionalRequests))

	results := make(map[uint16]*packets.MissingNodeRespVote)
	claims := make(map[uint16]*packets.MissingNodeClaimsVote)

	packet := packets.NewPhase2Packet(pulse.PulseNumber)
	err := packet.SetGlobuleHashSignature(state.GlobuleProof.Signature.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "[ NET Consensus phase-2.1 ] Failed to set pulse proof in Phase2Packet")
	}
	packet.SetBitSet(state.BitSet)

	voteAnswers, err := sp.Communicator.ExchangePhase21(ctx, state.ConsensusState, packet, additionalRequests)
	if err != nil {
		return nil, errors.Wrap(err, "[ NET Consensus phase-2.1 ] Failed to send additional requests")
	}

	if len(additionalRequests) == 0 {
		return state, nil
	}

	for _, vote := range voteAnswers {
		switch v := vote.(type) {
		case *packets.MissingNodeRespVote:
			results[v.NodeIndex] = v
		case *packets.MissingNodeClaimsVote:
			claims[v.NodeIndex] = v
		}
	}

	err = stats.RecordWithTags(ctx, []tag.Mutator{tag.Upsert(consensus.TagPhase, "phase 21")}, consensus.PacketsRecv.M(int64(len(results))))
	if err != nil {
		logger.Warn("[ NET Consensus phase-2.1 ] Failed to record received results metric: " + err.Error())
	}
	if len(results) != len(additionalRequests) {
		return nil, errors.Errorf("[ NET Consensus phase-2.1 ] Failed to receive enough MissingNodeRespVote responses,"+
			" received: %d/%d", len(results), len(additionalRequests))
	}

	origin := sp.NodeKeeper.GetOrigin().ID()
	bitsetChanges := make([]packets.BitSetCell, 0)
	for index, result := range results {
		claim := result.NodeClaimUnsigned
		node, err := node.ClaimToNode("", &claim)
		if err != nil {
			return nil, errors.Wrapf(err, "[ NET Consensus phase-2.1 ] Failed to convert claim to node, "+
				"ref: %s", claim.NodeRef)
		}

		merkleProof := &merkle.PulseProof{
			BaseProof: merkle.BaseProof{
				Signature: insolar.SignatureFromBytes(result.NodePulseProof.Signature()),
			},
			StateHash: result.NodePulseProof.StateHash(),
		}

		state.NodesMutator.AddWorkingNode(node)
		state.BitsetMapper.AddNode(node, index)
		err = state.ConsensusInfo.AddTemporaryMapping(claim.NodeRef, claim.ShortNodeID, claim.NodeAddress.Get())
		if err != nil {
			logger.Warn("Error adding temporary mapping: " + err.Error())
		}
		valid := validateProof(sp.Calculator, state.NodesMutator, state.PulseHash, node.ID(), merkleProof)
		if !valid {
			logger.Warnf("[ NET Consensus phase-2.1 ] Failed to validate proof from %s", node.ID())
			continue
		}

		err = state.Matrix.ReceivedProofFromNode(origin, node.ID())
		if err != nil {
			return nil, errors.Wrapf(err, "[ NET Consensus phase-2.1 ] Failed to assign proof from node %s "+
				"to state matrix", claim.NodeRef)
		}
		state.HashStorage.AddProof(node.ID(), &result.NodePulseProof)
		state.ValidProofs[node] = merkleProof
		bitsetChanges = append(bitsetChanges, packets.BitSetCell{NodeID: node.ID(), State: packets.Legit})
	}

	err = state.BitSet.ApplyChanges(bitsetChanges, state.BitsetMapper)
	if err != nil {
		return nil, errors.Wrap(err, "[ NET Consensus phase-2.1 ] Failed to apply changes to current bitset")
	}
	claimMap := make(map[insolar.Reference][]packets.ReferendumClaim)
	for index, claim := range claims {
		ref, err := state.BitsetMapper.IndexToRef(int(index))
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
	state.HashStorage.SetGlobuleHashSignature(origin, ghs)

	return state, nil
}

func (sp *SecondPhaseImpl) checkPacketSignature(packet *packets.Phase2Packet, recordRef insolar.Reference, accessor network.Accessor) error {
	activeNode := accessor.GetActiveNode(recordRef)
	if activeNode == nil {
		return errors.New("failed to get active node")
	}
	key := activeNode.PublicKey()
	return packet.Verify(sp.Cryptography, key)
}
