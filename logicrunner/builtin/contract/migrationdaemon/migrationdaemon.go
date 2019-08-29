/*
 *
 *  Copyright  2019. Insolar Technologies GmbH
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package migrationdaemon

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

// MigrationDaemonContract make migration procedure.
type MigrationDaemon struct {
	foundation.BaseContract
	IsActive              bool
	MigrationDaemonMember insolar.Reference
}

// Set status Migration daemon.
func (md *MigrationDaemon) SetActivationStatus(status bool) error {
	md.IsActive = status
	return nil
}

// Return status migration daemon.
// ins:immutable
func (md *MigrationDaemon) GetActivationStatus() (bool, error) {
	return md.IsActive, nil
}

// Return reference on migration daemon.
// ins:immutable
func (md *MigrationDaemon) GetMigrationDaemonMember(status bool) (insolar.Reference, error) {
	return md.MigrationDaemonMember, nil
}
