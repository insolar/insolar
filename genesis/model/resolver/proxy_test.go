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

func TestBaseProxy_GetClassID(t *testing.T) {
	proxy := &BaseProxy{}
	assert.Equal(t, "Proxy", proxy.GetClassID())
}

func TestBaseProxy_GetParent(t *testing.T) {
	proxy := &BaseProxy{}
	assert.Nil(t, proxy.GetParent())
}

func TestBaseProxy_SetReference(t *testing.T) {
	ref, _ := object.NewReference("1", "2", object.GlobalScope)
	proxy := &BaseProxy{}
	proxy.SetReference(ref)
	assert.Equal(t, ref, proxy.reference)
}

func TestBaseProxy_GetReference(t *testing.T) {
	ref, _ := object.NewReference("1", "2", object.GlobalScope)
	proxy := &BaseProxy{
		reference: ref,
	}
	assert.Equal(t, ref, proxy.GetReference())
}
