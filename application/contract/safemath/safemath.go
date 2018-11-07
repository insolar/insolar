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

func Div(a uint, b uint) (uint, error) {
	if b == 0 {
		return 0, errors.New("divisor cannot be zero")
	}

	return a / b, nil
}

func Sub(a uint, b uint) (uint, error) {
	if a < b {
		return 0, errors.New("subtrahend must be less than minuend")
	}
	return a - b, nil
}

func Add(a uint, b uint) (uint, error) {
	c := a + b

	if c < a {
		return 0, errors.New("overflow at addition")
	}

	return c, nil
}

func Mod(a uint, b uint) (uint, error) {
	if b == 0 {
		return 0, errors.New("divisor cannot be zero")
	}

	return a % b, nil
}
