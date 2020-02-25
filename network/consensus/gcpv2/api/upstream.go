// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package api

import (
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/pulse"
)

type UpstreamController interface {
	/* Called on a valid Pulse, but the pulse can yet be rolled back.
	Application should return immediately and start preparation of NodeState hash.
	NodeState should be sent into the channel when ready, but the channel will not be read if CancelPulseChange() has happened.

	The provided channel is guaranteed to have a buffer for one element.
	*/
	PreparePulseChange(report UpstreamReport, ch chan<- UpstreamState)

	/* Called on a confirmed Pulse and indicates final change of Pulse for the application.	*/
	CommitPulseChange(report UpstreamReport, pd pulse.Data, activeCensus census.Operational)

	/* Called on a rollback of Pulse and indicates continuation of the previous Pulse for the application. */
	CancelPulseChange()

	/* Consensus is finished. If expectedCensus == nil then this node was evicted from consensus.	*/
	ConsensusFinished(report UpstreamReport, expectedCensus census.Operational)

	/* Consensus was stopped abnormally */
	ConsensusAborted()
}

type UpstreamReport struct {
	PulseNumber pulse.Number
	MemberPower member.Power
	MemberMode  member.OpMode
	IsJoiner    bool
}

type UpstreamState struct {
	NodeState proofs.NodeStateHash
	// TODO ClaimFeeder
}
