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
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/insolar/insolar/genesis/model/object"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type mockChild struct {
	Reference      object.Reference
	ContextStorage storage.Storage
	parent         object.Parent
}

func (c *mockChild) GetClassID() string {
	return "mockChild"
}

func (c *mockChild) GetReference() object.Reference {
	return c.Reference
}

func (c *mockChild) GetParent() object.Parent {
	return c.parent
}

var child = &mockChild{}

type mockParent struct {
	Reference      object.Reference
	ContextStorage storage.Storage
	parent         object.Parent
}

func (p *mockParent) GetParent() object.Parent {
	return p.parent
}

func (p *mockParent) GetClassID() string {
	return "mockParent"
}

func (p *mockParent) GetReference() object.Reference {
	return p.Reference
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
	parent object.Parent
}

func (p *mockProxy) GetClassID() string {
	return "mockProxy"
}

func (p *mockProxy) GetReference() object.Reference {
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

func (f *mockFactory) GetReference() object.Reference {
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

func TestNewMemberDomain(t *testing.T) {
	parent := &mockParent{}
	mDomain := newMemberDomain(parent)

	assert.Equal(t, &memberDomain{
		BaseDomain: *domain.NewBaseDomain(parent, MemberDomainName),
	}, mDomain)
}

func TestMemberDomain_GetClassID(t *testing.T) {
	mDomain := newMemberDomain(nil)
	domainID := mDomain.GetClassID()
	assert.Equal(t, MemberDomainID, domainID)
}

func TestMemberDomain_CreateMember(t *testing.T) {
	parent := &mockParent{}
	mDomain := newMemberDomain(parent)

	factory := &mockFactory{}
	member, err := mDomain.CreateMember(factory)
	assert.NoError(t, err)

	_, err = uuid.FromString(member)
	assert.NoError(t, err)
}

func TestMemberDomain_CreateMember_WithError(t *testing.T) {
	parent := &mockParent{}
	mDomain := newMemberDomain(parent)

	factory := &mockFactoryError{}
	_, err := mDomain.CreateMember(factory)
	assert.EqualError(t, err, "factory create error")
}

func TestMemberDomain_CreateMember_WithNilError(t *testing.T) {
	parent := &mockParent{}
	mDomain := newMemberDomain(parent)

	factory := &mockFactoryNilError{}
	_, err := mDomain.CreateMember(factory)
	assert.EqualError(t, err, "factory returns nil")
}

func TestMemberDomain_GetMember(t *testing.T) {
	parent := &mockParent{}
	mDomain := newMemberDomain(parent)

	factory := &mockFactory{}
	member, err := mDomain.CreateMember(factory)
	assert.NoError(t, err)

	resolved, err := mDomain.GetMember(member)
	assert.NoError(t, err)

	assert.Equal(t, &mockProxy{
		parent: mDomain,
	}, resolved)
}

func TestMemberDomain_GetMember_IncorrectRef(t *testing.T) {
	parent := &mockParent{}
	mDomain := newMemberDomain(parent)

	_, err := mDomain.GetMember("1")
	assert.EqualError(t, err, "object with record 1 does not exist")
}

func TestNewMemberDomainProxy(t *testing.T) {
	parent := &mockParent{}
	mDomainProxy := newMemberDomainProxy(parent)

	assert.Equal(t, &memberDomainProxy{
		instance: newMemberDomain(parent),
	}, mDomainProxy)
}

func TestMemberDomainProxy_GetReference(t *testing.T) {
	parent := &mockParent{}
	mDomainProxy := newMemberDomainProxy(parent)

	reference := mDomainProxy.GetReference()
	// TODO should return actual reference
	assert.Nil(t, reference)
}

func TestMemberDomainProxy_GetParent(t *testing.T) {
	parent := &mockParent{}
	mDomainProxy := newMemberDomainProxy(parent)

	returnedParent := mDomainProxy.GetParent()
	assert.Equal(t, parent, returnedParent)
}

func TestMemberDomainProxy_GetClassID(t *testing.T) {
	parent := &mockParent{}
	mDomainProxy := newMemberDomainProxy(parent)

	id := mDomainProxy.GetClassID()
	assert.Equal(t, MemberDomainID, id)
}

func TestNewMemberDomainFactory(t *testing.T) {
	expected := &memberDomainFactory{}
	factory := NewMemberDomainFactory()

	assert.Equal(t, expected, factory)
}

func TestMemberDomainFactory_GetClassID(t *testing.T) {
	factory := NewMemberDomainFactory()
	id := factory.GetClassID()

	assert.Equal(t, MemberDomainID, id)
}

func TestMemberDomainFactory_GetReference(t *testing.T) {
	factory := NewMemberDomainFactory()
	reference := factory.GetReference()

	assert.Nil(t, reference)
}

func TestMemberDomainFactory_Create(t *testing.T) {
	parent := &mockParent{}
	factory := NewMemberDomainFactory()
	proxy, err := factory.Create(parent)
	mDomain := newMemberDomain(parent)

	assert.NoError(t, err)
	assert.Equal(t, &memberDomainProxy{
		instance: mDomain,
	}, proxy)
}

func TestMemberDomainFactory_CreateWithError(t *testing.T) {
	parent := &mockParentWithError{}
	factory := NewMemberDomainFactory()
	proxy, err := factory.Create(parent)

	assert.EqualError(t, err, "add child error")
	assert.Nil(t, proxy)
}
