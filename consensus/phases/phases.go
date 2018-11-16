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
	"math"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/merkle"
	"github.com/pkg/errors"
)

const ConsensusAtPercents = 0.66

func consensusReached(resultLen, participanstLen int) bool {
	minParticipants := int(math.Floor(ConsensusAtPercents*float64(participanstLen))) + 1

	return resultLen >= minParticipants
}

// FirstPhase is a first phase.
type FirstPhase struct {
	NodeNetwork  core.NodeNetwork  `inject:""`
	Calculator   merkle.Calculator `inject:""`
	Communicator Communicator      `inject:""`
}

// Execute do first phase
func (fp *FirstPhase) Execute(ctx context.Context, pulse *core.Pulse) (*FirstPhaseState, error) {
	entry := &merkle.PulseEntry{Pulse: pulse}
	pulseHash, pulseProof, err := fp.Calculator.GetPulseProof(entry)

	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to calculate pulse proof.")
	}

	packet := packets.Phase1Packet{}
	err = packet.SetPulseProof(pulseProof.StateHash, pulseProof.Signature.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to set pulse proof in Phase1Packet.")
	}

	activeNodes := fp.NodeNetwork.GetActiveNodes()
	proofSet, err := fp.Communicator.ExchangePhase1(ctx, activeNodes, packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to exchange results.")
	}

	nodeProofs := make(map[core.Node]*merkle.PulseProof)
	for ref, packet := range proofSet {
		node := fp.NodeNetwork.GetActiveNode(ref)
		rawProof := packet.GetPulseProof()
		proof := &merkle.PulseProof{
			BaseProof: merkle.BaseProof{
				Signature: core.SignatureFromBytes(rawProof.Signature()),
			},
			StateHash: rawProof.StateHash(),
		}

		if !fp.Calculator.IsValid(proof, pulseHash, node.PublicKey()) {
			nodeProofs[node] = proof
		}
	}

	if !consensusReached(len(nodeProofs), len(activeNodes)) {
		return nil, errors.New("[ Execute ] Consensus not reached")
	}

	return &FirstPhaseState{
		PulseEntry:    entry,
		PulseHash:     pulseHash,
		PulseProof:    pulseProof,
		PulseProofSet: nodeProofs,
	}, nil
}

// SecondPhase is a second phase.
type SecondPhase struct {
	NodeNetwork  core.NodeNetwork  `inject:""`
	Network      core.Network      `inject:""`
	Calculator   merkle.Calculator `inject:""`
	Communicator Communicator      `inject:""`
}

func (sp *SecondPhase) Execute(ctx context.Context, state *FirstPhaseState) (*SecondPhaseState, error) {
	prevCloudHash := sp.NodeNetwork.GetCloudHash()
	globuleID := sp.Network.GetGlobuleID()

	entry := &merkle.GlobuleEntry{
		PulseEntry:    state.PulseEntry,
		ProofSet:      state.PulseProofSet,
		PulseHash:     state.PulseHash,
		PrevCloudHash: prevCloudHash,
		GlobuleID:     globuleID,
	}
	globuleHash, globuleProof, err := sp.Calculator.GetGlobuleProof(entry)

	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to calculate pulse proof.")
	}

	packet := packets.Phase2Packet{}
	err = packet.SetGlobuleHashSignature(globuleProof.Signature.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to set pulse proof in Phase2Packet.")
	}

	activeNodes := sp.NodeNetwork.GetActiveNodes()
	proofSet, err := sp.Communicator.ExchangePhase2(ctx, activeNodes, packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to exchange results.")
	}

	nodeProofs := make(map[core.Node]*merkle.GlobuleProof)
	for ref, packet := range proofSet {
		node := sp.NodeNetwork.GetActiveNode(ref)
		proof := &merkle.GlobuleProof{
			BaseProof: merkle.BaseProof{
				Signature: core.SignatureFromBytes(packet.GetGlobuleHashSignature()),
			},
			PrevCloudHash: prevCloudHash,
			GlobuleID:     globuleProof.GlobuleID,
			NodeCount:     globuleProof.NodeCount,
			NodeRoot:      globuleProof.NodeRoot,
		}

		if !sp.Calculator.IsValid(proof, globuleHash, node.PublicKey()) {
			nodeProofs[node] = proof
		}
	}

	if !consensusReached(len(nodeProofs), len(activeNodes)) {
		return nil, errors.New("[ Execute ] Consensus not reached")
	}

	return &SecondPhaseState{
		FirstPhaseState: state,

		GlobuleEntry:    entry,
		GlobuleHash:     globuleHash,
		GlobuleProof:    globuleProof,
		GlobuleProofSet: nodeProofs,
	}, nil
}

// ThirdPhasePulse.
type ThirdPhasePulse struct {
	NodeNetwork core.NodeNetwork `inject:""`
	State       *ThirdPhasePulseState
}

func (tpp *ThirdPhasePulse) Execute(ctx context.Context, state *SecondPhaseState) error {
	// TODO: do something here
	return nil
}

// ThirdPhaseReferendum.
type ThirdPhaseReferendum struct {
	NodeNetwork core.NodeNetwork `inject:""`
	State       *ThirdPhaseReferendumState
}

func (tpr *ThirdPhaseReferendum) Execute(ctx context.Context, state *SecondPhaseState) error {
	// TODO: do something here
	return nil
}
