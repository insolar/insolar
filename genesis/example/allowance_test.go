/*
 *    Copyright 2018 INS Ecosystem
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

package example

import (
	"testing"

	"github.com/insolar/insolar/genesis/model/class"
	"github.com/stretchr/testify/assert"
)

var testAmount = 100500
var testSender = "test"

func TestNewAllowance(t *testing.T) {

	al := newAllowance(testSender, testAmount)

	assert.Equal(t, &allowance{
		amount: testAmount,
		sender: testSender,
		active: false,
	}, al)

	assert.Equal(t, testAmount, al.GetAmount())
	assert.Equal(t, testSender, al.GetSender())
}

func TestAllowance_GetAmount(t *testing.T) {
	al := newAllowance(testSender, testAmount)
	assert.Equal(t, testAmount, al.GetAmount())
}

func TestAllowance_GetSender(t *testing.T) {
	al := newAllowance(testSender, testAmount)
	assert.Equal(t, testSender, al.GetSender())
}

func TestAllowance_GetInterfaceKey(t *testing.T) {
	al := newAllowance(testSender, testAmount)
	assert.Equal(t, class.AllowanceID, al.GetInterfaceKey())
}
