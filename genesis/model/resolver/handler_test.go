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

package resolver

import (
	"fmt"
	"testing"

	"github.com/insolar/insolar/genesis/mock/storage"
	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/stretchr/testify/assert"
)

type invalidScopeReference struct{}

func (r *invalidScopeReference) GetClassID() string {
	return class.ReferenceID
}

func (r *invalidScopeReference) GetRecord() string {
	return "145"
}

func (r *invalidScopeReference) GetDomain() string {
	return "123"
}

func (r *invalidScopeReference) GetScope() object.ScopeType {
	return object.ScopeType(10000)
}

func (r *invalidScopeReference) String() string {
	return fmt.Sprintf("#%s.#%s", "145", "123")
}

func (r *invalidScopeReference) GetReference() object.Reference {
	return r
}

func (r *invalidScopeReference) SetReference(ref object.Reference) {
}

func (r *invalidScopeReference) GetParent() object.Parent {
	return nil
}

func TestNewHandler(t *testing.T) {
	mockParent := &mockParentProxy{}
	handler := NewHandler(mockParent)

	assert.Equal(t, &Handler{
		globalResolver: GlobalResolver,
		childResolver: &childResolver{
			parent: mockParent,
		},
		contextResolver: &contextResolver{
			parent: mockParent,
		},
	}, handler)
}

func TestHandler_GetObject_Not_Reference(t *testing.T) {
	mockParent := &mockParentProxy{}
	resolverHandler := NewHandler(mockParent)

	obj, err := resolverHandler.GetObject("not reference", "mockChild")

	assert.EqualError(t, err, "reference is not Reference class object")
	assert.Nil(t, obj)
}

func TestHandler_GetObject_GlobalScope(t *testing.T) {
	mockParent := &mockParentProxy{}
	resolverHandler := NewHandler(nil)
	newMap := make(map[string]Proxy)
	resolverHandler.InitGlobalMap(&newMap)

	ref, _ := object.NewReference("123", "1", object.GlobalScope)
	(*GlobalResolver.globalInstanceMap)["123"] = mockParent

	obj, err := resolverHandler.GetObject(ref, "mockChild")

	assert.NoError(t, err)
	assert.Equal(t, child, obj)
}

func TestHandler_GetObject_ChildScope(t *testing.T) {
	mockParent := &mockParentProxy{}
	resolverHandler := NewHandler(mockParent)
	ref, _ := object.NewReference("1", "1", object.ChildScope)

	obj, err := resolverHandler.GetObject(ref, "mockChild")

	assert.NoError(t, err)
	assert.Equal(t, child, obj)
}

func TestHandler_GetObject_ContextScope(t *testing.T) {
	contextStorage := storage.NewMapStorage()
	record, _ := contextStorage.Set(child)
	mockParent := &mockParentProxy{
		ContextStorage: contextStorage,
	}
	resolverHandler := NewHandler(mockParent)
	ref, _ := object.NewReference("1", record, object.ContextScope)

	obj, err := resolverHandler.GetObject(ref, "mockChild")

	assert.NoError(t, err)
	assert.Equal(t, child, obj)
}

func TestHandler_GetObject_default(t *testing.T) {
	mockParent := &mockParentProxy{}
	resolverHandler := NewHandler(mockParent)
	ref := &invalidScopeReference{}

	obj, err := resolverHandler.GetObject(ref, "mockChild")

	assert.EqualError(t, err, "unknown scope type: 10000")
	assert.Nil(t, obj)
}

func TestHandler_SetGlobalMap(t *testing.T) {
	resolverHandler := NewHandler(nil)
	resolverHandler.globalResolver.globalInstanceMap = nil

	newMap := make(map[string]Proxy)
	resolverHandler.InitGlobalMap(&newMap)

	assert.Equal(t, &newMap, resolverHandler.globalResolver.globalInstanceMap)
}
