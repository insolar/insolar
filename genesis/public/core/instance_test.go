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

package core

import (
	"fmt"
	"testing"

	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/contract"
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type mockProxy struct {
	reference object.Reference
}

func (p *mockProxy) GetClassID() string {
	return "mockProxy"
}

func (p *mockProxy) GetClass() object.Factory {
	return nil
}

func (p *mockProxy) GetReference() object.Reference {
	return p.reference
}

func (p *mockProxy) SetReference(reference object.Reference) {
	p.reference = reference
}

type mockFactory struct {
}

func (f *mockFactory) Create(parent object.Parent) (object.Proxy, error) {
	return &mockChildProxy{
		parent: parent,
	}, nil
}

func (f *mockFactory) GetClassID() string {
	return "mockFactory"
}

func (f *mockFactory) GetClass() object.Factory {
	return &mockFactory{}
}

func (f *mockFactory) GetReference() object.Reference {
	return nil
}

func (f *mockFactory) GetParent() object.Parent {
	return nil
}

func (f *mockFactory) SetReference(reference object.Reference) {

}

type mockFactoryError struct {
	mockFactory
}

func (f *mockFactoryError) Create(parent object.Parent) (object.Proxy, error) {
	return nil, fmt.Errorf("factory create error")
}

type mockFactoryNilError struct {
	mockFactory
}

func (f *mockFactoryNilError) Create(parent object.Parent) (object.Proxy, error) {
	return nil, nil
}

func TestNewInstanceDomain(t *testing.T) {
	factory := &mockFactory{}
	parent := &mockParent{}
	instDom, err := newInstanceDomain(parent, factory)

	assert.NoError(t, err)
	assert.Equal(t, &instanceDomain{
		BaseDomain: *domain.NewBaseDomain(parent, factory, InstanceDomainName),
	}, instDom)
}

func TestNewInstanceDomain_WithNilParent(t *testing.T) {
	factory := &mockFactory{}
	_, err := newInstanceDomain(nil, factory)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestInstanceDomain_GetClassID(t *testing.T) {
	factory := &mockFactory{}
	parent := &mockParent{}
	instDom, err := newInstanceDomain(parent, factory)
	assert.NoError(t, err)

	domainID := instDom.GetClassID()
	assert.Equal(t, class.InstanceDomainID, domainID)
}

func TestInstanceDomain_GetClass(t *testing.T) {
	factory := &mockFactory{}
	parent := &mockParent{}
	instDom, err := newInstanceDomain(parent, factory)
	assert.NoError(t, err)

	assert.Equal(t, factory, instDom.GetClass())
}

func TestInstanceDomain_CreateInstance(t *testing.T) {
	factory := &mockFactory{}
	parent := &mockParent{}
	instDom, err := newInstanceDomain(parent, factory)
	assert.NoError(t, err)

	registered, err := instDom.CreateInstance(factory)
	assert.NoError(t, err)

	_, err = uuid.FromString(registered)
	assert.NoError(t, err)
}

func TestInstanceDomain_CreateInstance_WithError(t *testing.T) {
	factory := &mockFactory{}
	parent := &mockParent{}
	instDom, err := newInstanceDomain(parent, factory)
	assert.NoError(t, err)

	factoryErr := &mockFactoryError{}
	_, err = instDom.CreateInstance(factoryErr)
	assert.EqualError(t, err, "factory create error")
}

func TestInstanceDomain_CreateInstance_WithNilError(t *testing.T) {
	factory := &mockFactory{}
	parent := &mockParent{}
	instDom, err := newInstanceDomain(parent, factory)
	assert.NoError(t, err)

	factoryErr := &mockFactoryNilError{}
	_, err = instDom.CreateInstance(factoryErr)
	assert.EqualError(t, err, "factory returns nil")
}

func TestInstanceDomain_GetInstance(t *testing.T) {
	factory := &mockFactory{}
	parent := &mockParent{}
	instDom, err := newInstanceDomain(parent, factory)
	assert.NoError(t, err)

	registered, err := instDom.CreateInstance(factory)
	assert.NoError(t, err)

	resolved, err := instDom.GetInstance(registered)
	assert.NoError(t, err)

	assert.Equal(t, &mockChildProxy{
		parent: instDom,
	}, resolved)
}

func TestInstanceDomain_GetInstance_IncorrectRef(t *testing.T) {
	factory := &mockFactory{}
	parent := &mockParent{}
	instDom, err := newInstanceDomain(parent, factory)
	assert.NoError(t, err)

	_, err = instDom.GetInstance("1")
	assert.EqualError(t, err, "object with record 1 does not exist")
}

func TestNewInstanceDomainProxy(t *testing.T) {
	factory := &mockFactory{}
	parent := &mockParent{}
	instDom, err := newInstanceDomain(parent, factory)
	assert.NoError(t, err)

	proxy, err := newInstanceDomainProxy(parent, factory)
	assert.NoError(t, err)

	assert.Equal(t, &instanceDomainProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: instDom,
		},
	}, proxy)
}

func TestNewInstanceDomainProxy_GetClass(t *testing.T) {
	factory := &mockFactory{}
	parent := &mockParent{}
	proxy, err := newInstanceDomainProxy(parent, factory)
	assert.NoError(t, err)
	assert.Equal(t, factory, proxy.GetClass())
}

func TestNewInstanceDomainProxy_WithNilParent(t *testing.T) {
	factory := &mockFactory{}
	_, err := newInstanceDomainProxy(nil, factory)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestInstanceDomainProxy_CreateInstance(t *testing.T) {
	factory := &mockFactory{}
	parent := &mockParent{}
	proxy, err := newInstanceDomainProxy(parent, factory)
	assert.NoError(t, err)

	registered, err := proxy.CreateInstance(factory)
	assert.NoError(t, err)

	_, err = uuid.FromString(registered)
	assert.NoError(t, err)
}

func TestInstanceDomainProxy_CreateInstance_WithError(t *testing.T) {
	factory := &mockFactory{}
	parent := &mockParent{}
	proxy, err := newInstanceDomainProxy(parent, factory)
	assert.NoError(t, err)

	factoryErr := &mockFactoryError{}
	_, err = proxy.CreateInstance(factoryErr)
	assert.EqualError(t, err, "factory create error")
}

func TestInstanceDomainProxy_CreateInstance_WithNilError(t *testing.T) {
	factory := &mockFactory{}
	parent := &mockParent{}
	proxy, err := newInstanceDomainProxy(parent, factory)
	assert.NoError(t, err)

	factoryErr := &mockFactoryNilError{}
	_, err = proxy.CreateInstance(factoryErr)
	assert.EqualError(t, err, "factory returns nil")
}

func TestInstanceDomainProxy_GetInstance(t *testing.T) {
	factory := &mockFactory{}
	parent := &mockParent{}
	proxy, err := newInstanceDomainProxy(parent, factory)
	assert.NoError(t, err)

	registered, err := proxy.CreateInstance(factory)
	assert.NoError(t, err)

	resolved, err := proxy.GetInstance(registered)
	assert.NoError(t, err)

	assert.Equal(t, &mockChildProxy{
		parent: proxy.Instance.(object.Parent),
	}, resolved)
}

func TestInstanceDomainProxy_GetInstance_IncorrectRef(t *testing.T) {
	factory := &mockFactory{}
	parent := &mockParent{}
	proxy, err := newInstanceDomainProxy(parent, factory)
	assert.NoError(t, err)

	_, err = proxy.GetInstance("1")
	assert.EqualError(t, err, "object with record 1 does not exist")
}

func TestNewInstanceDomainFactory(t *testing.T) {
	parent := &mockParent{}
	factory := NewInstanceDomainFactory(parent)
	assert.Equal(t, &instanceDomainFactory{parent: parent}, factory)
}

func TestInstanceDomainFactory_GetClassID(t *testing.T) {
	parent := &mockParent{}
	factory := NewInstanceDomainFactory(parent)
	assert.Equal(t, class.InstanceDomainID, factory.GetClassID())
}

func TestInstanceDomainFactory_GetClass(t *testing.T) {
	parent := &mockParent{}
	factory := NewInstanceDomainFactory(parent)
	assert.Equal(t, factory, factory.GetClass())
}

func TestInstanceDomainFactory_Create(t *testing.T) {
	parent := &mockParent{}
	factory := NewInstanceDomainFactory(parent)
	proxy, err := factory.Create(parent)
	instDom, _ := newInstanceDomain(parent, factory)

	assert.NoError(t, err)
	assert.Equal(t, &instanceDomainProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: instDom,
		},
	}, proxy)
}

func TestInstanceDomainFactory_CreateWithNoParent(t *testing.T) {
	parent := &mockParent{}
	factory := NewInstanceDomainFactory(parent)
	proxy, err := factory.Create(nil)

	assert.EqualError(t, err, "parent must not be nil")
	assert.Nil(t, proxy)
}

func TestInstanceDomainFactory_CreateWithError(t *testing.T) {
	parent := &mockParentWithError{}
	factory := NewInstanceDomainFactory(parent)
	proxy, err := factory.Create(parent)

	assert.EqualError(t, err, "add child error")
	assert.Nil(t, proxy)
}
