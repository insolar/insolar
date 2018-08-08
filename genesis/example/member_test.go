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
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/insolar/insolar/genesis/model/resolver"
	"github.com/stretchr/testify/assert"
)

type BaseComposite struct{}

func (c *BaseComposite) GetInterfaceKey() string {
	return "BaseComposite"
}

func (c *BaseComposite) GetClassID() string {
	return "BaseComposite"
}

type BaseCompositeFactory struct{}

func (cf *BaseCompositeFactory) Create() (object.Composite, error) {
	return &BaseComposite{}, nil
}

func TestNewMember(t *testing.T) {
	parent := &mockParent{}
	testMember, err := newMember(parent)

	assert.NoError(t, err)
	expectedMember := &member{}
	expectedMember.CompositeMap = make(map[string]object.Composite)
	assert.Equal(t, expectedMember, testMember)
}

func TestNewMember_WithNilParent(t *testing.T) {
	_, err := newMember(nil)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestMember_GetClassID(t *testing.T) {
	parent := &mockParent{}
	testMember, _ := newMember(parent)

	memberID := testMember.GetClassID()
	assert.Equal(t, class.MemberID, memberID)
}

func TestMember_GetUsername(t *testing.T) {
	parent := &mockParent{}
	testMember, _ := newMember(parent)

	username := testMember.GetUsername()
	assert.Equal(t, "", username)
}

func TestMember_GetPublicKey(t *testing.T) {
	parent := &mockParent{}
	testMember, _ := newMember(parent)

	publicKey := testMember.GetPublicKey()
	assert.Equal(t, "", publicKey)
}

func TestNewMemberProxy(t *testing.T) {
	parent := &mockParent{}
	_, err := newMember(parent)
	assert.NoError(t, err)

	proxy, err := newMemberProxy(parent)
	assert.NoError(t, err)

	expectedMember := &member{}
	expectedMember.CompositeMap = make(map[string]object.Composite)
	assert.Equal(t, &memberProxy{
		BaseProxy: resolver.BaseProxy{
			Instance: expectedMember,
		},
	}, proxy)
}

func TestNewMemberProxy_WithNilParent(t *testing.T) {
	_, err := newMemberProxy(nil)
	assert.EqualError(t, err, "parent must not be nil")
}

func TestMemberProxy_GetUsername(t *testing.T) {
	parent := &mockParent{}
	proxy, _ := newMemberProxy(parent)

	username := proxy.GetUsername()
	assert.Equal(t, "", username)
}

func TestMemberProxy_GetPublicKey(t *testing.T) {
	parent := &mockParent{}
	proxy, _ := newMemberProxy(parent)

	publicKey := proxy.GetPublicKey()
	assert.Equal(t, "", publicKey)
}

func TestMemberProxy_GetOrCreateComposite_Get(t *testing.T) {
	parent := &mockParent{}
	proxy, _ := newMemberProxy(parent)
	composite := &BaseComposite{}
	compositeFactory := &BaseCompositeFactory{}

	res, err := proxy.GetOrCreateComposite(composite.GetInterfaceKey(), compositeFactory)

	assert.NoError(t, err)
	assert.Equal(t, composite, res)
}

func TestMemberProxy_GetOrCreateComposite_Create(t *testing.T) {
	parent := &mockParent{}
	proxy, _ := newMemberProxy(parent)
	composite := &BaseComposite{}
	compositeFactory := &BaseCompositeFactory{}

	res, err := proxy.GetOrCreateComposite(composite.GetInterfaceKey(), compositeFactory)

	assert.Len(t, proxy.Instance.(*member).CompositeMap, 1)
	assert.Equal(t, proxy.Instance.(*member).CompositeMap[composite.GetInterfaceKey()], res)
	assert.NoError(t, err)
}

func TestNewMemberFactory(t *testing.T) {
	parent := &mockParent{}
	expected := &memberFactory{parent: parent}

	factory := NewMemberFactory(parent)

	assert.Equal(t, expected, factory)
}

func TestMemberFactory_GetClassID(t *testing.T) {
	parent := &mockParent{}
	factory := NewMemberFactory(parent)
	id := factory.GetClassID()

	assert.Equal(t, class.MemberID, id)
}

func TestMemberFactory_GetParent(t *testing.T) {
	parent := &mockParent{}
	factory := NewMemberFactory(parent)
	p := factory.GetParent()

	assert.Equal(t, parent, p)
}

func TestMemberFactory_Create(t *testing.T) {
	parent := &mockParent{}
	factory := NewMemberFactory(parent)

	proxy, err := factory.Create(parent)

	assert.NoError(t, err)

	expectedMember := &member{}
	expectedMember.CompositeMap = make(map[string]object.Composite)
	assert.Equal(t, &memberProxy{
		BaseProxy: resolver.BaseProxy{
			Instance: expectedMember,
		},
	}, proxy)
}

func TestMemberFactory_CreateWithNoParent(t *testing.T) {
	parent := &mockParent{}
	factory := NewMemberFactory(parent)
	proxy, err := factory.Create(nil)

	assert.EqualError(t, err, "parent must not be nil")
	assert.Nil(t, proxy)
}

func TestMemberFactory_CreateWithError(t *testing.T) {
	parent := &mockParentWithError{}
	factory := NewMemberFactory(parent)
	proxy, err := factory.Create(parent)

	assert.EqualError(t, err, "add child error")
	assert.Nil(t, proxy)
}
