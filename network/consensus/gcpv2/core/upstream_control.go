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

package core

import (
	"time"

	"github.com/insolar/insolar/network/consensus/gcpv2/census"

	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"

	"github.com/insolar/insolar/network/consensus/common"
)

type UpstreamPulseController interface {
	/* Called when pulse is expected soon.
	Application traffic should be throttled down a bit.
	*/
	PulseIsComing(anticipatedStart time.Time)

	/* Called on receiving seem-to-be-valid Pulsar or Phase0 packets.
	Application traffic should be stopped or throttled down severely for a limited time (1-2 secs).
	Restoration of application traffic should be done automatically, unless PulseDetected() is called again.
	Can be called multiple time in sequence.

	Application MUST NOT consider it as a new pulse.
	*/
	PulseDetected()

	/* Called on a valid Pulse, but the pulse can yet be rolled back. No additional implications on traffic.
	Application should return immediately and start preparation of NodeStateHash.
	NodeStateHash should be sent into the channel when ready.
	*/
	PreparePulseChange(report MembershipUpstreamReport) <-chan common2.NodeStateHash

	/* Called on a confirmed Pulse and indicates final change of Pulse for the application.
	Application traffic can be resumed, but should remain throttled.
	*/
	CommitPulseChange(report MembershipUpstreamReport, activeCensus census.OperationalCensus)

	/* Called on a rollback of Pulse and indicates continuation of the previous Pulse for the application.
	Application traffic can be resumed at full.
	*/
	CancelPulseChange()

	/* Consensus is finished and population for the next pulse is finalized
	Application traffic can be resumed at full
	*/
	ConsensusCompleted(report MembershipUpstreamReport, expectedCensus census.OperationalCensus)

	/*
		This node was unable to prove its identity, but it may not yet be excluded (is suspected).
		All new Operations should be suspended.
	*/
	ConsensusFailed()

	/* Node has joined a globula */
	MembershipAcquired()

	/* Node was expelled from a globula */
	MembershipLost()

	/* Application traffic should be stopped or throttled down severely for a limited time (1-2 secs). */
	SuspendTraffic()
	/* Application traffic can be resumed at full */
	ResumeTraffic()
}

type MembershipUpstreamReport struct {
	PulseNumber     common.PulseNumber
	MembershipState common2.MembershipState
	MemberPower     common2.MemberPower
}
