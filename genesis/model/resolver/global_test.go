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
	instanceMap := make(map[*object.Reference]object.Proxy)

	assert.Equal(t, &globalResolver{
		globalInstanceMap: &instanceMap,
	}, resolver)
}

func TestGlobalResolver_GetObject_No_Object(t *testing.T) {
	resolver := newGlobalResolver()
	ref, _ := object.NewReference("1", "1", object.GlobalScope)

	obj, err := resolver.GetObject(ref, "someClass")

	assert.Equal(t, "reference with address `#1.#1` not found", err.Error())
	assert.Nil(t, obj)
}

func TestGlobalResolver_GetObject_Wrong_classID(t *testing.T) {
	mockParent := &mockParent{}
	resolver := newGlobalResolver()
	ref, _ := object.NewReference("1", "1", object.GlobalScope)
	(*resolver.globalInstanceMap)[ref] = mockParent

	obj, err := resolver.GetObject(ref, "someClass")

	assert.Equal(t, "instance class is not `someClass`", err.Error())
	assert.Nil(t, obj)
}

func TestGlobalResolver_GetObject(t *testing.T) {
	mockParent := &mockParent{}
	resolver := newGlobalResolver()
	ref, _ := object.NewReference("1", "1", object.GlobalScope)
	(*resolver.globalInstanceMap)[ref] = mockParent

	obj, err := resolver.GetObject(ref, "mockParent")

	assert.Nil(t, err)
	assert.Equal(t, mockParent, obj)
}
