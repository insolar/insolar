/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package servicenetwork

import (
	"context"

	"github.com/insolar/insolar/consensus/claimhandler"
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/consensus/phases"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
)

type CommunicatorTestOpt int

const (
	PartialPositive1Phase = CommunicatorTestOpt(iota + 1)
	PartialNegative1Phase
	PartialPositive2Phase
	PartialNegative2Phase
	PartialPositive3Phase
	PartialNegative3Phase
	PartialPositive23Phase
	PartialNegative23Phase
)

type CommunicatorMock struct {
	communicator phases.Communicator
	ignoreFrom   core.RecordRef
	testOpt      CommunicatorTestOpt
}

func (cm *CommunicatorMock) ExchangePhase1(
	ctx context.Context,
	originClaim *packets.NodeAnnounceClaim,
	participants []core.Node,
	packet *packets.Phase1Packet,
) (map[core.RecordRef]*packets.Phase1Packet, error) {
	pckts, err := cm.communicator.ExchangePhase1(ctx, originClaim, participants, packet)
	if err != nil {
		return nil, err
	}
	switch cm.testOpt {
	case PartialNegative1Phase, PartialPositive1Phase:
		delete(pckts, cm.ignoreFrom)
	}
	return pckts, nil
}

func (cm *CommunicatorMock) ExchangePhase2(ctx context.Context, list network.UnsyncList, handler *claimhandler.ClaimHandler,
	participants []core.Node, packet *packets.Phase2Packet) (map[core.RecordRef]*packets.Phase2Packet, error) {

	pckts, err := cm.communicator.ExchangePhase2(ctx, list, handler, participants, packet)
	if err != nil {
		return nil, err
	}
	switch cm.testOpt {
	case PartialPositive2Phase, PartialNegative2Phase, PartialPositive23Phase, PartialNegative23Phase:
		delete(pckts, cm.ignoreFrom)
	}
	return pckts, nil
}

func (cm *CommunicatorMock) ExchangePhase21(ctx context.Context, list network.UnsyncList, handler *claimhandler.ClaimHandler,
	packet *packets.Phase2Packet, additionalRequests []*phases.AdditionalRequest) ([]packets.ReferendumVote, error) {

	return cm.communicator.ExchangePhase21(ctx, list, handler, packet, additionalRequests)
}

func (cm *CommunicatorMock) ExchangePhase3(ctx context.Context, participants []core.Node, packet *packets.Phase3Packet) (map[core.RecordRef]*packets.Phase3Packet, error) {
	pckts, err := cm.communicator.ExchangePhase3(ctx, participants, packet)
	if err != nil {
		return nil, err
	}
	switch cm.testOpt {
	case PartialPositive3Phase, PartialNegative3Phase, PartialPositive23Phase, PartialNegative23Phase:
		delete(pckts, cm.ignoreFrom)
	}
	return pckts, nil
}

func (cm *CommunicatorMock) Init(ctx context.Context) error {
	return cm.communicator.Init(ctx)
}
