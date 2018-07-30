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
	"testing"

	"github.com/insolar/insolar/genesis/mock/storage"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	mockParent := &mockParent{}
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
	mockParent := &mockParent{}
	resolverHandler := NewHandler(mockParent)

	obj, err := resolverHandler.GetObject("not reference", "mockChild")

	assert.EqualError(t, err, "reference is not Reference class object")
	assert.Nil(t, obj)
}

func TestHandler_GetObject_GlobalScope(t *testing.T) {
	mockParent := &mockParent{}
	resolverHandler := NewHandler(nil)
	newMap := make(map[string]object.Proxy)
	resolverHandler.InitGlobalMap(&newMap)

	ref, _ := object.NewReference("1", "123", object.GlobalScope)
	(*GlobalResolver.globalInstanceMap)["123"] = mockParent

	obj, err := resolverHandler.GetObject(ref, "mockChild")

	assert.NoError(t, err)
	assert.Equal(t, child, obj)
}

func TestHandler_GetObject_ChildScope(t *testing.T) {
	mockParent := &mockParent{}
	resolverHandler := NewHandler(mockParent)
	ref, _ := object.NewReference("1", "1", object.ChildScope)

	obj, err := resolverHandler.GetObject(ref, "mockChild")

	assert.NoError(t, err)
	assert.Equal(t, child, obj)
}

func TestHandler_GetObject_ContextScope(t *testing.T) {
	contextStorage := storage.NewMapStorage()
	record, _ := contextStorage.Set(child)
	mockParent := &mockParent{
		ContextStorage: contextStorage,
	}
	resolverHandler := NewHandler(mockParent)
	ref, _ := object.NewReference(record, "1", object.ContextScope)

	obj, err := resolverHandler.GetObject(ref, "mockChild")

	assert.NoError(t, err)
	assert.Equal(t, child, obj)
}

func TestHandler_GetObject_default(t *testing.T) {
	mockParent := &mockParent{}
	resolverHandler := NewHandler(mockParent)
	ref := &object.Reference{
		Scope: object.ScopeType(10000),
	}

	obj, err := resolverHandler.GetObject(ref, "mockChild")

	assert.EqualError(t, err, "unknown scope type: 10000")
	assert.Nil(t, obj)
}

func TestHandler_SetGlobalMap(t *testing.T) {
	resolverHandler := NewHandler(nil)
	resolverHandler.globalResolver.globalInstanceMap = nil

	newMap := make(map[string]object.Proxy)
	resolverHandler.InitGlobalMap(&newMap)

	assert.Equal(t, &newMap, resolverHandler.globalResolver.globalInstanceMap)
}
