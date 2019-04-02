/*
 *    Copyright 2019 Insolar Technologies
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

	"github.com/pkg/errors"

	"github.com/insolar/insolar"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/ledger/storage/nodes"
	"github.com/insolar/insolar/utils/entropy"
)

// JetCoordinator is responsible for all jet interactions
type JetCoordinator struct {
	NodeNet                    core.NodeNetwork                `inject:""`
	PlatformCryptographyScheme core.PlatformCryptographyScheme `inject:""`
	PulseStorage               core.PulseStorage               `inject:""`
	JetStorage                 storage.JetStorage              `inject:""`
	PulseTracker               storage.PulseTracker            `inject:""`
	Nodes                      nodes.Accessor                  `inject:""`

	lightChainLimit int
}

// NewJetCoordinator creates new coordinator instance.
func NewJetCoordinator(lightChainLimit int) *JetCoordinator {
	return &JetCoordinator{lightChainLimit: lightChainLimit}
}

// Hardcoded roles count for validation and execution
const (
	VirtualValidatorCount  = 3
	MaterialValidatorCount = 3

	VirtualExecutorCount  = 1
	MaterialExecutorCount = 1
)

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
		if objID.Pulse() == core.PulseNumberJet {
			node, err := jc.LightExecutorForJet(ctx, objID, pulse)
			if err != nil {
				return nil, err
			}
			return []core.RecordRef{*node}, nil
		}
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

// VirtualExecutorForObject returns list of VEs for a provided pulse and objID
func (jc *JetCoordinator) VirtualExecutorForObject(
	ctx context.Context, objID core.RecordID, pulse core.PulseNumber,
) (*core.RecordRef, error) {
	nodes, err := jc.virtualsForObject(ctx, objID, pulse, VirtualExecutorCount)
	if err != nil {
		return nil, err
	}
	return &nodes[0], nil
}

// VirtualValidatorsForObject returns list of VVs for a provided pulse and objID
func (jc *JetCoordinator) VirtualValidatorsForObject(
	ctx context.Context, objID core.RecordID, pulse core.PulseNumber,
) ([]core.RecordRef, error) {
	nodes, err := jc.virtualsForObject(ctx, objID, pulse, VirtualValidatorCount+VirtualExecutorCount)
	if err != nil {
		return nil, err
	}
	// Skipping `VirtualExecutorCount` for validators
	// because it will be selected as the executor(s) for the same pulse.
	return nodes[VirtualExecutorCount:], nil
}

// LightExecutorForJet returns list of LEs for a provided pulse and jetID
func (jc *JetCoordinator) LightExecutorForJet(
	ctx context.Context, jetID core.RecordID, pulse core.PulseNumber,
) (*core.RecordRef, error) {
	nodes, err := jc.lightMaterialsForJet(ctx, jetID, pulse, MaterialExecutorCount)
	if err != nil {
		return nil, err
	}
	inslogger.FromContext(ctx).Debug(
		"jet miss: node ", nodes[0].String(),
		" is LME for ", jetID.DebugString(),
		" in pulse ", pulse,
	)
	return &nodes[0], nil
}

// LightValidatorsForJet returns list of LVs for a provided pulse and jetID
func (jc *JetCoordinator) LightValidatorsForJet(
	ctx context.Context, jetID core.RecordID, pulse core.PulseNumber,
) ([]core.RecordRef, error) {
	nodes, err := jc.lightMaterialsForJet(ctx, jetID, pulse, MaterialValidatorCount+MaterialExecutorCount)
	if err != nil {
		return nil, err
	}
	// Skipping `MaterialExecutorCount` for validators
	// because it will be selected as the executor(s) for the same pulse.
	return nodes[MaterialExecutorCount:], nil
}

// LightExecutorForObject returns list of LEs for a provided pulse and objID
func (jc *JetCoordinator) LightExecutorForObject(
	ctx context.Context, objID core.RecordID, pulse core.PulseNumber,
) (*core.RecordRef, error) {
	jetID, _ := jc.JetStorage.FindJet(ctx, pulse, objID)
	return jc.LightExecutorForJet(ctx, *jetID, pulse)
}

// LightValidatorsForObject returns list of LVs for a provided pulse and objID
func (jc *JetCoordinator) LightValidatorsForObject(
	ctx context.Context, objID core.RecordID, pulse core.PulseNumber,
) ([]core.RecordRef, error) {
	jetID, _ := jc.JetStorage.FindJet(ctx, pulse, objID)
	return jc.LightValidatorsForJet(ctx, *jetID, pulse)
}

// Heavy returns *core.RecorRef to a heavy of specific pulse
func (jc *JetCoordinator) Heavy(ctx context.Context, pulse core.PulseNumber) (*core.RecordRef, error) {
	candidates, err := jc.Nodes.InRole(pulse, core.StaticRoleHeavyMaterial)
	if err == core.ErrNoNodes {
		return nil, err
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch active heavy nodes for pulse %v", pulse)
	}
	if len(candidates) == 0 {
		return nil, errors.New(fmt.Sprintf("no active heavy nodes for pulse %d", pulse))
	}
	ent, err := jc.entropy(ctx, pulse)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch entropy for pulse %v", pulse)
	}

	refs, err := getRefs(
		jc.PlatformCryptographyScheme,
		ent[:],
		candidates,
		1,
	)
	if err != nil {
		return nil, err
	}
	return &refs[0], nil
}

// IsBeyondLimit calculates if target pulse is behind clean-up limit
// or if currentPN|targetPN didn't found in in-memory pulse-storage.
func (jc *JetCoordinator) IsBeyondLimit(ctx context.Context, currentPN, targetPN core.PulseNumber) (bool, error) {
	currentPulse, err := jc.PulseTracker.GetPulse(ctx, currentPN)
	if err == core.ErrNotFound {
		return true, nil
	}
	if err != nil {
		return false, errors.Wrapf(err, "failed to fetch pulse %v", currentPN)
	}

	targetPulse, err := jc.PulseTracker.GetPulse(ctx, targetPN)
	if err == core.ErrNotFound {
		return true, nil
	}
	if err != nil {
		return false, errors.Wrapf(err, "failed to fetch pulse %v", targetPN)
	}

	if currentPulse.SerialNumber-targetPulse.SerialNumber < jc.lightChainLimit {
		return false, nil
	}

	return true, nil
}

// NodeForJet calculates a node (LME or heavy) for a specific jet for a specific pulseNumber
func (jc *JetCoordinator) NodeForJet(ctx context.Context, jetID core.RecordID, rootPN, targetPN core.PulseNumber) (*core.RecordRef, error) {
	toHeavy, err := jc.IsBeyondLimit(ctx, rootPN, targetPN)
	if err != nil {
		return nil, err
	}

	if toHeavy {
		return jc.Heavy(ctx, rootPN)
	}
	return jc.LightExecutorForJet(ctx, jetID, targetPN)
}

// NodeForObject calculates a node (LME or heavy) for a specific jet for a specific pulseNumber
func (jc *JetCoordinator) NodeForObject(ctx context.Context, objectID core.RecordID, rootPN, targetPN core.PulseNumber) (*core.RecordRef, error) {
	toHeavy, err := jc.IsBeyondLimit(ctx, rootPN, targetPN)
	if err != nil {
		return nil, err
	}

	if toHeavy {
		return jc.Heavy(ctx, rootPN)
	}
	return jc.LightExecutorForObject(ctx, objectID, targetPN)
}

func (jc *JetCoordinator) virtualsForObject(
	ctx context.Context, objID core.RecordID, pulse core.PulseNumber, count int,
) ([]core.RecordRef, error) {
	candidates, err := jc.Nodes.InRole(pulse, core.StaticRoleVirtual)
	if err == core.ErrNoNodes {
		return nil, err
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch active virtual nodes for pulse %v", pulse)
	}
	if len(candidates) == 0 {
		return nil, errors.New(fmt.Sprintf("no active virtual nodes for pulse %d", pulse))
	}

	ent, err := jc.entropy(ctx, pulse)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch entropy for pulse %v", pulse)
	}

	return getRefs(
		jc.PlatformCryptographyScheme,
		circleXOR(ent[:], objID.Hash()),
		candidates,
		count,
	)
}

func (jc *JetCoordinator) lightMaterialsForJet(
	ctx context.Context, jetID core.RecordID, pulse core.PulseNumber, count int,
) ([]core.RecordRef, error) {
	_, prefix := jet.Jet(jetID)

	candidates, err := jc.Nodes.InRole(pulse, core.StaticRoleLightMaterial)
	if err == core.ErrNoNodes {
		return nil, err
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch active light nodes for pulse %v", pulse)
	}
	if len(candidates) == 0 {
		return nil, core.ErrNoNodes
	}

	ent, err := jc.entropy(ctx, pulse)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch entropy for pulse %v", pulse)
	}

	return getRefs(
		jc.PlatformCryptographyScheme,
		circleXOR(ent[:], prefix),
		candidates,
		count,
	)
}

func (jc *JetCoordinator) entropy(ctx context.Context, pulse core.PulseNumber) (core.Entropy, error) {
	current, err := jc.PulseStorage.Current(ctx)
	if err != nil {
		return core.Entropy{}, errors.Wrap(err, "failed to get current pulse")
	}

	if current.PulseNumber == pulse {
		return current.Entropy, nil
	}

	older, err := jc.PulseTracker.GetPulse(ctx, pulse)
	if err != nil {
		return core.Entropy{}, errors.Wrapf(err, "failed to fetch pulse data for pulse %v", pulse)
	}

	return older.Pulse.Entropy, nil
}

func getRefs(
	scheme core.PlatformCryptographyScheme,
	e []byte,
	values []insolar.Node,
	count int,
) ([]core.RecordRef, error) {
	// TODO: remove sort when network provides sorted result from GetActiveNodesByRole (INS-890) - @nordicdyno 5.Dec.2018
	sort.SliceStable(values, func(i, j int) bool {
		v1 := values[i].ID
		v2 := values[j].ID
		return bytes.Compare(v1[:], v2[:]) < 0
	})
	in := make([]interface{}, 0, len(values))
	for _, value := range values {
		in = append(in, interface{}(value.ID))
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
