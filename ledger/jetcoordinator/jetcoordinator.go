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
	"fmt"
	"sort"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/utils/entropy"
	"github.com/pkg/errors"
)

// JetCoordinator is responsible for all jet interactions
type JetCoordinator struct {
	db                         *storage.DB
	roleCounts                 map[core.DynamicRole]int
	NodeNet                    core.NodeNetwork                `inject:""`
	PlatformCryptographyScheme core.PlatformCryptographyScheme `inject:""`
}

// NewJetCoordinator creates new coordinator instance.
func NewJetCoordinator(db *storage.DB, conf configuration.JetCoordinator) *JetCoordinator {
	jc := JetCoordinator{db: db}
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

// Me returns current node.
func (jc *JetCoordinator) Me() core.RecordRef {
	return jc.NodeNet.GetOrigin().ID()
}

// IsAuthorized checks for role on concrete pulse for the address.
func (jc *JetCoordinator) IsAuthorized(
	ctx context.Context,
	role core.DynamicRole,
	obj core.RecordID,
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
	objID core.RecordID,
	pulse core.PulseNumber,
) ([]core.RecordRef, error) {
	switch role {
	case core.DynamicRoleVirtualExecutor:
		node, err := jc.VirtualExecutorForObject(ctx, objID, pulse)
		if err != nil {
			return nil, err
		}
		return []core.RecordRef{*node}, nil

	case core.DynamicRoleVirtualValidator:
		return jc.VirtualValidatorsForObject(ctx, objID, pulse)

	case core.DynamicRoleLightExecutor:
		node, err := jc.LightExecutorForObject(ctx, objID, pulse)
		if err != nil {
			return nil, err
		}
		return []core.RecordRef{*node}, nil

	case core.DynamicRoleLightValidator:
		return jc.LightValidatorsForObject(ctx, objID, pulse)

	case core.DynamicRoleHeavyExecutor:
		node, err := jc.Heavy(ctx, pulse)
		if err != nil {
			return nil, err
		}
		return []core.RecordRef{*node}, nil
	}

	panic("unexpected role")
}

func (jc *JetCoordinator) VirtualExecutorForObject(
	ctx context.Context, objID core.RecordID, pulse core.PulseNumber,
) (*core.RecordRef, error) {
	nodes, err := jc.virtualsForObject(ctx, objID, pulse, 1)
	if err != nil {
		return nil, err
	}
	return &nodes[0], nil
}

func (jc *JetCoordinator) VirtualValidatorsForObject(
	ctx context.Context, objID core.RecordID, pulse core.PulseNumber,
) ([]core.RecordRef, error) {
	count, ok := jc.roleCounts[core.DynamicRoleVirtualValidator]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no candidates for role %d", core.DynamicRoleVirtualValidator))
	}
	nodes, err := jc.virtualsForObject(ctx, objID, pulse, count+1)
	if err != nil {
		return nil, err
	}
	return nodes[1:], nil
}

func (jc *JetCoordinator) LightExecutorForJet(
	ctx context.Context, jetID core.RecordID, pulse core.PulseNumber,
) (*core.RecordRef, error) {
	nodes, err := jc.lightMaterialsForJet(ctx, jetID, pulse, 1)
	if err != nil {
		return nil, err
	}
	return &nodes[0], nil
}

func (jc *JetCoordinator) LightValidatorsForJet(
	ctx context.Context, jetID core.RecordID, pulse core.PulseNumber,
) ([]core.RecordRef, error) {
	count, ok := jc.roleCounts[core.DynamicRoleLightValidator]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no candidates for role %d", core.DynamicRoleLightValidator))
	}
	nodes, err := jc.lightMaterialsForJet(ctx, jetID, pulse, count+1)
	if err != nil {
		return nil, err
	}
	return nodes[1:], nil
}

func (jc *JetCoordinator) LightExecutorForObject(
	ctx context.Context, objID core.RecordID, pulse core.PulseNumber,
) (*core.RecordRef, error) {
	tree, err := jc.db.GetJetTree(ctx, pulse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch jet tree")
	}
	jetID, _ := tree.Find(objID)
	return jc.LightExecutorForJet(ctx, *jetID, pulse)
}

func (jc *JetCoordinator) LightValidatorsForObject(
	ctx context.Context, objID core.RecordID, pulse core.PulseNumber,
) ([]core.RecordRef, error) {
	tree, err := jc.db.GetJetTree(ctx, pulse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch jet tree")
	}
	jetID, _ := tree.Find(objID)
	return jc.LightValidatorsForJet(ctx, *jetID, pulse)
}

func (jc *JetCoordinator) Heavy(ctx context.Context, pulse core.PulseNumber) (*core.RecordRef, error) {
	candidates, err := jc.db.GetActiveNodesByRole(pulse, core.StaticRoleHeavyMaterial)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch active nodes for pulse %v", pulse)
	}
	if len(candidates) == 0 {
		return nil, errors.New(fmt.Sprintf("no active nodes for pulse %d", pulse))
	}
	pulseData, err := jc.db.GetPulse(ctx, pulse)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch pulse data for pulse %v", pulse)
	}
	nodes, err := getRefs(
		jc.PlatformCryptographyScheme,
		pulseData.Pulse.Entropy[:],
		candidates,
		1,
	)
	if err != nil {
		return nil, err
	}
	return &nodes[0], nil
}

func (jc *JetCoordinator) virtualsForObject(
	ctx context.Context, objID core.RecordID, pulse core.PulseNumber, count int,
) ([]core.RecordRef, error) {
	candidates, err := jc.db.GetActiveNodesByRole(pulse, core.StaticRoleVirtual)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch active nodes for pulse %v", pulse)
	}
	if len(candidates) == 0 {
		return nil, errors.New(fmt.Sprintf("no active nodes for pulse %d", pulse))
	}
	pulseData, err := jc.db.GetPulse(ctx, pulse)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch pulse data for pulse %v", pulse)
	}
	return getRefs(
		jc.PlatformCryptographyScheme,
		circleXOR(pulseData.Pulse.Entropy[:], objID.Hash()),
		candidates,
		count,
	)
}

func (jc *JetCoordinator) lightMaterialsForJet(
	ctx context.Context, jetID core.RecordID, pulse core.PulseNumber, count int,
) ([]core.RecordRef, error) {
	_, prefix := jet.Jet(jetID)
	candidates, err := jc.db.GetActiveNodesByRole(pulse, core.StaticRoleLightMaterial)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch active nodes for pulse %v", pulse)
	}
	if len(candidates) == 0 {
		return nil, errors.New(fmt.Sprintf("no active nodes for pulse %d", pulse))
	}
	pulseData, err := jc.db.GetPulse(ctx, pulse)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch pulse data for pulse %v", pulse)
	}
	return getRefs(
		jc.PlatformCryptographyScheme,
		circleXOR(pulseData.Pulse.Entropy[:], prefix),
		candidates,
		count,
	)
}

func getRefs(
	scheme core.PlatformCryptographyScheme,
	e []byte,
	values []core.Node,
	count int,
) ([]core.RecordRef, error) {
	// TODO: remove sort when network provides sorted result from GetActiveNodesByRole (INS-890) - @nordicdyno 5.Dec.2018
	sort.SliceStable(values, func(i, j int) bool {
		v1 := values[i].ID()
		v2 := values[j].ID()
		return bytes.Compare(v1[:], v2[:]) < 0
	})
	in := make([]interface{}, 0, len(values))
	for _, value := range values {
		in = append(in, interface{}(value.ID()))
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

// CircleXOR performs XOR for 'value' and 'src'. The result is returned as new byte slice.
// If 'value' is smaller than 'dst', XOR starts from the beginning of 'src'.
func circleXOR(value, src []byte) []byte {
	result := make([]byte, len(value))
	srcLen := len(src)
	for i := range result {
		result[i] = value[i] ^ src[i%srcLen]
	}
	return result
}
