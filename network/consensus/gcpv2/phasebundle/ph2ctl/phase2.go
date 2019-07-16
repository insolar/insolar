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

package ph2ctl

import (
	"context"
	"math"
	"time"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/announce"

	"github.com/insolar/insolar/network/consensus/common/lazyhead"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"

	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/insolar/network/consensus/gcpv2/core"
)

func NewPhase2Controller(loopingMinimalDelay time.Duration, packetPrepareOptions transport.PacketSendOptions,
	queueNshReady <-chan *core.NodeAppearance) *Phase2Controller {

	return &Phase2Controller{
		packetPrepareOptions: packetPrepareOptions,
		queueNshReady:        queueNshReady,
		loopingMinimalDelay:  loopingMinimalDelay,
	}
}

var _ core.PhaseController = &Phase2Controller{}

type Phase2Controller struct {
	core.PhaseControllerTemplate
	core.MemberPacketDispatcher
	R                    *core.FullRealm
	packetPrepareOptions transport.PacketSendOptions
	queueNshReady        <-chan *core.NodeAppearance
	loopingMinimalDelay  time.Duration
}

type TrustUpdateSignal struct {
	NewTrustLevel member.TrustLevel
	UpdatedNode   *core.NodeAppearance
}

var pingSignal = TrustUpdateSignal{}

func (v *TrustUpdateSignal) IsPingSignal() bool {
	return v.UpdatedNode == nil
}

func (c *Phase2Controller) GetPacketType() []phases.PacketType {
	return []phases.PacketType{phases.PacketPhase2}
}

func (c *Phase2Controller) CreatePacketDispatcher(pt phases.PacketType, ctlIndex int, realm *core.FullRealm) (core.PacketDispatcher, core.PerNodePacketDispatcherFactory) {
	c.R = realm
	return c, nil
}

type resolvedNeighbour struct {
	neighbour *core.NodeAppearance
	ma        profiles.MembershipAnnouncement
}

func (c *Phase2Controller) DispatchMemberPacket(ctx context.Context, reader transport.MemberPacketReader, sender *core.NodeAppearance) error {

	p2 := reader.AsPhase2Packet()

	signalSent, err := announce.ApplyMemberAnnouncement(ctx, p2, p2.GetBriefIntroduction(), sender, c.R)
	if err != nil {
		return err
	}

	neighbourhood := p2.GetNeighbourhood()
	neighbours, err := c.verifyNeighbourhood(neighbourhood, sender)
	if err != nil {
		rep := misbehavior.FraudOf(err)
		if rep != nil {
			return sender.RegisterFraud(*rep)
		}
		return err
	}

	for i, nb := range neighbours {
		modified, err := nb.neighbour.ApplyNeighbourEvidence(sender, nb.ma)
		if err == nil && modified {
			signalSent = true

			if nb.neighbour.IsJoiner() {
				ja := neighbourhood[i].GetJoinerAnnouncement()
				err = c.R.AdvancePurgatoryNode(nb.neighbour.GetShortNodeID(), ja.GetBriefIntroduction(), nil, nb.neighbour)
			} else if !nb.ma.JoinerID.IsAbsent() {
				err = c.R.AdvancePurgatoryNode(nb.ma.JoinerID, nil, nil, nb.neighbour)
			}
		}
		if err != nil {
			inslogger.FromContext(ctx).Error(err)
			continue
		}
	}
	if !signalSent {
		sender.NotifyOnCustom(pingSignal)
	}

	return nil
}

func (c *Phase2Controller) verifyNeighbourhood(neighbourhood []transport.MembershipAnnouncementReader,
	n *core.NodeAppearance) ([]resolvedNeighbour, error) {

	hasThis := false
	hasSelf := false
	neighbours := make([]resolvedNeighbour, len(neighbourhood))
	nc := c.R.GetNodeCount()
	localID := c.R.GetSelfNodeID()

	for idx, nb := range neighbourhood {
		nid := nb.GetNodeID()
		if nid == n.GetShortNodeID() {
			hasSelf = true
		}
		neighbour := c.R.GetPopulation().GetNodeAppearance(nid)
		if neighbour == nil {
			return nil, n.Frauds().NewUnknownNeighbour(n.GetProfile())
		}
		nr := nb.GetNodeRank()

		if !profiles.MatchIntroAndRank(neighbour.GetProfile(), nc, nr) {
			return nil, n.Frauds().NewMismatchedNeighbourRank(n.GetProfile())
		}

		// TODO validate node proof - if fails, then fraud on sender
		// neighbourProfile.IsValidPacketSignature(nshEvidence.GetSignature())

		neighbours[idx].ma = announce.AnnouncementFromReader(nb)
		neighbours[idx].neighbour = neighbour

		if nid == localID {
			hasThis = true
		}
	}

	if !hasThis || hasSelf {
		return nil, n.Frauds().NewNeighbourMissingTarget(n.GetProfile())
	}
	if hasSelf {
		return nil, n.Frauds().NewNeighbourContainsRource(n.GetProfile())
	}
	return neighbours, nil
}

func (c *Phase2Controller) StartWorker(ctx context.Context, realm *core.FullRealm) {
	go c.workerPhase2(ctx)
}

func (c *Phase2Controller) workerPhase2(ctx context.Context) {

	// This duration is a soft timeout - the worker will attempt to send all data in the queue before stopping.
	timings := c.R.GetTimings()
	// endOfPhase := time.After(c.R.AdjustedAfter(timings.EndOfPhase2))
	weightScaler := NewNeighbourWeightScalerInt64(timings.EndOfPhase1.Nanoseconds())

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
	nodeQueue := lazyhead.NewHeadedLazySortedList(neighbourSize-1, core.LessByNeighbourWeightForNodeAppearance, remainingNodes>>1)
	joinQueue := lazyhead.NewHeadedLazySortedList(neighbourJoiners+joinersBoost, core.LessByNeighbourWeightForNodeAppearance, remainingJoiners>>1)

	idleLoop := false
	softTimeout := false
	for {
	inner:
		for {
			/*	This loop attempts to reads all messages from the channel before passing out
				Also it does some waiting on idle loops, but such waiting should be minimal, as queue weights are time-based.
			*/
			np, done := readQueueOrDone(ctx, idleLoop, c.loopingMinimalDelay, c.queueNshReady)
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

			nh := make([]*core.NodeAppearance, len(nhBuf))
			for i, np := range nhBuf {
				// don't create MembershipAnnouncementReader here to avoid hitting lock by this only process
				nh[i] = np.(*core.NodeAppearance)
			}

			go c.sendPhase2(ctx, nh)

			idleLoop = false
		}
	}
}

func (c *Phase2Controller) sendPhase2(ctx context.Context, neighbourhood []*core.NodeAppearance) {

	neighbourhoodAnnouncements := make([]transport.MembershipAnnouncementReader, len(neighbourhood))
	for i, np := range neighbourhood {
		neighbourhoodAnnouncements[i] = c.R.CreateAnnouncement(np)
	}

	p2 := c.R.GetPacketBuilder().PreparePhase2Packet(c.R.CreateLocalAnnouncement(), nil,
		neighbourhoodAnnouncements, c.packetPrepareOptions)

	p2.SendToMany(ctx, len(neighbourhood), c.R.GetPacketSender(),
		func(ctx context.Context, targetIdx int) (profiles.ActiveNode, transport.PacketSendOptions) {
			np := neighbourhood[targetIdx]
			np.SetPacketSent(phases.PacketPhase2)
			return np.GetProfile(), 0
		})
}

func availableInQueue(captured int, queue lazyhead.HeadedLazySortedList, remains int, maxWeight uint32) int {
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

	log.Info("Phase2 has started re-requesting Phase1")

	s := c.R.GetSelf()
	if s.IsNshRequired() {
		// we are close to end of Phase2 have no NSH - so missing Phase1 packets is the lesser problem
		return
	}

	pr1 := c.R.GetPacketBuilder().PreparePhase1Packet(c.R.CreateLocalAnnouncement(),
		c.R.GetOriginalPulse(), nil, transport.RequestForPhase1|c.packetPrepareOptions)

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
