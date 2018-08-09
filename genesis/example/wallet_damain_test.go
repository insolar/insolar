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
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewWalletDomain_WithNilParent(t *testing.T) {
	_, err := newWalletDomain(nil)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestNewWalletDomain(t *testing.T) {
	parent := &mockParent{}
	proxy, err := newWalletDomain(parent)

	assert.NoError(t, err)
	assert.Equal(t, &walletDomain{
		BaseDomain: *domain.NewBaseDomain(parent, WalletDomainName),
	}, proxy)
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

	factory := &mockFactory{}

	record, err := wallet.CreateWallet(factory)
	assert.NoError(t, err)

	_, err = uuid.FromString(record)
	assert.NoError(t, err)
}

func TestWalletDomain_CreateWallet_WithError(t *testing.T) {
	parent := &mockParent{}
	wdomain, err := newWalletDomain(parent)
	assert.NoError(t, err)

	factory := &mockFactoryError{}

	_, err = wdomain.CreateWallet(factory)

	assert.EqualError(t, err, "factory create error")
}

func TestWalletDomain_CreateWallet_WithNilError(t *testing.T) {
	parent := &mockParent{}
	wdomain, err := newWalletDomain(parent)
	assert.NoError(t, err)

	factory := &mockFactoryNilError{}

	_, err = wdomain.CreateWallet(factory)

	assert.EqualError(t, err, "factory returns nil")
}

func TestWalletDomain_GetClassID(t *testing.T) {
	parent := &mockParent{}
	wdomain, err := newWalletDomain(parent)

	assert.NoError(t, err)
	assert.Equal(t, class.WalletDomainID, wdomain.GetClassID())
}

func TestWalletDomain_GetWallet_NoSuchRecord(t *testing.T) {
	parent := &mockParent{}
	wdomain, err := newWalletDomain(parent)
	assert.NoError(t, err)

	_, err = wdomain.GetWallet("test")
	assert.EqualError(t, err, "object with record test does not exist")
}

func TestWalletDomain_GetWallet(t *testing.T) {
	parent := &mockParent{}
	wdomain, err := newWalletDomain(parent)
	assert.NoError(t, err)

	factory := &mockFactory{}
	record, err := wdomain.CreateWallet(factory)
	assert.NoError(t, err)

	proxy, err := wdomain.GetWallet(record)
	assert.NoError(t, err)

	assert.Equal(t, &mockProxy{
		parent: wdomain,
	}, proxy)
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

	wdomain, _ := newWalletDomain(parent)

	proxy, err := factory.Create(parent)
	assert.NoError(t, err)

	assert.Equal(t, &walletDomainProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: wdomain,
		},
	}, proxy)
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

	factory := &mockFactory{}
	record, err := proxy.CreateWallet(factory)
	assert.NoError(t, err)

	_, err = uuid.FromString(record)
	assert.NoError(t, err)
}

func TestWalletDomainProxy_GetWallet(t *testing.T) {
	parent := &mockParent{}
	proxyD, err := newWalletDomainProxy(parent)
	assert.NoError(t, err)

	factory := &mockFactory{}
	record, err := proxyD.CreateWallet(factory)
	assert.NoError(t, err)

	proxyW, err := proxyD.GetWallet(record)
	assert.NoError(t, err)

	assert.Equal(t, &mockProxy{
		parent: proxyD.Instance.(object.Parent),
	}, proxyW)
}
