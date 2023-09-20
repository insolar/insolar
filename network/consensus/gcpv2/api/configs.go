package api

import (
	"fmt"
	"time"
)

type RoundTimings struct {
	// Time to wait since NSH is requested before starting Phase0.
	StartPhase0At time.Duration

	// When Phase2 can be finished sooner by number of covered nodes, termination of Phase2 will be delayed
	// by BeforeInPhase2ChasingDelay after every Phase2 packet. No extra delays when = 0.
	// Total Phase2 time can NOT exceed EndOfPhase2
	BeforeInPhase2ChasingDelay time.Duration

	// When Phase3 can be finished sooner by number of covered nodes, termination of Phase2 will be delayed
	// by BeforeInPhase3ChasingDelay after every Phase3 packet. No extra delays when = 0.
	// Total Phase3 time can NOT exceed EndOfPhase3
	BeforeInPhase3ChasingDelay time.Duration

	// Time to finish receiving of Phase1 packets from other nodes and to finish producing Phase2 packets as well
	// since start of the consensus round
	EndOfPhase1 time.Duration

	// Time to wait before re-sending Phase1 packets (marked as requests) to missing nodes
	// since start of the consensus round. No retries when = 0
	StartPhase1RetryAt time.Duration

	// Time to finish receiving Phase2 packets from other nodes and START producing Phase3 packets
	// Phase3 can start sooner if there is enough number of nodes covered by Phase2
	EndOfPhase2 time.Duration

	// Time to finish receiving Phase3 packets from other nodes
	EndOfPhase3 time.Duration

	// Hard stop for all consensus operations
	EndOfConsensus time.Duration
}

func (t RoundTimings) String() string {
	return fmt.Sprintf("EndOfConsensus: %s, EndOfPhase1: %s, EndOfPhase2: %s, EndOfPhase3: %s",
		t.EndOfConsensus,
		t.EndOfPhase1,
		t.EndOfPhase2,
		t.EndOfPhase3,
	)
}
