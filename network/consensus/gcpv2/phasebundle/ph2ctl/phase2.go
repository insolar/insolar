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
	"runtime"
	"time"

	"github.com/insolar/insolar/network/consensus/gcpv2/core/population"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/endpoints"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/announce"

	"github.com/insolar/insolar/network/consensus/common/lazyhead"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"

	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/insolar/network/consensus/gcpv2/core"
)

func NewPhase2Controller(loopingMinimalDelay time.Duration, packetPrepareOptions transport.PacketPrepareOptions,
	queueNshReady <-chan *population.NodeAppearance, lockOSThread bool) *Phase2Controller {

	return &Phase2Controller{
		packetPrepareOptions: packetPrepareOptions,
		queueNshReady:        queueNshReady,
		loopingMinimalDelay:  loopingMinimalDelay,
		lockOSThread:         lockOSThread,
	}
}

var _ core.PhaseController = &Phase2Controller{}

type Phase2Controller struct {
	core.PhaseControllerTemplate
	R                    *core.FullRealm
	packetPrepareOptions transport.PacketPrepareOptions
	queueNshReady        <-chan *population.NodeAppearance
	loopingMinimalDelay  time.Duration
	lockOSThread         bool
}

type Phase2PacketDispatcher struct {
	core.MemberPacketDispatcherTemplate
	ctl      *Phase2Controller
	isCapped bool
}

type UpdateSignal struct {
	NewTrustLevel member.TrustLevel
	UpdatedNode   *population.NodeAppearance
	DynNode       bool
}

func NewTrustUpdateSignal(n *population.NodeAppearance, newLevel member.TrustLevel) UpdateSignal {
	return UpdateSignal{UpdatedNode: n, NewTrustLevel: newLevel}
}

func NewDynamicNodeCreated(n *population.NodeAppearance) UpdateSignal {
	return UpdateSignal{UpdatedNode: n, DynNode: true, NewTrustLevel: member.UnknownTrust}
}

func NewDynamicNodeReady(n *population.NodeAppearance) UpdateSignal {
	return UpdateSignal{UpdatedNode: n, DynNode: true, NewTrustLevel: member.TrustBySome}
}

var pingSignal = UpdateSignal{}

func (v *UpdateSignal) IsPingSignal() bool {
	return v.UpdatedNode == nil && v.NewTrustLevel == 0 && !v.DynNode
}

func (c *Phase2Controller) GetPacketType() []phases.PacketType {
	return []phases.PacketType{phases.PacketPhase2, phases.PacketExtPhase2}
}

func (c *Phase2Controller) CreatePacketDispatcher(pt phases.PacketType, ctlIndex int,
	realm *core.FullRealm) (population.PacketDispatcher, core.PerNodePacketDispatcherFactory) {

	return &Phase2PacketDispatcher{ctl: c, isCapped: pt != phases.PacketPhase2}, nil
}

func (c *Phase2PacketDispatcher) DispatchMemberPacket(ctx context.Context, reader transport.MemberPacketReader, sender *population.NodeAppearance) error {

	p2 := reader.AsPhase2Packet()
	realm := c.ctl.R

	signalSent, announcedJoiner, err := announce.ApplyMemberAnnouncement(ctx, p2,
		p2.GetBriefIntroduction(), false, sender, realm)

	if err != nil {
		return err
	}

	neighbourhood := p2.GetNeighbourhood()
	neighbours, err := announce.VerifyNeighbourhood(ctx, neighbourhood, sender, announcedJoiner, realm)

	if err != nil {
		rep := misbehavior.FraudOf(err) // TODO unify approach to fraud registration
		if rep != nil {
			return sender.RegisterFraud(*rep)
		}
		return err
	}

	purgatory := realm.GetPurgatory()
	//senderID := sender.GetNodeID()

	for i, nb := range neighbours {
		modified := false
		if nb.Neighbour == nil {
			rank := neighbourhood[i].GetNodeRank()
			if rank.IsJoiner() && nb.Announcement.Joiner.JoinerProfile == nil {
				if announcedJoiner == nil || announcedJoiner.GetStaticNodeID() != nb.Announcement.MemberID {
					panic("unexpected")
				}
				// skip joiner that was introduced by the sender
				continue
			}
			err = purgatory.UnknownFromNeighbourhood(ctx, rank, nb.Announcement, c.isCapped)
		} else {
			modified, err = nb.Neighbour.ApplyNeighbourEvidence(sender, nb.Announcement, c.isCapped, nil)
		}
		if err != nil {
			inslogger.FromContext(ctx).Error(err)
		} else if modified {
			signalSent = true
		}
	}
	if !signalSent {
		sender.NotifyOnCustom(pingSignal)
	}

	return nil
}

func (c *Phase2PacketDispatcher) TriggerUnknownMember(ctx context.Context, memberID insolar.ShortNodeID,
	packet transport.MemberPacketReader, from endpoints.Inbound) (bool, error) {

	p2 := packet.AsPhase2Packet()

	// TODO check endpoint and PK

	return announce.ApplyUnknownAnnouncement(ctx, memberID, p2, p2.GetBriefIntroduction(), false, c.ctl.R)
}

func (c *Phase2Controller) StartWorker(ctx context.Context, realm *core.FullRealm) {
	c.R = realm
	go c.workerPhase2(ctx)
}

func (c *Phase2Controller) workerPhase2(ctx context.Context) {

	if c.lockOSThread {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
	}

	log := inslogger.FromContext(ctx)

	// This duration is a soft timeout - the worker will attempt to send all data in the queue before stopping.
	timings := c.R.GetTimings()
	// endOfPhase := time.After(c.R.AdjustedAfter(timings.EndOfPhase2))
	weightScaler := NewScalerInt64(timings.EndOfPhase1.Nanoseconds())
	startedAt := c.R.GetStartedAt()

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

	processedJoiners := 0
	processedNodes := 0

	if c.R.IsJoiner() { // exclude self
		processedJoiners++
	} else {
		processedNodes++
	}

	/*
		Is is safe to use an unsafe core.LessByNeighbourWeightForNodeAppearance as all objects have passed
		through a channel (sync) after neighbourWeight field was modified.
	*/
	nodeQueue := lazyhead.NewHeadedLazySortedList(neighbourSize-1, population.LessByNeighbourWeightForNodeAppearance, 1+c.R.GetPopulation().GetIndexedCount()>>1)
	joinQueue := lazyhead.NewHeadedLazySortedList(neighbourJoiners+joinersBoost, population.LessByNeighbourWeightForNodeAppearance, 1+nodeQueue.Len()>>1)

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
				log.Debug(">>>>workerPhase2: Done")
				return
			case np == nil:
				switch {
				// there is actually no need for early exit here
				case softTimeout && idleLoop:
					log.Debug(">>>>workerPhase2: timeout + idle")
					return
				case joinQueue.Len() > 0 || nodeQueue.Len() > 0 || !softTimeout:
					break inner
				case softTimeout:
					idleLoop = true
				}
			case np.GetNodeID() == c.R.GetSelfNodeID():
				continue
			case np.IsJoiner():
				//c.R.CreateAnnouncement(np, false) // sanity check
				joinQueue.Add(np)
			default:
				//c.R.CreateAnnouncement(np, false) // sanity check
				nodeQueue.Add(np)
			}
		}
		idleLoop = true
		if c.R.GetSelf().IsNSHRequired() {
			// we can't send anything yet - NSH wasn't provided yet
			continue
		}

		maxWeight := weightScaler.ScaleInt64(time.Since(startedAt).Nanoseconds())

		pop := c.R.GetPopulation()
		nodeCount, isComplete := pop.GetCountAndCompleteness(false)
		nodeCapacity, _ := pop.GetSealedCapacity()
		remainingNodes := nodeCount - processedNodes
		// TODO there are may be a racing when joiners are added?
		remainingJoiners := pop.GetJoinersCount() - processedJoiners

		/*
			Here we do scaling based on a number of nodes, to enable earlier processing for faster networks
			We can allow either all nodes to be processed ahead of time or not more 3/4 of population to ensure that
			there are enough members will be available to late members.
		*/
		available := nodeQueue.Len() + joinQueue.Len()
		switch {
		case isComplete && available >= remainingNodes+remainingJoiners:
			maxWeight = math.MaxUint32
		case nodeCapacity > neighbourSize<<2: // only works for large enough populations, > 4*neighbourhood size
			coverage := processedNodes + processedJoiners + available

			if coverage > nodeCapacity>>1 { // only if more half of members arrived
				const zeroCorrection = math.MaxUint32 >> 2 // zero correction value - scale is limited by 3/4

				k := uint64(math.MaxUint32) * uint64(coverage) / uint64(nodeCapacity)
				newWeight := uint32(0)
				switch {
				case k >= math.MaxUint32:
					newWeight = math.MaxUint32 - zeroCorrection
				case k >= zeroCorrection:
					newWeight = uint32(k) - zeroCorrection
				}
				if newWeight > maxWeight {
					maxWeight = newWeight
				}
			}
		}

		if maxWeight == math.MaxUint32 { // time is up
			softTimeout = true
		}

		takeJoiners := availableInQueue(joinersBoost, joinQueue, maxWeight)
		takeNodes := availableInQueue(takeJoiners, nodeQueue, maxWeight)

		if joinersBoost > 0 && takeNodes == 0 && joinQueue.Len() > takeJoiners {
			// when no normal nodes are zero then Ph2 can carry more joiners, lets unblock the boost
			takeJoiners = availableInQueue(0, joinQueue, maxWeight)
		}

		// NB! There is no reason to send Ph2 packet to only 1 non-joining node - the data will be the same as for Phase1
		if takeNodes > 0 || takeJoiners > 0 {
			nhBuf := neighbourhoodBuf[0:0]
			nhBuf = joinQueue.CutOffHeadByLenInto(takeJoiners, nhBuf)
			nhBuf = nodeQueue.CutOffHeadByLenInto(takeNodes, nhBuf)
			processedJoiners += takeJoiners
			processedNodes += takeNodes

			remainingJoiners -= takeJoiners
			remainingNodes -= takeNodes

			nh := make([]*population.NodeAppearance, len(nhBuf))
			for i, np := range nhBuf {
				// don't create MembershipAnnouncementReader here to avoid hitting lock by this only process
				nh[i] = np.(*population.NodeAppearance)
			}

			go c.sendPhase2(ctx, nh)

			idleLoop = false
		}
	}
}

func availableInQueue(captured int, queue lazyhead.HeadedLazySortedList, maxWeight uint32) int {
	if maxWeight == math.MaxUint32 {
		return queue.GetAvailableHeadLen(captured)
	}

	if queue.HasFullHead(captured) && queue.GetReversedHead(captured).(*population.NodeAppearance).GetNeighbourWeight() <= maxWeight {
		return queue.GetAvailableHeadLen(captured)
	}
	return 0
}

func readQueueOrDone(ctx context.Context, needsSleep bool, sleep time.Duration,
	q <-chan *population.NodeAppearance) (np *population.NodeAppearance, done bool) {

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

func (c *Phase2Controller) sendPhase2(ctx context.Context, neighbourhood []*population.NodeAppearance) {

	neighbourhoodAnnouncements := make([]transport.MembershipAnnouncementReader, len(neighbourhood))
	for i, np := range neighbourhood {
		neighbourhoodAnnouncements[i] = c.R.CreateAnnouncement(np, false)
	}

	p2 := c.R.GetPacketBuilder().PreparePhase2Packet(c.R.CreateLocalAnnouncement(), nil,
		neighbourhoodAnnouncements, c.packetPrepareOptions)

	sendOptions := c.packetPrepareOptions.AsSendOptions()

	p2.SendToMany(ctx, len(neighbourhood), c.R.GetPacketSender(),
		func(ctx context.Context, targetIdx int) (transport.TargetProfile, transport.PacketSendOptions) {
			np := neighbourhood[targetIdx]
			np.SetPacketSent(phases.PacketPhase2)
			return np, sendOptions
		})
}

func (c *Phase2Controller) workerRetryOnMissingNodes(ctx context.Context) {
	log := inslogger.FromContext(ctx)

	log.Debug("Phase2 has started re-requesting Phase1")

	s := c.R.GetSelf()
	if s.IsNSHRequired() {
		// we are close to end of Phase2 have no NSH - so missing Phase1 packets is the lesser problem
		return
	}

	pr1 := c.R.GetPacketBuilder().PreparePhase1Packet(c.R.CreateLocalAnnouncement(),
		c.R.GetOriginalPulse(), c.R.GetWelcomePackage(),
		transport.AlternativePhasePacket|transport.PrepareWithoutPulseData|c.packetPrepareOptions)

	sendOptions := c.packetPrepareOptions.AsSendOptions()

	for _, v := range c.R.GetPopulation().GetShuffledOtherNodes() {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if !v.IsNSHRequired() {
			continue
		}
		pr1.SendTo(ctx, v, sendOptions, c.R.GetPacketSender())
	}
}
