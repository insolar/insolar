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

package insolar

import (
	"context"
)

// DynamicRole is number representing a node role.
type DynamicRole int

const (
	// DynamicRoleUndefined is used for special cases.
	DynamicRoleUndefined = DynamicRole(iota)
	// DynamicRoleVirtualExecutor is responsible for current pulse CPU operations.
	DynamicRoleVirtualExecutor
	// DynamicRoleVirtualValidator is responsible for previous pulse CPU operations.
	DynamicRoleVirtualValidator
	// DynamicRoleLightExecutor is responsible for current pulse Disk operations.
	DynamicRoleLightExecutor
	// DynamicRoleLightValidator is responsible for previous pulse Disk operations.
	DynamicRoleLightValidator
	// DynamicRoleHeavyExecutor is responsible for permanent Disk operations.
	DynamicRoleHeavyExecutor
)

// IsVirtualRole checks if node role is virtual (validator or executor).
func (r DynamicRole) IsVirtualRole() bool {
	switch r {
	case DynamicRoleVirtualExecutor:
		return true
	case DynamicRoleVirtualValidator:
		return true
	}
	return false
}

//go:generate minimock -i github.com/insolar/insolar/insolar.PulseManager -o ../testutils -s _mock.go

// PulseManager provides Ledger's methods related to Pulse.
type PulseManager interface {
	// Set set's new pulse and closes current jet drop. If dry is true, nothing will be saved to storage.
	Set(ctx context.Context, pulse Pulse, persist bool) error
}

//go:generate minimock -i github.com/insolar/insolar/insolar.JetCoordinator -o ../testutils -s _mock.go

// JetCoordinator provides methods for calculating Jet affinity
// (e.g. to which Jet a message should be sent).
type JetCoordinator interface {
	// Me returns current node.
	Me() Reference

	// IsAuthorized checks for role on concrete pulse for the address.
	IsAuthorized(ctx context.Context, role DynamicRole, obj ID, pulse PulseNumber, node Reference) (bool, error)

	// QueryRole returns node refs responsible for role bound operations for given object and pulse.
	QueryRole(ctx context.Context, role DynamicRole, obj ID, pulse PulseNumber) ([]Reference, error)

	VirtualExecutorForObject(ctx context.Context, objID ID, pulse PulseNumber) (*Reference, error)
	VirtualValidatorsForObject(ctx context.Context, objID ID, pulse PulseNumber) ([]Reference, error)

	LightExecutorForObject(ctx context.Context, objID ID, pulse PulseNumber) (*Reference, error)
	LightValidatorsForObject(ctx context.Context, objID ID, pulse PulseNumber) ([]Reference, error)
	// LightExecutorForJet calculates light material executor for provided jet.
	LightExecutorForJet(ctx context.Context, jetID ID, pulse PulseNumber) (*Reference, error)
	LightValidatorsForJet(ctx context.Context, jetID ID, pulse PulseNumber) ([]Reference, error)

	Heavy(ctx context.Context, pulse PulseNumber) (*Reference, error)

	IsBeyondLimit(ctx context.Context, currentPN, targetPN PulseNumber) (bool, error)
	NodeForJet(ctx context.Context, jetID ID, rootPN, targetPN PulseNumber) (*Reference, error)

	// NodeForObject calculates a node (LME or heavy) for a specific jet for a specific pulseNumber
	NodeForObject(ctx context.Context, objectID ID, rootPN, targetPN PulseNumber) (*Reference, error)
}

// KV is a generic key/value struct.
type KV struct {
	K []byte
	V []byte
}

// KVSize returns size of key/value array in bytes.
func KVSize(kvs []KV) (amount int64) {
	for _, kv := range kvs {
		amount += int64(len(kv.K) + len(kv.V))
	}
	return
}

// StorageExportResult represents storage data view.
type StorageExportResult struct {
	Data     map[string]interface{}
	NextFrom *PulseNumber
	Size     int
}

var (
	DomainID = *NewID(0, nil)
)
