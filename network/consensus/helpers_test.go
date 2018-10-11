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

package consensus

import (
	"context"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
)

func newActiveNode(ref byte, pulse int) *core.ActiveNode {
	return &core.ActiveNode{
		NodeID:    core.RecordRef{ref},
		PulseNum:  core.PulseNumber(pulse),
		State:     core.NodeActive,
		Role:      core.RoleUnknown,
		PublicKey: []byte{0, 0, 0},
	}
}

type TestNode struct {
	self            Participant
	allParticipants []Participant
	consensus       Consensus
	ctx             context.Context
}

type TestParticipant struct {
	node   *core.ActiveNode
	holder mockUnsyncHolder
}

func NewParticipant(ref byte, list []*core.ActiveNode) *TestParticipant {
	return &TestParticipant{node: newActiveNode(ref, 0),
		holder: mockUnsyncHolder{list}}
}

func (p *TestParticipant) GetID() core.RecordRef {
	return p.node.NodeID
}

func (p *TestParticipant) GetActiveNode() *core.ActiveNode {
	return p.node
}

func (m *TestParticipant) GetUnsync() []*core.ActiveNode {
	return m.holder.GetUnsync()
}

func (m *TestParticipant) GetPulse() core.PulseNumber {
	return m.holder.GetPulse()
}

func (m *TestParticipant) SetHash([]*NodeUnsyncHash) {
}

func (TestParticipant) GetHash(blockTimeout time.Duration) ([]*NodeUnsyncHash, error) {
	return nil, nil
}

// =====

type mockUnsyncHolder struct {
	list []*core.ActiveNode
}

func (m *mockUnsyncHolder) GetUnsync() []*core.ActiveNode {
	return m.list
}

func (mockUnsyncHolder) GetPulse() core.PulseNumber {
	return 0
}

func (mockUnsyncHolder) SetHash([]*NodeUnsyncHash) {
}

func (mockUnsyncHolder) GetHash(blockTimeout time.Duration) ([]*NodeUnsyncHash, error) {
	return nil, nil
}

type testCommunicator struct {
	self Participant
}

func (c *testCommunicator) ExchangeData(ctx context.Context, pulse core.PulseNumber, p Participant, data []*core.ActiveNode) ([]*core.ActiveNode, error) {
	log.Infof("returns data: %v", data)
	tp := p.(*TestParticipant)
	return tp.holder.GetUnsync(), nil
}

func (c *testCommunicator) ExchangeHash(ctx context.Context, pulse core.PulseNumber, p Participant, data []*NodeUnsyncHash) ([]*NodeUnsyncHash, error) {
	return []*NodeUnsyncHash{}, nil
}
