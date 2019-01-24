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

package pulsemanager

import (
	"context"

	"github.com/pkg/errors"
)

func (m *PulseManager) initJetSyncState(ctx context.Context) error {
	allJets, err := m.ReplicaStorage.GetAllNonEmptySyncClientJets(ctx)
	if err != nil {
		return errors.Wrap(err, "failed get heavy client jets' sync state")
	}
	for jetID, pulses := range allJets {
		m.syncClientsPool.AddPulsesToSyncClient(ctx, jetID, false, pulses...)
	}
	return nil
}
