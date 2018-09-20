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

package event

import (
	"io"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

type GetCode struct {
	code        core.RecordRef     // nolint
	machinePref []core.MachineType // nolint
}

func (e *GetCode) Serialize() (io.Reader, error) {
	return serialize(e, TypeGetCode)
}

func (e *GetCode) GetReference() core.RecordRef {
	return e.code
}

func (e *GetCode) GetOperatingRole() core.JetRole {
	return core.RoleLightExecutor
}

func (e *GetCode) React(c core.Components) (core.Reaction, error) {
	ledgerComponent, exists := c["core.Ledger"]
	if !exists {
		return nil, errors.New("no core.EventBus in components")
	}
	ledger, ok := ledgerComponent.(core.Ledger)
	if !ok {
		return nil, errors.New("EventBus assertion failed")
	}

	return ledger.GetArtifactManager().HandleEvent(e)
}
