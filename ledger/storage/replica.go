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

package storage

import (
	"context"

	"github.com/insolar/insolar/core"
)

// SetReplicatedPulse saves last pulse succesfully replicated to 'heavy material' node.
func (db *DB) SetReplicatedPulse(ctx context.Context, pulsenum core.PulseNumber) error {
	return db.Update(ctx, func(tx *TransactionManager) error {
		return tx.set(ctx, prefixkey(scopeIDSystem, []byte{sysReplicatedPulse}), pulsenum.Bytes())
	})
}

// GetReplicatedPulse returns last pulse succesfully replicated to 'heavy material' node.
func (db *DB) GetReplicatedPulse(ctx context.Context) (core.PulseNumber, error) {
	buf, err := db.get(ctx, prefixkey(scopeIDSystem, []byte{sysReplicatedPulse}))
	if err != nil {
		if err == ErrNotFound {
			err = nil
		}
		return 0, err
	}
	return core.NewPulseNumber(buf), nil
}

// SetLastPulseAsLightMaterial saves last pulse then node had a 'light material' role.
func (db *DB) SetLastPulseAsLightMaterial(ctx context.Context, pulsenum core.PulseNumber) error {
	return db.Update(ctx, func(tx *TransactionManager) error {
		return tx.set(ctx, prefixkey(scopeIDSystem, []byte{sysLastPulseAsLightMaterial}), pulsenum.Bytes())
	})
}

// GetLastPulseAsLightMaterial returns last pulse then node had a 'light material' role.
func (db *DB) GetLastPulseAsLightMaterial(ctx context.Context) (core.PulseNumber, error) {
	buf, err := db.get(ctx, prefixkey(scopeIDSystem, []byte{sysLastPulseAsLightMaterial}))
	if err != nil {
		if err == ErrNotFound {
			err = nil
		}
		return 0, err
	}
	return core.NewPulseNumber(buf), nil
}
