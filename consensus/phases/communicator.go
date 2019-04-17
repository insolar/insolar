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
	"sync/atomic"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/consensus"
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
)

//go:generate minimock -i github.com/insolar/insolar/consensus/phases.Communicator -o . -s _mock.go

// Communicator interface provides methods to exchange data between nodes
type Communicator interface {
	// ExchangePhase1 used in first consensus step to exchange data between participants
	ExchangePhase1(
		ctx context.Context,
		originClaim *packets.NodeAnnounceClaim,
		participants []insolar.NetworkNode,
		packet *packets.Phase1Packet,
	) (map[insolar.Reference]*packets.Phase1Packet, error)
	// ExchangePhase2 used in second consensus step to exchange data between participants
	ExchangePhase2(ctx context.Context, state *ConsensusState,
		participants []insolar.NetworkNode, packet *packets.Phase2Packet) (map[insolar.Reference]*packets.Phase2Packet, error)
	// ExchangePhase21 is used between phases 2 and 3 of consensus to send additional MissingNode requests
	ExchangePhase21(ctx context.Context, state *ConsensusState,
		packet *packets.Phase2Packet, additionalRequests []*AdditionalRequest) ([]packets.ReferendumVote, error)
	// ExchangePhase3 used in third consensus step to exchange data between participants
	ExchangePhase3(ctx context.Context,
		participants []insolar.NetworkNode, packet *packets.Phase3Packet) (map[insolar.Reference]*packets.Phase3Packet, error)

	component.Initer
}

type phase1Result struct {
	id     insolar.Reference
	packet *packets.Phase1Packet
}

type phase2Result struct {
	id     insolar.Reference
	packet *packets.Phase2Packet
}

type phase3Result struct {
	id     insolar.Reference
	packet *packets.Phase3Packet
}

// ConsensusCommunicator is simple Communicator implementation which communicates with each participants
type ConsensusCommunicator struct {
	ConsensusNetwork network.ConsensusNetwork    `inject:""`
	PulseHandler     network.PulseHandler        `inject:""`
	Cryptography     insolar.CryptographyService `inject:""`
	NodeKeeper       network.NodeKeeper          `inject:""`

	phase1result chan phase1Result
	phase2result chan phase2Result
	phase3result chan phase3Result

	currentPulseNumber uint32
}

// NewCommunicator constructor creates new ConsensusCommunicator
func NewCommunicator() *ConsensusCommunicator {
	return &ConsensusCommunicator{}
}

// Init method implements Initer interface
func (nc *ConsensusCommunicator) Init(ctx context.Context) error {
	nc.phase1result = make(chan phase1Result)
	nc.phase2result = make(chan phase2Result)
	nc.phase3result = make(chan phase3Result)
	nc.ConsensusNetwork.RegisterPacketHandler(packets.Phase1, nc.phase1DataHandler)
	nc.ConsensusNetwork.RegisterPacketHandler(packets.Phase2, nc.phase2DataHandler)
	nc.ConsensusNetwork.RegisterPacketHandler(packets.Phase3, nc.phase3DataHandler)
	return nil
}

func (nc *ConsensusCommunicator) getPulseNumber() insolar.PulseNumber {
	pulseNumber := atomic.LoadUint32(&nc.currentPulseNumber)
	return insolar.PulseNumber(pulseNumber)
}

func (nc *ConsensusCommunicator) setPulseNumber(new insolar.PulseNumber) bool {
	old := nc.getPulseNumber()
	return old < new && atomic.CompareAndSwapUint32(&nc.currentPulseNumber, uint32(old), uint32(new))
}

func (nc *ConsensusCommunicator) sendRequestToNodes(ctx context.Context, participants []insolar.NetworkNode, packet packets.ConsensusPacket) {
	for _, node := range participants {
		if node.ID().Equal(nc.NodeKeeper.GetOrigin().ID()) {
			continue
		}

		go func(ctx context.Context, n insolar.NetworkNode, packet packets.ConsensusPacket) {
			logger := inslogger.FromContext(ctx)
			logger.Debugf("Send %s request to %s", packet.GetType(), n.ID())
			err := nc.ConsensusNetwork.SignAndSendPacket(packet, n.ID(), nc.Cryptography)
			if err != nil {
				logger.Errorf("Failed to send %s request: %s", packet.GetType(), err.Error())
				return
			}
			err = stats.RecordWithTags(context.Background(), []tag.Mutator{tag.Upsert(consensus.TagPhase, packet.GetType().String())}, consensus.PacketsSent.M(1))
			if err != nil {
				logger.Warn("Failed to record metric of sent phase1 requests")
			}
		}(ctx, node, packet.Clone())
	}
}

func (nc *ConsensusCommunicator) sendRequestToNodesWithOrigin(ctx context.Context, originClaim *packets.NodeAnnounceClaim,
	participants []insolar.NetworkNode, packet *packets.Phase1Packet) error {

	requests := make(map[insolar.Reference]packets.ConsensusPacket)
	for _, participant := range participants {
		if participant.ID().Equal(nc.NodeKeeper.GetOrigin().ID()) {
			continue
		}

		err := originClaim.Update(participant.ID(), nc.Cryptography)
		if err != nil {
			return errors.Wrap(err, "Failed to update claims before sending in phase1")
		}
		requests[participant.ID()] = packet.Clone()
	}

	for ref, req := range requests {
		go func(ctx context.Context, node insolar.Reference, consensusPacket packets.ConsensusPacket) {
			logger := inslogger.FromContext(ctx)
			logger.Debug("Send phase1 request with origin to %s", node)
			err := nc.ConsensusNetwork.SignAndSendPacket(consensusPacket, node, nc.Cryptography)
			if err != nil {
				logger.Error("Failed to send phase1 request with origin: " + err.Error())
				return
			}
			err = stats.RecordWithTags(context.Background(), []tag.Mutator{tag.Upsert(consensus.TagPhase, consensusPacket.GetType().String())}, consensus.PacketsSent.M(1))
			if err != nil {
				logger.Warn("Failed to record metric of sent phase1 requests")
			}
		}(ctx, ref, req)
	}
	return nil
}

func (nc *ConsensusCommunicator) generatePhase2Response(ctx context.Context, origReq, req *packets.Phase2Packet,
	state *ConsensusState) (*packets.Phase2Packet, error) {

	logger := inslogger.FromContext(ctx)
	answers := make([]packets.ReferendumVote, 0)
	for _, vote := range req.GetVotes() {
		if vote.Type() != packets.TypeMissingNode {
			continue
		}
		v, ok := vote.(*packets.MissingNode)
		if !ok {
			logger.Warnf("Phase 2 MissingNode request type mismatch")
			continue
		}
		ref, err := state.BitsetMapper.IndexToRef(int(v.NodeIndex))
		if err != nil {
			logger.Warnf("Phase 2 MissingNode requested index: %d, error: %s", v.NodeIndex, err.Error())
			continue
		}
		node := state.NodesMutator.GetActiveNode(ref)
		if node == nil {
			logger.Warnf("Phase 2 MissingNode requested index: %d; mapped ref %s not found", v.NodeIndex, ref)
			continue
		}
		claim, err := packets.NodeToClaim(node)
		if err != nil {
			logger.Warnf("Phase 2 MissingNode requested index: %d, mapped ref: %s, convertation node -> claim error: %s",
				v.NodeIndex, ref, err.Error())
			continue
		}
		proof := state.HashStorage.GetProof(ref)
		if proof == nil {
			logger.Warnf("Phase 2 MissingNode requested index: %d, mapped ref: %s, proof not found", v.NodeIndex, ref)
			continue
		}
		answer := packets.MissingNodeSupplementaryVote{
			NodeIndex:         v.NodeIndex,
			NodePulseProof:    *proof,
			NodeClaimUnsigned: *claim,
		}
		answers = append(answers, &answer)
		claims := state.ClaimHandler.GetClaimsFromNode(ref)
		for _, claim := range claims {
			claimAnswer := packets.MissingNodeClaim{NodeIndex: v.NodeIndex, Claim: claim}
			answers = append(answers, &claimAnswer)
		}
	}
	response := packets.NewPhase2Packet(origReq.GetPulseNumber())
	response.SetBitSet(origReq.GetBitSet())
	ghs := origReq.GetGlobuleHashSignature()
	err := response.SetGlobuleHashSignature(ghs[:])
	if err != nil {
		return nil, errors.Wrap(err, "Failed to set globule hash in phase2 response")
	}
	for _, answer := range answers {
		response.AddVote(answer)
	}
	return response, nil
}

// ExchangePhase1 used in first consensus phase to exchange data between participants
func (nc *ConsensusCommunicator) ExchangePhase1(
	ctx context.Context,
	originClaim *packets.NodeAnnounceClaim,
	participants []insolar.NetworkNode,
	packet *packets.Phase1Packet,
) (map[insolar.Reference]*packets.Phase1Packet, error) {
	_, span := instracer.StartSpan(ctx, "Communicator.ExchangePhase1")
	span.AddAttributes(trace.Int64Attribute("pulse", int64(packet.GetPulseNumber())))
	defer span.End()
	logger := inslogger.FromContext(ctx)

	result := make(map[insolar.Reference]*packets.Phase1Packet, len(participants))
	result[nc.ConsensusNetwork.GetNodeID()] = packet
	nc.setPulseNumber(packet.GetPulse().PulseNumber)

	var request *packets.Phase1Packet

	type none struct{}
	sentRequests := make(map[insolar.Reference]none)

	// TODO: awful, need rework
	if originClaim == nil {
		request = packet
	} else {
		request = packet.Clone().(*packets.Phase1Packet)
		request.RemoveAnnounceClaim()
	}

	if originClaim == nil {
		nc.sendRequestToNodes(ctx, participants, request)
	} else {
		err := nc.sendRequestToNodesWithOrigin(ctx, originClaim, participants, packet)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to send requests")
		}
	}
	for _, p := range participants {
		sentRequests[p.ID()] = none{}
	}

	shouldSendResponse := func(ref insolar.Reference) bool {
		_, ok := sentRequests[ref]
		return !ok
	}
	response := request

	for {
		select {
		case res := <-nc.phase1result:
			logger.Debugf("Got phase1 request from %s", res.id)
			currentPulse := nc.getPulseNumber()
			if res.packet.GetPulseNumber() != currentPulse {
				logger.Debugf("Filtered phase1 packet, packet pulse %d != %d (current pulse)",
					res.packet.GetPulseNumber(), currentPulse)
				continue
			}

			if res.id.IsEmpty() {
				logger.Debug("Got unknown phase1 request, try to get routing info from announce claim")
				claim := res.packet.GetAnnounceClaim()
				if claim == nil {
					logger.Warn("Could not get announce claim from phase1 packet")
					continue
				}
				res.id = claim.NodeRef
				err := nc.NodeKeeper.GetConsensusInfo().AddTemporaryMapping(claim.NodeRef, claim.ShortNodeID, claim.NodeAddress.Get())
				if err != nil {
					logger.Warn("Error adding temporary mapping: " + err.Error())
					continue
				}
			}

			if shouldSendResponse(res.id) {
				// send response
				logger.Debugf("Send phase1 response to %s", res.id)
				err := nc.ConsensusNetwork.SignAndSendPacket(response, res.id, nc.Cryptography)
				if err != nil {
					logger.Error("Error sending phase1 response: " + err.Error())
				}
			}
			if !res.id.IsEmpty() {
				sentRequests[res.id] = none{}
				result[res.id] = res.packet
			}

			// FIXME: early return is commented to have synchronized length of phases on all nodes
			// if len(result) == len(participants) {
			// 	return result, nil
			// }
		case <-ctx.Done():
			return result, nil
		}
	}
}

// ExchangePhase2 used in second consensus phase to exchange data between participants
func (nc *ConsensusCommunicator) ExchangePhase2(ctx context.Context, state *ConsensusState,
	participants []insolar.NetworkNode, packet *packets.Phase2Packet) (map[insolar.Reference]*packets.Phase2Packet, error) {
	_, span := instracer.StartSpan(ctx, "Communicator.ExchangePhase2")
	span.AddAttributes(trace.Int64Attribute("pulse", int64(packet.GetPulseNumber())))
	defer span.End()
	logger := inslogger.FromContext(ctx)

	result := make(map[insolar.Reference]*packets.Phase2Packet, len(participants))

	result[nc.ConsensusNetwork.GetNodeID()] = packet

	nc.sendRequestToNodes(ctx, participants, packet)

	type none struct{}
	sentRequests := make(map[insolar.Reference]none)

	shouldSendResponse := func(p *phase2Result) bool {
		_, ok := sentRequests[p.id]
		return !ok || p.packet.ContainsRequests()
	}

	var err error
	for {
		select {
		case res := <-nc.phase2result:
			logger.Debugf("Got phase2 request from %s", res.id)
			currentPulse := nc.getPulseNumber()
			if res.packet.GetPulseNumber() != currentPulse {
				logger.Debugf("Filtered phase2 packet, packet pulse %d != %d (current pulse)",
					res.packet.GetPulseNumber(), currentPulse)
				continue
			}

			if shouldSendResponse(&res) {
				logger.Debugf("Send phase2 response to %s", res.id)
				// send response
				response := packet
				if res.packet.ContainsRequests() {
					response, err = nc.generatePhase2Response(ctx, packet, res.packet, state)
					if err != nil {
						logger.Warnf("Failed to generate phase 2 response packet: %s", err.Error())
						continue
					}
				}
				err := nc.ConsensusNetwork.SignAndSendPacket(response, res.id, nc.Cryptography)
				if err != nil {
					logger.Error("Error sending phase2 response: " + err.Error())
				}
			}
			result[res.id] = res.packet
			sentRequests[res.id] = none{}

			// FIXME: early return is commented to have synchronized length of phases on all nodes
			// if len(result) == len(participants) {
			// 	return result, nil
			// }
		case <-ctx.Done():
			return result, nil
		}
	}
}

func selectCandidate(candidates []insolar.Reference) insolar.Reference {
	// TODO: make it random
	if len(candidates) == 0 {
		panic("candidates list should have at least 1 candidate (check consensus state matrix calculation routines)")
	}
	return candidates[0]
}

func (nc *ConsensusCommunicator) sendAdditionalRequests(ctx context.Context, origReq *packets.Phase2Packet,
	additionalRequests []*AdditionalRequest) error {

	logger := inslogger.FromContext(ctx)
	for _, req := range additionalRequests {
		newReq := *origReq
		newReq.AddVote(&packets.MissingNode{NodeIndex: uint16(req.RequestIndex)})
		receiver := selectCandidate(req.Candidates)
		err := nc.ConsensusNetwork.SignAndSendPacket(&newReq, receiver, nc.Cryptography)
		if err != nil {
			return errors.Wrapf(err, "Failed to send additional phase 2.1 request for index %d to node %s", req.RequestIndex, receiver)
		}
		err = stats.RecordWithTags(context.Background(), []tag.Mutator{tag.Upsert(consensus.TagPhase, origReq.GetType().String())}, consensus.PacketsSent.M(1))
		if err != nil {
			logger.Warn("Failed to record metric of sent phase2.1 additional requests")
		}
	}

	return nil
}

// ExchangePhase21 used in second consensus phase to exchange data between participants
func (nc *ConsensusCommunicator) ExchangePhase21(ctx context.Context, state *ConsensusState,
	packet *packets.Phase2Packet, additionalRequests []*AdditionalRequest) ([]packets.ReferendumVote, error) {
	_, span := instracer.StartSpan(ctx, "Communicator.ExchangePhase21")
	span.AddAttributes(trace.Int64Attribute("pulse", int64(packet.GetPulseNumber())))
	defer span.End()
	logger := inslogger.FromContext(ctx)

	type none struct{}
	incoming := make(map[insolar.Reference]none)
	responsesFilter := make(map[int]none)
	for _, req := range additionalRequests {
		responsesFilter[req.RequestIndex] = none{}
	}

	shouldSendResponse := func(p *phase2Result) bool {
		_, ok := incoming[p.id]
		return !ok || p.packet.ContainsRequests()
	}

	result := make([]packets.ReferendumVote, 0)

	appendResult := func(index uint16, vote packets.ReferendumVote) {
		_, ok := responsesFilter[int(index)]
		if !ok {
			return
		}
		result = append(result, vote)
	}

	err := nc.sendAdditionalRequests(ctx, packet, additionalRequests)
	if err != nil {
		return nil, errors.Wrap(err, "[ ExchangePhase2.1 ] Failed to send additional phase 2.1 requests")
	}

	for {
		select {
		case res := <-nc.phase2result:
			logger.Debugf("Got phase2 request from %s", res.id)
			currentPulse := nc.getPulseNumber()
			if res.packet.GetPulseNumber() != currentPulse {
				logger.Debugf("Filtered phase2 packet, packet pulse %d != %d (current pulse)",
					res.packet.GetPulseNumber(), currentPulse)
				continue
			}

			if shouldSendResponse(&res) {
				logger.Debugf("Send phase2 response to %s", res.id)
				// send response
				response := packet
				if res.packet.ContainsRequests() {
					response, err = nc.generatePhase2Response(ctx, packet, res.packet, state)
					if err != nil {
						logger.Warnf("Failed to generate phase 2 response packet: %s", err.Error())
						continue
					}
				}
				err := nc.ConsensusNetwork.SignAndSendPacket(response, res.id, nc.Cryptography)
				if err != nil {
					logger.Error("Error sending phase2 response: " + err.Error())
				}
			}

			if res.packet.ContainsResponses() {
				voteAnswers := res.packet.GetVotes()
				for _, vote := range voteAnswers {
					switch v := vote.(type) {
					case *packets.MissingNodeSupplementaryVote:
						appendResult(v.NodeIndex, v)
					case *packets.MissingNodeClaim:
						appendResult(v.NodeIndex, v)
					}
				}
			}

			incoming[res.id] = none{}

			// FIXME: early return is commented to have synchronized length of phases on all nodes
			// if len(result) == len(participants) {
			// 	return result, nil
			// }
		case <-ctx.Done():
			return result, nil
		}
	}
}

// ExchangePhase3 used in third consensus step to exchange data between participants
func (nc *ConsensusCommunicator) ExchangePhase3(ctx context.Context, participants []insolar.NetworkNode, packet *packets.Phase3Packet) (map[insolar.Reference]*packets.Phase3Packet, error) {
	result := make(map[insolar.Reference]*packets.Phase3Packet, len(participants))
	_, span := instracer.StartSpan(ctx, "Communicator.ExchangePhase3")
	span.AddAttributes(trace.Int64Attribute("pulse", int64(packet.GetPulseNumber())))
	defer span.End()
	logger := inslogger.FromContext(ctx)

	result[nc.ConsensusNetwork.GetNodeID()] = packet

	nc.sendRequestToNodes(ctx, participants, packet)

	type none struct{}
	sentRequests := make(map[insolar.Reference]none)

	shouldSendResponse := func(p *phase3Result) bool {
		_, ok := sentRequests[p.id]
		return !ok
	}

	for {
		select {
		case res := <-nc.phase3result:
			logger.Debugf("Got phase3 request from %s", res.id)
			currentPulse := nc.getPulseNumber()
			if res.packet.GetPulseNumber() != currentPulse {
				logger.Debugf("Filtered phase3 packet, packet pulse %d != %d (current pulse)",
					res.packet.GetPulseNumber(), currentPulse)
				continue
			}

			if shouldSendResponse(&res) {
				logger.Debugf("Send phase3 response to %s", res.id)
				// send response
				err := nc.ConsensusNetwork.SignAndSendPacket(packet, res.id, nc.Cryptography)
				if err != nil {
					logger.Error("Error sending phase3 response: " + err.Error())
				}
			}
			result[res.id] = res.packet
			sentRequests[res.id] = none{}

			// FIXME: early return is commented to have synchronized length of phases on all nodes
			// if len(result) == len(participants) {
			// 	return result, nil
			// }
		case <-ctx.Done():
			return result, nil
		}
	}
}

func (nc *ConsensusCommunicator) phase1DataHandler(packet packets.ConsensusPacket, sender insolar.Reference) {
	p, ok := packet.(*packets.Phase1Packet)
	if !ok {
		log.Error("invalid Phase1Packet")
		return
	}

	newPulse := p.GetPulse()

	if newPulse.PulseNumber < nc.getPulseNumber() {
		log.Warn("ignore old pulse Phase1Packet")
		return
	}

	if nc.setPulseNumber(newPulse.PulseNumber) {
		go nc.PulseHandler.HandlePulse(context.Background(), newPulse)
	}

	nc.phase1result <- phase1Result{id: sender, packet: p}
}

func (nc *ConsensusCommunicator) phase2DataHandler(packet packets.ConsensusPacket, sender insolar.Reference) {
	p, ok := packet.(*packets.Phase2Packet)
	if !ok {
		log.Error("invalid Phase2Packet")
		return
	}

	pulseNumber := p.GetPulseNumber()

	if pulseNumber < nc.getPulseNumber() {
		log.Warn("ignore old pulse Phase2Packet")
		return
	}

	nc.phase2result <- phase2Result{id: sender, packet: p}
}

func (nc *ConsensusCommunicator) phase3DataHandler(packet packets.ConsensusPacket, sender insolar.Reference) {
	p, ok := packet.(*packets.Phase3Packet)
	if !ok {
		log.Warn("failed to cast a type 3 packet to phase3packet")
	}
	nc.phase3result <- phase3Result{id: sender, packet: p}
}
