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

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/merkle"
	"github.com/pkg/errors"
)

type SecondPhase interface {
	Execute(ctx context.Context, state *FirstPhaseState) (*SecondPhaseState, error)
}

func NewSecondPhase() SecondPhase {
	return &secondPhase{}
}

type secondPhase struct {
	NodeKeeper   network.NodeKeeper       `inject:""`
	Calculator   merkle.Calculator        `inject:""`
	Communicator Communicator             `inject:""`
	Cryptography core.CryptographyService `inject:""`
}

func (sp *secondPhase) Execute(ctx context.Context, state *FirstPhaseState) (*SecondPhaseState, error) {
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
		return nil, errors.Wrap(err, "[ Execute ] Failed to calculate pulse proof.")
	}

	packet := packets.Phase2Packet{}
	err = packet.SetGlobuleHashSignature(globuleProof.Signature.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to set pulse proof in Phase2Packet.")
	}
	bitset, err := generatePhase2Bitset(state.UnsyncList, state.ValidProofs)
	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to generate bitset for Phase2Packet")
	}
	packet.SetBitSet(bitset)
	err = sp.signPhase2Packet(&packet)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign a packet")
	}
	activeNodes := state.UnsyncList.GetActiveNodes()
	packets, err := sp.Communicator.ExchangePhase2(ctx, activeNodes, &packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to exchange results.")
	}

	nodeProofs := make(map[core.Node]*merkle.GlobuleProof)

	for ref, packet := range packets {
		signIsCorrect, err := sp.isSignPhase2PacketRight(packet, ref)
		if err != nil {
			log.Warn("failed to check a sign: ", err.Error())
		} else if !signIsCorrect {
			log.Warn("recieved a bad sign packet: ", err.Error())
		}
		node := state.UnsyncList.GetActiveNode(ref)
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

func generatePhase2Bitset(list network.UnsyncList, proofs map[core.Node]*merkle.PulseProof) (packets.BitSet, error) {
	bitset, err := packets.NewBitSet(list.Length())
	if err != nil {
		return nil, err
	}
	cells := make([]packets.BitSetCell, 0)
	for node := range proofs {
		cells = append(cells, packets.BitSetCell{NodeID: node.ID(), State: packets.Legit})
	}
	err = bitset.ApplyChanges(cells, list)
	if err != nil {
		return nil, err
	}
	return bitset, nil
}

func (sp *secondPhase) signPhase2Packet(p *packets.Phase2Packet) error {
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

func (sp *secondPhase) isSignPhase2PacketRight(packet *packets.Phase2Packet, recordRef core.RecordRef) (bool, error) {
	key := sp.NodeKeeper.GetActiveNode(recordRef).PublicKey()

	raw, err := packet.RawFirstPart()
	if err != nil {
		return false, errors.Wrap(err, "failed to serialize")
	}

	return sp.Cryptography.Verify(key, core.SignatureFromBytes(raw), raw), nil
}
