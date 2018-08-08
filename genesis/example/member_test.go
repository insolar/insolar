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

	"github.com/insolar/insolar/genesis/mock/storage"
	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/insolar/insolar/genesis/model/resolver"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type mockCallable struct {
	reference object.Reference
}

func (c *mockCallable) GetReference() object.Reference {
	return c.reference
}

func (c *mockCallable) SetReference(reference object.Reference) {
	c.reference = reference
}

type mockChild struct {
	mockCallable
	ContextStorage storage.Storage
	parent         object.Parent
}

func (c *mockChild) GetClassID() string {
	return "mockChild"
}

func (c *mockChild) GetParent() object.Parent {
	return c.parent
}

var child = &mockChild{}

type mockParent struct {
	mockCallable
	ContextStorage storage.Storage
	parent         object.Parent
}

func (p *mockParent) GetParent() object.Parent {
	return p.parent
}

func (p *mockParent) GetClassID() string {
	return "mockParent"
}

func (p *mockParent) GetChildStorage() storage.Storage {
	return nil
}

func (p *mockParent) AddChild(child object.Child) (string, error) {
	return "", nil
}

func (p *mockParent) GetChild(key string) (object.Child, error) {
	return child, nil
}

func (p *mockParent) GetContext() []string {
	return []string{}
}

func (p *mockParent) GetContextStorage() storage.Storage {
	return p.ContextStorage
}

type mockParentWithError struct {
	mockParent
}

func (p *mockParentWithError) AddChild(child object.Child) (string, error) {
	return "", fmt.Errorf("add child error")
}

type mockProxy struct {
	mockCallable
	parent object.Parent
}

func (p *mockProxy) GetClassID() string {
	return "mockProxy"
}

func (p *mockProxy) GetParent() object.Parent {
	return p.parent
}

type mockFactory struct {
	mockCallable
}

func (f *mockFactory) Create(parent object.Parent) (resolver.Proxy, error) {
	return &mockProxy{
		parent: parent,
	}, nil
}

func (f *mockFactory) GetClassID() string {
	return "mockFactory"
}

func (f *mockFactory) GetParent() object.Parent {
	return nil
}

type mockFactoryError struct {
	mockFactory
}

func (f *mockFactoryError) Create(parent object.Parent) (resolver.Proxy, error) {
	return nil, fmt.Errorf("factory create error")
}

type mockFactoryNilError struct {
	mockFactory
}

func (f *mockFactoryNilError) Create(parent object.Parent) (resolver.Proxy, error) {
	return nil, nil
}

func TestNewMemberDomain(t *testing.T) {
	parent := &mockParent{}
	mDomain, err := newMemberDomain(parent)

	assert.NoError(t, err)
	assert.NotNil(t, mDomain)
	assert.NotNil(t, mDomain.(*memberDomain).memberFactoryRecord)
	assert.NotEmpty(t, mDomain.(*memberDomain).ChildStorage.GetKeys())
}

func TestNewMemberDomain_WithNilParent(t *testing.T) {
	_, err := newMemberDomain(nil)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestMemberDomain_GetClassID(t *testing.T) {
	parent := &mockParent{}
	mDomain, _ := newMemberDomain(parent)
	domainID := mDomain.GetClassID()
	assert.Equal(t, class.MemberDomainID, domainID)
}

func TestMemberDomain_CreateMember(t *testing.T) {
	parent := &mockParent{}
	mDomain, _ := newMemberDomain(parent)

	recordID, err := mDomain.CreateMember()
	assert.NoError(t, err)

	_, err = uuid.FromString(recordID)
	assert.NoError(t, err)
}

func TestMemberDomain_CreateMember_WithError(t *testing.T) {
	parent := &mockParent{}
	factoryWithError := &mockFactoryError{}

	mDomain := &memberDomain{
		BaseDomain: *domain.NewBaseDomain(parent, MemberDomainName),
	}
	record, _ := mDomain.AddChild(factoryWithError)
	mDomain.memberFactoryRecord = record

	_, err := mDomain.CreateMember()
	assert.EqualError(t, err, "factory create error")
}

func TestMemberDomain_CreateMember_WithNilError(t *testing.T) {
	parent := &mockParent{}
	factoryWithNil := &mockFactoryNilError{}

	mDomain := &memberDomain{
		BaseDomain: *domain.NewBaseDomain(parent, MemberDomainName),
	}
	record, _ := mDomain.AddChild(factoryWithNil)
	mDomain.memberFactoryRecord = record
	_, err := mDomain.CreateMember()
	assert.EqualError(t, err, "factory returns nil")
}

func TestMemberDomain_CreateMember_NoMemberFactoryRecord(t *testing.T) {
	parent := &mockParent{}

	mDomain := &memberDomain{
		BaseDomain: *domain.NewBaseDomain(parent, MemberDomainName),
	}
	mDomain.memberFactoryRecord = "unexistedRecord"
	_, err := mDomain.CreateMember()
	assert.EqualError(t, err, "object with record unexistedRecord does not exist")
}

func TestMemberDomain_CreateMember_NotFactory(t *testing.T) {
	parent := &mockParent{}

	mDomain := &memberDomain{
		BaseDomain: *domain.NewBaseDomain(parent, MemberDomainName),
	}
	record, _ := mDomain.AddChild(parent)
	mDomain.memberFactoryRecord = record
	_, err := mDomain.CreateMember()
	assert.EqualError(t, err, fmt.Sprintf("child by record `%s` is not Factory instance", record))
}
func TestMemberDomain_CreateMember_NotMember(t *testing.T) {
	parent := &mockParent{}
	factory := &mockFactory{}

	mDomain := &memberDomain{
		BaseDomain: *domain.NewBaseDomain(parent, MemberDomainName),
	}
	record, _ := mDomain.AddChild(factory)
	mDomain.memberFactoryRecord = record
	_, err := mDomain.CreateMember()
	assert.EqualError(t, err, "factory returns not Member instance")
}

func TestMemberDomain_GetMember(t *testing.T) {
	parent := &mockParent{}
	mDomain, _ := newMemberDomain(parent)

	recordID, err := mDomain.CreateMember()
	assert.NoError(t, err)

	resolved, err := mDomain.GetMember(recordID)
	assert.NoError(t, err)

	_, ok := resolved.(*memberProxy)
	assert.True(t, ok)
}

func TestMemberDomain_GetMember_IncorrectRef(t *testing.T) {
	parent := &mockParent{}
	mDomain, _ := newMemberDomain(parent)

	_, err := mDomain.GetMember("1")
	assert.EqualError(t, err, "object with record 1 does not exist")
}

func TestNewMemberDomainProxy(t *testing.T) {
	parent := &mockParent{}

	proxy, err := newMemberDomainProxy(parent)

	assert.NoError(t, err)
	assert.NotNil(t, proxy)
	_, ok := proxy.Instance.(*memberDomain)
	assert.True(t, ok)
}

func TestNewMemberDomainProxy_WithNilParent(t *testing.T) {
	_, err := newMemberDomainProxy(nil)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestMemberDomainProxy_CreateMember(t *testing.T) {
	parent := &mockParent{}
	proxy, _ := newMemberDomainProxy(parent)

	recordID, err := proxy.CreateMember()
	assert.NoError(t, err)

	_, err = uuid.FromString(recordID)
	assert.NoError(t, err)
}

func TestMemberDomainProxy_GetMember(t *testing.T) {
	parent := &mockParent{}
	proxy, _ := newMemberDomainProxy(parent)

	recordID, err := proxy.CreateMember()
	assert.NoError(t, err)

	resolved, err := proxy.GetMember(recordID)
	assert.NoError(t, err)

	_, ok := resolved.(*memberProxy)
	assert.True(t, ok)
}

func TestNewMemberDomainFactory(t *testing.T) {
	parent := &mockParent{}
	expected := &memberDomainFactory{parent: parent}

	factory := NewMemberDomainFactory(parent)

	assert.Equal(t, expected, factory)
}

func TestMemberDomainFactory_GetClassID(t *testing.T) {
	parent := &mockParent{}
	factory := NewMemberDomainFactory(parent)
	id := factory.GetClassID()

	assert.Equal(t, class.MemberDomainID, id)
}

func TestMemberDomainFactory_GetParent(t *testing.T) {
	parent := &mockParent{}
	factory := NewMemberDomainFactory(parent)
	reference := factory.GetParent()

	assert.Nil(t, reference)
}

func TestMemberDomainFactory_Create(t *testing.T) {
	parent := &mockParent{}
	factory := NewMemberDomainFactory(parent)
	proxy, err := factory.Create(parent)

	assert.NotNil(t, proxy)
	assert.NoError(t, err)
}

func TestMemberDomainFactory_CreateWithNoParent(t *testing.T) {
	parent := &mockParent{}
	factory := NewMemberDomainFactory(parent)
	proxy, err := factory.Create(nil)

	assert.EqualError(t, err, "parent must not be nil")
	assert.Nil(t, proxy)
}

func TestMemberDomainFactory_CreateWithError(t *testing.T) {
	parent := &mockParentWithError{}
	factory := NewMemberDomainFactory(parent)
	proxy, err := factory.Create(parent)

	assert.EqualError(t, err, "add child error")
	assert.Nil(t, proxy)
}
