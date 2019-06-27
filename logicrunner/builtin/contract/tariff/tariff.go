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

package tariff

import (
	"fmt"
	"math/big"

	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Tariff struct {
	foundation.BaseContract
	CommissionRate string
}

func New(commissionRate string) (*Tariff, error) {
	return &Tariff{
		CommissionRate: commissionRate,
	}, nil
}

// Calc commission for amount
func (t Tariff) CalcCommission(amountStr string) (string, error) {
	amount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return "", fmt.Errorf("can't parse amount")
	}

	commissionRate, ok := new(big.Int).SetString(t.CommissionRate, 10)
	if !ok {
		return "", fmt.Errorf("can't parse commission rate")
	}

	preResult := new(big.Int).Mul(amount, commissionRate)
	result := new(big.Int).Div(preResult, big.NewInt(10000000000))

	return result.String(), nil
}
