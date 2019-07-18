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

package ph3ctl

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/insolar/insolar/network/consensus/common/chaser"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/statevector"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/nodeset"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/ph2ctl"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/insolar/network/consensus/gcpv2/core"
)

func NewPhase3Controller(loopingMinimalDelay time.Duration, packetPrepareOptions transport.PacketSendOptions, queueTrustUpdated <-chan ph2ctl.TrustUpdateSignal,
	consensusStrategy ConsensusSelectionStrategy, inspectionFactory VectorInspectorFactory) *Phase3ControllerV2 {
	return &Phase3ControllerV2{
		packetPrepareOptions: packetPrepareOptions,
		queueTrustUpdated:    queueTrustUpdated,
		consensusStrategy:    consensusStrategy,
		loopingMinimalDelay:  loopingMinimalDelay,
		inspectionFactory:    inspectionFactory,
	}
}

var _ core.PhaseController = &Phase3ControllerV2{}

type Phase3ControllerV2 struct {
	core.PhaseControllerTemplate
	core.MemberPacketDispatcherTemplate
	packetPrepareOptions transport.PacketSendOptions
	queueTrustUpdated    <-chan ph2ctl.TrustUpdateSignal
	queuePh3Recv         chan InspectedVector
	consensusStrategy    ConsensusSelectionStrategy
	inspectionFactory    VectorInspectorFactory
	loopingMinimalDelay  time.Duration
	R                    *core.FullRealm

	rw        sync.RWMutex
	inspector VectorInspector
	// packetHandler to Worker channel
}

func (c *Phase3ControllerV2) CreatePacketDispatcher(pt phases.PacketType, ctlIndex int, realm *core.FullRealm) (core.PacketDispatcher, core.PerNodePacketDispatcherFactory) {
	c.R = realm
	return c, nil
}

func (*Phase3ControllerV2) GetPacketType() []phases.PacketType {
	return []phases.PacketType{phases.PacketPhase3}
}

func (c *Phase3ControllerV2) DispatchMemberPacket(ctx context.Context, reader transport.MemberPacketReader, n *core.NodeAppearance) error {

	p3 := reader.AsPhase3Packet()

	// TODO validations

	iv := c.getInspector().InspectVector(n, statevector.NewVector(p3.GetBitset(),
		statevector.NewSubVector(p3.GetTrustedGlobulaAnnouncementHash(), p3.GetTrustedGlobulaStateSignature(), p3.GetTrustedExpectedRank()),
		statevector.NewSubVector(p3.GetDoubtedGlobulaAnnouncementHash(), p3.GetDoubtedGlobulaStateSignature(), p3.GetDoubtedExpectedRank())))

	if iv == nil {
		panic("illegal state")
	}
	c.queuePh3Recv <- iv

	return nil
}

func (c *Phase3ControllerV2) getInspector() VectorInspector {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.inspector
}

func (c *Phase3ControllerV2) setInspector(inspector VectorInspector) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.inspector = inspector
}

func (c *Phase3ControllerV2) StartWorker(ctx context.Context, realm *core.FullRealm) {
	c.queuePh3Recv = make(chan InspectedVector, c.R.GetNodeCount())
	c.inspector = NewBypassInspector()

	go c.workerPhase3(ctx)
}

func (c *Phase3ControllerV2) workerPhase3(ctxRound context.Context) {

	ctx, cancel := context.WithDeadline(ctxRound, time.Now().Add(c.R.AdjustedAfter(c.R.GetTimings().EndOfPhase3)))
	defer cancel()

	if !c.workerPrePhase3(ctx) {
		// context was stopped in a hard way, we are dead in terms of consensus
		// TODO should wait for further packets to decide if we need to turn ourselves into suspended state
		// c.R.StopRoundByTimeout()
		return
	}

	vectorHelper := c.R.GetPopulation().CreateVectorHelper()
	localProjection := vectorHelper.CreateProjection()
	localInspector := c.inspectionFactory.CreateInspector(&localProjection, c.R.GetDigestFactory(), c.R.GetSelfNodeID())

	// it also finalizes internal state to allow later parallel use
	localHashedVector := localInspector.CreateVector(c.R.GetSigner())

	c.setInspector(localInspector)

	go c.workerSendPhase3(ctx, localHashedVector)

	if !c.workerRecvPhase3(ctx, localInspector) {
		// context was stopped in a hard way or we have left a consensus
		return
	}
	// TODO should wait for further packets to decide if we need to turn ourselves into suspended state
	// c.R.StopRoundByTimeout()

	// avoid any links to controllers for this flusher
	go workerQueueFlusher(ctxRound, c.queuePh3Recv, c.queueTrustUpdated)
}

func workerQueueFlusher(ctxRound context.Context, q0 chan InspectedVector, q1 <-chan ph2ctl.TrustUpdateSignal) {
	for {
		select {
		case <-ctxRound.Done():
			return
		case _, ok := <-q0:
			if ok {
				continue
			}
			if q1 == nil {
				return
			}
			q0 = nil
		case _, ok := <-q1:
			if ok {
				continue
			}
			if q0 == nil {
				return
			}
			q1 = nil
		}
	}
}

func (c *Phase3ControllerV2) workerPrePhase3(ctx context.Context) bool {
	log := inslogger.FromContext(ctx)

	log.Debug(">>>>workerPrePhase3: begin")

	timings := c.R.GetTimings()
	startOfPhase3 := time.After(c.R.AdjustedAfter(timings.EndOfPhase2))
	chasingDelayTimer := chaser.NewChasingTimer(timings.BeforeInPhase2ChasingDelay)

	var countFraud = 0
	var countHasNsh = 0
	var countTrustBySome = 0
	var countTrustByNeighbors = 0

outer:
	for {
		select {
		case <-ctx.Done():
			log.Debug(">>>>workerPrePhase3: ctx.Done")
			return false // ctx.Err() ?
		case <-chasingDelayTimer.Channel():
			log.Debug(">>>>workerPrePhase3: chaseExpired")
			break outer
		case <-startOfPhase3:
			log.Debug(">>>>workerPrePhase3: startOfPhase3")
			break outer
		case sig := <-c.queueTrustUpdated:
			switch {
			case sig.IsPingSignal(): // ping indicates arrival of Phase2 packet, to support chasing
				// TODO chasing
				continue
			case sig.NewTrustLevel < 0:
				countFraud++
				continue // no chasing delay on fraud
			case sig.NewTrustLevel == member.UnknownTrust:
				countHasNsh++
				// if countHasNsh >= R.othersCount {
				// 	// we have answers from all
				// 	break outer
				// }
			case sig.NewTrustLevel >= member.TrustByNeighbors:
				countTrustByNeighbors++
				fallthrough
			default:
				countTrustBySome++

				pop := c.R.GetPopulation()
				// We have some-trusted from all nodes, and the majority of them are well-trusted
				if countTrustBySome >= pop.GetOthersCount() && countTrustByNeighbors >= pop.GetBftMajorityCount() {
					log.Debug(">>>>workerPrePhase3: all")
					break outer
				}

				if chasingDelayTimer.IsEnabled() {
					// We have answers from all nodes, and the majority of them are well-trusted
					if countHasNsh >= pop.GetOthersCount() && countTrustByNeighbors >= pop.GetBftMajorityCount() {
						chasingDelayTimer.RestartChase()
						log.Debug(">>>>workerPrePhase3: chaseStartedAll")
					} else if countTrustBySome-countFraud >= pop.GetBftMajorityCount() {
						// We can start chasing-timeout after getting answers from majority of some-trusted nodes
						chasingDelayTimer.RestartChase()
						log.Debug(">>>>workerPrePhase3: chaseStartedSome")
					}
				}
			}
		}
	}

	/* Ensure that NSH is available, otherwise we can't normally build packets */
	for c.R.GetSelf().IsNshRequired() {
		select {
		case <-ctx.Done():
			log.Debug(">>>>workerPrePhase3: ctx.Done")
			return false // ctx.Err() ?
		case <-c.queueTrustUpdated:
		case <-time.After(c.loopingMinimalDelay):
		}
	}
	return true
}

func (c *Phase3ControllerV2) workerRescanForMissing(ctx context.Context, missing chan InspectedVector) {
	for {
		select {
		case <-ctx.Done():
			return
		case sig := <-c.queueTrustUpdated:
			if sig.IsPingSignal() {
				continue
			}
			// TODO
		case <-missing:
			panic("not implemented")
			//TODO
		}
	}
}

func (c *Phase3ControllerV2) workerSendPhase3(ctx context.Context, selfData statevector.Vector) {

	otherNodes := c.R.GetPopulation().GetShuffledOtherNodes()

	p3 := c.R.GetPacketBuilder().PreparePhase3Packet(c.R.CreateLocalAnnouncement(), selfData,
		c.packetPrepareOptions)

	p3.SendToMany(ctx, len(otherNodes), c.R.GetPacketSender(),
		func(ctx context.Context, targetIdx int) (profiles.StaticProfile, transport.PacketSendOptions) {
			np := otherNodes[targetIdx]
			np.SetPacketSent(phases.PacketPhase3)
			return np.GetProfile().GetStatic(), 0
		})

	// TODO send to shuffled joiners as well?
}

func (c *Phase3ControllerV2) workerRecvPhase3(ctx context.Context, localInspector VectorInspector) bool {

	log := inslogger.FromContext(ctx)

	var queueMissing chan InspectedVector

	timings := c.R.GetTimings()
	softDeadline := time.After(c.R.AdjustedAfter(timings.EndOfPhase3))
	chasingDelayTimer := chaser.NewChasingTimer(timings.BeforeInPhase3ChasingDelay)

	verifiedStatTbl := nodeset.NewConsensusStatTable(c.R.GetNodeCount())
	originalStatTbl := nodeset.NewConsensusStatTable(c.R.GetNodeCount())

	// should it be updatable?
	{
		selfIndex := c.R.GetSelf().GetIndex().AsInt()
		localStat := nodeset.StateToConsensusStatRow(localInspector.GetBitset())
		localStatCopy := localStat
		verifiedStatTbl.PutRow(selfIndex, &localStat)
		originalStatTbl.PutRow(selfIndex, &localStatCopy)
	}

	remainingNodes := c.R.GetPopulation().GetOthersCount()

	// TODO detect nodes produced similar bitmaps, but different GSH
	// even if we wont have all NSH, we can let to know these nodes on such collision
	// bitsetMatcher := make(map[gcpv2.StateBitset])

	// hasher := nodeset.NewFilteredSequenceHasher(c.R.GetDigestFactory(), localVector)

	alteredDoubtedGshCount := 0
	var consensusSelection ConsensusSelection

outer:
	for {
		select {
		case <-ctx.Done():
			log.Debug("Phase3 cancelled")
			return false
		case <-softDeadline:
			log.Debug("Phase3 deadline")
			consensusSelection = c.consensusStrategy.SelectOnStopped(&verifiedStatTbl, true, c.R)
			break outer
		case <-chasingDelayTimer.Channel():
			log.Debug("Phase3 chasing expired")
			consensusSelection = c.consensusStrategy.SelectOnStopped(&verifiedStatTbl, true, c.R)
			break outer
		case d := <-c.queuePh3Recv:
			switch {
			case d.HasMissingMembers():
				if queueMissing == nil {
					queueMissing = make(chan InspectedVector, len(c.queuePh3Recv))
					go c.workerRescanForMissing(ctx, queueMissing)
				}
				queueMissing <- d
				// do chasing
			case !d.IsInspected():
				d = d.Reinspect(localInspector)
				if !d.IsInspected() {
					if d.HasMissingMembers() {
						c.queuePh3Recv <- d
					}
					// TODO heavy inspection with hash recalculations should be running on a limited pool
					go func() {
						d.Inspect()
						if d.IsInspected() {
							c.queuePh3Recv <- d
						} else {
							inslogger.FromContext(ctx).Errorf("unable to inspect vector: %v", d)
						}
					}()
					break // do chasing
				}
				fallthrough
			default:
				nodeStats, vr := d.GetInspectionResults()
				if log.Is(insolar.DebugLevel) {
					var logMsg interface{}
					switch {
					case d.HasSenderFault() || nodeStats == nil:
						logMsg = "fault"
					case nodeStats.HasAllValues(0):
						logMsg = "missed"
					default:
						logMsg = "added"
					}

					na := d.GetNode()
					log.Debugf(
						"%s: s:%v t:%v idx:%d left:%d\n Here:%v\nThere:%v\n Comp:%v\nStats:%v\n",
						logMsg, na.GetNodeID(), c.R.GetSelf().GetNodeID(), na.GetIndex(), remainingNodes,
						localInspector.GetBitset(), d.GetBitset(), d, nodeStats,
					)
				}

				if nodeStats == nil {
					break
				}

				{
					nodeIndex := d.GetNode().GetIndex().AsInt()
					verifiedStatTbl.PutRow(nodeIndex, nodeStats)
					originalStat := nodeset.StateToConsensusStatRow(d.GetBitset())
					originalStatTbl.PutRow(nodeIndex, &originalStat)
				}

				remainingNodes--

				if vr.AnyOf(nodeset.NvrDoubtedAlteredNodeSet) {
					alteredDoubtedGshCount++
				}

				if remainingNodes <= 0 {
					consensusSelection = c.consensusStrategy.SelectOnStopped(&verifiedStatTbl, false, c.R)
					log.Debug("Phase3 done all")
					break outer
				}

				consensusSelection = c.consensusStrategy.TrySelectOnAdded(
					&verifiedStatTbl, d.GetNode().GetProfile(), nodeStats, c.R)

				if consensusSelection == nil {
					continue
				}
			}
			if chasingDelayTimer.IsEnabled() && consensusSelection.CanBeImproved() {
				log.Debug("Phase3 (re)start chasing")
				chasingDelayTimer.RestartChase()
				continue
			}
			log.Debug("Phase3 done by strategy")
			break outer
		}
	}

	if log.Is(insolar.DebugLevel) {
		tblHeader := fmt.Sprintf("Consensus Node View (%%s): %v", c.R.GetSelfNodeID())
		typeHeader := "Original=Verified"
		if !originalStatTbl.EqualsTyped(&verifiedStatTbl) {
			log.Debug(originalStatTbl.TableFmt(fmt.Sprintf(tblHeader, "Original"), nodeset.FmtConsensusStat))
			typeHeader = "Verified"
		}
		log.Debug(verifiedStatTbl.TableFmt(fmt.Sprintf(tblHeader, typeHeader), nodeset.FmtConsensusStat))
	}

	if consensusSelection == nil {
		panic("illegal state")
	}

	sameWithActive := false
	selectionSet := (*nodeset.ConsensusBitsetRow)(nil)

	if consensusSelection.IsSameWithActive() {
		sameWithActive = true
	} else {
		selectionSet = consensusSelection.GetConsensusNodes()
		sameWithActive = selectionSet.Len() == c.R.GetNodeCount() && selectionSet.HasAllValues(nodeset.CbsIncluded)
	}

	if sameWithActive {
		selectionSet = nil
		log.Info("Consensus is finished as same")
	} else {
		log.Infof("Consensus is finished as different, %v", selectionSet)
		// TODO update population and/or start Phase 4
	}

	b := c.R.CreateNextCensusBuilder()
	if c.buildNextPopulation(b.GetPopulationBuilder(), selectionSet) {
		// TODO HACK
		priming := c.R.GetPrimingCloudHash()
		b.SetGlobulaStateHash(priming)
		b.SealCensus()
		c.R.FinishRound(b, priming)
		return true
	}
	log.Info("Node has left")
	c.R.FinishRound(b, nil)
	return false
}

func (c *Phase3ControllerV2) buildNextPopulation(pb census.PopulationBuilder, nodeset *nodeset.ConsensusBitsetRow) bool {

	// pop := c.R.GetPopulation()
	// count := 0
	// for _, na := pop.GetIndexedNodes() {
	//
	// }
	//
	// for
	//
	// if isLeaver, leaveReason, _, _, _ := c.R.GetSelf().GetRequestedState(); isLeaver {
	//	//we are leaving, no need to build population, but lets make it look nice
	//	pb.RemoveOthers()
	//	lp := pb.GetLocalProfile()
	//	lp.SetIndex(0)
	//	lp.SetOpModeAndLeaveReason(leaveReason)
	//	return false
	// }
	//
	//
	// //if pb.GetLocalProfile().GetOpMode().IsEvicted() /* TODO and local is still evicted */ {
	// //	//this node was evicted, so we can have a consensus with ourselves
	// //	pb.RemoveOthers()
	// //	return
	// //}
	//
	// pop := c.R.GetPopulation()
	// for _, np := range pb.GetUnorderedProfiles() {
	//	opm := np.GetOpMode()
	//	if np.IsJoiner() || opm.IsEvicted() { panic("illegal state") }
	//
	//	idx := np.GetIndex()
	//	//TODO if nodeset.
	//
	//	na := pop.GetNodeAppearanceByIndex(idx)
	//	//TODO MUST use cached vector values
	//	isLeaver, _, _, _, _ := na.GetRequestedState()
	//	if isLeaver {
	//		np.SetOpMode(common2.MemberModeEvictedGracefully)
	//	}
	// }
	return true
}
