//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package jetcoordinator

import (
	"context"
	"fmt"
	"sort"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/pulse"
	"github.com/insolar/insolar/utils/entropy"
)

// Coordinator is responsible for all jet interactions
type Coordinator struct {
	OriginProvider             network.OriginProvider             `inject:""`
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`

	PulseAccessor   insolarPulse.Accessor   `inject:""`
	PulseCalculator insolarPulse.Calculator `inject:""`

	JetAccessor jet.Accessor  `inject:""`
	Nodes       node.Accessor `inject:""`

	lightChainLimit int
}

// NewJetCoordinator creates new coordinator instance.
func NewJetCoordinator(lightChainLimit int) *Coordinator {
	return &Coordinator{lightChainLimit: lightChainLimit}
}

// Hardcoded roles count for validation and execution
const (
	VirtualValidatorCount  = 3
	MaterialValidatorCount = 3

	VirtualExecutorCount  = 1
	MaterialExecutorCount = 1
)

// Me returns current node.
func (jc *Coordinator) Me() insolar.Reference {
	return jc.OriginProvider.GetOrigin().ID()
}

// IsAuthorized checks for role on concrete pulse for the address.
func (jc *Coordinator) IsAuthorized(
	ctx context.Context,
	role insolar.DynamicRole,
	obj insolar.ID,
	pulse insolar.PulseNumber,
	node insolar.Reference,
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

// IsMeAuthorizedNow checks role of the current node in the current pulse for the address.
// Wrapper around IsAuthorized.
func (jc *Coordinator) IsMeAuthorizedNow(
	ctx context.Context,
	role insolar.DynamicRole,
	obj insolar.ID,
) (bool, error) {
	p, err := jc.PulseAccessor.Latest(ctx)
	if err != nil {
		return false, errors.Wrap(err, "couldn't get pulse")
	}
	return jc.IsAuthorized(ctx, role, obj, p.PulseNumber, jc.Me())
}

// QueryRole returns node refs responsible for role bound operations for given object and pulse.
func (jc *Coordinator) QueryRole(
	ctx context.Context,
	role insolar.DynamicRole,
	objID insolar.ID,
	pulseNumber insolar.PulseNumber,
) ([]insolar.Reference, error) {
	switch role {
	case insolar.DynamicRoleVirtualExecutor:
		n, err := jc.VirtualExecutorForObject(ctx, objID, pulseNumber)
		if err != nil {
			return nil, errors.Wrapf(err, "calc DynamicRoleVirtualExecutor for object %v failed", objID.String())
		}
		return []insolar.Reference{*n}, nil

	case insolar.DynamicRoleVirtualValidator:
		return jc.VirtualValidatorsForObject(ctx, objID, pulseNumber)

	case insolar.DynamicRoleLightExecutor:
		if objID.Pulse() == pulse.Jet {
			n, err := jc.LightExecutorForJet(ctx, objID, pulseNumber)
			if err != nil {
				return nil, errors.Wrapf(err, "calc DynamicRoleLightExecutor for object %v failed", objID.String())
			}
			return []insolar.Reference{*n}, nil
		}
		n, err := jc.LightExecutorForObject(ctx, objID, pulseNumber)
		if err != nil {
			return nil, errors.Wrapf(err, "calc LightExecutorForObject for object %v failed", objID.String())
		}
		return []insolar.Reference{*n}, nil

	case insolar.DynamicRoleLightValidator:
		ref, err := jc.LightValidatorsForObject(ctx, objID, pulseNumber)
		if err != nil {
			return nil, errors.Wrapf(err, "calc DynamicRoleLightValidator for object %v failed", objID.String())
		}
		return ref, nil

	case insolar.DynamicRoleHeavyExecutor:
		n, err := jc.Heavy(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "calc DynamicRoleHeavyExecutor for pulse %v failed", pulseNumber.String())
		}
		return []insolar.Reference{*n}, nil
	}

	inslogger.FromContext(ctx).Panicf("unexpected role %v", role.String())
	return nil, nil
}

// VirtualExecutorForObject returns list of VEs for a provided pulse and objID
func (jc *Coordinator) VirtualExecutorForObject(
	ctx context.Context, objID insolar.ID, pulse insolar.PulseNumber,
) (*insolar.Reference, error) {
	nodes, err := jc.virtualsForObject(ctx, objID, pulse, VirtualExecutorCount)
	if err != nil {
		return nil, err
	}
	return &nodes[0], nil
}

// VirtualValidatorsForObject returns list of VVs for a provided pulse and objID
func (jc *Coordinator) VirtualValidatorsForObject(
	ctx context.Context, objID insolar.ID, pulse insolar.PulseNumber,
) ([]insolar.Reference, error) {
	nodes, err := jc.virtualsForObject(ctx, objID, pulse, VirtualValidatorCount+VirtualExecutorCount)
	if err != nil {
		return nil, errors.Wrapf(err, "calc VirtualValidatorsForObject for object %v failed", objID.String())
	}
	// Skipping `VirtualExecutorCount` for validators
	// because it will be selected as the executor(s) for the same pulse.
	return nodes[VirtualExecutorCount:], nil
}

// LightExecutorForJet returns list of LEs for a provided pulse and jetID
func (jc *Coordinator) LightExecutorForJet(
	ctx context.Context, jetID insolar.ID, pulse insolar.PulseNumber,
) (*insolar.Reference, error) {
	nodes, err := jc.lightMaterialsForJet(ctx, jetID, pulse, MaterialExecutorCount)
	if err != nil {
		return nil, err
	}
	return &nodes[0], nil
}

// LightValidatorsForJet returns list of LVs for a provided pulse and jetID
func (jc *Coordinator) LightValidatorsForJet(
	ctx context.Context, jetID insolar.ID, pulse insolar.PulseNumber,
) ([]insolar.Reference, error) {
	nodes, err := jc.lightMaterialsForJet(ctx, jetID, pulse, MaterialValidatorCount+MaterialExecutorCount)
	if err != nil {
		return nil, err
	}
	// Skipping `MaterialExecutorCount` for validators
	// because it will be selected as the executor(s) for the same pulse.
	return nodes[MaterialExecutorCount:], nil
}

// LightExecutorForObject returns list of LEs for a provided pulse and objID
func (jc *Coordinator) LightExecutorForObject(
	ctx context.Context, objID insolar.ID, pulse insolar.PulseNumber,
) (*insolar.Reference, error) {
	jetID, _ := jc.JetAccessor.ForID(ctx, pulse, objID)
	return jc.LightExecutorForJet(ctx, insolar.ID(jetID), pulse)
}

// LightValidatorsForObject returns list of LVs for a provided pulse and objID
func (jc *Coordinator) LightValidatorsForObject(
	ctx context.Context, objID insolar.ID, pulse insolar.PulseNumber,
) ([]insolar.Reference, error) {
	jetID, _ := jc.JetAccessor.ForID(ctx, pulse, objID)
	return jc.LightValidatorsForJet(ctx, insolar.ID(jetID), pulse)
}

// Heavy returns *insolar.RecorRef to heavy
func (jc *Coordinator) Heavy(ctx context.Context) (*insolar.Reference, error) {
	latest, err := jc.PulseAccessor.Latest(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch pulse")
	}

	candidates, err := jc.Nodes.InRole(latest.PulseNumber, insolar.StaticRoleHeavyMaterial)
	if err == node.ErrNoNodes {
		return nil, err
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch active heavy nodes for pulse %v", latest.PulseNumber)
	}
	if len(candidates) == 0 {
		return nil, errors.New(fmt.Sprintf("no active heavy nodes for pulse %d", latest.PulseNumber))
	}
	ent, err := jc.entropy(ctx, latest.PulseNumber)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch entropy for pulse %v", latest.PulseNumber)
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
func (jc *Coordinator) IsBeyondLimit(ctx context.Context, targetPN insolar.PulseNumber) (bool, error) {
	// Genesis case. When there is no any data on a lme
	if targetPN <= insolar.GenesisPulse.PulseNumber {
		return true, nil
	}

	latest, err := jc.PulseAccessor.Latest(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to fetch pulse")
	}

	// Out target on the latest pulse. It's within limit.
	if latest.PulseNumber <= targetPN {
		return false, nil
	}

	iter := latest.PulseNumber
	for i := 1; i <= jc.lightChainLimit; i++ {
		stepBack, err := jc.PulseCalculator.Backwards(ctx, latest.PulseNumber, i)
		// We could not reach our target and ran out of known pulses. It means it's beyond limit.
		if err == insolarPulse.ErrNotFound {
			return true, nil
		}
		if err != nil {
			return false, errors.Wrap(err, "failed to calculate pulse")
		}
		// We reached our target. It's within limit.
		if iter <= targetPN {
			return false, nil
		}

		iter = stepBack.PulseNumber
	}
	// We iterated limit back. It means our data is further back and beyond limit.
	return true, nil
}

// NodeForJet calculates a node (LME or heavy) for a specific jet for a specific pulseNumber
func (jc *Coordinator) NodeForJet(ctx context.Context, jetID insolar.ID, targetPN insolar.PulseNumber) (*insolar.Reference, error) {
	toHeavy, err := jc.IsBeyondLimit(ctx, targetPN)
	if err != nil {
		return nil, errors.Wrapf(err, "[IsBeyondLimit] failed, targetPN - %v", targetPN)
	}

	if toHeavy {
		return jc.Heavy(ctx)
	}
	return jc.LightExecutorForJet(ctx, jetID, targetPN)
}

// NodeForObject calculates a node (LME or heavy) for a specific jet for a specific pulseNumber
func (jc *Coordinator) NodeForObject(ctx context.Context, objectID insolar.ID, targetPN insolar.PulseNumber) (*insolar.Reference, error) {
	toHeavy, err := jc.IsBeyondLimit(ctx, targetPN)
	if err != nil {
		return nil, err
	}

	if toHeavy {
		return jc.Heavy(ctx)
	}
	return jc.LightExecutorForObject(ctx, objectID, targetPN)
}

func (jc *Coordinator) virtualsForObject(
	ctx context.Context, objID insolar.ID, pulse insolar.PulseNumber, count int,
) ([]insolar.Reference, error) {
	candidates, err := jc.Nodes.InRole(pulse, insolar.StaticRoleVirtual)
	if err == node.ErrNoNodes {
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
		utils.CircleXOR(ent[:], objID.Hash()),
		candidates,
		count,
	)
}

func (jc *Coordinator) lightMaterialsForJet(
	ctx context.Context, jetID insolar.ID, pulse insolar.PulseNumber, count int,
) ([]insolar.Reference, error) {
	prefix := insolar.JetID(jetID).Prefix()

	candidates, err := jc.Nodes.InRole(pulse, insolar.StaticRoleLightMaterial)
	if err == node.ErrNoNodes {
		return nil, err
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch active light nodes for pulse %v", pulse)
	}
	if len(candidates) == 0 {
		return nil, node.ErrNoNodes
	}

	ent, err := jc.entropy(ctx, pulse)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch entropy for pulse %v", pulse)
	}

	return getRefs(
		jc.PlatformCryptographyScheme,
		utils.CircleXOR(ent[:], prefix),
		candidates,
		count,
	)
}

func (jc *Coordinator) entropy(ctx context.Context, pulse insolar.PulseNumber) (insolar.Entropy, error) {
	current, err := jc.PulseAccessor.Latest(ctx)
	if err != nil {
		return insolar.Entropy{}, errors.Wrap(err, "failed to get current pulse")
	}

	if current.PulseNumber == pulse {
		return current.Entropy, nil
	}

	older, err := jc.PulseAccessor.ForPulseNumber(ctx, pulse)
	if err != nil {
		return insolar.Entropy{}, errors.Wrapf(err, "failed to fetch pulse data for pulse %v", pulse)
	}

	return older.Entropy, nil
}

func getRefs(
	scheme insolar.PlatformCryptographyScheme,
	e []byte,
	values []insolar.Node,
	count int,
) ([]insolar.Reference, error) {
	sort.SliceStable(values, func(i, j int) bool {
		return values[i].ID.Compare(values[j].ID) < 0
	})
	in := make([]interface{}, 0, len(values))
	for _, value := range values {
		in = append(in, interface{}(value.ID))
	}

	res, err := entropy.SelectByEntropy(scheme, e, in, count)
	if err != nil {
		return nil, err
	}
	out := make([]insolar.Reference, 0, len(res))
	for _, value := range res {
		out = append(out, value.(insolar.Reference))
	}
	return out, nil
}
