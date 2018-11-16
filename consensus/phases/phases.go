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
	"github.com/insolar/insolar/network"
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
	NodeKeeper   network.NodeKeeper `inject:""`
	Calculator   merkle.Calculator  `inject:""`
	Communicator Communicator       `inject:""`
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

	if fp.NodeKeeper.NodesJoinedDuringPreviousPulse() {
		err = packet.AddClaim(fp.NodeKeeper.GetOriginClaim())
		if err != nil {
			return nil, errors.Wrap(err, "[ Execute ] Failed to add origin claim in Phase1Packet.")
		}
	}
	// TODO: add other claims

	activeNodes := fp.NodeKeeper.GetActiveNodes()
	resultPackets, err := fp.Communicator.ExchangePhase1(ctx, activeNodes, packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to exchange results.")
	}

	proofSet := make(map[core.RecordRef]*merkle.PulseProof)
	claimSet := make([]packets.ReferendumClaim, 0)
	for ref, packet := range resultPackets {
		rawProof := packet.GetPulseProof()
		proofSet[ref] = &merkle.PulseProof{
			BaseProof: merkle.BaseProof{
				Signature: core.SignatureFromBytes(rawProof.Signature()),
			},
			StateHash: rawProof.StateHash(),
		}
		claims := packet.GetClaims() // TODO: build set of claims
		claimSet = append(claimSet, claims...)
	}

	fp.processClaims(ctx, claimSet)

	// Get active nodes again for case when current node just joined to network and don't have full active list
	activeNodes = fp.NodeKeeper.GetActiveNodes()
	deviantsNodes := make([]core.Node, 0)
	timedOutNodes := make([]core.Node, 0)
	validProofs := make(map[core.Node]*merkle.PulseProof)

	for _, node := range activeNodes {
		proof, ok := proofSet[node.ID()]
		if !ok {
			timedOutNodes = append(timedOutNodes, node)
		}

		if !fp.Calculator.IsValid(proof, pulseHash, node.PublicKey()) {
			validProofs[node] = proof
		} else {
			deviantsNodes = append(deviantsNodes, node)
		}
	}

	return &FirstPhaseState{
		PulseEntry:    entry,
		PulseHash:     pulseHash,
		PulseProof:    pulseProof,
		PulseProofSet: validProofs,
		TimedOutNodes: timedOutNodes,
		DeviantNodes:  deviantsNodes,
	}, nil
}

func (fp *FirstPhase) processClaims(ctx context.Context, claims []packets.ReferendumClaim) {
	var nodes []core.Node
	var unsyncClaims []packets.ReferendumClaim

	for _, genericClaim := range claims {
		switch claim := genericClaim.(type) {
		case *packets.NodeAnnounceClaim:
			nodes = append(nodes, claim.Node())
		case *packets.NodeJoinClaim:
			panic("Not implemented yet") // TODO: authorize node here
		case *packets.NodeLeaveClaim:
			unsyncClaims = append(unsyncClaims, claim)
		default:
			panic("Not implemented yet")
		}
	}

	fp.NodeKeeper.AddActiveNodes(nodes)
	fp.NodeKeeper.AddUnsyncClaims(unsyncClaims)
}

// SecondPhase is a second phase.
type SecondPhase struct {
	NodeKeeper   network.NodeKeeper `inject:""`
	Network      core.Network       `inject:""`
	Calculator   merkle.Calculator  `inject:""`
	Communicator Communicator       `inject:""`
}

func (sp *SecondPhase) Execute(ctx context.Context, state *FirstPhaseState) (*SecondPhaseState, error) {
	prevCloudHash := sp.NodeKeeper.GetCloudHash()
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

	activeNodes := sp.NodeKeeper.GetActiveNodes()
	proofSet, err := sp.Communicator.ExchangePhase2(ctx, activeNodes, packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to exchange results.")
	}

	nodeProofs := make(map[core.Node]*merkle.GlobuleProof)

	var deviants []core.Node
	deviants = append(deviants, state.TimedOutNodes...)
	deviants = append(deviants, state.DeviantNodes...)

	for ref, packet := range proofSet {
		node := sp.NodeKeeper.GetActiveNode(ref)
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

	// TODO: check
	if !consensusReached(len(nodeProofs), len(activeNodes)) {
		return nil, errors.New("[ Execute ] Consensus not reached")
	}

	sp.NodeKeeper.Sync(deviants)

	return &SecondPhaseState{
		FirstPhaseState: state,

		GlobuleEntry:    entry,
		GlobuleHash:     globuleHash,
		GlobuleProof:    globuleProof,
		GlobuleProofSet: nodeProofs,
	}, nil
}

func (sp *SecondPhase) processTimedOutNodes(timedOutNodes []core.Node) {
	// TODO: process
}

func (sp *SecondPhase) calculateListForNextPulse() (uint16, []byte) {
	// TODO: calculate
	return 1337, []byte("1337")
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
