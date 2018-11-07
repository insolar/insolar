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

package safemath

import (
	"errors"
)

// Mul multiplies two uints, error on overflow.
func Mul(a uint, b uint) (uint, error) {
	if a == 0 {
		return 0, nil
	}

	c := a * b
	if c/a != b {
		return 0, errors.New("multiplication overflow")
	}
	return c, nil
}

// Div is integer division of two uints truncating the quotient, reverts on division by zero.
func Div(a uint, b uint) (uint, error) {
	if b == 0 {
		return 0, errors.New("divisor cannot be zero")
	}

	return a / b, nil
}

// Sub subtracts two uints, reverts on overflow (i.e. if subtrahend is greater than minuend).
func Sub(a uint, b uint) (uint, error) {
	if a < b {
		return 0, errors.New("subtrahend must be smaller than minuend")
	}
	return a - b, nil
}

// Add adds two uints, reverts on overflow.
func Add(a uint, b uint) (uint, error) {
	c := a + b

	if c < a {
		return 0, errors.New("overflow at addition")
	}

	return c, nil
}

// Mod divides two uints and returns the remainder (unsigned integer modulo), reverts when dividing by zero.
func Mod(a uint, b uint) (uint, error) {
	if b == 0 {
		return 0, errors.New("divisor cannot be zero")
	}

	return a % b, nil
}
