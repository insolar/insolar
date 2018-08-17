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
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/stretchr/testify/assert"
)

func TestNewWallet_WithNilParent(t *testing.T) {
	_, err := newWallet(nil)

	assert.EqualError(t, err, "parent must not be nil")
}

func TestNewWallet(t *testing.T) {
	parent := &mockParent{}
	wActual, err := newWallet(parent)

	assert.NoError(t, err)
	assert.Equal(t, &wallet{
		BaseSmartContract: *contract.NewBaseSmartContract(parent),
	}, wActual)
}

func TestNewWalletProxy_WithNilParent(t *testing.T) {
	_, err := newWalletProxy(nil)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestNewWalletProxy(t *testing.T) {
	parent := &mockParent{}
	nWallet, err := newWallet(parent)
	assert.NoError(t, err)

	proxy, err := newWalletProxy(parent)
	assert.NoError(t, err)

	assert.Equal(t, &walletProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: nWallet,
		},
	}, proxy)
}

func TestWalletProxy_GetBalance(t *testing.T) {

	parent := &mockParent{}
	testBalance := 42

	mParent, _ := newMember(parent)

	w := wallet{
		BaseSmartContract: *contract.NewBaseSmartContract(mParent.(object.Parent)),
		balance:           testBalance,
	}

	proxy := walletProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: &w,
		},
	}

	tBalance, err := proxy.GetBalance()
	assert.NoError(t, err)
	assert.Equal(t, testBalance, tBalance)

}

func TestWalletProxy_GetBalance_WrongParent(t *testing.T) {
	parent := &mockParent{}

	w := wallet{
		BaseSmartContract: *contract.NewBaseSmartContract(parent),
	}

	proxy := walletProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: &w,
		},
	}

	_, err := proxy.GetBalance()
	assert.EqualError(t, err, "parent must be wallet")
}

type memberBadComposite struct {
	member
}

func (*memberBadComposite) GetOrCreateComposite(compositeFactory object.CompositeFactory) (object.Composite, error) {
	return &allowance{}, nil
}

func TestWalletProxy_GetBalance_AllowanceIsNotCollection(t *testing.T) {

	m := memberBadComposite{}

	w := wallet{
		BaseSmartContract: *contract.NewBaseSmartContract(&m),
	}

	proxy := walletProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: &w,
		},
	}

	_, err := proxy.GetBalance()
	assert.EqualError(t, err, "allowance must be composite collection")
}

func TestWalletProxy_GetBalance_With_Allowances(t *testing.T) {

	parent := &mockParent{}
	testBalance := 42
	testAllowance := 100500

	mParent, _ := newMember(parent)

	w := wallet{
		BaseSmartContract: *contract.NewBaseSmartContract(mParent.(object.Parent)),
		balance:           testBalance,
	}

	factory, _ := NewAllowanceFactory(parent)
	composite, err := mParent.GetOrCreateComposite(factory)
	assert.NoError(t, err)
	alCollection := composite.(object.CompositeCollection)

	nAllowance, err := newAllowanceWithParams(parent, testSender, testReceiver, testAllowance)
	assert.NoError(t, err)

	alCollection.Add(nAllowance)
	alCollection.Add(nAllowance)

	proxy := walletProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: &w,
		},
	}

	tBalance, err := proxy.GetBalance()
	assert.NoError(t, err)
	assert.Equal(t, testBalance+testAllowance*2, tBalance)

}

func TestWalletProxy_GetInterfaceKey(t *testing.T) {
	parent := &mockParent{}
	proxy, err := newWalletProxy(parent)
	assert.NoError(t, err)
	assert.Equal(t, class.WalletID, proxy.GetInterfaceKey())
}

func TestNewWalletFactory(t *testing.T) {
	parent := &mockParent{}
	expected := &walletFactory{parent: parent}

	factory := NewWalletFactory(parent)

	assert.Equal(t, expected, factory)
}

func TestWallet_GetBalance(t *testing.T) {
	parent := &mockParent{}
	test_balance := 42

	m, err := newMember(parent)
	assert.NoError(t, err)

	w := wallet{
		BaseSmartContract: *contract.NewBaseSmartContract(m.(object.Parent)),
		balance:           test_balance,
	}

	tBalance, err := w.GetBalance()
	assert.NoError(t, err)

	assert.Equal(t, tBalance, test_balance)
}

func TestWallet_GetClassID(t *testing.T) {
	parent := &mockParent{}
	w, err := newWallet(parent)
	assert.NoError(t, err)

	assert.Equal(t, class.WalletID, w.GetClassID())
}

func TestWallet_GetInterfaceKey(t *testing.T) {
	parent := &mockParent{}
	w, err := newWallet(parent)
	assert.NoError(t, err)

	assert.Equal(t, class.WalletID, w.GetInterfaceKey())
}

func TestWalletFactory_Create_WithNilParent(t *testing.T) {
	parent := &mockParent{}
	factory := NewWalletFactory(parent)
	_, err := factory.Create(nil)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestWalletFactory_CreateWithError(t *testing.T) {
	parent := &mockParentWithError{}
	factory := NewWalletFactory(parent)
	_, err := factory.Create(parent)

	assert.EqualError(t, err, "add child error")

}

func TestWalletFactory_Create(t *testing.T) {
	parent := &mockParent{}
	factory := NewWalletFactory(parent)
	proxy, err := factory.Create(parent)
	assert.NoError(t, err)

	expecatedWallet := wallet{
		BaseSmartContract: *contract.NewBaseSmartContract(parent),
	}

	assert.Equal(t, &walletProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: &expecatedWallet,
		},
	}, proxy)
}

func TestWalletFactory_GetClassID(t *testing.T) {
	parent := &mockParent{}
	factory := NewWalletFactory(parent)

	assert.Equal(t, class.WalletID, factory.GetClassID())
}

func TestWalletFactory_GetInterfaceKey(t *testing.T) {
	parent := &mockParent{}
	factory := NewWalletFactory(parent)

	assert.Equal(t, class.WalletID, factory.GetInterfaceKey())
}

func TestWalletFactory_GetParent(t *testing.T) {
	parent := &mockParent{}
	factory := NewWalletFactory(parent)

	assert.Equal(t, parent, factory.GetParent())
}

func TestWalletFactory_InterfaceKeyEqualClassID(t *testing.T) {
	parent := &mockParent{}
	factory := NewWalletFactory(parent)

	assert.Equal(t, factory.GetInterfaceKey(), factory.GetClassID())
}
