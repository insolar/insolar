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
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/merkle"
	"github.com/pkg/errors"
)

const BFTPercent = 2.0 / 3.0
const MajorityPercent = 0.5

func consensusReachedBFT(resultLen, participanstLen int) bool {
	return consensusReachedWithPercent(resultLen, participanstLen, BFTPercent)
}

func consensusReachedMajority(resultLen, participanstLen int) bool {
	return consensusReachedWithPercent(resultLen, participanstLen, MajorityPercent)
}

func consensusReachedWithPercent(resultLen, participanstLen int, percent float64) bool {
	minParticipants := int(math.Floor(percent*float64(participanstLen))) + 1
	return resultLen >= minParticipants
}

// FirstPhase is a first phase.
type FirstPhase struct {
	NodeNetwork  core.NodeNetwork         `inject:""`
	Calculator   merkle.Calculator        `inject:""`
	Communicator Communicator             `inject:""`
	Cryptography core.CryptographyService `inject:""`
	NodeKeeper   network.NodeKeeper       `inject:""`
	State        *FirstPhaseState
	UnsyncList   network.UnsyncList
}

// Execute do first phase
func (fp *FirstPhase) Execute(ctx context.Context, pulse *core.Pulse) (*FirstPhaseState, error) {
	entry := &merkle.PulseEntry{Pulse: pulse}
	pulseHash, pulseProof, err := fp.Calculator.GetPulseProof(entry)
	if fp.NodeKeeper.GetState() == network.Ready {
		fp.UnsyncList = fp.NodeKeeper.GetUnsyncList()
	}

	if err != nil {
		return nil, errors.Wrap(err, "[ FirstPhase ] Failed to calculate pulse proof.")
	}

	packet := packets.Phase1Packet{}
	err = packet.SetPulseProof(pulseProof.StateHash, pulseProof.Signature.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "[ FirstPhase ] Failed to set pulse proof in Phase1Packet.")
	}

	var success bool
	if fp.NodeKeeper.NodesJoinedDuringPreviousPulse() {
		originClaim, err := fp.NodeKeeper.GetOriginAnnounceClaim()
		if err != nil {
			return nil, errors.Wrap(err, "[ FirstPhase ] Failed to get origin claim")
		}
		success = packet.AddClaim(originClaim)
		if !success {
			return nil, errors.Wrap(err, "[ FirstPhase ] Failed to add origin claim in Phase1Packet.")
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
		return nil, errors.Wrap(err, "[ FirstPhase ] failed to sign a packet")
	}
	resultPackets, err := fp.Communicator.ExchangePhase1(ctx, activeNodes, &packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ FirstPhase ] Failed to exchange results.")
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
		claimMap[ref] = fp.filterClaims(ref, packet.GetClaims())
	}

	if fp.NodeKeeper.GetState() == network.Waiting {
		length, err := detectSparseBitsetLength(claimMap)
		if err != nil {
			return nil, errors.Wrapf(err, "[ FirstPhase ] Failed to detect bitset length")
		}
		fp.UnsyncList = fp.NodeKeeper.GetSparseUnsyncList(length)
	}

	fp.UnsyncList.AddClaims(claimMap)

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

func (fp *FirstPhase) signPhase1Packet(packet *packets.Phase1Packet) error {
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

func (fp *FirstPhase) isSignPhase1PacketRight(packet *packets.Phase1Packet, recordRef core.RecordRef) (bool, error) {
	key := fp.NodeNetwork.GetActiveNode(recordRef).PublicKey()
	raw, err := packet.RawBytes()

	if err != nil {
		return false, errors.Wrap(err, "failed to serialize packet")
	}
	return fp.Cryptography.Verify(key, core.SignatureFromBytes(raw), raw), nil
}

func (fp *FirstPhase) validateProofs(
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

func (fp *FirstPhase) validateProof(pulseHash merkle.OriginHash, nodeID core.RecordRef, proof *merkle.PulseProof) bool {
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

func (fp *FirstPhase) filterClaims(nodeID core.RecordRef, claims []packets.ReferendumClaim) []packets.ReferendumClaim {
	result := make([]packets.ReferendumClaim, 0)
	for _, claim := range claims {
		signedClaim, ok := claim.(packets.SignedClaim)
		if ok {
			err := fp.checkClaimSignature(signedClaim)
			if err != nil {
				log.Error("[ filterClaims ] failed to check a claim sign")
				continue
			}
		}
		supClaim, ok := claim.(packets.ClaimSupplementary)
		if ok {
			supClaim.AddSupplementaryInfo(nodeID)
		}
		result = append(result, claim)
	}
	return result
}

func (fp *FirstPhase) checkClaimSignature(claim packets.SignedClaim) error {
	key, err := claim.GetPublicKey()
	if err != nil {
		return errors.Wrap(err, "[ checkClaimSignature ] failed to import a key")
	}
	rawClaim, err := claim.SerializeRaw()
	if err != nil {
		return errors.Wrap(err, "[ checkClaimSignature ] failed to serialize a claim")
	}
	success := fp.Cryptography.Verify(key, core.SignatureFromBytes(claim.GetSignature()), rawClaim)
	if !success {
		return errors.New("Signature verification failed")
	}
	return nil
}
