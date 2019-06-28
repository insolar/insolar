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
