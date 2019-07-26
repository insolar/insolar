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

	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

type Tariff struct {
	foundation.BaseContract
}

// New creates new tariff.
func New() (*Tariff, error) {
	return &Tariff{}, nil
}

func calcFeeRate(amountStr string) (string, error) {
	amount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return "", fmt.Errorf("can't parse amount")
	}

	if amount.Cmp(big.NewInt(1000*1000*1000)) >= 0 {
		return "1000000000", nil // 1000 * 1000 * 1000 = 10%
	}
	if amount.Cmp(big.NewInt(1000*1000)) >= 0 {
		return "2000000000", nil // 2 * 1000 * 1000 * 1000 = 20%
	}
	if amount.Cmp(big.NewInt(1000)) >= 0 {
		return "3000000000", nil // 3 * 1000 * 1000 * 1000 = 30%
	}
	return "4000000000", nil // 4 * 1000 * 1000 * 1000 = 40%

}

// CalcFee calculates fee for amount. Returns fee.
func (t Tariff) CalcFee(amountStr string) (string, error) {
	amount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return "", fmt.Errorf("can't parse amount")
	}

	commissionRateStr, err := calcFeeRate(amountStr)
	if err != nil {
		return "", fmt.Errorf("failed to calc fee rate")
	}

	commissionRate, ok := new(big.Int).SetString(commissionRateStr, 10)
	if !ok {
		return "", fmt.Errorf("can't parse commission rate")
	}

	preResult := new(big.Int).Mul(amount, commissionRate)

	capacity := big.NewInt(10 * 1000 * 1000 * 1000)
	result := new(big.Int).Div(preResult, capacity)

	mod := new(big.Int).Mod(preResult, capacity)
	if mod.Cmp(big.NewInt(0)) == 1 {
		result = new(big.Int).Add(result, big.NewInt(1))
	}

	return result.String(), nil
}
