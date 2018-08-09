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

//<<<<<<< HEAD
//type mockCallable struct {
//reference object.Reference
//}

/*func (c *mockCallable) GetReference() object.Reference {
	return c.reference
}

func (c *mockCallable) SetReference(reference object.Reference) {
	c.reference = reference
}*/

/*type mockChild struct {
	//mockCallable
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
	//mockCallable
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
	//mockCallable
	parent    object.Parent
	reference object.Reference
}

func (p *mockProxy) GetClassID() string {
	return "mockProxy"
}

func (p *mockProxy) GetParent() object.Parent {
	return p.parent
}

func (p *mockProxy) GetReference() object.Reference {
	return p.reference
}

func (p *mockProxy) SetReference(reference object.Reference) {
	p.reference = reference
}

type mockFactory struct {
	//mockCallable
	reference object.Reference
}

func (f *mockFactory) Create(parent object.Parent) (resolver.Proxy, error) {
	return &mockProxy{
		parent: parent,
	}, nil
}
=======*/
type BaseComposite struct{}

//>>>>>>> c58cdcfc979b429cafb9cc4fad009c35a8c990ff

func (c *BaseComposite) GetInterfaceKey() string {
	return "BaseComposite"
}

func (c *BaseComposite) GetClassID() string {
	return "BaseComposite"
}

/*<<<<<<< HEAD
func (f *mockFactory) GetReference() object.Reference {
	return f.reference
}

func (f *mockFactory) SetReference(reference object.Reference) {
	f.reference = reference
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
=======*/
type BaseCompositeFactory struct{}

//>>>>>>> c58cdcfc979b429cafb9cc4fad009c35a8c990ff

func (cf *BaseCompositeFactory) Create() (object.Composite, error) {
	return &BaseComposite{}, nil
}

func TestNewMember(t *testing.T) {
	parent := &mockParent{}
	testMember, err := newMember(parent)

	assert.NoError(t, err)
	expectedMember := &member{
		BaseSmartContract: *contract.NewBaseSmartContract(parent),
	}
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

	expectedMember := &member{
		BaseSmartContract: *contract.NewBaseSmartContract(parent),
	}
	expectedMember.CompositeMap = make(map[string]object.Composite)
	assert.Equal(t, &memberProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
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

	/*<<<<<<< HEAD
		proxy, err := newMemberDomainProxy(parent)
		assert.NoError(t, err)

		assert.Equal(t, &memberDomainProxy{
			BaseSmartContractProxy: contract.BaseSmartContractProxy{
				Instance: mDomain,
			},
		}, proxy)
	=======*/
	username := proxy.GetUsername()
	assert.Equal(t, "", username)
	//>>>>>>> c58cdcfc979b429cafb9cc4fad009c35a8c990ff
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
	/*<<<<<<< HEAD
		assert.Equal(t, &memberDomainProxy{
			BaseSmartContractProxy: contract.BaseSmartContractProxy{
				Instance: mDomain,
	=======*/

	expectedMember := &member{
		BaseSmartContract: *contract.NewBaseSmartContract(parent),
	}
	expectedMember.CompositeMap = make(map[string]object.Composite)
	assert.Equal(t, &memberProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: expectedMember,
			//>>>>>>> c58cdcfc979b429cafb9cc4fad009c35a8c990ff
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
