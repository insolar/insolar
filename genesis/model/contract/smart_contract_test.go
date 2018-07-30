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

package contract

import (
	"testing"

	"github.com/insolar/insolar/genesis/mock/storage"
	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/insolar/insolar/genesis/model/resolver"
	"github.com/stretchr/testify/assert"
)

type mockChild struct {
	Reference *object.Reference
	parent    object.Parent
}

func (c *mockChild) GetClassID() string {
	return "mockChild"
}

func (c *mockChild) GetReference() *object.Reference {
	return c.Reference
}

func (c *mockChild) GetParent() object.Parent {
	return c.parent
}

type mockParent struct {
	Reference      *object.Reference
	ContextStorage storage.Storage
}

func (p *mockParent) GetClassID() string {
	return "mockParent"
}

func (p *mockParent) GetReference() *object.Reference {
	return p.Reference
}

func (p *mockParent) GetChildStorage() storage.Storage {
	return nil
}

func (p *mockParent) AddChild(child object.Child) (string, error) {
	return "", nil
}

func (p *mockParent) GetChild(key string) (object.Child, error) {
	return nil, nil
}

func (p *mockParent) GetContext() []string {
	return []string{}
}

func (p *mockParent) GetContextStorage() storage.Storage {
	return p.ContextStorage
}

type BaseComposite struct{}

func (c *BaseComposite) GetInterfaceKey() string {
	return "BaseComposite"
}

func (c *BaseComposite) GetClassID() string {
	return "BaseComposite"
}

type anotherBaseComposite struct{}

func (c *anotherBaseComposite) GetInterfaceKey() string {
	return "anotherBaseComposite"
}

func (c *anotherBaseComposite) GetClassID() string {
	return "anotherBaseComposite"
}

type BaseCompositeFactory struct{}

func (cf *BaseCompositeFactory) Create() object.Composite {
	return &BaseComposite{}
}

func TestNewBaseSmartContract(t *testing.T) {
	parent := &mockParent{}
	childStorage := storage.NewMapStorage()
	sc := NewBaseSmartContract(parent)

	assert.Equal(t, &BaseSmartContract{
		CompositeMap: make(map[string]object.Composite),
		ChildStorage: childStorage,
		Parent:       parent,
		resolver:     resolver.NewHandler(sc),
	}, sc)
}

func TestSmartContract_GetClassID(t *testing.T) {
	parent := &mockParent{}
	sc := NewBaseSmartContract(parent)

	classID := sc.GetClassID()

	assert.Equal(t, class.SmartContractID, classID)
}

func TestSmartContract_GetReference(t *testing.T) {
	parent := &mockParent{}
	sc := NewBaseSmartContract(parent)
	ref := sc.GetReference()

	assert.Nil(t, ref)
}

func TestSmartContract_CreateComposite(t *testing.T) {
	parent := &mockParent{}
	sc := NewBaseSmartContract(parent)
	compositeFactory := BaseCompositeFactory{}

	composite, err := sc.CreateComposite(&compositeFactory)

	assert.Len(t, sc.CompositeMap, 1)
	assert.Equal(t, sc.CompositeMap[composite.GetInterfaceKey()], composite)
	assert.NoError(t, err)
}

func TestSmartContract_CreateComposite_Error(t *testing.T) {
	parent := &mockParent{}
	sc := NewBaseSmartContract(parent)
	compositeFactory := BaseCompositeFactory{}
	sc.CreateComposite(&compositeFactory)

	res, err := sc.CreateComposite(&compositeFactory)

	assert.Nil(t, res)
	assert.EqualError(t, err, "delegate with name BaseComposite already exist")
}

func TestSmartContract_GetComposite(t *testing.T) {
	parent := &mockParent{}
	sc := NewBaseSmartContract(parent)
	compositeFactory := BaseCompositeFactory{}
	composite, _ := sc.CreateComposite(&compositeFactory)

	res, err := sc.GetComposite(composite.GetInterfaceKey())

	assert.Equal(t, composite, res)
	assert.NoError(t, err)
}

func TestSmartContract_GetComposite_Error(t *testing.T) {
	parent := &mockParent{}
	sc := NewBaseSmartContract(parent)
	composite := BaseComposite{}

	res, err := sc.GetComposite(composite.GetInterfaceKey())

	assert.Nil(t, res)
	assert.EqualError(t, err, "delegate with name BaseComposite does not exist")
}

func TestSmartContract_GetOrCreateComposite_Get(t *testing.T) {
	parent := &mockParent{}
	sc := NewBaseSmartContract(parent)
	composite := &BaseComposite{}
	compositeFactory := &BaseCompositeFactory{}

	res, err := sc.GetOrCreateComposite(composite.GetInterfaceKey(), compositeFactory)

	assert.NoError(t, err)
	assert.Equal(t, composite, res)
}

func TestSmartContract_GetOrCreateComposite_Create(t *testing.T) {
	parent := &mockParent{}
	sc := NewBaseSmartContract(parent)
	composite := &BaseComposite{}
	compositeFactory := &BaseCompositeFactory{}

	res, err := sc.GetOrCreateComposite(composite.GetInterfaceKey(), compositeFactory)

	assert.Len(t, sc.CompositeMap, 1)
	assert.Equal(t, sc.CompositeMap[composite.GetInterfaceKey()], res)
	assert.NoError(t, err)
}

func TestSmartContract_GetOrCreateComposite_Error(t *testing.T) {
	parent := &mockParent{}
	sc := NewBaseSmartContract(parent)
	composite := &anotherBaseComposite{}
	compositeFactory := &BaseCompositeFactory{}
	sc.CreateComposite(compositeFactory)

	res, err := sc.GetOrCreateComposite(composite.GetInterfaceKey(), compositeFactory)

	assert.Nil(t, res)
	assert.EqualError(t, err, "delegate with name BaseComposite already exist")
}

func TestSmartContract_GetChildStorage(t *testing.T) {
	parent := &mockParent{}
	sc := NewBaseSmartContract(parent)

	res := sc.GetChildStorage()

	assert.Equal(t, sc.ChildStorage, res)
}

func TestSmartContract_AddChild(t *testing.T) {
	parent := &mockParent{}
	sc := NewBaseSmartContract(parent)
	child := &mockChild{}

	res, err := sc.AddChild(child)

	assert.NoError(t, err)
	assert.Len(t, sc.ChildStorage.GetKeys(), 1)
	assert.Equal(t, sc.ChildStorage.GetKeys()[0], res)
}

func TestSmartContract_GetChild(t *testing.T) {
	parent := &mockParent{}
	sc := NewBaseSmartContract(parent)
	child := &mockChild{}
	key, _ := sc.AddChild(child)

	res, err := sc.GetChild(key)

	assert.NoError(t, err)
	assert.Equal(t, child, res)
}

func TestSmartContract_GetChild_Error(t *testing.T) {
	parent := &mockParent{}
	sc := NewBaseSmartContract(parent)

	res, err := sc.GetChild("someKey")

	assert.Nil(t, res)
	assert.EqualError(t, err, "object with record someKey does not exist")
}

func TestSmartContract_GetContextStorage(t *testing.T) {
	parent := &mockParent{}
	sc := NewBaseSmartContract(parent)

	res := sc.GetContextStorage()

	assert.Equal(t, sc.ContextStorage, res)
}

func TestSmartContract_GetContext(t *testing.T) {
	parent := &mockParent{}
	contextStorage := storage.NewMapStorage()
	sc := NewBaseSmartContract(parent)
	sc.ContextStorage = contextStorage

	res := sc.GetContext()

	assert.Equal(t, contextStorage.GetKeys(), res)
}

func TestSmartContract_GetParent(t *testing.T) {
	parent := &mockParent{}
	sc := NewBaseSmartContract(parent)

	res := sc.GetParent()

	assert.Equal(t, sc.Parent, res)
}

func TestSmartContract_GetResolver(t *testing.T) {
	parent := &mockParent{}
	sc := BaseSmartContract{
		CompositeMap: make(map[string]object.Composite),
		ChildStorage: storage.NewMapStorage(),
		Parent:       parent,
	}
	assert.Nil(t, sc.resolver)
	sc.GetResolver()

	assert.NotNil(t, sc.resolver)
}
