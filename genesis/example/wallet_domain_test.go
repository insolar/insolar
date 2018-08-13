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
	"fmt"
	"testing"

	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/stretchr/testify/assert"
)

func TestNewWalletDomain_WithNilParent(t *testing.T) {
	_, err := newWalletDomain(nil)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestNewWalletDomain(t *testing.T) {
	parent := &mockParent{}
	wDomain, err := newWalletDomain(parent)

	assert.NoError(t, err)
	assert.NotNil(t, wDomain)
	assert.NotEmpty(t, wDomain.ChildStorage.GetKeys())
}

func TestNewWalletDomainFactory(t *testing.T) {
	parent := &mockParent{}
	factory := NewWalletDomainFactory(parent)

	expected := &walletDomainFactory{parent: parent}

	assert.Equal(t, expected, factory)
}

func TestWalletDomain_CreateWallet(t *testing.T) {
	parent := &mockParent{}
	wallet, err := newWalletDomain(parent)
	assert.NoError(t, err)

	mProxy, err := newMemberProxy(parent)
	assert.NoError(t, err)

	// TODO: Check inserted composite
	err = wallet.CreateWallet(mProxy)

	assert.NoError(t, err)
}

func TestWalletDomain_CreateWallet_NoMemberFactoryRecord(t *testing.T) {
	parent := &mockParent{}

	wDomain := &walletDomain{
		BaseDomain: *domain.NewBaseDomain(parent, WalletDomainName),
	}
	wDomain.walletFactoryReference, _ = object.NewReference("", "unexistedRecord", object.ChildScope)

	mProxy, err := newMemberProxy(parent)
	assert.NoError(t, err)
	err = wDomain.CreateWallet(mProxy)

	assert.EqualError(t, err, "object with record unexistedRecord does not exist")
}

type mockWalletNotFactory struct {
	mockProxy
}

func (f *mockWalletNotFactory) GetClassID() string {
	return class.WalletID
}

func (f *mockWalletNotFactory) GetParent() object.Parent {
	return nil
}

func TestWalletDomain_CreateWallet_NotFactory(t *testing.T) {
	parent := &mockParent{}
	notFactory := &mockWalletNotFactory{}

	wDomain := &walletDomain{
		BaseDomain: *domain.NewBaseDomain(parent, WalletDomainName),
	}
	record, _ := wDomain.AddChild(notFactory)
	mProxy, err := newMemberProxy(parent)
	assert.NoError(t, err)

	wDomain.walletFactoryReference, _ = object.NewReference("", record, object.ChildScope)
	err = wDomain.CreateWallet(mProxy)
	assert.EqualError(t, err, fmt.Sprintf("child by reference `#.#%s` is not CompositeFactory instance", record))
}

func TestWalletDomain_GetClassID(t *testing.T) {
	parent := &mockParent{}
	wdomain, err := newWalletDomain(parent)

	assert.NoError(t, err)
	assert.Equal(t, class.WalletDomainID, wdomain.GetClassID())
}

func TestWalletDomainFactory_Create_WithNoParent(t *testing.T) {
	parent := &mockParent{}
	factory := NewWalletDomainFactory(parent)
	_, err := factory.Create(nil)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestWalletDomainFactory_Create(t *testing.T) {
	parent := &mockParent{}
	factory := NewWalletDomainFactory(parent)

	proxy, err := factory.Create(parent)
	assert.NoError(t, err)
	assert.NotNil(t, proxy)
}

func TestWalletDomainFactory_Create_WithError(t *testing.T) {
	parent := &mockParentWithError{}
	factory := NewWalletDomainFactory(parent)
	_, err := factory.Create(parent)

	assert.EqualError(t, err, "add child error")
}

func TestWalletDomainFactory_GetClassID(t *testing.T) {
	parent := &mockParentWithError{}
	factory := NewWalletDomainFactory(parent)
	assert.Equal(t, class.WalletDomainID, factory.GetClassID())
}

func TestWalletDomainFactory_GetParent(t *testing.T) {
	parent := &mockParentWithError{}
	factory := NewWalletDomainFactory(parent)
	actual := factory.GetParent()

	assert.Nil(t, actual)
}

func TestWalletDomainProxy_CreateWallet_WithNoParent(t *testing.T) {
	_, err := newWalletDomainProxy(nil)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestWalletDomainProxy_CreateWallet(t *testing.T) {
	parent := &mockParent{}
	proxy, err := newWalletDomainProxy(parent)
	assert.NoError(t, err)

	mProxy, err := newMemberProxy(parent)
	err = proxy.CreateWallet(mProxy)

	// TODO: Check inserted composite

	assert.NoError(t, err)
}
