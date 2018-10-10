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

package example

import (
	"testing"

	"github.com/insolar/insolar/application/model/class"
	"github.com/insolar/insolar/application/model/contract"
	"github.com/insolar/insolar/application/model/object"
	"github.com/stretchr/testify/assert"
)

func MakeTestReference(record string) object.Reference {
	result, _ := object.NewReference("test", record, object.ContextScope)
	return result
}

var testAmount = 100500
var testSender = MakeTestReference("sender")
var testReceiver = MakeTestReference("receiver")

func TestNewAllowance(t *testing.T) {

	compositeFactory := &MockBaseCompositeFactory{}
	parent := &mockParent{}
	al, err := newAllowance(parent, compositeFactory)
	assert.NoError(t, err)

	assert.Equal(t, &allowance{
		amount:            0,
		completed:         false,
		BaseSmartContract: *contract.NewBaseSmartContract(parent, compositeFactory),
	}, al)
}

func TestAllowance_IsCompleted(t *testing.T) {
	al := allowance{
		completed: false,
	}

	assert.Equal(t, false, al.IsCompleted())

	al.completed = true
	assert.Equal(t, true, al.IsCompleted())
}

func TestAllowance_MarkCompleted(t *testing.T) {
	al := allowance{
		completed: false,
	}
	assert.Equal(t, false, al.completed)

	al.MarkCompleted()
	assert.Equal(t, true, al.completed)
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

func TestAllowance_GetReceiver(t *testing.T) {
	al := allowance{
		receiver: testReceiver,
	}
	assert.Equal(t, testReceiver, al.GetReceiver())
}

func TestAllowance_GetInterfaceKey(t *testing.T) {
	compositeFactory := &MockBaseCompositeFactory{}
	parent := &mockParent{}
	al, err := newAllowance(parent, compositeFactory)
	assert.NoError(t, err)
	assert.Equal(t, class.AllowanceID, al.GetInterfaceKey())
}

func TestNewAllowanceFactory(t *testing.T) {
	parent := &mockParent{}
	factory, _ := NewAllowanceFactory(parent)

	expected := &allowanceFactory{
		parent: parent,
	}

	assert.Equal(t, expected, factory)

}

func TestAllowanceFactory_GetClassID(t *testing.T) {
	parent := &mockParent{}
	factory, _ := NewAllowanceFactory(parent)
	assert.Equal(t, class.AllowanceID, factory.GetClassID())
}

func TestAllowanceFactory_GetInterfaceKey(t *testing.T) {
	parent := &mockParent{}
	factory, _ := NewAllowanceFactory(parent)
	assert.Equal(t, class.AllowanceID, factory.GetInterfaceKey())
}

func TestAllowanceFactory_InterfaceKeyEqualClassID(t *testing.T) {
	parent := &mockParent{}
	factory, _ := NewAllowanceFactory(parent)
	assert.Equal(t, factory.GetInterfaceKey(), factory.GetClassID())
}

func TestAllowanceFactory_Create_WithNilParet(t *testing.T) {
	parent := &mockParent{}
	factory, _ := NewAllowanceFactory(parent)
	_, err := factory.Create(nil)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestAllowanceFactory_CreateWithError(t *testing.T) {
	parent := &mockParentWithError{}
	factory, _ := NewAllowanceFactory(parent)
	_, err := factory.Create(parent)

	assert.EqualError(t, err, "add child error")
}

func TestAllowanceFactory_Create(t *testing.T) {
	parent := &mockParent{}
	factory, _ := NewAllowanceFactory(parent)
	proxy, err := factory.Create(parent)
	assert.NoError(t, err)

	expectedAllowance := AllowanceCompositeCollection{
		parent: parent,
		class:  factory,
	}

	assert.Equal(t, &contract.BaseCompositeCollectionProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: &expectedAllowance,
		},
	}, proxy)
}

func TestAllowanceFactory_GetParent(t *testing.T) {
	parent := &mockParent{}
	factory, _ := NewAllowanceFactory(parent)

	assert.Equal(t, parent, factory.GetParent())
}

func TestNewAllowanceProxy_WithNilParent(t *testing.T) {
	compositeFactory := &MockBaseCompositeFactory{}
	_, err := newAllowanceProxy(nil, compositeFactory)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestNewAllowanceProxy(t *testing.T) {
	compositeFactory := &MockBaseCompositeFactory{}
	parent := &mockParent{}
	proxy, err := newAllowanceProxy(parent, compositeFactory)
	assert.NoError(t, err)

	nAllowance, err := newAllowance(parent, compositeFactory)
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

func TestAllowanceProxy_GetReceiver(t *testing.T) {
	al := allowance{
		receiver: testReceiver,
	}

	proxy := allowanceProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: &al,
		},
	}

	assert.Equal(t, proxy.GetReceiver(), testReceiver)
}

func TestAllowanceProxy_GetInterfaceKey(t *testing.T) {
	compositeFactory := &MockBaseCompositeFactory{}
	parent := &mockParent{}
	proxy, err := newAllowanceProxy(parent, compositeFactory)
	assert.NoError(t, err)
	assert.Equal(t, class.AllowanceID, proxy.GetInterfaceKey())
}

func TestAllowanceProxy_MarkCompleted(t *testing.T) {
	compositeFactory := &MockBaseCompositeFactory{}
	parent := &mockParent{}
	proxy, err := newAllowanceProxy(parent, compositeFactory)
	assert.NoError(t, err)
	assert.Equal(t, false, proxy.Instance.(*allowance).completed)
	proxy.MarkCompleted()
	assert.Equal(t, true, proxy.Instance.(*allowance).completed)
}

func TestAllowanceProxy_IsCompleted(t *testing.T) {
	compositeFactory := &MockBaseCompositeFactory{}
	parent := &mockParent{}
	proxy, err := newAllowanceProxy(parent, compositeFactory)
	assert.NoError(t, err)
	assert.Equal(t, false, proxy.Instance.(*allowance).IsCompleted())

	proxy.Instance.(*allowance).completed = true
	assert.Equal(t, true, proxy.Instance.(*allowance).IsCompleted())
}
