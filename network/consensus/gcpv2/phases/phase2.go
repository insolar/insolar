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
	"fmt"
	"math"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

func NewPhase2Controller(packetPrepareOptions core.PacketSendOptions, queueNshReady <-chan *core.NodeAppearance) *Phase2Controller {
	return &Phase2Controller{
		packetPrepareOptions: packetPrepareOptions,
		queueNshReady:        queueNshReady,
		//queueTrustUpdated:    queueTrustUpdated,
	}
}

var _ core.PhaseController = &Phase2Controller{}

type Phase2Controller struct {
	PhaseControllerWithJoinersTemplate
	packetPrepareOptions core.PacketSendOptions
	queueNshReady        <-chan *core.NodeAppearance
	//queueTrustUpdated    chan<- TrustUpdateSignal // small enough to be sent as values
}

type TrustUpdateSignal struct {
	NewTrustLevel packets.NodeTrustLevel
	UpdatedNode   *core.NodeAppearance
}

var pingSignal = TrustUpdateSignal{}

func (v *TrustUpdateSignal) IsPingSignal() bool {
	return v.UpdatedNode == nil
}

func (*Phase2Controller) GetPacketType() packets.PacketType {
	return packets.PacketPhase2
}

func (c *Phase2Controller) CreatePerNodePacketHandler(ctlIndex int, node *core.NodeAppearance,
	realm *core.FullRealm, sharedNodeContext context.Context) (core.PhasePerNodePacketFunc, context.Context) {

	return c.createPerNodePacketHandler(ctlIndex, node, realm, sharedNodeContext, c.handleJoinerPacket)
}

func (c *Phase2Controller) HandleMemberPacket(ctx context.Context, reader packets.MemberPacketReader, n *core.NodeAppearance) error {

	p2 := reader.AsPhase2Packet()
	err := n.SetReceivedWithDupCheck(c.GetPacketType())
	if err != nil {
		// allows to avoid cheating by boosting up trust with additional Phase2 packets
		return err
	}

	// for _, ni := range p2.GetIntroductions() {
	// TODO I'm joiner, lets add another joiner's profile
	// and this enable the added other node for sending Ph2 by this node
	// }

	// TODO Verify neighbourhood first (incl presence of self) before applying it
	var signalSent = false
	for _, nb := range p2.GetNeighbourhood() {

		nid := nb.GetNodeID()
		neighbour := c.R.GetPopulation().GetNodeAppearance(nid)
		if neighbour == nil {
			// TODO unknown node - blame sender
			panic(fmt.Errorf("unlisted neighbour node: %v", nid))
			// R.GetBlameFactory().NewMultipleNshBlame()
		}

		nr := nb.GetNodeRank()
		mp := common2.NewMembershipProfile(nr.GetIndex(), nr.GetPower(), nb.GetNodeStateHashEvidence(),
			nb.GetAnnouncementSignature(), nb.GetRequestedPower())

		// TODO validate node proof - if fails, then fraud on sender
		// neighbourProfile.IsValidPacketSignature(nshEvidence.GetSignature())
		// TODO check NodeRank also
		// if p1.HasSelfIntro() {
		//	// TODO register protocol misbehavior - IntroClaim was not expected
		// }

		var modified bool
		nc := c.R.GetNodeCount()
		if nc != int(nr.GetTotalCount()) {
			modified, _ = n.RegisterFraud(n.Frauds().NewMismatchedMembershipNodeCount(n.GetProfile(), mp, nc))
		} else {
			modified, _ = neighbour.ApplyNeighbourEvidence(n, mp)
		}
		if modified {
			signalSent = true
		}
	}
	if !signalSent {
		n.NotifyOnCustom(pingSignal)
	}

	return nil
}

func (c *Phase2Controller) handleJoinerPacket(ctx context.Context, reader packets.MemberPacketReader, from *JoinerController) error {
	panic("unsupported")
}

func (c *Phase2Controller) StartWorker(ctx context.Context) {
	go c.workerPhase2(ctx)
}

func (c *Phase2Controller) workerPhase2(ctx context.Context) {

	// This duration is a soft timeout - the worker will attempt to send all data in the queue before stopping.
	timings := c.R.GetTimings()
	// endOfPhase := time.After(c.R.AdjustedAfter(timings.EndOfPhase2))
	weightScaler := common2.NewNeighbourWeightScalerInt64(timings.EndOfPhase1.Nanoseconds())

	if timings.StartPhase1RetryAt > 0 {
		timer := time.AfterFunc(c.R.AdjustedAfter(timings.StartPhase1RetryAt), func() {
			c.workerRetryOnMissingNodes(ctx)
		})
		defer timer.Stop()
	}

	neighbourSizes := c.R.GetNeighbourhoodSizes()
	neighbourSizes.VerifySizes()
	neighbourSize := neighbourSizes.NeighbourhoodSize
	neighbourJoiners := neighbourSizes.JoinersPerNeighbourhood
	joinersBoost := neighbourSizes.JoinersBoost

	neighbourhoodBuf := make([]interface{}, 0, neighbourSize-1)

	remainingJoiners := c.R.GetPopulation().GetJoinersCount()
	remainingNodes := c.R.GetNodeCount() - remainingJoiners

	if c.R.IsJoiner() { // exclude self
		neighbourJoiners--
		remainingJoiners--
	} else {
		remainingNodes--
	}

	/*
		Is is safe to use an unsafe core.LessByNeighbourWeightForNodeAppearance as all objects have passed
		through a channel (sync) after neighbourWeight field was modified.
	*/
	nodeQueue := common.NewHeadedLazySortedList(neighbourSize-1, core.LessByNeighbourWeightForNodeAppearance, remainingNodes>>1)
	joinQueue := common.NewHeadedLazySortedList(neighbourJoiners+joinersBoost, core.LessByNeighbourWeightForNodeAppearance, remainingJoiners>>1)

	idleLoop := false
	softTimeout := false
	for {
	inner:
		for {
			/*	This loop attempts to reads all messages from the channel before passing out
				Also it does some waiting on idle loops, but such waiting should be minimal, as queue weights are time-based.
			*/
			np, done := readQueueOrDone(ctx, idleLoop, loopingMinimalDelay, c.queueNshReady)
			switch {
			case done:
				return
			case np == nil:
				switch {
				case softTimeout && idleLoop:
					return
				case joinQueue.Len() > 0 || nodeQueue.Len() > 0 || !softTimeout:
					break inner
				case softTimeout:
					idleLoop = true
				}
			case np.IsJoiner():
				joinQueue.Add(np)
			default:
				nodeQueue.Add(np)
			}
		}
		idleLoop = true
		if c.R.GetSelf().IsNshRequired() {
			// we can't send anything yet - NSH wasn't provided yet
			continue
		}

		maxWeight := weightScaler.ScaleInt64(time.Since(c.R.GetStartedAt()).Nanoseconds())
		if maxWeight == math.MaxUint32 { // time is up
			softTimeout = true
		}

		takeJoiners := availableInQueue(joinersBoost, joinQueue, remainingJoiners, maxWeight)
		takeNodes := availableInQueue(takeJoiners, nodeQueue, remainingNodes, maxWeight)

		if joinersBoost > 0 && takeNodes == 0 && joinQueue.Len() > takeJoiners {
			// when no normal nodes are zero then Ph2 can carry more joiners, lets unblock the boost
			takeJoiners = availableInQueue(0, joinQueue, remainingJoiners, maxWeight)
		}

		// NB! There is no reason to send Ph2 packet to only 1 non-joining node - the data will be the same as for Phase1
		if takeNodes > 1 || takeJoiners > 0 {
			nhBuf := neighbourhoodBuf[0:0]
			nhBuf = joinQueue.CutOffHeadByLenInto(takeJoiners, nhBuf)
			nhBuf = nodeQueue.CutOffHeadByLenInto(takeNodes, nhBuf)
			remainingJoiners -= takeJoiners
			remainingNodes -= takeNodes

			nh := make([]*core.NodeAppearance, 1, len(nhBuf)+1)
			nh[0] = c.R.GetSelf()
			for _, np := range nhBuf {
				nh = append(nh, np.(*core.NodeAppearance))
			}

			go c.sendPhase2(ctx, nh, takeJoiners)

			idleLoop = false
		}
	}
}

func (c *Phase2Controller) sendPhase2(ctx context.Context, neighbourhood []*core.NodeAppearance, joinerCount int) {

	neighbourhoodAnnouncements := make([]packets.MembershipAnnouncementReader, len(neighbourhood))
	//introductions := make([]common2.NodeIntroduction, 0, joinerCount)

	for i, np := range neighbourhood {
		neighbourhoodAnnouncements[i] = c.R.CreateAnnouncement(np)
		//node := np.GetProfile()
		//if node.IsJoiner() {
		//	introductions = append(introductions, node.GetIntroduction())
		//}
	}

	p2 := c.R.GetPacketBuilder().PreparePhase2Packet(c.R.CreateLocalAnnouncement(),
		neighbourhoodAnnouncements, c.packetPrepareOptions)

	for _, np := range neighbourhood[1:] { // start from 1 to skip sending to self
		p2.SendTo(ctx, np.GetProfile(), 0, c.R.GetPacketSender())
		np.SetSentByPacketType(c.GetPacketType())
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func availableInQueue(captured int, queue common.HeadedLazySortedList, remains int, maxWeight uint32) int {
	if maxWeight == math.MaxUint32 {
		return queue.GetAvailableHeadLen(captured)
	}

	if queue.HasFullHead(captured) || (queue.Len() > 0 && queue.Len() >= remains) {
		if queue.GetReversedHead(captured).(*core.NodeAppearance).GetNeighbourWeight() <= maxWeight {
			return queue.GetAvailableHeadLen(captured)
		}
	}
	return 0
}

func readQueueOrDone(ctx context.Context, needsSleep bool, sleep time.Duration,
	q <-chan *core.NodeAppearance) (np *core.NodeAppearance, done bool) {

	if needsSleep {
		select {
		case <-ctx.Done():
			return nil, true // ctx.Err() ?
		case np = <-q:
			return np, false
		case <-time.After(sleep):
			return nil, false
		}
	} else {
		select {
		case <-ctx.Done():
			return nil, true // ctx.Err() ?
		case np = <-q:
			return np, false
		default:
			return nil, false
		}
	}
}

func (c *Phase2Controller) workerRetryOnMissingNodes(ctx context.Context) {
	log := inslogger.FromContext(ctx)

	log.Infof("Phase2 has started re-requesting Phase1")

	s := c.R.GetSelf()
	if s.IsNshRequired() {
		// we are close to end of Phase2 have no NSH - so missing Phase1 packets is the lesser problem
		return
	}

	pr1 := c.R.GetPacketBuilder().PreparePhase1Packet(c.R.CreateLocalAnnouncement(),
		c.R.GetOriginalPulse(), core.RequestForPhase1|c.packetPrepareOptions)

	for _, v := range c.R.GetPopulation().GetShuffledOtherNodes() {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if !v.IsNshRequired() {
			continue
		}
		pr1.SendTo(ctx, v.GetProfile(), 0, c.R.GetPacketSender())
	}
}

// func (R *ConsensusRoundController) preparePhase1Packet(po gcpv2.PacketSendOptions) gcpv2.PreparedSendPacket {
//
// 	var selfIntro gcpv2.NodeIntroduction = nil
// 	if c.R.joinersCount > 0 {
// 		selfIntro = c.R.self.nodeProfile.GetLastPublishedIntroduction()
// 	}
//
// 	return c.R.callback.PreparePhase1Packet(c.R.activePulse.originalPacket, selfIntro,
// 		c.R.self.membership, c.R.strategy.GetPacketOptions(gcpv2.Phase1)|po)
// }
