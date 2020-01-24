// Copyright 2020 Insolar Network Ltd.
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

package deposit

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	days       = 1096
	multiplier = 3
)

func TestVestingCoeffs(t *testing.T) {
	exp := make([]float64, days)
	var expTotal float64
	for i := 0; i < days; i++ {
		exp[i] = math.Pow(float64(days*multiplier), float64(i+1)/float64(days)) / float64(days*multiplier)
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
		// fmt.Printf("\"%d\",\n", coeffInt)
		assert.Truef(t, expectedCoeff.Cmp(coeffInt) == 0, "step: %d, expected: %s actual: %s", i, expectedCoeff.String(), coeffInt.String())
	}
}

func TestVestingCoeffs_Table(t *testing.T) {
	denom, ok := new(big.Int).SetString(VestingCoeffDenominator, 10)
	require.True(t, ok)
	amount, ok := new(big.Int).SetString("5000000000000000000", 10)
	require.True(t, ok)
	xnsDenom, ok := new(big.Int).SetString("10000000000", 10)
	require.True(t, ok)

	var b strings.Builder
	fmt.Fprintf(&b, "%4s %25s %25s %25s | %25s %25s %25s\n", "Step", "Coefficient", "Released", "Daily release", "Fraction", "Released XNS", "Daily XNS")

	lastReleased := big.NewInt(0)
	for i, c := range VestingCoeffs {
		coeff, ok := new(big.Int).SetString(c, 10)
		require.True(t, ok)
		frac := new(big.Float).Quo(
			new(big.Float).SetInt(coeff),
			new(big.Float).SetInt(denom),
		)
		released := new(big.Int).Quo(
			new(big.Int).Mul(coeff, amount),
			denom,
		)
		daily := new(big.Int).Sub(released, lastReleased)
		lastReleased = released

		releasedXNS := new(big.Float).Quo(
			new(big.Float).SetInt(released),
			new(big.Float).SetInt(xnsDenom),
		)
		dailyXNS := new(big.Float).Quo(
			new(big.Float).SetInt(daily),
			new(big.Float).SetInt(xnsDenom),
		)

		fmt.Fprintf(&b, "%4d %25d %25d %25d | %25.20f %25.10f %25.10f\n", i, coeff, released, daily, frac, releasedXNS, dailyXNS)
	}

	// fmt.Println(b.String())
}

func TestVestedByNow_min_amount(t *testing.T) {
	amount := big.NewInt(1)
	zero := big.NewInt(0)
	for i := uint64(0); i <= 1094; i++ {
		assert.Equal(t, zero, VestedByNow(amount, i, 1096))
	}
	assert.Equal(t, amount, VestedByNow(amount, 1095, 1096))
}
