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

package consensusv1

import (
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	// TagPhase is a tag for consensus metrics.
	TagPhase = insmetrics.MustTagKey("phase")
)

var (
	// PacketsSent urrent consensus transport packets sent counter.
	PacketsSent = stats.Int64("consensus/packets/sent", "Current consensus transport packets sent counter", stats.UnitDimensionless)
	// PacketsRecv current consensus transport packets recv counter.
	PacketsRecv = stats.Int64("consensus/packets/recv", "Current consensus transport packets recv counter", stats.UnitDimensionless)
	// DeclinedClaims consensus claims declined counter.
	DeclinedClaims = stats.Int64("consensus/claims/declined", "Consensus claims declined counter", stats.UnitDimensionless)
	// SentSize consensus sent packets size.
	SentSize = stats.Int64("consensus/packets/sent/bytes", "Consensus sent packets size", stats.UnitDimensionless)
	// RecvSize consensus received packets size.
	RecvSize = stats.Int64("consensus/packets/recv/bytes", "Consensus received packets size", stats.UnitDimensionless)
	// FailedCheckProof consensus validate proof fails.
	FailedCheckProof = stats.Int64("consensus/proof/failed", "Consensus validate proof fails", stats.UnitDimensionless)
	// Phase2TimedOuts timed out nodes on phase 2.
	Phase2TimedOuts = stats.Int64("consensus/phase2/timedout", "Timed out nodes on phase 2", stats.UnitDimensionless)
	// Phase21Exec phase 21 execution counter.
	Phase21Exec = stats.Int64("consensus/phase21/exec", "Phase 21 execution counter", stats.UnitDimensionless)
	// Phase3Exec phase 3 execution counter
	Phase3Exec = stats.Int64("consensus/phase3/exec", "Phase 3 execution counter", stats.UnitDimensionless)
	// ActiveNodes active nodes count after consensus.
	ActiveNodes = stats.Int64("consensus/activenodes/count", "Active nodes count after consensus", stats.UnitDimensionless)
)

func init() {
	commontags := []tag.Key{TagPhase}
	err := view.Register(
		&view.View{
			Name:        PacketsSent.Name(),
			Description: PacketsSent.Description(),
			Measure:     PacketsSent,
			Aggregation: view.Count(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        PacketsRecv.Name(),
			Description: PacketsRecv.Description(),
			Measure:     PacketsRecv,
			Aggregation: view.Count(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        DeclinedClaims.Name(),
			Description: DeclinedClaims.Description(),
			Measure:     DeclinedClaims,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        SentSize.Name(),
			Description: SentSize.Description(),
			Measure:     SentSize,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        RecvSize.Name(),
			Description: RecvSize.Description(),
			Measure:     RecvSize,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        FailedCheckProof.Name(),
			Description: FailedCheckProof.Description(),
			Measure:     FailedCheckProof,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        Phase2TimedOuts.Name(),
			Description: Phase2TimedOuts.Description(),
			Measure:     Phase2TimedOuts,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        Phase3Exec.Name(),
			Description: Phase3Exec.Description(),
			Measure:     Phase3Exec,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        Phase21Exec.Name(),
			Description: Phase21Exec.Description(),
			Measure:     Phase21Exec,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        ActiveNodes.Name(),
			Description: ActiveNodes.Description(),
			Measure:     ActiveNodes,
			Aggregation: view.LastValue(),
		},
	)
	if err != nil {
		panic(err)
	}
}
