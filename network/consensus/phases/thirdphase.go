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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/consensus/packets"
	"github.com/insolar/insolar/network/merkle"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

type ThirdPhase interface {
	Execute(ctx context.Context, pulse *insolar.Pulse, state *SecondPhaseState) (*ThirdPhaseState, error)
}

func NewThirdPhase() ThirdPhase {
	return &ThirdPhaseImpl{}
}

type ThirdPhaseImpl struct {
	Cryptography insolar.CryptographyService `inject:""`
	Communicator Communicator                `inject:""`
	NodeKeeper   network.NodeKeeper          `inject:""`
	Calculator   merkle.Calculator           `inject:""`
}

func (tp *ThirdPhaseImpl) Execute(ctx context.Context, pulse *insolar.Pulse, state *SecondPhaseState) (*ThirdPhaseState, error) {
	ctx, span := instracer.StartSpan(ctx, "ThirdPhase.Execute")
	span.AddAttributes(trace.Int64Attribute("pulse", int64(state.PulseEntry.Pulse.PulseNumber)))
	defer span.End()
	stats.Record(ctx, consensus.Phase3Exec.M(1))

	logger := inslogger.FromContext(ctx)
	totalCount := state.BitsetMapper.Length()

	var gSign [packets.SignatureLength]byte
	copy(gSign[:], state.GlobuleProof.Signature.Bytes()[:packets.SignatureLength])
	packet := packets.NewPhase3Packet(pulse.PulseNumber, gSign, state.BitSet)

	nodes := make([]insolar.NetworkNode, 0)
	for _, node := range state.MatrixState.Active {
		nodes = append(nodes, state.NodesMutator.GetActiveNode(node))
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

	for ref, packet := range responses {
		err = nil
		if !ref.Equal(tp.NodeKeeper.GetOrigin().ID()) {
			err = tp.checkPacketSignature(packet, ref, state.NodesMutator)
		}
		if err != nil {
			logger.Warnf("[ NET Consensus phase-3 ] Failed to check phase3 packet signature from %s: %s", ref, err.Error())
			continue
		}
		// not needed until we implement fraud detection
		// cells, err := packet.GetBitset().GetCells(state.UnsyncList)

		state.HashStorage.SetGlobuleHashSignature(ref, packet.GetGlobuleHashSignature())
	}

	prevCloudHash := tp.NodeKeeper.GetCloudHash()
	validNodes := 0
	for _, node := range nodes {
		ghs, ok := state.HashStorage.GetGlobuleHashSignature(node.ID())
		if !ok {
			log.Warnf("[ NET Consensus phase-3 ] No globule hash signature for node %s", node.ID())
			continue
		}
		proof := &merkle.GlobuleProof{
			BaseProof: merkle.BaseProof{
				Signature: insolar.SignatureFromBytes(ghs[:]),
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

	claimSplit := state.ClaimHandler.FilterClaims(state.MatrixState.Active, pulse.Entropy)

	return &ThirdPhaseState{
		ActiveNodes:    nodes,
		GlobuleProof:   state.GlobuleProof,
		ApprovedClaims: claimSplit.ApprovedClaims,
		ReconnectTo:    checkReconnectClaim(claimSplit.ApprovedClaims),
	}, nil
}

func checkReconnectClaim(claims []packets.ReferendumClaim) string {
	for _, claim := range claims {
		if claim.Type() == packets.TypeChangeNetworkClaim {
			return claim.(*packets.ChangeNetworkClaim).Address
		}
	}
	return ""
}

func (tp *ThirdPhaseImpl) checkPacketSignature(packet *packets.Phase3Packet, recordRef insolar.Reference, accessor network.Accessor) error {
	activeNode := accessor.GetActiveNode(recordRef)
	if activeNode == nil {
		return errors.New("failed to get active node")
	}
	key := activeNode.PublicKey()
	return packet.Verify(tp.Cryptography, key)
}
