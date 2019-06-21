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
	"time"

	"github.com/insolar/insolar/network/consensus/common"

	"github.com/insolar/insolar/network/consensus/gcpv2/nodeset"

	"github.com/insolar/insolar/network/consensus/gcpv2/packets"

	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"

	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/consensus/gcpv2/stats"
)

func NewPhase3Controller(packetPrepareOptions core.PacketSendOptions, queueTrustUpdated <-chan TrustUpdateSignal,
	consensusStrategy ConsensusSelectionStrategy) *Phase3Controller {
	return &Phase3Controller{
		packetPrepareOptions: packetPrepareOptions,
		queueTrustUpdated:    queueTrustUpdated,
		consensusStrategy:    consensusStrategy,
	}
}

type ConsensusSelection interface {
	/* When false - disables chasing timeout */
	CanBeImproved() bool
	IsSameWithActive() bool
	/* This bitset only allows values of NbsConsensus[*] */
	GetConsensusNodes() *nodeset.ConsensusBitsetRow
}

type ConsensusSelectionStrategy interface {
	/* Result can be nil - it means no-decision */
	TrySelectOnAdded(globulaStats *stats.StatTable, addedNode common2.NodeProfile,
		nodeStats *stats.Row, realm *core.FullRealm) ConsensusSelection
	SelectOnStopped(globulaStats *stats.StatTable, timeIsOut bool, realm *core.FullRealm) ConsensusSelection
}

var _ core.PhaseController = &Phase3Controller{}

type Phase3Controller struct {
	core.PhaseControllerPerMemberTemplate
	packetPrepareOptions core.PacketSendOptions
	queueTrustUpdated    <-chan TrustUpdateSignal
	queuePh3Recv         chan ph3Data
	consensusStrategy    ConsensusSelectionStrategy
	// packetHandler to Worker channel
}

type ph3Data struct {
	np     *core.NodeAppearance
	vector nodeset.HashedNodeVector
}

func (*Phase3Controller) GetPacketType() packets.PacketType {
	return packets.PacketPhase3
}

func (c *Phase3Controller) HandleMemberPacket(reader packets.MemberPacketReader, n *core.NodeAppearance) error {

	p3 := reader.AsPhase3Packet()

	err := n.SetReceivedWithDupCheck(c.GetPacketType())
	if err != nil {
		return err
	}
	bs := p3.GetBitset()
	gshTrusted := p3.GetTrustedGsh()
	gshDoubted := p3.GetDoubtedGsh()
	// TODO ClaimHashes as well

	c.queuePh3Recv <- ph3Data{
		np: n,
		vector: nodeset.HashedNodeVector{
			Bitset:        bs,
			DoubtedVector: gshDoubted,
			TrustedVector: gshTrusted,
		},
	}

	return nil
}

func (c *Phase3Controller) StartWorker(ctx context.Context) {
	c.queuePh3Recv = make(chan ph3Data, c.R.GetNodeCount())

	go c.workerPhase3(ctx)
}

func (c *Phase3Controller) workerPhase3(ctxRound context.Context) {

	ctx, cancel := context.WithDeadline(ctxRound, time.Now().Add(c.R.AdjustedAfter(c.R.GetTimings().EndOfPhase3)))
	defer cancel()

	if !c.workerPrePhase3(ctx) {
		return
	}
	d := c.calcGshPair()

	go c.workerSendPhase3(ctx, d)

	if !c.workerRecvPhase3(ctx, d) {
		//context was stopped in a hard way, we are dead in terms of consensus
		//TODO should wait for further packets to decide if we need to turn ourself into suspended state
		//c.R.StopRoundByTimeout()
	}

	go workerQueueFlusher(ctxRound, c.queuePh3Recv, c.queueTrustUpdated)
}

func workerQueueFlusher(ctxRound context.Context, q0 chan ph3Data, q1 <-chan TrustUpdateSignal) {
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

func (c *Phase3Controller) workerPrePhase3(ctx context.Context) bool {
	logger := c.R.Log()

	logger.Debug(">>>>workerPrePhase3: begin")

	timings := c.R.GetTimings()
	startOfPhase3 := time.After(c.R.AdjustedAfter(timings.EndOfPhase2))
	chasingDelayTimer := common.NewChasingTimer(timings.BeforeInPhase2ChasingDelay)

	var countFraud = 0
	var countHasNsh = 0
	var countTrustBySome = 0
	var countTrustByNeighbors = 0

outer:
	for {
		select {
		case <-ctx.Done():
			logger.Debug(">>>>workerPrePhase3: ctx.Done")
			return false // ctx.Err() ?
		case <-chasingDelayTimer.Channel():
			logger.Debug(">>>>workerPrePhase3: chaseExpired")
			break outer
		case <-startOfPhase3:
			logger.Debug(">>>>workerPrePhase3: startOfPhase3")
			break outer
		case sig := <-c.queueTrustUpdated:
			if sig.IsPingSignal() { // ping indicates arrival of Phase2 packet, to support chasing
				break
			}
			switch {
			case sig.NewTrustLevel < 0:
				countFraud++
				continue // no chasing delay on fraud
			case sig.NewTrustLevel == packets.UnknownTrust:
				countHasNsh++
				// if countHasNsh >= R.othersCount {
				// 	// we have answers from all
				// 	break outer
				// }
			case sig.NewTrustLevel >= packets.TrustByNeighbors:
				countTrustByNeighbors++
				fallthrough
			default:
				countTrustBySome++

				// We have some-trusted from all nodes, and the majority of them are well-trusted
				if countTrustBySome >= c.R.GetOthersCount() && countTrustByNeighbors >= c.R.GetBftMajorityCount() {
					logger.Debug(">>>>workerPrePhase3: all")
					break outer
				}

				if chasingDelayTimer.IsEnabled() {
					// We have answers from all nodes, and the majority of them are well-trusted
					if countHasNsh >= c.R.GetOthersCount() && countTrustByNeighbors >= c.R.GetBftMajorityCount() {
						chasingDelayTimer.RestartChase()
						logger.Debug(">>>>workerPrePhase3: chaseStartedAll")
					} else if countTrustBySome-countFraud >= c.R.GetBftMajorityCount() {
						// We can start chasing-timeout after getting answers from majority of some-trusted nodes
						chasingDelayTimer.RestartChase()
						logger.Debug(">>>>workerPrePhase3: chaseStartedSome")
					}
				}
			}
		}
	}

	/* Ensure that NSH is available, otherwise we can't normally build packets */
	for c.R.GetSelf().IsNshRequired() {
		select {
		case <-ctx.Done():
			logger.Debug(">>>>workerPrePhase3: ctx.Done")
			return false // ctx.Err() ?
		case <-c.queueTrustUpdated:
		case <-time.After(loopingMinimalDelay):
		}
	}
	return true
}

func (c *Phase3Controller) calcGshPair() nodeset.HashedNodeVector {

	/*
		NB! SequenceDigester requires at least one hash to be added. So to avoid errors, local node MUST always
		have trust level set high enough to get bitset[i].IsTrusted() == true
	*/
	aggTrusted := c.R.GetDigestFactory().GetGshDigester()
	var aggDoubted common.SequenceDigester

	bitset := make(nodeset.NodeBitset, c.R.GetNodeCount())
	nodes := c.R.GetIndexedNodes()
	for i := range nodes {
		membership, trust := nodes[i].GetNodeMembershipAndTrust()
		switch {
		case membership.IsEmpty():
			bitset[i] = nodeset.NbsTimeout
			continue
		case trust.IsNegative():
			bitset[i] = nodeset.NbsFraud
		case trust == packets.UnknownTrust:
			bitset[i] = nodeset.NbsBaselineTrust
		case trust < packets.TrustByNeighbors:
			bitset[i] = nodeset.NbsLimitedTrust
		default:
			bitset[i] = nodeset.NbsHighTrust
		}
		if bitset[i].IsTrusted() {
			aggTrusted.AddNext(membership.Nsh)
			if aggDoubted == nil {
				continue
			}
		} else if aggDoubted == nil {
			aggDoubted = aggTrusted.ForkSequence()
		}
		aggDoubted.AddNext(membership.Nsh)
	}

	res := nodeset.HashedNodeVector{Bitset: bitset}
	res.TrustedVector = aggTrusted.FinishSequence().AsDigestHolder()
	if aggDoubted != nil {
		res.DoubtedVector = aggDoubted.FinishSequence().AsDigestHolder()
	}
	return res
}

func (c *Phase3Controller) workerSendPhase3(ctx context.Context, selfData nodeset.HashedNodeVector) {

	p3 := c.R.GetPacketBuilder().PreparePhase3Packet(c.R.GetLocalProfile(), c.R.GetPulseData(),
		selfData.Bitset, selfData.TrustedVector, selfData.DoubtedVector, c.packetPrepareOptions)

	for _, np := range c.R.GetShuffledOtherNodes() {
		select {
		case <-ctx.Done():
			return
		default:
		}

		p3.SendTo(np.GetProfile(), 0, c.R.GetPacketSender())
		np.SetSentByPacketType(c.GetPacketType())
	}
}

func (c *Phase3Controller) workerRecvPhase3(ctx context.Context, selfData nodeset.HashedNodeVector) bool {
	logger := c.R.Log()

	timings := c.R.GetTimings()
	softDeadline := time.After(c.R.AdjustedAfter(timings.EndOfPhase3))
	chasingDelayTimer := common.NewChasingTimer(timings.BeforeInPhase3ChasingDelay)

	statTbl := nodeset.NewConsensusStatTable(c.R.GetNodeCount())
	statTbl.PutRow(c.R.GetSelf().GetIndex(), selfData.Bitset.LocalToConsensusStatRow())

	remainingNodes := c.R.GetOthersCount()

	// TODO detect nodes produced similar bitmaps, but different GSH
	// even if we wont have all NSH, we can let to know these nodes on such collision
	// bitsetMatcher := make(map[gcpv2.NodeBitset])

	alteredDoubtedGshCount := 0

	hasher := nodeset.NewFilteredSequenceHasher(c.R.GetDigestFactory(), c.handleNodeHashing)
	var consensusSelection ConsensusSelection

outer:
	for {
		select {
		case <-ctx.Done():
			logger.Debug("Phase3 cancelled")
			return false
		case <-softDeadline:
			logger.Debug("Phase3 deadline")
			consensusSelection = c.consensusStrategy.SelectOnStopped(&statTbl, true, c.R)
			break outer
		case <-chasingDelayTimer.Channel():
			logger.Debug("Phase3 chasing expired")
			consensusSelection = c.consensusStrategy.SelectOnStopped(&statTbl, true, c.R)
			break outer
		case d := <-c.queuePh3Recv:
			nodeStats := statTbl.NewRow()

			vr := nodeset.ClassifyByNodeGsh(selfData, d.vector, nodeStats, hasher)

			logMsg := "add"
			if nodeStats.HasAllValues(0) {
				logMsg = "missed"
			}

			logger.Debugf(
				"%s: s:%v, t:%v, %d, %d, %d: %v",
				logMsg,
				d.np.GetShortNodeID(),
				c.R.GetSelf().GetShortNodeID(),
				d.np.GetIndex(),
				remainingNodes,
				vr,
				nodeStats,
			)

			statTbl.PutRow(d.np.GetIndex(), nodeStats)
			remainingNodes--

			if vr.AnyOf(nodeset.NvrDoubtedAlteredNodeSet) {
				alteredDoubtedGshCount++
			}

			if remainingNodes <= 0 {
				consensusSelection = c.consensusStrategy.SelectOnStopped(&statTbl, false, c.R)
				logger.Debug("Phase3 done all")
				break outer
			} else {
				consensusSelection = c.consensusStrategy.TrySelectOnAdded(&statTbl, d.np.GetProfile(), nodeStats, c.R)
				if consensusSelection == nil {
					continue
				}
				if chasingDelayTimer.IsEnabled() && consensusSelection.CanBeImproved() {
					logger.Debug("Phase3 (re)start chasing")
					chasingDelayTimer.RestartChase()
				} else {
					logger.Debug("Phase3 done strategy")
					break outer
				}
			}
		}
	}

	// TODO do table generation only when it is needed for logging
	// if c.R.Log().IsInfo() {
	tblHeader := fmt.Sprintf("Consensus Node View: %v", c.R.GetSelfNodeID())
	logger.Debug(statTbl.TableFmt(tblHeader, nodeset.FmtConsensusStat))
	// }

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
		logger.Infof("Consensus is finished as same")
	} else {
		logger.Infof("Consensus is finished as different, %v", selectionSet)
		// TODO update population and/or start Phase 4
	}

	b := c.R.CreateNextPopulationBuilder()
	priming := c.R.GetPrimingCloudHash()
	//population is to be updated somewhere here
	b.SetGlobulaStateHash(priming)
	b.SealCensus()
	c.R.FinishRound(b, priming)

	return true
}

func (c *Phase3Controller) handleNodeHashing(index int, digester common.SequenceDigester) {
	nsh := c.R.GetNodeAppearanceByIndex(index).GetNodeStateHashEvidence().GetNodeStateHash()
	digester.AddNext(nsh)
}
