// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package safemath

import (
	"errors"
	"math/big"
)

// Sub subtracts two uints, reverts on overflow (i.e. if subtrahend is greater than minuend).
func Sub(a *big.Int, b *big.Int) (*big.Int, error) {
	result := new(big.Int)

	if a.Cmp(b) == -1 {
		return result, errors.New("subtrahend must be smaller than minuend")
	}

	return result.Sub(a, b), nil
}

// Add adds two uints, reverts on overflow.
func Add(a *big.Int, b *big.Int) (*big.Int, error) {
	result := new(big.Int)
	result.Add(a, b)

	if a.Cmp(result) == 1 {
		return result, errors.New("overflow at addition")
	}

	return result, nil
}
