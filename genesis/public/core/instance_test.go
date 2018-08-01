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
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type mockProxy struct {
	parent object.Parent
}

func (p *mockProxy) GetClassID() string {
	return "mockProxy"
}

func (p *mockProxy) GetReference() *object.Reference {
	return nil
}

func (p *mockProxy) GetParent() object.Parent {
	return p.parent
}

type mockFactory struct{}

func (f *mockFactory) Create(parent object.Parent) (object.Proxy, error) {
	return &mockProxy{
		parent: parent,
	}, nil
}

func (f *mockFactory) GetClassID() string {
	return "mockFactory"
}

func (f *mockFactory) GetReference() *object.Reference {
	return nil
}

func (f *mockFactory) GetParent() object.Parent {
	return nil
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
	parent := &mockParent{}
	instDom, err := newInstanceDomain(parent)

	assert.NoError(t, err)
	assert.Equal(t, &instanceDomain{
		BaseDomain: *domain.NewBaseDomain(parent, InstanceDomainName),
	}, instDom)
}

func TestNewInstanceDomain_WithNilParent(t *testing.T) {
	_, err := newInstanceDomain(nil)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestInstanceDomain_GetClassID(t *testing.T) {
	parent := &mockParent{}
	instDom, err := newInstanceDomain(parent)
	assert.NoError(t, err)

	domainID := instDom.GetClassID()
	assert.Equal(t, class.InstanceDomainID, domainID)
}

func TestInstanceDomain_CreateInstance(t *testing.T) {
	parent := &mockParent{}
	instDom, err := newInstanceDomain(parent)
	assert.NoError(t, err)

	factory := &mockFactory{}
	registered, err := instDom.CreateInstance(factory)
	assert.NoError(t, err)

	_, err = uuid.FromString(registered)
	assert.NoError(t, err)
}

func TestInstanceDomain_CreateInstance_WithError(t *testing.T) {
	parent := &mockParent{}
	instDom, err := newInstanceDomain(parent)
	assert.NoError(t, err)

	factory := &mockFactoryError{}
	_, err = instDom.CreateInstance(factory)
	assert.EqualError(t, err, "factory create error")
}

func TestInstanceDomain_CreateInstance_WithNilError(t *testing.T) {
	parent := &mockParent{}
	instDom, err := newInstanceDomain(parent)
	assert.NoError(t, err)

	factory := &mockFactoryNilError{}
	_, err = instDom.CreateInstance(factory)
	assert.EqualError(t, err, "factory returns nil")
}

func TestInstanceDomain_GetInstance(t *testing.T) {
	parent := &mockParent{}
	instDom, err := newInstanceDomain(parent)
	assert.NoError(t, err)

	factory := &mockFactory{}
	registered, err := instDom.CreateInstance(factory)
	assert.NoError(t, err)

	resolved, err := instDom.GetInstance(registered)
	assert.NoError(t, err)

	assert.Equal(t, &mockProxy{
		parent: instDom,
	}, resolved)
}

func TestInstanceDomain_GetInstance_IncorrectRef(t *testing.T) {
	parent := &mockParent{}
	instDom, err := newInstanceDomain(parent)
	assert.NoError(t, err)

	_, err = instDom.GetInstance("1")
	assert.EqualError(t, err, "object with record 1 does not exist")
}

func TestNewInstanceDomainProxy(t *testing.T) {
	parent := &mockParent{}
	instDom, err := newInstanceDomain(parent)
	assert.NoError(t, err)

	proxy, err := newInstanceDomainProxy(parent)
	assert.NoError(t, err)

	assert.Equal(t, &instanceDomainProxy{
		instance: instDom,
	}, proxy)
}

func TestNewInstanceDomainProxy_WithNilParent(t *testing.T) {
	_, err := newInstanceDomainProxy(nil)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestInstanceDomainProxy_CreateInstance(t *testing.T) {
	parent := &mockParent{}
	proxy, err := newInstanceDomainProxy(parent)
	assert.NoError(t, err)

	factory := &mockFactory{}
	registered, err := proxy.CreateInstance(factory)
	assert.NoError(t, err)

	_, err = uuid.FromString(registered)
	assert.NoError(t, err)
}

func TestInstanceDomainProxy_CreateInstance_WithError(t *testing.T) {
	parent := &mockParent{}
	proxy, err := newInstanceDomainProxy(parent)
	assert.NoError(t, err)

	factory := &mockFactoryError{}
	_, err = proxy.CreateInstance(factory)
	assert.EqualError(t, err, "factory create error")
}

func TestInstanceDomainProxy_CreateInstance_WithNilError(t *testing.T) {
	parent := &mockParent{}
	proxy, err := newInstanceDomainProxy(parent)
	assert.NoError(t, err)

	factory := &mockFactoryNilError{}
	_, err = proxy.CreateInstance(factory)
	assert.EqualError(t, err, "factory returns nil")
}

func TestInstanceDomainProxy_GetInstance(t *testing.T) {
	parent := &mockParent{}
	proxy, err := newInstanceDomainProxy(parent)
	assert.NoError(t, err)

	factory := &mockFactory{}
	registered, err := proxy.CreateInstance(factory)
	assert.NoError(t, err)

	resolved, err := proxy.GetInstance(registered)
	assert.NoError(t, err)

	assert.Equal(t, &mockProxy{
		parent: proxy.instance,
	}, resolved)
}

func TestInstanceDomainProxy_GetReference(t *testing.T) {
	parent := &mockParent{}
	proxy, err := newInstanceDomainProxy(parent)
	assert.NoError(t, err)

	reference := proxy.GetReference()
	// TODO should return actual reference
	assert.Nil(t, reference)
}

func TestInstanceDomainProxy_GetInstance_IncorrectRef(t *testing.T) {
	parent := &mockParent{}
	proxy, err := newInstanceDomainProxy(parent)
	assert.NoError(t, err)

	_, err = proxy.GetInstance("1")
	assert.EqualError(t, err, "object with record 1 does not exist")
}

func TestInstanceDomainProxy_GetParent(t *testing.T) {
	parent := &mockParent{}
	proxy, err := newInstanceDomainProxy(parent)
	assert.NoError(t, err)

	returnedParent := proxy.GetParent()
	assert.Equal(t, parent, returnedParent)
}

func TestInstanceDomainProxy_GetClassID(t *testing.T) {
	parent := &mockParent{}
	proxy, err := newInstanceDomainProxy(parent)
	assert.NoError(t, err)

	assert.Equal(t, class.InstanceDomainID, proxy.GetClassID())
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

func TestInstanceDomainFactory_GetReference(t *testing.T) {
	parent := &mockParent{}
	factory := NewInstanceDomainFactory(parent)
	assert.Nil(t, factory.GetReference())
}

func TestInstanceDomainFactory_Create(t *testing.T) {
	parent := &mockParent{}
	factory := NewInstanceDomainFactory(parent)
	proxy, err := factory.Create(parent)
	instDom, _ := newInstanceDomain(parent)

	assert.NoError(t, err)
	assert.Equal(t, &instanceDomainProxy{
		instance: instDom,
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
