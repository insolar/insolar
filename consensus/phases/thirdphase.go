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
	"fmt"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/merkle"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

type ThirdPhase interface {
	Execute(ctx context.Context, pulse *core.Pulse, state *SecondPhaseState) (*ThirdPhaseState, error)
}

func NewThirdPhase() ThirdPhase {
	return &thirdPhase{}
}

type thirdPhase struct {
	Cryptography core.CryptographyService `inject:""`
	Communicator Communicator             `inject:""`
	NodeKeeper   network.NodeKeeper       `inject:""`
	Calculator   merkle.Calculator        `inject:""`
}

func (tp *thirdPhase) Execute(ctx context.Context, pulse *core.Pulse, state *SecondPhaseState) (*ThirdPhaseState, error) {
	ctx, span := instracer.StartSpan(ctx, "ThirdPhase.Execute")
	span.AddAttributes(trace.Int64Attribute("pulse", int64(state.PulseEntry.Pulse.PulseNumber)))
	defer span.End()
	metrics.ConsensusPhase3Exec.Inc()
	var gSign [packets.SignatureLength]byte
	copy(gSign[:], state.GlobuleProof.Signature.Bytes()[:packets.SignatureLength])
	packet := packets.NewPhase3Packet(pulse.PulseNumber, gSign, state.BitSet)

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

func (tp *thirdPhase) checkPacketSignature(packet *packets.Phase3Packet, recordRef core.RecordRef, unsyncList network.UnsyncList) error {
	activeNode := unsyncList.GetActiveNode(recordRef)
	if activeNode == nil {
		return errors.New("failed to get active node")
	}
	key := activeNode.PublicKey()
	return packet.Verify(tp.Cryptography, key)
}
