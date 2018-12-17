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

	err := tp.signPhase3Packet(&packet)

	if err != nil {
		return nil, errors.Wrap(err, "[ Phase 3 ] Failed to sign phase 3 packet")
	}

	nodes := state.FirstPhaseState.UnsyncList.GetActiveNodes()
	responses, err := tp.Communicator.ExchangePhase3(ctx, nodes, &packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ Phase 3 ] Failed exchange packets on phase 3")
	}

	for ref, packet := range responses {
		signed, err := tp.isSignPhase3PacketRight(packet, ref)
		if err != nil {
			inslogger.FromContext(ctx).Warnf("Failed to check phase3 packet signature from %s: %s", ref, err.Error())
			continue
		} else if !signed {
			inslogger.FromContext(ctx).Warnf("Received phase3 packet from %s with bad signature", ref)
			continue
		}
		// not needed until we implement fraud detection
		// cells, err := packet.GetBitset().GetCells(state.UnsyncList)

		state.UnsyncList.GlobuleHashSignatures()[ref] = packet.GetGlobuleHashSignature()
	}

	totalCount := state.UnsyncList.Length()
	validCount := 0
	prevCloudHash := tp.NodeKeeper.GetCloudHash()
	for ref, ghs := range state.UnsyncList.GlobuleHashSignatures() {
		node := state.UnsyncList.GetActiveNode(ref)
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
			validCount++
		}
	}

	if !consensusReachedBFT(validCount, totalCount) {
		return nil, errors.New(fmt.Sprintf("[ Phase 3 ] Failed to pass BFT consensus: %d/%d", validCount, totalCount))
	}

	// cloudEntry := &merkle.CloudEntry{
	//
	// }
	// cloudHash, _, _ := tp.Calculator.GetCloudProof(cloudEntry)

	return &ThirdPhaseState{}, nil
}

func (tp *ThirdPhase) signPhase3Packet(p *packets.Phase3Packet) error {
	data, err := p.RawBytes()
	if err != nil {
		return errors.Wrap(err, "failed to get raw bytes")
	}
	sign, err := tp.Cryptography.Sign(data)
	if err != nil {
		return errors.Wrap(err, "failed to sign a phase 2 packet")
	}

	copy(p.SignatureHeaderSection1[:], sign.Bytes())
	return nil
}

func (tp *ThirdPhase) isSignPhase3PacketRight(packet *packets.Phase3Packet, recordRef core.RecordRef) (bool, error) {
	key := tp.NodeNetwork.GetActiveNode(recordRef).PublicKey()

	raw, err := packet.RawBytes()
	if err != nil {
		return false, errors.Wrap(err, "failed to serialize")
	}

	return tp.Cryptography.Verify(key, core.SignatureFromBytes(raw), raw), nil
}
