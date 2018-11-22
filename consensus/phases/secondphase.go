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

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/merkle"
	"github.com/pkg/errors"
)

// SecondPhase is a second phase.
type SecondPhase struct {
	NodeNetwork  core.NodeNetwork         `inject:""`
	NodeKeeper   network.NodeKeeper       `inject:""`
	Network      core.Network             `inject:""`
	Calculator   merkle.Calculator        `inject:""`
	Communicator Communicator             `inject:""`
	Cryptography core.CryptographyService `inject:""`
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
	err = sp.signPhase2Packet(&packet)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign a packet")
	}
	proofSet, err := sp.Communicator.ExchangePhase2(ctx, activeNodes, packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to exchange results.")
	}

	nodeProofs := make(map[core.Node]*merkle.GlobuleProof)

	var deviants []core.Node
	deviants = append(deviants, state.TimedOutNodes...)
	deviants = append(deviants, state.DeviantNodes...)

	for ref, packet := range proofSet {
		signIsCorrect, err := sp.isSignPhase2PacketRight(packet, ref)
		if err != nil {
			log.Warn("failed to check a sign: ", err.Error())
		} else if !signIsCorrect {
			log.Warn("recieved a bad sign packet: ", err.Error())
		}
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

	// TODO: timeouts, deviants, etc.
	sp.NodeKeeper.Sync(state.UnsyncList)

	return &SecondPhaseState{
		FirstPhaseState: state,

		GlobuleEntry:    entry,
		GlobuleHash:     globuleHash,
		GlobuleProof:    globuleProof,
		GlobuleProofSet: nodeProofs,
	}, nil
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
	key := sp.NodeNetwork.GetActiveNode(recordRef).PublicKey()

	raw, err := packet.RawFirstPart()
	if err != nil {
		return false, errors.Wrap(err, "failed to serialize")
	}

	return sp.Cryptography.Verify(key, core.SignatureFromBytes(raw), raw), nil
}

func (sp *SecondPhase) processTimedOutNodes(timedOutNodes []core.Node) {
	// TODO: process
}

func (sp *SecondPhase) calculateListForNextPulse() (uint16, []byte) {
	// TODO: calculate
	return 1337, []byte("1337")
}
