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
	"math"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/merkle"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

const ConsensusAtPercents = 2.0 / 3.0

func consensusReached(resultLen, participanstLen int) bool {
	minParticipants := int(math.Floor(ConsensusAtPercents*float64(participanstLen))) + 1

	return resultLen >= minParticipants
}

type FirstPhase interface {
	Execute(ctx context.Context, pulse *core.Pulse) (*FirstPhaseState, error)
}

func NewFirstPhase() FirstPhase {
	return &firstPhase{}
}

type firstPhase struct {
	Calculator   merkle.Calculator        `inject:""`
	Communicator Communicator             `inject:""`
	Cryptography core.CryptographyService `inject:""`
	NodeKeeper   network.NodeKeeper       `inject:""`
	State        *FirstPhaseState
	UnsyncList   network.UnsyncList
}

// Execute do first phase
func (fp *firstPhase) Execute(ctx context.Context, pulse *core.Pulse) (*FirstPhaseState, error) {
	entry := &merkle.PulseEntry{Pulse: pulse}
	pulseHash, pulseProof, err := fp.Calculator.GetPulseProof(entry)
	if fp.NodeKeeper.GetState() == network.Ready {
		fp.UnsyncList = fp.NodeKeeper.GetUnsyncList()
	}

	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to calculate pulse proof.")
	}

	packet := packets.Phase1Packet{}
	err = packet.SetPulseProof(pulseProof.StateHash, pulseProof.Signature.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to set pulse proof in Phase1Packet.")
	}

	var success bool
	if fp.NodeKeeper.NodesJoinedDuringPreviousPulse() {
		originClaim, err := fp.NodeKeeper.GetOriginClaim()
		if err != nil {
			return nil, errors.Wrap(err, "[ Execute ] Failed to get origin claim")
		}
		success = packet.AddClaim(originClaim)
		if !success {
			return nil, errors.Wrap(err, "[ Execute ] Failed to add origin claim in Phase1Packet.")
		}
	}
	for {
		success = packet.AddClaim(fp.NodeKeeper.GetClaimQueue().Front())
		if !success {
			break
		}
		_ = fp.NodeKeeper.GetClaimQueue().Pop()
	}

	activeNodes := fp.NodeKeeper.GetActiveNodes()
	err = fp.signPhase1Packet(&packet)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign a packet")
	}
	resultPackets, addressMap, err := fp.Communicator.ExchangePhase1(ctx, activeNodes, &packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to exchange results.")
	}

	proofSet := make(map[core.RecordRef]*merkle.PulseProof)
	claimMap := make(map[core.RecordRef][]packets.ReferendumClaim)
	for ref, packet := range resultPackets {
		signIsCorrect, err := fp.isSignPhase1PacketRight(packet, ref)
		if err != nil {
			log.Warn("failed to check a sign: ", err.Error())
		} else if !signIsCorrect {
			log.Warn("recieved a bad sign packet: ", err.Error())
		}
		rawProof := packet.GetPulseProof()
		proofSet[ref] = &merkle.PulseProof{
			BaseProof: merkle.BaseProof{
				Signature: core.SignatureFromBytes(rawProof.Signature()),
			},
			StateHash: rawProof.StateHash(),
		}
		claimMap[ref] = fp.getSignedClaims(packet.GetClaims())
	}

	if fp.NodeKeeper.GetState() == network.Waiting {
		length, err := detectSparseBitsetLength(claimMap)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to detect bitset length")
		}
		fp.UnsyncList = fp.NodeKeeper.GetSparseUnsyncList(length)
	}

	fp.UnsyncList.AddClaims(claimMap, addressMap)

	valid, fault := fp.validateProofs(pulseHash, proofSet)

	return &FirstPhaseState{
		PulseEntry:  entry,
		PulseHash:   pulseHash,
		PulseProof:  pulseProof,
		ValidProofs: valid,
		FaultProofs: fault,
		UnsyncList:  fp.UnsyncList,
	}, nil
}

func (fp *firstPhase) signPhase1Packet(packet *packets.Phase1Packet) error {
	data, err := packet.RawBytes()
	if err != nil {
		return errors.Wrap(err, "failed to get raw bytes")
	}
	sign, err := fp.Cryptography.Sign(data)
	if err != nil {
		return errors.Wrap(err, "failed to sign a phase 2 packet")
	}
	copy(packet.Signature[:], sign.Bytes())
	return nil
}

func (fp *firstPhase) isSignPhase1PacketRight(packet *packets.Phase1Packet, recordRef core.RecordRef) (bool, error) {
	key := fp.NodeKeeper.GetActiveNode(recordRef).PublicKey()
	raw, err := packet.RawBytes()

	if err != nil {
		return false, errors.Wrap(err, "failed to serialize packet")
	}
	return fp.Cryptography.Verify(key, core.SignatureFromBytes(raw), raw), nil
}

func (fp *firstPhase) validateProofs(
	pulseHash merkle.OriginHash,
	proofs map[core.RecordRef]*merkle.PulseProof,
) (valid map[core.Node]*merkle.PulseProof, fault map[core.RecordRef]*merkle.PulseProof) {

	validProofs := make(map[core.Node]*merkle.PulseProof)
	faultProofs := make(map[core.RecordRef]*merkle.PulseProof)
	for nodeID, proof := range proofs {
		valid := fp.validateProof(pulseHash, nodeID, proof)
		if valid {
			validProofs[fp.UnsyncList.GetActiveNode(nodeID)] = proof
		} else {
			faultProofs[nodeID] = proof
		}
	}
	return validProofs, faultProofs
}

func (fp *firstPhase) validateProof(pulseHash merkle.OriginHash, nodeID core.RecordRef, proof *merkle.PulseProof) bool {
	node := fp.UnsyncList.GetActiveNode(nodeID)
	if node == nil {
		return false
	}
	return fp.Calculator.IsValid(proof, pulseHash, node.PublicKey())
}

func detectSparseBitsetLength(claims map[core.RecordRef][]packets.ReferendumClaim) (int, error) {
	// TODO: NETD18-47
	for _, claimList := range claims {
		for _, claim := range claimList {
			if claim.Type() == packets.TypeNodeAnnounceClaim {
				announceClaim, ok := claim.(*packets.NodeAnnounceClaim)
				if !ok {
					continue
				}
				return int(announceClaim.NodeCount), nil
			}
		}
	}
	return 0, errors.New("no announce claims were received")
}

func (fp *firstPhase) getSignedClaims(claims []packets.ReferendumClaim) []packets.ReferendumClaim {
	result := make([]packets.ReferendumClaim, 0)
	for _, claim := range claims {
		joinClaim, ok := claim.(*packets.NodeJoinClaim)
		if ok {
			signConfirmed, err := fp.claimSignIsOk(joinClaim)
			if err != nil {
				log.Error("[ getSignedClaims ] failed to check a claim sign")
				continue
			}
			if !signConfirmed {
				log.Error("[ getSginedClaims ] sign is unconfirmed")
				continue
			}
		}
		result = append(result, claim)
	}
	return result
}

func (fp *firstPhase) claimSignIsOk(claim *packets.NodeJoinClaim) (bool, error) {
	keyProc := platformpolicy.NewKeyProcessor()
	key, err := keyProc.ImportPublicKeyPEM(claim.NodePK[:])
	if err != nil {
		return false, errors.Wrap(err, "[ claimSignIsOk ] failed to import a key")
	}
	rawClaim, err := claim.SerializeWithoutSign()
	if err != nil {
		return false, errors.Wrap(err, "[ claimSignIsOk ] failed to serialize a claim")
	}
	return fp.Cryptography.Verify(key, core.SignatureFromBytes(claim.Signature[:]), rawClaim), nil
}
