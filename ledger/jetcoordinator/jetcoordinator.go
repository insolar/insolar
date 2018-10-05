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
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/pkg/errors"
)

// JetCoordinator is responsible for all jet interactions
type JetCoordinator struct {
	db             *storage.DB
	rootJetNode    *JetNode
	roleCandidates map[core.JetRole][]core.RecordRef
	roleCounts     map[core.JetRole]int
}

// NewJetCoordinator creates new coordinator instance.
func NewJetCoordinator(db *storage.DB, conf configuration.JetCoordinator) (*JetCoordinator, error) {
	jc := JetCoordinator{
		db: db,
		rootJetNode: &JetNode{
			ref: core.RecordRef{},
			left: &JetNode{
				left:  &JetNode{ref: core.RecordRef{}},
				right: &JetNode{ref: core.RecordRef{}},
			},
			right: &JetNode{
				left:  &JetNode{ref: core.RecordRef{}},
				right: &JetNode{ref: core.RecordRef{}},
			},
		},
	}
	jc.loadConfig(conf)

	return &jc, nil
}

func (jc *JetCoordinator) loadConfig(conf configuration.JetCoordinator) {
	jc.roleCandidates = map[core.JetRole][]core.RecordRef{}
	jc.roleCounts = map[core.JetRole]int{}

	for intRole, candidates := range conf.RoleCandidates {
		role := core.JetRole(intRole)
		jc.roleCandidates[role] = []core.RecordRef{}
		for _, cand := range candidates {
			jc.roleCandidates[role] = append(jc.roleCandidates[role], core.NewRefFromBase58(cand))
		}
	}

	for intRole, count := range conf.RoleCounts {
		role := core.JetRole(intRole)
		jc.roleCounts[role] = count
	}
}

// IsAuthorized checks for role on concrete pulse for the address.
func (jc *JetCoordinator) IsAuthorized(
	role core.JetRole, obj core.RecordRef, pulse core.PulseNumber, node core.RecordRef,
) (bool, error) {
	nodes, err := jc.QueryRole(role, obj, pulse)
	if err != nil {
		return false, err
	}
	for _, n := range nodes {
		if n == node {
			return true, nil
		}
	}
	return false, nil
}

// QueryRole returns node refs responsible for role bound operations for given object and pulse.
func (jc *JetCoordinator) QueryRole(
	role core.JetRole, obj core.RecordRef, pulse core.PulseNumber,
) ([]core.RecordRef, error) {
	entropy, err := jc.db.GetEntropy(pulse)
	if err != nil {
		return nil, err
	}

	candidates, ok := jc.roleCandidates[role]
	if !ok {
		return nil, errors.New("no candidates for this role")
	}
	count, ok := jc.roleCounts[role]
	if !ok {
		return nil, errors.New("no candidate count for this role")
	}

	selected, err := selectByEntropy(*entropy, candidates, count)
	if err != nil {
		return nil, err
	}

	return selected, nil
}

// CreateDrop creates jet drop for provided pulse number.
func (jc *JetCoordinator) CreateDrop(pulse core.PulseNumber) (*jetdrop.JetDrop, error) {
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

func (jc *JetCoordinator) jetRef(objRef core.RecordRef) *core.RecordRef { // nolint: megacheck
	return jc.rootJetNode.GetContaining(&objRef)
}
