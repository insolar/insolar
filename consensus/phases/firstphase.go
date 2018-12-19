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
	"github.com/insolar/insolar/instrumentation/inslogger"
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
	UnsyncList   network.UnsyncList
}

// Execute do first phase
func (fp *FirstPhase) Execute(ctx context.Context, pulse *core.Pulse) (*FirstPhaseState, error) {
	entry := &merkle.PulseEntry{Pulse: pulse}

	if fp.Calculator == nil {
		panic("calculator is nil")
	}

	if fp.NodeNetwork == nil {
		panic("NodeNetwork is nil")
	}
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
	var originClaim *packets.NodeAnnounceClaim
	if fp.NodeKeeper.NodesJoinedDuringPreviousPulse() {
		originClaim, err = fp.NodeKeeper.GetOriginAnnounceClaim(fp.UnsyncList)
		if err != nil {
			return nil, errors.Wrap(err, "[ FirstPhase ] Failed to get origin claim")
		}
		success = packet.AddClaim(originClaim)
		if !success {
			return nil, errors.Wrap(err, "[ FirstPhase ] Failed to add origin claim in Phase1Packet.")
		}
	}
	for {
		claim := fp.NodeKeeper.GetClaimQueue().Front()
		if claim == nil {
			break
		}
		success = packet.AddClaim(claim)
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
	resultPackets, err := fp.Communicator.ExchangePhase1(ctx, originClaim, activeNodes, &packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ FirstPhase ] Failed to exchange results.")
	}

	proofSet := make(map[core.RecordRef]*merkle.PulseProof)
	rawProofs := make(map[core.RecordRef]*packets.NodePulseProof)
	claimMap := make(map[core.RecordRef][]packets.ReferendumClaim)
	for ref, packet := range resultPackets {
		signIsCorrect, err := fp.isSignPhase1PacketRight(packet, ref)
		if err != nil {
			inslogger.FromContext(ctx).Warnf("Failed to check phase1 packet signature from %s: %s", ref, err.Error())
			continue
		}
		if !signIsCorrect {
			inslogger.FromContext(ctx).Warnf("Received phase1 packet from %s with bad signature", ref)
			continue
		}
		rawProof := packet.GetPulseProof()
		rawProofs[ref] = rawProof
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

	valid, fault := validateProofs(fp.Calculator, fp.UnsyncList, pulseHash, proofSet)
	for node := range valid {
		fp.UnsyncList.AddProof(node.ID(), rawProofs[node.ID()])
	}
	for nodeID := range fault {
		inslogger.FromContext(ctx).Warnf("Failed to validate proof from %s", nodeID)
		delete(claimMap, nodeID)
	}
	err = fp.UnsyncList.AddClaims(claimMap)
	if err != nil {
		return nil, errors.Wrapf(err, "[ FirstPhase ] Failed to add claims")
	}

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
	activeNode := fp.NodeNetwork.GetActiveNode(recordRef)
	if activeNode == nil {
		return false, errors.New("failed to get active node")
	}

	key := activeNode.PublicKey()
	raw, err := packet.RawBytes()

	if err != nil {
		return false, errors.Wrap(err, "failed to serialize packet")
	}
	return fp.Cryptography.Verify(key, core.SignatureFromBytes(packet.Signature[:]), raw), nil
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
		return errors.Wrap(err, "[ checkClaimSignature ] Failed to import a key")
	}
	rawClaim, err := claim.SerializeRaw()
	if err != nil {
		return errors.Wrap(err, "[ checkClaimSignature ] Failed to serialize a claim")
	}
	success := fp.Cryptography.Verify(key, core.SignatureFromBytes(claim.GetSignature()), rawClaim)
	if !success {
		return errors.New("[ checkClaimSignature ] Signature verification failed")
	}
	return nil
}
