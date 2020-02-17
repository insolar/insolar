// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
)

var (
	// ConsensusPacketsSent consensus sent packets size.
	ConsensusPacketsSent = stats.Int64("consensus_packets_sent", "Consensus sent packets size", stats.UnitBytes)
	// ConsensusPacketsRecv consensus received packets size.
	ConsensusPacketsRecv = stats.Int64("consensus_packets_recv", "Consensus received packets size", stats.UnitBytes)
	// ConsensusPacketsRecvBad consensus received packets size.
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
			Measure:     ConsensusPacketsSent,
			Aggregation: view.Sum(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        ConsensusPacketsSent.Name() + "_count",
			Measure:     ConsensusPacketsSent,
			Aggregation: view.Count(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        ConsensusPacketsRecv.Name() + "_bytes",
			Measure:     ConsensusPacketsRecv,
			Aggregation: view.Sum(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        ConsensusPacketsRecv.Name() + "_count",
			Measure:     ConsensusPacketsRecv,
			Aggregation: view.Count(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        FailedCheckProof.Name(),
			Measure:     FailedCheckProof,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        ActiveNodes.Name(),
			Measure:     ActiveNodes,
			Aggregation: view.LastValue(),
		},
	)
	if err != nil {
		panic(err)
	}
}
