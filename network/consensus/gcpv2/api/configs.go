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
