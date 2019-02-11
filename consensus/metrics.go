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
	TagPhase = insmetrics.MustTagKey("phase")
)

var (
	ConsensusPacketsSent      = stats.Int64("consensus/packets/sent", "Current consensus transport packets sent counter", stats.UnitDimensionless)
	ConsensusPacketsRecv      = stats.Int64("consensus/packets/recv", "Current consensus transport packets recv counter", stats.UnitDimensionless)
	ConsensusDeclinedClaims   = stats.Int64("consensus/claims/declined", "Consensus claims declined counter", stats.UnitDimensionless)
	ConsensusSentSize         = stats.Int64("consensus/packets/sent/bytes", "Consensus sent packets size", stats.UnitDimensionless)
	ConsensusRecvSize         = stats.Int64("consensus/packets/recv/bytes", "Consensus received packets size", stats.UnitDimensionless)
	ConsensusFailedCheckProof = stats.Int64("consensus/proof/failed", "Consensus validate proof fails", stats.UnitDimensionless)
	ConsensusPhase2TimedOuts  = stats.Int64("consensus/phase2/timedout", "Timed out nodes on phase 2", stats.UnitDimensionless)
	ConsensusPhase21Exec      = stats.Int64("consensus/phase21/exec", "Phase 21 execution counter", stats.UnitDimensionless)
	ConsensusPhase3Exec       = stats.Int64("consensus/phase3/exec", "Phase 3 execution counter", stats.UnitDimensionless)
	ConsensusActiveNodes      = stats.Int64("consensus/activenodes/count", "Active nodes count after consensus", stats.UnitDimensionless)
)

func init() {
	commontags := []tag.Key{TagPhase}
	err := view.Register(
		&view.View{
			Name:        ConsensusPacketsSent.Name(),
			Description: ConsensusPacketsSent.Description(),
			Measure:     ConsensusPacketsSent,
			Aggregation: view.Count(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        ConsensusPacketsRecv.Name(),
			Description: ConsensusPacketsRecv.Description(),
			Measure:     ConsensusPacketsRecv,
			Aggregation: view.Count(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        ConsensusDeclinedClaims.Name(),
			Description: ConsensusDeclinedClaims.Description(),
			Measure:     ConsensusDeclinedClaims,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        ConsensusSentSize.Name(),
			Description: ConsensusSentSize.Description(),
			Measure:     ConsensusSentSize,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        ConsensusRecvSize.Name(),
			Description: ConsensusRecvSize.Description(),
			Measure:     ConsensusRecvSize,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        ConsensusFailedCheckProof.Name(),
			Description: ConsensusFailedCheckProof.Description(),
			Measure:     ConsensusFailedCheckProof,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        ConsensusPhase2TimedOuts.Name(),
			Description: ConsensusPhase2TimedOuts.Description(),
			Measure:     ConsensusPhase2TimedOuts,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        ConsensusPhase3Exec.Name(),
			Description: ConsensusPhase3Exec.Description(),
			Measure:     ConsensusPhase3Exec,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        ConsensusPhase21Exec.Name(),
			Description: ConsensusPhase21Exec.Description(),
			Measure:     ConsensusPhase21Exec,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        ConsensusActiveNodes.Name(),
			Description: ConsensusActiveNodes.Description(),
			Measure:     ConsensusActiveNodes,
			Aggregation: view.LastValue(),
		},
	)
	if err != nil {
		panic(err)
	}
}
