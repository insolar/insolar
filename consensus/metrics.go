/*
 *    Copyright 2019 Insolar
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

package consensus

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
