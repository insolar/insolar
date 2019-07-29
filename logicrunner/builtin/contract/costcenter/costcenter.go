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

package costcenter

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

// CostCenter provides access to tariffs and wallet for commission.
type CostCenter struct {
	foundation.BaseContract
	CommissionWallet insolar.Reference
	CurrentTariff    insolar.Reference
	Tariffs          []insolar.Reference
}

// New creates new cost center.
func NewCostCenter(commissionWallet insolar.Reference, currentTariff insolar.Reference) (*CostCenter, error) {
	return &CostCenter{
		CommissionWallet: commissionWallet,
		CurrentTariff:    currentTariff,
	}, nil
}

// SetTariffs sets tariffs.
func (cc CostCenter) SetTariffs(tariffs []insolar.Reference) error {
	cc.Tariffs = tariffs
	return nil
}

// GetTariffs gets tariffs.
func (cc CostCenter) GetTariffs() ([]insolar.Reference, error) {
	return cc.Tariffs, nil
}

// SetCurrentTariff sets current tariff.
func (cc CostCenter) SetCurrentTariff(currentTariff insolar.Reference) error {
	cc.CurrentTariff = currentTariff
	return nil
}

// GetCurrentTariff gets current tariff.
func (cc CostCenter) GetCurrentTariff() (insolar.Reference, error) {
	return cc.CurrentTariff, nil
}
