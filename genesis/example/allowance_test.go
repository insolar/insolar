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
	"github.com/insolar/insolar/genesis/model/contract"
	"github.com/stretchr/testify/assert"
)

var testAmount = 100500
var testSender = "test"
var testReciever = "testReciever"

func TestNewAllowance(t *testing.T) {

	parent := &mockParent{}
	al, err := newAllowance(parent)
	assert.NoError(t, err)

	assert.Equal(t, &allowance{
		amount:            0,
		sender:            "",
		completed:         false,
		BaseSmartContract: *contract.NewBaseSmartContract(parent),
	}, al)
}

func TestAllowance_GetAmount(t *testing.T) {

	al := allowance{
		amount: testAmount,
	}

	assert.Equal(t, testAmount, al.GetAmount())
}

func TestAllowance_GetSender(t *testing.T) {
	al := allowance{
		sender: testSender,
	}
	assert.Equal(t, testSender, al.GetSender())
}

func TestAllowance_GetReciever(t *testing.T) {
	al := allowance{
		reciever: testReciever,
	}
	assert.Equal(t, testReciever, al.GetReciever())
}

func TestAllowance_GetInterfaceKey(t *testing.T) {
	parent := &mockParent{}
	al, err := newAllowance(parent)
	assert.NoError(t, err)
	assert.Equal(t, class.AllowanceID, al.GetInterfaceKey())
}

func TestNewAllowanceFactory(t *testing.T) {
	parent := &mockParent{}
	factory := NewAllowanceFactory(parent)

	expected := &allowanceFactory{
		parent: parent,
	}

	assert.Equal(t, expected, factory)

}

func TestAllowanceFactory_GetClassID(t *testing.T) {
	parent := &mockParent{}
	factory := NewAllowanceFactory(parent)
	assert.Equal(t, class.AllowanceID, factory.GetClassID())
}

func TestAllowanceFactory_GetInterfaceKey(t *testing.T) {
	parent := &mockParent{}
	factory := NewAllowanceFactory(parent)
	assert.Equal(t, class.AllowanceID, factory.GetInterfaceKey())
}

func TestAllowanceFactory_InterfaceKeyEqualClassID(t *testing.T) {
	parent := &mockParent{}
	factory := NewAllowanceFactory(parent)
	assert.Equal(t, factory.GetInterfaceKey(), factory.GetClassID())
}

func TestAllowanceFactory_Create_WithNilParet(t *testing.T) {
	parent := &mockParent{}
	factory := NewAllowanceFactory(parent)
	_, err := factory.Create(nil)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestAllowanceFactory_CreateWithError(t *testing.T) {
	parent := &mockParentWithError{}
	factory := NewAllowanceFactory(parent)
	_, err := factory.Create(parent)

	assert.EqualError(t, err, "add child error")
}

func TestAllowanceFactory_Create(t *testing.T) {
	parent := &mockParent{}
	factory := NewAllowanceFactory(parent)
	proxy, err := factory.Create(parent)
	assert.NoError(t, err)

	expecatedAllowance := allowance{
		BaseSmartContract: *contract.NewBaseSmartContract(parent),
	}

	assert.Equal(t, &allowanceProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: &expecatedAllowance,
		},
	}, proxy)
}

func TestAllowanceFactory_GetParent(t *testing.T) {
	parent := &mockParent{}
	factory := NewAllowanceFactory(parent)

	assert.Equal(t, parent, factory.GetParent())
}

func TestNewAllowanceProxy_WithNilParent(t *testing.T) {
	_, err := newAllowanceProxy(nil)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestNewAllowanceProxy(t *testing.T) {
	parent := &mockParent{}
	proxy, err := newAllowanceProxy(parent)
	assert.NoError(t, err)

	nAllowance, err := newAllowance(parent)
	assert.NoError(t, err)

	assert.Equal(t, &allowanceProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: nAllowance,
		},
	}, proxy)

}

func TestAllowanceProxy_GetAmount(t *testing.T) {
	al := allowance{
		amount: testAmount,
	}

	proxy := allowanceProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: &al,
		},
	}

	assert.Equal(t, proxy.GetAmount(), testAmount)
}

func TestAllowanceProxy_GetSender(t *testing.T) {
	al := allowance{
		sender: testSender,
	}

	proxy := allowanceProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: &al,
		},
	}

	assert.Equal(t, proxy.GetSender(), testSender)
}

func TestAllowanceProxy_GetReciever(t *testing.T) {
	al := allowance{
		reciever: testReciever,
	}

	proxy := allowanceProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: &al,
		},
	}

	assert.Equal(t, proxy.GetReciever(), testReciever)
}

func TestAllowanceProxy_GetInterfaceKey(t *testing.T) {
	parent := &mockParent{}
	proxy, err := newAllowanceProxy(parent)
	assert.NoError(t, err)
	assert.Equal(t, class.AllowanceID, proxy.GetInterfaceKey())
}

func TestAllowanceProxy_MarkCompleted(t *testing.T) {
	parent := &mockParent{}
	proxy, err := newAllowanceProxy(parent)
	assert.NoError(t, err)
	assert.Equal(t, false, proxy.Instance.(*allowance).completed)
}
