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

package deposit

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const days = 1827

func TestVestingCoeffs(t *testing.T) {
	exp := make([]float64, days)
	var expTotal float64
	for i := 0; i < days; i++ {
		exp[i] = math.Pow(float64(days), float64(i+1)/float64(days))
		expTotal += exp[i]
	}

	frac := make([]float64, days)
	for i := 0; i < days; i++ {
		var expSum float64
		for j := 0; j <= i; j++ {
			expSum += exp[j]
		}
		frac[i] = expSum / expTotal
	}

	denom, ok := new(big.Int).SetString(VestingCoeffDenominator, 10)
	require.True(t, ok, "failed to parse vesting denominator")
	for i, f := range frac {
		fr := big.NewFloat(f)
		coeff := new(big.Float).Mul(fr, new(big.Float).SetInt(denom))

		coeffInt, _ := coeff.Int(nil)
		expectedCoeff, ok := new(big.Int).SetString(VestingCoeffs[i], 10)
		require.True(t, ok)
		assert.Truef(t, expectedCoeff.Cmp(coeffInt) == 0, "step: %d, expected: %s actual: %s", i, expectedCoeff.String(), coeffInt.String())
	}
}
