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
	"sync/atomic"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/pkg/errors"
)

// Communicator interface provides methods to exchange data between nodes
//go:generate minimock -i github.com/insolar/insolar/consensus/phases.Communicator -o . -s _mock.go
type Communicator interface {
	// ExchangePhase1 used in first consensus step to exchange data between participants
	ExchangePhase1(
		ctx context.Context,
		originClaim *packets.NodeAnnounceClaim,
		participants []core.Node,
		packet *packets.Phase1Packet,
	) (map[core.RecordRef]*packets.Phase1Packet, error)
	// ExchangePhase2 used in second consensus step to exchange data between participants
	ExchangePhase2(ctx context.Context, list network.UnsyncList, participants []core.Node, packet *packets.Phase2Packet) (map[core.RecordRef]*packets.Phase2Packet, error)
	// ExchangePhase21 is used between phases 2 and 3 of consensus to send additional MissingNode requests
	ExchangePhase21(ctx context.Context, list network.UnsyncList, packet *packets.Phase2Packet, additionalRequests []*AdditionalRequest) ([]packets.ReferendumVote, error)
	// ExchangePhase3 used in third consensus step to exchange data between participants
	ExchangePhase3(ctx context.Context, participants []core.Node, packet *packets.Phase3Packet) (map[core.RecordRef]*packets.Phase3Packet, error)
}

type phase1Result struct {
	id     core.RecordRef
	packet *packets.Phase1Packet
}

type phase2Result struct {
	id     core.RecordRef
	packet *packets.Phase2Packet
}

type phase3Result struct {
	id     core.RecordRef
	packet *packets.Phase3Packet
}

// NaiveCommunicator is simple Communicator implementation which communicates with each participants
type NaiveCommunicator struct {
	ConsensusNetwork network.ConsensusNetwork `inject:""`
	PulseHandler     network.PulseHandler     `inject:""`
	Cryptography     core.CryptographyService `inject:""`
	NodeKeeper       network.NodeKeeper       `inject:""`

	phase1result chan phase1Result
	phase2result chan phase2Result
	phase3result chan phase3Result

	currentPulseNumber uint32
}

// NewNaiveCommunicator constructor creates new NaiveCommunicator
func NewNaiveCommunicator() *NaiveCommunicator {
	return &NaiveCommunicator{}
}

// Start method implements Starter interface
func (nc *NaiveCommunicator) Start(ctx context.Context) error {
	nc.phase1result = make(chan phase1Result)
	nc.phase2result = make(chan phase2Result)
	nc.phase3result = make(chan phase3Result)
	nc.ConsensusNetwork.RegisterPacketHandler(packets.Phase1, nc.phase1DataHandler)
	nc.ConsensusNetwork.RegisterPacketHandler(packets.Phase2, nc.phase2DataHandler)
	nc.ConsensusNetwork.RegisterPacketHandler(packets.Phase3, nc.phase3DataHandler)
	return nil
}

func (nc *NaiveCommunicator) getPulseNumber() core.PulseNumber {
	pulseNumber := atomic.LoadUint32(&nc.currentPulseNumber)
	return core.PulseNumber(pulseNumber)
}

func (nc *NaiveCommunicator) setPulseNumber(new core.PulseNumber) bool {
	old := nc.getPulseNumber()
	return old < new && atomic.CompareAndSwapUint32(&nc.currentPulseNumber, uint32(old), uint32(new))
}

func (nc *NaiveCommunicator) sendRequestToNodes(participants []core.Node, packet packets.ConsensusPacket) {
	for _, node := range participants {
		if node.ID().Equal(nc.NodeKeeper.GetOrigin().ID()) {
			continue
		}

		go func(n core.Node) {
			err := nc.ConsensusNetwork.SignAndSendPacket(packet, n.ID(), nc.Cryptography)
			if err != nil {
				log.Errorln(err.Error())
			}
		}(node)
	}
}

func (nc *NaiveCommunicator) sendRequestToNodesWithOrigin(originClaim *packets.NodeAnnounceClaim,
	participants []core.Node, packet *packets.Phase1Packet) error {

	requests := make(map[core.RecordRef]packets.ConsensusPacket)
	for _, participant := range participants {
		err := originClaim.Update(participant.ID(), nc.Cryptography)
		if err != nil {
			return errors.Wrap(err, "Failed to update claims before sending in phase1")
		}
		req := *packet
		requests[participant.ID()] = &req
	}

	for ref, req := range requests {
		if ref.Equal(nc.NodeKeeper.GetOrigin().ID()) {
			continue
		}

		go func(node core.RecordRef, consensusPacket packets.ConsensusPacket) {
			err := nc.ConsensusNetwork.SignAndSendPacket(consensusPacket, node, nc.Cryptography)
			if err != nil {
				log.Errorln(err.Error())
			}
		}(ref, req)
	}
	return nil
}

func (nc *NaiveCommunicator) generatePhase2Response(origReq, req *packets.Phase2Packet, list network.UnsyncList) (*packets.Phase2Packet, error) {
	answers := make([]packets.ReferendumVote, 0)
	for _, vote := range req.GetVotes() {
		if vote.Type() != packets.TypeMissingNode {
			continue
		}
		v, ok := vote.(*packets.MissingNode)
		if !ok {
			log.Warnf("Phase 2 MissingNode request type mismatch")
			continue
		}
		ref, err := list.IndexToRef(int(v.NodeIndex))
		if err != nil {
			log.Warnf("Phase 2 MissingNode requested index: %d, error: %s", v.NodeIndex, err.Error())
			continue
		}
		node := list.GetActiveNode(ref)
		if node == nil {
			log.Warnf("Phase 2 MissingNode requested index: %d; mapped ref %s not found", v.NodeIndex, ref)
			continue
		}
		claim, err := packets.NodeToClaim(node)
		if err != nil {
			log.Warnf("Phase 2 MissingNode requested index: %d, mapped ref: %s, convertation node -> claim error: %s",
				v.NodeIndex, ref, err.Error())
			continue
		}
		proof := list.GetProof(ref)
		if proof == nil {
			log.Warnf("Phase 2 MissingNode requested index: %d, mapped ref: %s, proof not found", v.NodeIndex, ref)
			continue
		}
		ghs, ok := list.GlobuleHashSignatures()[ref]
		if !ok {
			log.Warnf("Phase 2 MissingNode requested index: %d, mapped ref: %s, GHS not found", v.NodeIndex, ref)
			continue
		}
		answer := packets.MissingNodeSupplementaryVote{
			NodeIndex:            v.NodeIndex,
			NodePulseProof:       *proof,
			GlobuleHashSignature: ghs,
			NodeClaimUnsigned:    *claim,
		}
		answers = append(answers, &answer)
		claims := list.GetClaims(ref)
		for _, claim := range claims {
			claimAnswer := packets.MissingNodeClaim{NodeIndex: v.NodeIndex, Claim: claim}
			answers = append(answers, &claimAnswer)
		}
	}
	response := packets.Phase2Packet{}
	response.SetBitSet(origReq.GetBitSet())
	ghs := origReq.GetGlobuleHashSignature()
	err := response.SetGlobuleHashSignature(ghs[:])
	if err != nil {
		return nil, errors.Wrap(err, "[ generatePhase2Response ] failed to set globule hash")
	}
	for _, answer := range answers {
		response.AddVote(answer)
	}
	return &response, nil
}

// ExchangePhase1 used in first consensus phase to exchange data between participants
func (nc *NaiveCommunicator) ExchangePhase1(
	ctx context.Context,
	originClaim *packets.NodeAnnounceClaim,
	participants []core.Node,
	packet *packets.Phase1Packet,
) (map[core.RecordRef]*packets.Phase1Packet, error) {
	result := make(map[core.RecordRef]*packets.Phase1Packet, len(participants))

	result[nc.ConsensusNetwork.GetNodeID()] = packet

	nc.setPulseNumber(packet.GetPulse().PulseNumber)

	var request *packets.Phase1Packet

	// TODO: awful, need rework
	if originClaim == nil {
		request = packet
	} else {
		request = &packets.Phase1Packet{}
		*request = *packet
		request.RemoveAnnounceClaim()
	}

	if originClaim == nil {
		nc.sendRequestToNodes(participants, request)
	} else {
		err := nc.sendRequestToNodesWithOrigin(originClaim, participants, packet)
		if err != nil {
			return nil, errors.Wrap(err, "[ExchangePhase1] Failed to send requests")
		}
	}

	shouldSendResponse := func(ref core.RecordRef) bool {
		val, ok := result[ref]
		return !ok || val == nil
	}
	response := request

	inslogger.FromContext(ctx).Infof("result len %d", len(result))
	for {
		select {
		case res := <-nc.phase1result:
			if res.packet.GetPulseNumber() != core.PulseNumber(nc.currentPulseNumber) {
				continue
			}

			if res.id.IsEmpty() {
				claim := res.packet.GetAnnounceClaim()
				if claim == nil {
					continue
				}
				res.id = claim.NodeRef
				err := nc.NodeKeeper.AddTemporaryMapping(claim.NodeRef, claim.ShortNodeID, claim.NodeAddress.Get())
				if err != nil {
					inslogger.FromContext(ctx).Warn("Error adding temporary mapping: " + err.Error())
					continue
				}
			}

			if shouldSendResponse(res.id) {
				// send response
				err := nc.ConsensusNetwork.SignAndSendPacket(response, res.id, nc.Cryptography)
				if err != nil {
					log.Errorln(err.Error())
				}
			}
			result[res.id] = res.packet

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
func (nc *NaiveCommunicator) ExchangePhase2(ctx context.Context, list network.UnsyncList,
	participants []core.Node, packet *packets.Phase2Packet) (map[core.RecordRef]*packets.Phase2Packet, error) {

	result := make(map[core.RecordRef]*packets.Phase2Packet, len(participants))

	result[nc.ConsensusNetwork.GetNodeID()] = packet

	nc.sendRequestToNodes(participants, packet)

	shouldSendResponse := func(p *phase2Result) bool {
		val, ok := result[p.id]
		firstResultReceive := !ok || val == nil
		packetContainsVoteRequests := p.packet.ContainsRequests()
		return firstResultReceive || packetContainsVoteRequests
	}

	inslogger.FromContext(ctx).Infof("result len %d", len(result))
	var err error
	for {
		select {
		case res := <-nc.phase2result:
			if res.packet.GetPulseNumber() != core.PulseNumber(nc.currentPulseNumber) {
				continue
			}

			if shouldSendResponse(&res) {
				// send response
				response := packet
				if res.packet.ContainsRequests() {
					response, err = nc.generatePhase2Response(packet, res.packet, list)
					if err != nil {
						log.Warnf("Failed to generate phase 2 response packet: %s", err.Error())
						continue
					}
				}
				err := nc.ConsensusNetwork.SignAndSendPacket(response, res.id, nc.Cryptography)
				if err != nil {
					log.Errorln(err.Error())
				}
			}
			result[res.id] = res.packet

			// FIXME: early return is commented to have synchronized length of phases on all nodes
			// if len(result) == len(participants) {
			// 	return result, nil
			// }
		case <-ctx.Done():
			return result, nil
		}
	}

	return result, nil
}

func selectCandidate(candidates []core.RecordRef) core.RecordRef {
	// TODO: make it random
	if len(candidates) == 0 {
		panic("candidates list should have at least 1 candidate (check consensus state matrix calculation routines)")
	}
	return candidates[0]
}

func (nc *NaiveCommunicator) sendAdditionalRequests(origReq *packets.Phase2Packet, additionalRequests []*AdditionalRequest) error {
	for _, req := range additionalRequests {
		newReq := *origReq
		newReq.AddVote(&packets.MissingNode{NodeIndex: uint16(req.RequestIndex)})
		receiver := selectCandidate(req.Candidates)
		err := nc.ConsensusNetwork.SignAndSendPacket(&newReq, receiver, nc.Cryptography)
		if err != nil {
			return errors.Wrapf(err, "Failed to send additional phase 2.1 request for index %d to node %s", req.RequestIndex, receiver)
		}
	}

	return nil
}

// ExchangePhase21 used in second consensus phase to exchange data between participants
func (nc *NaiveCommunicator) ExchangePhase21(ctx context.Context, list network.UnsyncList, packet *packets.Phase2Packet,
	additionalRequests []*AdditionalRequest) ([]packets.ReferendumVote, error) {

	type none struct{}
	incoming := make(map[core.RecordRef]none)
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

	err := nc.sendAdditionalRequests(packet, additionalRequests)
	if err != nil {
		return nil, errors.Wrap(err, "[ ExchangePhase2.1 ] Failed to send additional phase 2.1 requests")
	}

	for {
		select {
		case res := <-nc.phase2result:
			if res.packet.GetPulseNumber() != core.PulseNumber(nc.currentPulseNumber) {
				continue
			}

			if shouldSendResponse(&res) {
				// send response
				response := packet
				if res.packet.ContainsRequests() {
					response, err = nc.generatePhase2Response(packet, res.packet, list)
					if err != nil {
						log.Warnf("Failed to generate phase 2 response packet: %s", err.Error())
						continue
					}
				}
				err := nc.ConsensusNetwork.SignAndSendPacket(response, res.id, nc.Cryptography)
				if err != nil {
					log.Errorln(err.Error())
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
func (nc *NaiveCommunicator) ExchangePhase3(ctx context.Context, participants []core.Node, packet *packets.Phase3Packet) (map[core.RecordRef]*packets.Phase3Packet, error) {
	result := make(map[core.RecordRef]*packets.Phase3Packet, len(participants))

	result[nc.ConsensusNetwork.GetNodeID()] = packet

	nc.sendRequestToNodes(participants, packet)

	shouldSendResponse := func(p *phase3Result) bool {
		val, ok := result[p.id]
		return !ok || val == nil
	}

	inslogger.FromContext(ctx).Infof("result len %d", len(result))
	for {
		select {
		case res := <-nc.phase3result:
			if shouldSendResponse(&res) {
				// send response
				err := nc.ConsensusNetwork.SignAndSendPacket(packet, res.id, nc.Cryptography)
				if err != nil {
					log.Errorln(err.Error())
				}
			}
			result[res.id] = res.packet

			// FIXME: early return is commented to have synchronized length of phases on all nodes
			// if len(result) == len(participants) {
			// 	return result, nil
			// }
		case <-ctx.Done():
			return result, nil
		}
	}
}

func (nc *NaiveCommunicator) phase1DataHandler(packet packets.ConsensusPacket, sender core.RecordRef) {
	p, ok := packet.(*packets.Phase1Packet)
	if !ok {
		log.Errorln("invalid Phase1Packet")
		return
	}

	newPulse := p.GetPulse()

	if newPulse.PulseNumber < nc.getPulseNumber() {
		log.Warnln("ignore old pulse Phase1Packet")
		return
	}

	if nc.setPulseNumber(newPulse.PulseNumber) {
		go nc.PulseHandler.HandlePulse(context.Background(), newPulse)
	}

	nc.phase1result <- phase1Result{id: sender, packet: p}
}

func (nc *NaiveCommunicator) phase2DataHandler(packet packets.ConsensusPacket, sender core.RecordRef) {
	p, ok := packet.(*packets.Phase2Packet)
	if !ok {
		log.Errorln("invalid Phase2Packet")
		return
	}

	pulseNumber := p.GetPulseNumber()

	if pulseNumber < nc.getPulseNumber() {
		log.Warnln("ignore old pulse Phase2Packet")
		return
	}

	nc.phase2result <- phase2Result{id: sender, packet: p}
}

func (nc *NaiveCommunicator) phase3DataHandler(packet packets.ConsensusPacket, sender core.RecordRef) {
	p, ok := packet.(*packets.Phase3Packet)
	if !ok {
		log.Warn("failed to cast a type 3 packet to phase3packet")
	}
	nc.phase3result <- phase3Result{id: sender, packet: p}
}
