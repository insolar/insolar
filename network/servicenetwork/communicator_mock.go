/*
 *    Copyright 2018 Insolar
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

package servicenetwork

import (
	"context"

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
	case PartialNegative1Phase:
		fallthrough
	case PartialPositive1Phase:
		delete(pckts, cm.ignoreFrom)
	}
	return pckts, nil
}

func (cm *CommunicatorMock) ExchangePhase2(ctx context.Context, list network.UnsyncList, participants []core.Node, packet *packets.Phase2Packet) (map[core.RecordRef]*packets.Phase2Packet, error) {
	pckts, err := cm.communicator.ExchangePhase2(ctx, list, participants, packet)
	if err != nil {
		return nil, err
	}
	switch cm.testOpt {
	case PartialPositive2Phase:
		fallthrough
	case PartialNegative2Phase:
		delete(pckts, cm.ignoreFrom)
	}
	return pckts, nil
}

func (cm *CommunicatorMock) ExchangePhase21(ctx context.Context, list network.UnsyncList, packet *packets.Phase2Packet, additionalRequests []*phases.AdditionalRequest) ([]packets.ReferendumVote, error) {
	return cm.communicator.ExchangePhase21(ctx, list, packet, additionalRequests)
}

func (cm *CommunicatorMock) ExchangePhase3(ctx context.Context, participants []core.Node, packet *packets.Phase3Packet) (map[core.RecordRef]*packets.Phase3Packet, error) {
	return cm.communicator.ExchangePhase3(ctx, participants, packet)
}

func (cm *CommunicatorMock) Start(ctx context.Context) error {
	return cm.communicator.Start(ctx)
}
