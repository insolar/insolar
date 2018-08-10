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
	test_balance := 42

	w := wallet{
		BaseSmartContract: *contract.NewBaseSmartContract(parent),
		balance:           test_balance,
	}
	proxy := walletProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: &w,
		},
	}

	assert.Equal(t, proxy.GetBalance(), test_balance)

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

	w := wallet{
		BaseSmartContract: *contract.NewBaseSmartContract(parent),
		balance:           test_balance,
	}

	assert.Equal(t, w.GetBalance(), test_balance)
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

func TestWalletFactory_GetParent(t *testing.T) {
	parent := &mockParent{}
	factory := NewWalletFactory(parent)

	assert.Equal(t, parent, factory.GetParent())
}
