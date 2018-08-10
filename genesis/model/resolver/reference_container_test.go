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

	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/stretchr/testify/assert"
)

func TestNewReferenceContainer(t *testing.T) {
	domain := "134"
	record := "156"
	ref, _ := object.NewReference(domain, record, object.GlobalScope)
	container := NewReferenceContainer(ref)
	assert.Equal(t, ref, container.reference)

}

func TestReferenceContainer_GetStoredReference(t *testing.T) {
	domain := "134"
	record := "156"
	ref, _ := object.NewReference(domain, record, object.GlobalScope)
	container := &ReferenceContainer{
		reference: ref,
	}

	assert.Equal(t, ref, container.GetStoredReference())
}

func TestReferenceContainer_GetClassID(t *testing.T) {
	domain := "134"
	record := "156"
	ref, _ := object.NewReference(domain, record, object.GlobalScope)
	container := &ReferenceContainer{
		reference: ref,
	}

	refID := container.GetClassID()
	assert.Equal(t, class.ReferenceID, refID)
}
