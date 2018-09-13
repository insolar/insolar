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

package jetcoordinator

import (
	"fmt"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
)

type mockHolder struct {
	virtualExecutor core.RecordRef
	lightExecutor   core.RecordRef
	heavyExecutor   core.RecordRef

	virtualValidators []core.RecordRef
	lightValidators   []core.RecordRef
}

// JetCoordinator is responsible for all jet interactions
type JetCoordinator struct {
	db          *storage.DB
	rootJetNode *JetNode

	mock *mockHolder // TODO: remove after actual implementation is ready
}

// NewJetCoordinator creates new coordinator instance.
func NewJetCoordinator(db *storage.DB, conf configuration.JetCoordinator) (*JetCoordinator, error) {
	mock, err := createMock(conf)
	if err != nil {
		return nil, err
	}

	rootJetNode := &JetNode{
		ref: core.RecordRef{},
		left: &JetNode{
			left:  &JetNode{ref: core.RecordRef{}},
			right: &JetNode{ref: core.RecordRef{}},
		},
		right: &JetNode{
			left:  &JetNode{ref: core.RecordRef{}},
			right: &JetNode{ref: core.RecordRef{}},
		},
	}

	return &JetCoordinator{
		db:          db,
		mock:        mock,
		rootJetNode: rootJetNode,
	}, nil
}

func createMock(conf configuration.JetCoordinator) (*mockHolder, error) {
	virtualExecutor := core.String2Ref(conf.VirtualExecutor)
	lightExecutor := core.String2Ref(conf.LightExecutor)
	heavyExecutor := core.String2Ref(conf.HeavyExecutor)

	virtualValidators := make([]core.RecordRef, len(conf.VirtualValidators))
	for i, vv := range conf.VirtualValidators {
		virtualValidators[i] = core.String2Ref(vv)
	}

	lightValidators := make([]core.RecordRef, len(conf.LightValidators))
	for i, lv := range conf.VirtualValidators {
		lightValidators[i] = core.String2Ref(lv)
	}

	return &mockHolder{
		virtualExecutor: virtualExecutor,
		lightExecutor:   lightExecutor,
		heavyExecutor:   heavyExecutor,

		virtualValidators: virtualValidators,
		lightValidators:   lightValidators,
	}, nil
}

func (jc *JetCoordinator) IsAuthorized(role core.JetRole, obj core.RecordRef, pulse core.PulseNumber, node core.RecordRef) bool {
	nodes := jc.QueryRole(role, obj, pulse)
	for _, n := range nodes {
		if n == node {
			return true
		}
	}
	return false
}

func (jc *JetCoordinator) QueryRole(role core.JetRole, obj core.RecordRef, pulse core.PulseNumber) []core.RecordRef {
	switch role {
	case core.RoleVirtualExecutor:
		return []core.RecordRef{jc.mock.virtualExecutor}
	case core.RoleLightExecutor:
		return []core.RecordRef{jc.mock.lightExecutor}
	case core.RoleHeavyExecutor:
		return []core.RecordRef{jc.mock.heavyExecutor}
	case core.RoleVirtualValidator:
		return jc.mock.virtualValidators
	case core.RoleLightValidator:
		return jc.mock.lightValidators
	default:
		panic("Unknown role")
	}
}

// Pulse creates new jet drop and ends current slot.
// This should be called when receiving a new pulse from pulsar.
func (jc *JetCoordinator) Pulse(new record.PulseNum) (*jetdrop.JetDrop, error) {
	current := jc.db.GetCurrentPulse()
	if new-current != 1 {
		panic(fmt.Sprintf("Wrong pulse, got %v, but current is %v\n", new, current))
	}

	// TODO: stop serving all requests (next node will be storage)

	drop, err := jc.createDrop(current)
	if err != nil {
		return nil, err
	}
	// nextExecutor, err := jc.getNextExecutor([][]byte{}) // TODO: fetch candidates from config
	// if err != nil {
	// 	return nil, err
	// }
	// nextValidators, err := jc.getNextValidators([][]byte{}, 3) // TODO: fetch candidates and count from config
	// if err != nil {
	// 	return nil, err
	// }

	// TODO: select next executor and validators. Send jet drop to current validators.

	jc.db.SetCurrentPulse(new)

	return drop, nil
}

// CreateDrop creates jet drop for provided pulse number.
func (jc *JetCoordinator) createDrop(pulse record.PulseNum) (*jetdrop.JetDrop, error) {
	prevDrop, err := jc.db.GetDrop(pulse - 1)
	if err != nil {
		return nil, err
	}
	newDrop, err := jc.db.SetDrop(pulse, prevDrop)
	if err != nil {
		return nil, err
	}
	return newDrop, nil
}

func (jc *JetCoordinator) getCurrentEntropy() ([]byte, error) { // nolint: megacheck
	return jc.db.GetEntropy(jc.db.GetCurrentPulse())
}

// TODO: real signature unknown
func (jc *JetCoordinator) getNextExecutor(candidates [][]byte) ([]byte, error) { // nolint: megacheck
	entropy, err := jc.getCurrentEntropy()
	if err != nil {
		return nil, err
	}
	idx, err := selectByEntropy(entropy, candidates, 1)
	if err != nil {
		return nil, err
	}

	return candidates[idx[0]], nil
}

// TODO: real signature unknown
func (jc *JetCoordinator) getNextValidators(candidates [][]byte, count int) ([][]byte, error) { // nolint: megacheck
	entropy, err := jc.getCurrentEntropy()
	if err != nil {
		return nil, err
	}
	idx, err := selectByEntropy(entropy, candidates, 1)
	if err != nil {
		return nil, err
	}
	selected := make([][]byte, 0, count)
	for _, i := range idx {
		selected = append(selected, candidates[i])
	}
	return selected, nil
}

func (jc *JetCoordinator) jetRef(objRef core.RecordRef) *core.RecordRef { // nolint: megacheck
	return jc.rootJetNode.GetContaining(&objRef)
}
