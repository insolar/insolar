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
	"bytes"
	"context"
	"sort"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/utils/entropy"
	"github.com/pkg/errors"
)

// JetCoordinator is responsible for all jet interactions
type JetCoordinator struct {
	db                         *storage.DB
	rootJetNode                *JetNode
	roleCounts                 map[core.DynamicRole]int
	NodeNet                    core.NodeNetwork                `inject:""`
	PlatformCryptographyScheme core.PlatformCryptographyScheme `inject:""`
}

// NewJetCoordinator creates new coordinator instance.
func NewJetCoordinator(db *storage.DB, conf configuration.JetCoordinator) *JetCoordinator {
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

	return &jc
}

func (jc *JetCoordinator) loadConfig(conf configuration.JetCoordinator) {
	jc.roleCounts = map[core.DynamicRole]int{}

	for intRole, count := range conf.RoleCounts {
		role := core.DynamicRole(intRole)
		jc.roleCounts[role] = count
	}
}

// IsAuthorized checks for role on concrete pulse for the address.
func (jc *JetCoordinator) IsAuthorized(
	ctx context.Context,
	role core.DynamicRole,
	obj *core.RecordRef,
	pulse core.PulseNumber,
	node core.RecordRef,
) (bool, error) {
	nodes, err := jc.QueryRole(ctx, role, obj, pulse)
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
	ctx context.Context,
	role core.DynamicRole,
	obj *core.RecordRef,
	pulse core.PulseNumber,
) ([]core.RecordRef, error) {
	pulseData, err := jc.db.GetPulse(ctx, pulse)
	if err != nil {
		return nil, err
	}
	candidates := jc.NodeNet.GetActiveNodesByRole(role)
	if len(candidates) == 0 {
		return nil, errors.New("no candidates for this role")
	}

	if obj == nil {
		count, ok := jc.roleCounts[role]
		if !ok {
			return nil, errors.New("no candidate count for this role")
		}
		return refsByEntropy(jc.PlatformCryptographyScheme, pulseData.Pulse.Entropy[:], candidates, count)
	}

	h := jc.PlatformCryptographyScheme.ReferenceHasher()
	_, err = h.Write(pulseData.Pulse.Entropy[:])
	if err != nil {
		return nil, err
	}
	_, err = h.Write(obj[:])
	if err != nil {
		return nil, err
	}
	objEntropy := h.Sum(nil)

	if role == core.DynamicRoleLightExecutor {

	}

	return refsByEntropy(jc.PlatformCryptographyScheme, objEntropy, candidates, 1)
}

func (jc *JetCoordinator) jetRef(objRef core.RecordRef) *core.RecordRef { // nolint: megacheck
	return jc.rootJetNode.GetContaining(&objRef)
}

// GetActiveNodes return active nodes for specified pulse.
func (jc *JetCoordinator) GetActiveNodes(pulse core.PulseNumber) ([]core.Node, error) {
	return jc.db.GetActiveNodes(pulse)
}

func (jc *JetCoordinator) getNodeViaJet(
	ctx context.Context, pulse core.Pulse, ent []byte, candidates []core.RecordRef,
) (*core.RecordRef, error) {
	jetTree, err := jc.db.GetJetTree(ctx, pulse.PulseNumber)

	// Select a jet via entropy.
	// TODO: select jet via binary tree.
	in := make([]interface{}, 0, len(jetTree.Jets))
	for _, jet := range jetTree.Jets {
		in = append(in, interface{}(jet.ID))
	}
	selectedJets, err := entropy.SelectByEntropy(jc.PlatformCryptographyScheme, ent, in, 1)
	if err != nil {
		return nil, err
	}
	selectedJet := selectedJets[0].(core.RecordID)

	// Select a node via jet.
	h := jc.PlatformCryptographyScheme.ReferenceHasher()
	_, err = h.Write(ent)
	if err != nil {
		return nil, err
	}
	_, err = h.Write(selectedJet[:])
	if err != nil {
		return nil, err
	}
	nodes, err := refsByEntropy(jc.PlatformCryptographyScheme, h.Sum(nil), candidates, 1)
	if err != nil {
		return nil, err
	}

	return &nodes[0], nil
}

func refsByEntropy(
	scheme core.PlatformCryptographyScheme,
	e []byte,
	values []core.RecordRef,
	count int,
) ([]core.RecordRef, error) {
	// TODO: remove sort when network provides sorted result from GetActiveNodesByRole (INS-890) - @nordicdyno 5.Dec.2018
	sort.SliceStable(values, func(i, j int) bool {
		return bytes.Compare(values[i][:], values[j][:]) < 0
	})
	in := make([]interface{}, 0, len(values))
	for _, value := range values {
		in = append(in, interface{}(value))
	}

	res, err := entropy.SelectByEntropy(scheme, e, in, count)
	if err != nil {
		return nil, err
	}
	out := make([]core.RecordRef, 0, len(res))
	for _, value := range res {
		out = append(out, value.(core.RecordRef))
	}
	return out, nil
}
