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

	"github.com/insolar/insolar/genesis/model/object"
	"github.com/stretchr/testify/assert"
)

func TestNewGlobalResolver(t *testing.T) {
	resolver := newGlobalResolver()

	assert.Equal(t, &globalResolver{
		globalInstanceMap: nil,
	}, resolver)
}

func TestGlobalResolver_GetObject_Not_Reference(t *testing.T) {
	resolver := newGlobalResolver()

	obj, err := resolver.GetObject("not reference", "someClass")

	assert.EqualError(t, err, "reference is not Reference class object")
	assert.Nil(t, obj)
}

func TestGlobalResolver_GetObject_No_Object(t *testing.T) {
	resolver := newGlobalResolver()
	newMap := make(map[string]object.Proxy)
	resolver.InitGlobalMap(&newMap)
	ref, _ := object.NewReference("1", "123", object.GlobalScope)

	obj, err := resolver.GetObject(ref, "someClass")

	assert.EqualError(t, err, "reference with address `#123.#1` not found")
	assert.Nil(t, obj)
}

func TestGlobalResolver_GetObject_Not_Parent(t *testing.T) {
	mockChild := &mockChild{}
	resolver := newGlobalResolver()
	newMap := make(map[string]object.Proxy)
	resolver.InitGlobalMap(&newMap)
	ref, _ := object.NewReference("1", "123", object.GlobalScope)
	(*resolver.globalInstanceMap)["123"] = mockChild

	obj, err := resolver.GetObject(ref, "mockChild")

	assert.EqualError(t, err, "object with domain `123` can not have children")
	assert.Nil(t, obj)
}

func TestGlobalResolver_GetObject_No_Child(t *testing.T) {
	mockParent := &mockParentWithError{}
	resolver := newGlobalResolver()
	newMap := make(map[string]object.Proxy)
	resolver.InitGlobalMap(&newMap)
	ref, _ := object.NewReference("1", "123", object.GlobalScope)
	(*resolver.globalInstanceMap)["123"] = mockParent

	obj, err := resolver.GetObject(ref, "someClass")

	assert.EqualError(t, err, "object with record 1 does not exist")
	assert.Nil(t, obj)
}

func TestGlobalResolver_GetObject_ClassID_Not_Str(t *testing.T) {
	mockParent := &mockParent{}
	resolver := newGlobalResolver()
	newMap := make(map[string]object.Proxy)
	resolver.InitGlobalMap(&newMap)
	ref, _ := object.NewReference("1", "123", object.GlobalScope)
	(*resolver.globalInstanceMap)["123"] = mockParent

	obj, err := resolver.GetObject(ref, ref)

	assert.EqualError(t, err, "classID is not string")
	assert.Nil(t, obj)
}

func TestGlobalResolver_GetObject_Wrong_classID(t *testing.T) {
	mockParent := &mockParent{}
	resolver := newGlobalResolver()
	newMap := make(map[string]object.Proxy)
	resolver.InitGlobalMap(&newMap)
	ref, _ := object.NewReference("1", "123", object.GlobalScope)
	(*resolver.globalInstanceMap)["123"] = mockParent

	obj, err := resolver.GetObject(ref, "someClass")

	assert.EqualError(t, err, "instance class is not `someClass`")
	assert.Nil(t, obj)
}

func TestGlobalResolver_GetObject(t *testing.T) {
	mockParent := &mockParent{}
	resolver := newGlobalResolver()
	newMap := make(map[string]object.Proxy)
	resolver.InitGlobalMap(&newMap)
	ref, _ := object.NewReference("1", "123", object.GlobalScope)
	(*resolver.globalInstanceMap)["123"] = mockParent

	obj, err := resolver.GetObject(ref, "mockChild")

	assert.NoError(t, err)
	assert.Equal(t, child, obj)
}

func TestGlobalResolver_InitGlobalMap(t *testing.T) {
	resolver := newGlobalResolver()

	assert.Nil(t, resolver.globalInstanceMap)

	newMap := make(map[string]object.Proxy)
	resolver.InitGlobalMap(&newMap)

	assert.Equal(t, &newMap, resolver.globalInstanceMap)
}
