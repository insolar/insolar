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

	"github.com/insolar/insolar/consensus"
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/merkle"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/jbenet/go-base58"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
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

type FirstPhase interface {
	Execute(ctx context.Context, pulse *core.Pulse) (*FirstPhaseState, error)
}

func NewFirstPhase() FirstPhase {
	return &FirstPhaseImpl{}
}

type FirstPhaseImpl struct {
	Calculator   merkle.Calculator        `inject:""`
	Communicator Communicator             `inject:""`
	Cryptography core.CryptographyService `inject:""`
	NodeKeeper   network.NodeKeeper       `inject:""`
}

// Execute do first phase
func (fp *FirstPhaseImpl) Execute(ctx context.Context, pulse *core.Pulse) (*FirstPhaseState, error) {
	entry := &merkle.PulseEntry{Pulse: pulse}
	logger := inslogger.FromContext(ctx)
	ctx, span := instracer.StartSpan(ctx, "FirstPhase.Execute")
	span.AddAttributes(trace.Int64Attribute("pulse", int64(pulse.PulseNumber)))
	defer span.End()

	var unsyncList network.UnsyncList

	pulseHash, pulseProof, err := fp.Calculator.GetPulseProof(entry)
	if fp.NodeKeeper.GetState() == core.ReadyNodeNetworkState {
		unsyncList = fp.NodeKeeper.GetUnsyncList()
	}

	logger.Infof("[ NET Consensus phase-1 ] Calculated pulse proof: %s", base58.Encode(pulseHash))

	if err != nil {
		return nil, errors.Wrap(err, "[ NET Consensus phase-1 ] Failed to calculate pulse proof")
	}

	packet := packets.NewPhase1Packet(*pulse)
	err = packet.SetPulseProof(pulseProof.StateHash, pulseProof.Signature.Bytes())
	if err != nil {
		return nil, errors.Wrapf(err, "[ NET Consensus phase-1 ] Failed to set pulse proof in Phase1Packet")
	}

	var success bool
	var originClaim *packets.NodeAnnounceClaim
	if fp.NodeKeeper.NodesJoinedDuringPreviousPulse() {
		log.Debugf("[ NET Consensus phase-1 ] Add origin announce claim to consensus phase1 packet")
		originClaim, err = fp.NodeKeeper.GetOriginAnnounceClaim(unsyncList)
		if err != nil {
			return nil, errors.Wrap(err, "[ NET Consensus phase-1 ] Failed to get origin claim")
		}
		success = packet.AddClaim(originClaim)
		if !success {
			return nil, errors.Wrap(err, "[ NET Consensus phase-1 ] Failed to add origin claim in Phase1Packet")
		}
		log.Debug("[ NET Consensus phase-1 ] Added origin claim in Phase1Packet")
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
		log.Debugf("[ NET Consensus phase-1 ] Added claim %s to Phase1Packet", claim.Type())
	}
	log.Infof("[ NET Consensus phase-1 ] Phase1Packet claims count: %d", len(packet.GetClaims()))

	activeNodes := fp.NodeKeeper.GetActiveNodes()
	resultPackets, err := fp.Communicator.ExchangePhase1(ctx, originClaim, activeNodes, packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ NET Consensus phase-1 ] Failed to exchange results")
	}
	if len(resultPackets) < 2 && fp.NodeKeeper.GetState() == core.WaitingNodeNetworkState {
		return nil, errors.New("[ NET Consensus phase-1 ] Failed to receive enough packets from other nodes")
	}
	if fp.NodeKeeper.GetState() == core.WaitingNodeNetworkState {
		logger.Infof("[ NET Consensus phase-1 ] received packets: %d", len(resultPackets))
	} else {
		logger.Infof("[ NET Consensus phase-1 ] received packets: %d/%d", len(resultPackets), len(activeNodes))
	}
	err = stats.RecordWithTags(ctx, []tag.Mutator{tag.Upsert(consensus.TagPhase, "phase 1")}, consensus.PacketsRecv.M(int64(len(resultPackets))))
	if err != nil {
		logger.Warn("[ NET Consensus phase-1 ] Failed to record received packets metric: " + err.Error())
	}

	proofSet := make(map[core.RecordRef]*merkle.PulseProof)
	rawProofs := make(map[core.RecordRef]*packets.NodePulseProof)
	claimMap := make(map[core.RecordRef][]packets.ReferendumClaim)
	for ref, packet := range resultPackets {
		err = nil
		if !ref.Equal(fp.NodeKeeper.GetOrigin().ID()) {
			err = fp.checkPacketSignature(packet, ref)
		}
		if err != nil {
			logger.Warnf("[ NET Consensus phase-1 ] Failed to check phase1 packet signature from %s: %s", ref, err.Error())
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

	if fp.NodeKeeper.GetState() == core.WaitingNodeNetworkState {
		length, err := detectSparseBitsetLength(claimMap, fp.NodeKeeper)
		if err != nil {
			return nil, errors.Wrap(err, "[ NET Consensus phase-1 ] Failed to detect bitset length")
		}
		logger.Debugf("[ NET Consensus phase-1 ] Bitset length: %d", length)
		unsyncList = fp.NodeKeeper.GetSparseUnsyncList(length)
	}

	err = unsyncList.AddClaims(claimMap)
	if err != nil {
		return nil, errors.Wrap(err, "[ NET Consensus phase-1 ] Failed to add claims")
	}
	valid, fault := validateProofs(fp.Calculator, unsyncList, pulseHash, proofSet)
	for node := range valid {
		unsyncList.AddProof(node.ID(), rawProofs[node.ID()])
	}
	for nodeID := range fault {
		logger.Warnf("[ NET Consensus phase-1 ] Failed to validate proof from %s", nodeID)
		unsyncList.RemoveNode(nodeID)
	}
	logger.Infof("[ NET Consensus phase-1 ] Valid proofs after phase: %d/%d", len(valid), unsyncList.Length())

	return &FirstPhaseState{
		PulseEntry:  entry,
		PulseHash:   pulseHash,
		PulseProof:  pulseProof,
		ValidProofs: valid,
		FaultProofs: fault,
		UnsyncList:  unsyncList,
	}, nil
}

func (fp *FirstPhaseImpl) checkPacketSignature(packet *packets.Phase1Packet, recordRef core.RecordRef) error {
	if fp.NodeKeeper.GetState() == core.WaitingNodeNetworkState {
		return fp.checkPacketSignatureFromClaim(packet, recordRef)
	}

	activeNode := fp.NodeKeeper.GetActiveNode(recordRef)
	if activeNode == nil {
		return errors.New("failed to get active node")
	}
	key := activeNode.PublicKey()
	return packet.Verify(fp.Cryptography, key)
}

func (fp *FirstPhaseImpl) checkPacketSignatureFromClaim(packet *packets.Phase1Packet, recordRef core.RecordRef) error {
	announceClaim := packet.GetAnnounceClaim()
	if announceClaim == nil {
		return errors.New("could not find announce claim")
	}
	pk, err := platformpolicy.NewKeyProcessor().ImportPublicKeyBinary(announceClaim.NodePK[:])
	if err != nil {
		return errors.Wrap(err, "could not import public key from announce claim")
	}
	return packet.Verify(fp.Cryptography, pk)
}

func detectSparseBitsetLength(claims map[core.RecordRef][]packets.ReferendumClaim, nk network.NodeKeeper) (int, error) {
	// TODO: NETD18-47
	for _, claimList := range claims {
		for _, claim := range claimList {
			if claim.Type() == packets.TypeNodeAnnounceClaim {
				announceClaim, ok := claim.(*packets.NodeAnnounceClaim)
				if !ok {
					continue
				}

				nk.SetCloudHash(announceClaim.CloudHash[:])
				return int(announceClaim.NodeCount), nil
			}
		}
	}
	return 0, errors.New("no announce claims were received")
}

func (fp *FirstPhaseImpl) filterClaims(nodeID core.RecordRef, claims []packets.ReferendumClaim) []packets.ReferendumClaim {
	result := make([]packets.ReferendumClaim, 0)
	for _, claim := range claims {
		signedClaim, ok := claim.(packets.SignedClaim)
		if ok && !nodeID.Equal(fp.NodeKeeper.GetOrigin().ID()) {
			err := fp.checkClaimSignature(signedClaim)
			if err != nil {
				stats.Record(context.Background(), consensus.DeclinedClaims.M(1))
				log.Error("failed to check claim signature: " + err.Error())
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

func (fp *FirstPhaseImpl) checkClaimSignature(claim packets.SignedClaim) error {
	key, err := claim.GetPublicKey()
	if err != nil {
		return errors.Wrap(err, "failed to import a key")
	}
	rawClaim, err := claim.SerializeRaw()
	if err != nil {
		return errors.Wrap(err, "failed to serialize a claim")
	}
	success := fp.Cryptography.Verify(key, core.SignatureFromBytes(claim.GetSignature()), rawClaim)
	if !success {
		return errors.New("signature verification failed")
	}
	return nil
}
