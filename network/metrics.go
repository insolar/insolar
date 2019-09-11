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

package network

import (
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	// TagPhase is a tag for consensus metrics.
	TagPhase = insmetrics.MustTagKey("phase")
	TagRole  = insmetrics.MustTagKey("role")
)

var (
	// ConsensusPacketsSentBytes consensus sent packets size.
	ConsensusPacketsSent = stats.Int64("consensus_packets_sent", "Consensus sent packets size", stats.UnitBytes)
	// ConsensusPacketsRecvBytes consensus received packets size.
	ConsensusPacketsRecv = stats.Int64("consensus_packets_recv", "Consensus received packets size", stats.UnitBytes)
	// ConsensusPacketsRecvBytes consensus received packets size.
	ConsensusPacketsRecvBad = stats.Int64("consensus_packets_recv_bad", "Consensus received packets size", stats.UnitBytes)

	// DeclinedClaims consensus claims declined counter.
	DeclinedClaims = stats.Int64("consensus_claims_declined", "Consensus claims declined counter", stats.UnitDimensionless)
	// FailedCheckProof consensus validate proof fails.
	FailedCheckProof = stats.Int64("consensus_proof_failed", "Consensus validate proof fails", stats.UnitDimensionless)
	// ActiveNodes active nodes count after consensus.
	ActiveNodes = stats.Int64("consensus_active_nodes_count", "Active nodes count after consensus", stats.UnitDimensionless)
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
			Name:        DeclinedClaims.Name(),
			Description: DeclinedClaims.Description(),
			Measure:     DeclinedClaims,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        ConsensusPacketsSent.Name() + "_bytes",
			Description: ConsensusPacketsSent.Description(),
			Measure:     ConsensusPacketsSent,
			Aggregation: view.Sum(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        ConsensusPacketsSent.Name() + "_count",
			Description: ConsensusPacketsSent.Description(),
			Measure:     ConsensusPacketsSent,
			Aggregation: view.Count(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        ConsensusPacketsRecv.Name() + "_bytes",
			Description: ConsensusPacketsRecv.Description(),
			Measure:     ConsensusPacketsRecv,
			Aggregation: view.Sum(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        ConsensusPacketsRecv.Name() + "_count",
			Description: ConsensusPacketsRecv.Description(),
			Measure:     ConsensusPacketsRecv,
			Aggregation: view.Count(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        FailedCheckProof.Name(),
			Description: FailedCheckProof.Description(),
			Measure:     FailedCheckProof,
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
