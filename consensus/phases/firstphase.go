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

const ConsensusAtPercents = 2.0 / 3.0

func consensusReached(resultLen, participanstLen int) bool {
	minParticipants := int(math.Floor(ConsensusAtPercents*float64(participanstLen))) + 1

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
		return nil, errors.Wrap(err, "[ Execute ] Failed to calculate pulse proof.")
	}

	packet := packets.Phase1Packet{}
	err = packet.SetPulseProof(pulseProof.StateHash, pulseProof.Signature.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to set pulse proof in Phase1Packet.")
	}

	var success bool
	if fp.NodeKeeper.NodesJoinedDuringPreviousPulse() {
		success = packet.AddClaim(fp.NodeKeeper.GetOriginClaim())
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
		claimMap[ref] = packet.GetClaims()
	}

	if fp.NodeKeeper.GetState() == network.Waiting {
		length, err := detectSparseBitsetLength(claimMap)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to detect bitset length")
		}
		fp.UnsyncList = fp.NodeKeeper.GetSparseUnsyncList(length)
	}

	fp.processClaims(ctx, claimMap, addressMap)

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
		UnsyncList:    fp.UnsyncList,
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

func (fp *FirstPhase) processClaims(ctx context.Context, claims map[core.RecordRef][]packets.ReferendumClaim, addressMap map[core.RecordRef]string) {
	fp.UnsyncList.AddClaims(claims, addressMap)
}
