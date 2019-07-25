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

	"github.com/insolar/insolar/network/consensus/common/consensuskit"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/consensus"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/inspectors"

	"github.com/insolar/insolar/network/consensus/common/chaser"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/statevector"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/nodeset"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/ph2ctl"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/insolar/network/consensus/gcpv2/core"
)

func NewPhase3Controller(loopingMinimalDelay time.Duration, packetPrepareOptions transport.PacketSendOptions,
	queueTrustUpdated <-chan ph2ctl.TrustUpdateSignal, consensusStrategy consensus.SelectionStrategy,
	inspectionFactory inspectors.VectorInspectorFactory, enabledFast bool) *Phase3Controller {
	return &Phase3Controller{
		packetPrepareOptions: packetPrepareOptions,
		queueTrustUpdated:    queueTrustUpdated,
		consensusStrategy:    consensusStrategy,
		loopingMinimalDelay:  loopingMinimalDelay,
		inspectionFactory:    inspectionFactory,
		isFastEnabled:        enabledFast,
	}
}

var _ core.PhaseController = &Phase3Controller{}

type Phase3Controller struct {
	core.PhaseControllerTemplate
	packetPrepareOptions transport.PacketSendOptions
	queueTrustUpdated    <-chan ph2ctl.TrustUpdateSignal
	queuePh3Recv         chan inspectors.InspectedVector
	consensusStrategy    consensus.SelectionStrategy
	inspectionFactory    inspectors.VectorInspectorFactory
	loopingMinimalDelay  time.Duration
	isFastEnabled        bool
	R                    *core.FullRealm

	rw        sync.RWMutex
	inspector inspectors.VectorInspector

	// packetHandler to Worker channel
}

type Phase3PacketDispatcher struct {
	core.MemberPacketDispatcherTemplate
	ctl           *Phase3Controller
	customOptions uint32
}

const outOfOrderPhase3 = 1

func (c *Phase3Controller) CreatePacketDispatcher(pt phases.PacketType, ctlIndex int, realm *core.FullRealm) (core.PacketDispatcher, core.PerNodePacketDispatcherFactory) {
	customOptions := uint32(0)
	if pt != phases.PacketPhase3 {
		customOptions = outOfOrderPhase3
	}
	return &Phase3PacketDispatcher{ctl: c, customOptions: customOptions}, nil
}

func (*Phase3Controller) GetPacketType() []phases.PacketType {
	return []phases.PacketType{phases.PacketPhase3, phases.PacketFastPhase3}
}

func (c *Phase3PacketDispatcher) DispatchMemberPacket(ctx context.Context, reader transport.MemberPacketReader, n *core.NodeAppearance) error {

	p3 := reader.AsPhase3Packet()

	// TODO validations

	iv := c.ctl.getInspector().InspectVector(ctx, n, c.customOptions, statevector.NewVector(p3.GetBitset(),
		statevector.NewSubVector(p3.GetTrustedGlobulaAnnouncementHash(), p3.GetTrustedGlobulaStateSignature(), p3.GetTrustedExpectedRank()),
		statevector.NewSubVector(p3.GetDoubtedGlobulaAnnouncementHash(), p3.GetDoubtedGlobulaStateSignature(), p3.GetDoubtedExpectedRank())))

	if iv == nil {
		panic("illegal state")
	}
	c.ctl.queuePh3Recv <- iv

	return nil
}

func (c *Phase3Controller) getInspector() inspectors.VectorInspector {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.inspector
}

func (c *Phase3Controller) setInspector(inspector inspectors.VectorInspector) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.inspector = inspector
}

func (c *Phase3Controller) StartWorker(ctx context.Context, realm *core.FullRealm) {
	c.R = realm
	c.queuePh3Recv = make(chan inspectors.InspectedVector, c.R.GetNodeCount())
	c.inspector = inspectors.NewBypassInspector()

	go c.workerPhase3(ctx)
}

func (c *Phase3Controller) workerPhase3(ctx context.Context) {

	if !c.workerPrePhase3(ctx) {
		// context was stopped in a hard way, we are dead in terms of consensus
		// TODO should wait for further packets to decide if we need to turn ourselves into suspended state
		// c.R.StopRoundByTimeout()
		return
	}

	vectorHelper := c.R.GetPopulation().CreateVectorHelper()
	localProjection := vectorHelper.CreateProjection()
	localInspector := c.inspectionFactory.CreateInspector(&localProjection, c.R.GetDigestFactory(), c.R.GetSelfNodeID())

	// enables parallel use
	localInspector.PrepareForInspection(ctx)
	c.setInspector(localInspector)

	if !c.R.IsJoiner() {
		// joiner has no vote in consensus, hence there is no reason to send Phase3 from it
		localHashedVector := localInspector.CreateVector(c.R.GetSigner())
		go c.workerSendPhase3(ctx, localHashedVector, c.packetPrepareOptions)
	}

	if !c.workerRecvPhase3(ctx, localInspector) {
		// context was stopped in a hard way or we have left a consensus
		return
	}
	// TODO should wait for further packets to decide if we need to turn ourselves into suspended state
	// c.R.StopRoundByTimeout()

	workerQueueFlusher(c.R, c.queuePh3Recv, c.queueTrustUpdated)
}

func workerQueueFlusher(realm *core.FullRealm, q0 chan inspectors.InspectedVector, q1 <-chan ph2ctl.TrustUpdateSignal) {
	realm.AddPoll(func(ctx context.Context) bool {
		select {
		case _, ok := <-q0:
			return ok
		default:
			return q0 != nil
		}
	})

	realm.AddPoll(func(ctx context.Context) bool {
		select {
		case _, ok := <-q1:
			return ok
		default:
			return q1 != nil
		}
	})
}

func (c *Phase3Controller) workerPrePhase3(ctx context.Context) bool {
	log := inslogger.FromContext(ctx)

	log.Debug(">>>>workerPrePhase3: begin")

	timings := c.R.GetTimings()
	startOfPhase3 := time.After(c.R.AdjustedAfter(timings.EndOfPhase2))
	chasingDelayTimer := chaser.NewChasingTimer(timings.BeforeInPhase2ChasingDelay)

	var countFraud = 0
	var countHasNsh = 0
	var countTrustBySome = 0
	var countTrustByNeighbors = 0

	pop := c.R.GetPopulation()
	didFastPhase3 := false

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
			case sig.NewTrustLevel >= member.TrustByNeighbors:
				countTrustByNeighbors++
				fallthrough
			default:
				countTrustBySome++

				indexedCount, isComplete := pop.GetCountAndCompleteness(false)
				bftMajority := consensuskit.BftMajority(indexedCount)

				// We have some-trusted from all nodes, and the majority of them are well-trusted
				if isComplete && countTrustBySome >= indexedCount && countTrustByNeighbors >= bftMajority {
					log.Debug(">>>>workerPrePhase3: all and complete")
					break outer
				}

				if chasingDelayTimer.IsEnabled() {
					// We have answers from all nodes, and the majority of them are well-trusted
					if countHasNsh >= indexedCount && countTrustByNeighbors >= bftMajority {
						chasingDelayTimer.RestartChase()
						log.Debug(">>>>workerPrePhase3: chaseStartedAll")
					} else if countTrustBySome-countFraud >= bftMajority {
						// We can start chasing-timeout after getting answers from majority of some-trusted nodes
						chasingDelayTimer.RestartChase()
						log.Debug(">>>>workerPrePhase3: chaseStartedSome")
					}
				}
			}

			// if we didn't went for a full phase3 sending, but we have all nodes, then should try a shortcut
			if c.isFastEnabled {
				indexedCount, isComplete := pop.GetCountAndCompleteness(false)
				if isComplete && countHasNsh >= indexedCount && !didFastPhase3 {
					didFastPhase3 = true
					go c.workerSendFastPhase3(ctx)
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

func (c *Phase3Controller) workerRescanForMissing(ctx context.Context, missing chan inspectors.InspectedVector) {
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
			// TODO - rescan vector and send results
			// c.queuePh3Recv <- d
		}
	}
}

func (c *Phase3Controller) workerSendFastPhase3(ctx context.Context) {

	// TODO vector calculation for fast options
	// handling of fast phase3 may also require a separate vector inspector
	// c.workerSendPhase3(ctx, nil, c.packetPrepareOptions|transport.AlternativePhasePacket)
}

func (c *Phase3Controller) workerSendPhase3(ctx context.Context, selfData statevector.Vector, options transport.PacketSendOptions) {

	p3 := c.R.GetPacketBuilder().PreparePhase3Packet(c.R.CreateLocalAnnouncement(), selfData, options)

	selfID := c.R.GetSelfNodeID()
	otherNodes := c.R.GetPopulation().GetAnyNodes(true, true)

	p3.SendToMany(ctx, len(otherNodes), c.R.GetPacketSender(),
		func(ctx context.Context, targetIdx int) (transport.TargetProfile, transport.PacketSendOptions) {
			np := otherNodes[targetIdx]
			if np.GetNodeID() == selfID {
				return nil, 0
			}
			np.SetPacketSent(phases.PacketPhase3)
			return np, 0
		})
}

func (c *Phase3Controller) workerRecvPhase3(ctx context.Context, localInspector inspectors.VectorInspector) bool {

	log := inslogger.FromContext(ctx)

	var queueMissing chan inspectors.InspectedVector

	timings := c.R.GetTimings()
	softDeadline := time.After(c.R.AdjustedAfter(timings.EndOfPhase3))
	chasingDelayTimer := chaser.NewChasingTimer(timings.BeforeInPhase3ChasingDelay)

	verifiedStatTbl := nodeset.NewConsensusStatTable(c.R.GetNodeCount())
	originalStatTbl := nodeset.NewConsensusStatTable(c.R.GetNodeCount())

	processedNodesFlawlessly := 0

	if !c.R.IsJoiner() {
		selfIndex := c.R.GetSelf().GetIndex().AsInt()
		// should it be updatable?
		localStat := nodeset.StateToConsensusStatRow(localInspector.GetBitset())
		localStatCopy := localStat
		verifiedStatTbl.PutRow(selfIndex, &localStat)
		originalStatTbl.PutRow(selfIndex, &localStatCopy)
		processedNodesFlawlessly++
	}

	population := c.R.GetPopulation()

	// TODO detect nodes produced similar bitmaps, but different GSH
	// even if we wont have all NSH, we can let to know these nodes on such collision
	// bitsetMatcher := make(map[gcpv2.StateBitset])

	// hasher := nodeset.NewFilteredSequenceHasher(c.R.GetDigestFactory(), localVector)

	// alteredDoubtedGshCount := 0
	var consensusSelection consensus.Selection

outer:
	for {
		popCount, popCompleteness := population.GetCountAndCompleteness(false)
		if popCompleteness && popCount <= processedNodesFlawlessly {
			consensusSelection = c.consensusStrategy.SelectOnStopped(&verifiedStatTbl, false,
				population.GetBftMajorityCount())

			log.Debug("Phase3 done all")
			break outer
		}

		select {
		case <-ctx.Done():
			log.Debug("Phase3 cancelled")
			return false
		case <-softDeadline:
			log.Debug("Phase3 deadline")
			consensusSelection = c.consensusStrategy.SelectOnStopped(&verifiedStatTbl, true, population.GetBftMajorityCount())
			break outer
		case <-chasingDelayTimer.Channel():
			log.Debug("Phase3 chasing expired")
			consensusSelection = c.consensusStrategy.SelectOnStopped(&verifiedStatTbl, true, population.GetBftMajorityCount())
			break outer
		case d := <-c.queuePh3Recv:
			switch {
			case d.HasMissingMembers():
				if queueMissing == nil {
					queueMissing = make(chan inspectors.InspectedVector, len(c.queuePh3Recv))
					go c.workerRescanForMissing(ctx, queueMissing)
				}
				queueMissing <- d
				// do chasing
			case !d.IsInspected():
				d = d.Reinspect(ctx, localInspector)
				if !d.IsInspected() {
					if d.HasMissingMembers() {
						// loop it back to be picked by "case d.HasMissingMembers()"
						c.queuePh3Recv <- d
					}
					// TODO heavy inspection with hash recalculations should be running on a limited pool
					// go func() {
					d.Inspect(ctx)
					if !d.IsInspected() {
						inslogger.FromContext(ctx).Errorf("unable to inspect vector: %v", d)
						break
						// } else {
						//	c.queuePh3Recv <- d
					}
					// }()
					// break // do chasing
				}
				fallthrough
			default:
				inspectedNode := d.GetNode()
				nodeIndex := -1
				if inspectedNode.IsJoiner() {
					panic("not implemented")
				} else {
					nodeIndex = inspectedNode.GetIndex().AsInt()
				}

				nodeStats, vr := d.GetInspectionResults()

				if log.Is(insolar.DebugLevel) {
					popLimit, popSealed := population.GetSealedLimit()
					remains := popLimit - originalStatTbl.RowCount() - 1
					logMsg := "validated"
					switch {
					case nodeStats == nil:
						remains++
						fallthrough
					case d.HasSenderFault() || nodeStats == nil:
						logMsg = "fault"
					case nodeStats.HasAllValues(nodeset.ConsensusStatUnknown):
						logMsg = "missed"
					case !vr.AnyOf(nodeset.NvrTrustedValid | nodeset.NvrDoubtedValid):
						if vr.AnyOf(nodeset.NvrTrustedFraud|nodeset.NvrDoubtedFraud) || !c.R.IsJoiner() {
							logMsg = "failed"
						} else {
							logMsg = "received"
						}
					}
					completenessMark := ' '
					if !popSealed {
						completenessMark = '+'
					}

					na := d.GetNode()
					log.Debugf(
						"%s: idx:%d remains:%d%c\n Here(%04d):%v\nThere(%04d):%v\n     Result:%v\n Comparison:%v\n",
						//													    Compared
						logMsg, na.GetIndex(), remains, completenessMark, c.R.GetSelf().GetNodeID(), localInspector.GetBitset(),
						na.GetNodeID(), d.GetBitset(),
						nodeStats, d,
					)
				}

				if nodeStats != nil {
					currentRow, _ := verifiedStatTbl.GetRow(nodeIndex)
					if currentRow != nil && currentRow.GetCustomOptions() == outOfOrderPhase3 && nodeStats.GetCustomOptions() != outOfOrderPhase3 {
						// TODO do something more efficient
						originalStatTbl.RemoveRow(nodeIndex)
						verifiedStatTbl.RemoveRow(nodeIndex)

						if currentRow.HasAllValuesOf(nodeset.ConsensusStatTrusted, nodeset.ConsensusStatDoubted) {
							processedNodesFlawlessly--
						}
					}

					originalStat := nodeset.StateToConsensusStatRow(d.GetBitset())
					originalStatTbl.PutRow(nodeIndex, &originalStat)
					verifiedStatTbl.PutRow(nodeIndex, nodeStats)
					if nodeStats.HasAllValuesOf(nodeset.ConsensusStatTrusted, nodeset.ConsensusStatDoubted) {
						processedNodesFlawlessly++
					}
				} else {
					break
				}

				consensusSelection = c.consensusStrategy.TrySelectOnAdded(&verifiedStatTbl,
					d.GetNode().GetProfile().GetStatic(), nodeStats)

				// remainingNodes--

				// if vr.AnyOf(nodeset.NvrDoubtedAlteredNodeSet) {
				//	alteredDoubtedGshCount++
				// }
			}

			if consensusSelection != nil {
				if !consensusSelection.CanBeImproved() || !chasingDelayTimer.IsEnabled() {
					log.Debug("Phase3 done earlier by strategy")
					break outer
				}

				log.Debug("Phase3 (re)start chasing")
				chasingDelayTimer.RestartChase()
			}
		}
	}

	if log.Is(insolar.DebugLevel) {

		limit, sealed := population.GetSealedLimit()
		limitStr := ""
		if sealed {
			limitStr = fmt.Sprintf("%d", limit)
		} else {
			limitStr = fmt.Sprintf("%d+", limit)
		}
		tblHeader := fmt.Sprintf("%%sConsensus Node View (%%s): ID=%v Members=%d/%s Joiners=%d",
			c.R.GetSelfNodeID(), population.GetIndexedCount(), limitStr, population.GetJoinersCount())
		typeHeader := "Original, Verified"
		prev := ""
		if !originalStatTbl.EqualsTyped(&verifiedStatTbl) {
			prev = originalStatTbl.TableFmt(fmt.Sprintf(tblHeader, prev, "Original"), nodeset.FmtConsensusStat)
			typeHeader = "Verified"
		}
		log.Debug(verifiedStatTbl.TableFmt(fmt.Sprintf(tblHeader, prev, typeHeader), nodeset.FmtConsensusStat))
	}

	if consensusSelection == nil {
		panic("illegal state")
	}

	selectionSet := consensusSelection.GetConsensusVector()

	if selectionSet.HasValues(nodeset.CbsExcluded) {
		log.Infof("Consensus is finished as different, %v", selectionSet)
		// TODO update population and/or start Phase 4
	} else {
		log.Info("Consensus is finished as expected")
	}

	popRanks, csh, gsh := localInspector.CreateNextPopulation(selectionSet)
	c.R.BuildNextPopulation(ctx, popRanks, gsh, csh)

	return true
}
