/*
 *    Copyright 2018 Insolar
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

package object

import (
	"testing"

	"github.com/insolar/insolar/application/model/class"
	"github.com/stretchr/testify/assert"
)

func TestNewReferenceContainer(t *testing.T) {
	domain := "134"
	record := "156"
	ref, _ := NewReference(domain, record, GlobalScope)
	container := NewReferenceContainer(ref)
	assert.Equal(t, &ReferenceContainer{
		storedReference: ref,
	}, container)

}

func TestReferenceContainer_GetStoredReference(t *testing.T) {
	domain := "134"
	record := "156"
	ref, _ := NewReference(domain, record, GlobalScope)
	container := &ReferenceContainer{
		storedReference: ref,
	}

	assert.Equal(t, ref, container.GetStoredReference())
}

func TestReferenceContainer_GetClassID(t *testing.T) {
	domain := "134"
	record := "156"
	ref, _ := NewReference(domain, record, GlobalScope)
	container := &ReferenceContainer{
		storedReference: ref,
	}

	refID := container.GetClassID()
	assert.Equal(t, class.ReferenceID, refID)
}
