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

// +build networktest

package servicenetwork

import (
	"context"
	"time"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/consensus/phases"
	"github.com/insolar/insolar/insolar"
)

type CommunicationPolicy int

const (
	PartialPositive1Phase = CommunicationPolicy(iota + 1)
	PartialNegative1Phase
	PartialPositive2Phase
	PartialNegative2Phase
	PartialPositive3Phase
	PartialNegative3Phase
	PartialPositive23Phase
	PartialNegative23Phase
	FullTimeout
)

type CommunicatorMock struct {
	communicator phases.Communicator
	ignoreFrom   insolar.Reference
	policy       CommunicationPolicy
}

func (cm *CommunicatorMock) ExchangePhase1(
	ctx context.Context,
	originClaim *packets.NodeAnnounceClaim,
	participants []insolar.NetworkNode,
	packet *packets.Phase1Packet,
) (map[insolar.Reference]*packets.Phase1Packet, error) {
	pckts, err := cm.communicator.ExchangePhase1(ctx, originClaim, participants, packet)
	if err != nil {
		return nil, err
	}
	switch cm.policy {
	case PartialNegative1Phase, PartialPositive1Phase:
		delete(pckts, cm.ignoreFrom)
	}
	return pckts, nil
}

func (cm *CommunicatorMock) ExchangePhase2(ctx context.Context, state *phases.ConsensusState,
	participants []insolar.NetworkNode, packet *packets.Phase2Packet) (map[insolar.Reference]*packets.Phase2Packet, error) {

	pckts, err := cm.communicator.ExchangePhase2(ctx, state, participants, packet)
	if err != nil {
		return nil, err
	}
	switch cm.policy {
	case PartialPositive2Phase, PartialNegative2Phase, PartialPositive23Phase, PartialNegative23Phase:
		delete(pckts, cm.ignoreFrom)
	}
	return pckts, nil
}

func (cm *CommunicatorMock) ExchangePhase21(ctx context.Context, state *phases.ConsensusState,
	packet *packets.Phase2Packet, additionalRequests []*phases.AdditionalRequest) ([]packets.ReferendumVote, error) {

	return cm.communicator.ExchangePhase21(ctx, state, packet, additionalRequests)
}

func (cm *CommunicatorMock) ExchangePhase3(ctx context.Context, participants []insolar.NetworkNode, packet *packets.Phase3Packet) (map[insolar.Reference]*packets.Phase3Packet, error) {
	pckts, err := cm.communicator.ExchangePhase3(ctx, participants, packet)
	if err != nil {
		return nil, err
	}
	switch cm.policy {
	case PartialPositive3Phase, PartialNegative3Phase, PartialPositive23Phase, PartialNegative23Phase:
		delete(pckts, cm.ignoreFrom)
	}
	return pckts, nil
}

func (cm *CommunicatorMock) Init(ctx context.Context) error {
	return cm.communicator.Init(ctx)
}

type FullTimeoutPhaseManager struct {
}

func (ftpm *FullTimeoutPhaseManager) OnPulse(ctx context.Context, pulse *insolar.Pulse, pulseStartTime time.Time) error {
	return nil
}
