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
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/merkle"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/pkg/errors"
)

// SecondPhase is a second phase.
type SecondPhase struct {
	NodeKeeper   network.NodeKeeper       `inject:""`
	Calculator   merkle.Calculator        `inject:""`
	Communicator Communicator             `inject:""`
	Cryptography core.CryptographyService `inject:""`
}

func (sp *SecondPhase) Execute(ctx context.Context, state *FirstPhaseState) (*SecondPhaseState, error) {
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
		return nil, errors.Wrap(err, "[ SecondPhase ] Failed to calculate pulse proof.")
	}

	packet := packets.Phase2Packet{}
	err = packet.SetGlobuleHashSignature(globuleProof.Signature.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "[ SecondPhase ] Failed to set pulse proof in Phase2Packet.")
	}
	bitset, err := sp.generatePhase2Bitset(state.UnsyncList, state.ValidProofs)
	if err != nil {
		return nil, errors.Wrap(err, "[ SecondPhase ] Failed to generate bitset for Phase2Packet")
	}
	packet.SetBitSet(bitset)
	err = sp.signPhase2Packet(&packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ SecondPhase ] Failed to sign a packet")
	}
	activeNodes := state.UnsyncList.GetActiveNodes()
	packets, err := sp.Communicator.ExchangePhase2(ctx, state.UnsyncList, activeNodes, &packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ SecondPhase ] Failed to exchange results.")
	}

	nodeProofs := make(map[core.Node]*GlobuleProofValidated)
	stateMatrix := NewStateMatrix(state.UnsyncList)

	for ref, packet := range packets {
		signIsCorrect, err := sp.isSignPhase2PacketRight(packet, ref)
		if err != nil {
			log.Warn("failed to check a sign: ", err.Error())
		} else if !signIsCorrect {
			log.Warn("recieved a bad sign packet: ", err.Error())
		}
		err = stateMatrix.ApplyBitSet(ref, packet.GetBitSet())
		if err != nil {
			log.Warnf("[ SecondPhase ] could not apply bitset from node %s", ref)
			continue
		}

		ghs := packet.GetGlobuleHashSignature()
		state.UnsyncList.SetGlobuleHashSignature(ref, ghs)
		node := state.UnsyncList.GetActiveNode(ref)
		proof := &merkle.GlobuleProof{
			BaseProof: merkle.BaseProof{
				Signature: core.SignatureFromBytes(ghs[:]),
			},
			PrevCloudHash: prevCloudHash,
			GlobuleID:     globuleProof.GlobuleID,
			NodeCount:     globuleProof.NodeCount,
			NodeRoot:      globuleProof.NodeRoot,
		}
		valid := sp.Calculator.IsValid(proof, globuleHash, node.PublicKey())
		nodeProofs[node] = &GlobuleProofValidated{Proof: proof, Valid: valid}
	}

	matrixCalculation, err := stateMatrix.CalculatePhase2(sp.NodeKeeper.GetOrigin().ID())
	if err != nil {
		return nil, errors.Wrap(err, "[ SecondPhase ] Failed to calculate bitset matrix consensus result")
	}

	return &SecondPhaseState{
		FirstPhaseState: state,
		Matrix:          stateMatrix,
		MatrixState:     matrixCalculation,
		BitSet:          bitset,

		GlobuleEntry:    entry,
		GlobuleHash:     globuleHash,
		GlobuleProof:    globuleProof,
		GlobuleProofSet: nodeProofs,
	}, nil
}

func (sp *SecondPhase) Execute21(ctx context.Context, state *SecondPhaseState) (*SecondPhaseState, error) {
	additionalRequests := state.MatrixState.AdditionalRequestsPhase2
	if len(additionalRequests) == 0 {
		return state, nil
	}

	count := len(additionalRequests)
	results := make(map[uint16]*packets.MissingNodeSupplementaryVote)
	claims := make(map[uint16]*packets.MissingNodeClaim)

	packet := packets.Phase2Packet{}
	err := packet.SetGlobuleHashSignature(state.GlobuleProof.Signature.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase 2.1 ] Failed to set pulse proof in Phase2Packet.")
	}
	packet.SetBitSet(state.BitSet)
	err = sp.signPhase2Packet(&packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase 2.1 ] Failed to sign a packet")
	}

	voteAnswers, err := sp.Communicator.ExchangePhase21(ctx, state.UnsyncList, &packet, additionalRequests)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase 2.1 ] Failed to send additional requests.")
	}

	for _, vote := range voteAnswers {
		switch v := vote.(type) {
		case *packets.MissingNodeSupplementaryVote:
			results[v.NodeIndex] = v
		case *packets.MissingNodeClaim:
			claims[v.NodeIndex] = v
		}
	}

	if len(results) != count {
		return nil, errors.New(fmt.Sprintf("[ Phase 2.1 ] Failed to receive enough MissingNodeSupplementaryVote responses: %d/%d", len(results), count))
	}

	for index, result := range results {
		node, err := nodenetwork.ClaimToNode("", &result.NodeClaimUnsigned)
		if err != nil {
			return nil, errors.Wrapf(err, "[ Phase 2.1 ] Failed to convert claim to node, ref: %s", result.NodeClaimUnsigned.NodeRef)
		}
		state.UnsyncList.AddNode(node, index)
		state.UnsyncList.AddProof(node.ID(), &result.NodePulseProof)
		state.UnsyncList.SetGlobuleHashSignature(node.ID(), result.GlobuleHashSignature)
	}
	claimMap := make(map[core.RecordRef][]packets.ReferendumClaim)
	for index, claim := range claims {
		ref, err := state.UnsyncList.IndexToRef(int(index))
		if err != nil {
			return nil, errors.Wrapf(err, "[ Phase 2.1 ] Failed to map index %d to ref", index)
		}
		list := claimMap[ref]
		if list == nil {
			list = make([]packets.ReferendumClaim, 0)
		}
		list = append(list, claim.Claim)
		claimMap[ref] = list
	}
	state.UnsyncList.AddClaims(claimMap)

	// cloudEntry := &merkle.CloudEntry{
	//
	// }

	// cloudHash, _, _ := sp.Calculator.GetCloudProof(cloudEntry)

	return state, nil
}

func (sp *SecondPhase) generatePhase2Bitset(list network.UnsyncList, proofs map[core.Node]*merkle.PulseProof) (packets.BitSet, error) {
	bitset, err := packets.NewBitSet(list.Length())
	if err != nil {
		return nil, err
	}
	cells := make([]packets.BitSetCell, 0)
	for node := range proofs {
		cells = append(cells, packets.BitSetCell{NodeID: node.ID(), State: packets.Legit})
	}
	cells = append(cells, packets.BitSetCell{NodeID: sp.NodeKeeper.GetOrigin().ID(), State: packets.Legit})
	err = bitset.ApplyChanges(cells, list)
	if err != nil {
		return nil, err
	}
	return bitset, nil
}

func (sp *SecondPhase) signPhase2Packet(p *packets.Phase2Packet) error {
	data, err := p.RawFirstPart()
	if err != nil {
		return errors.Wrap(err, "failed to get raw bytes")
	}
	sign, err := sp.Cryptography.Sign(data)
	if err != nil {
		return errors.Wrap(err, "failed to sign a phase 2 packet")
	}

	copy(p.SignatureHeaderSection1[:], sign.Bytes())
	// TODO: sign a second part after claim addition
	return nil
}

func (sp *SecondPhase) isSignPhase2PacketRight(packet *packets.Phase2Packet, recordRef core.RecordRef) (bool, error) {
	key := sp.NodeKeeper.GetActiveNode(recordRef).PublicKey()

	raw, err := packet.RawFirstPart()
	if err != nil {
		return false, errors.Wrap(err, "failed to serialize")
	}

	return sp.Cryptography.Verify(key, core.SignatureFromBytes(raw), raw), nil
}
