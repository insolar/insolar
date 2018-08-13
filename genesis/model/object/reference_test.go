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

package object

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewReference(t *testing.T) {
	domain := "134"
	record := "156"
	ref, err := NewReference(domain, record, GlobalScope)

	assert.NoError(t, err)
	assert.Equal(t, &reference{
		domain: domain,
		record: record,
		scope:  GlobalScope,
	}, ref)
}

func TestNewReference_Error(t *testing.T) {
	domain := "134"
	record := "156"
	unknownScope := ScopeType(100)
	ref, err := NewReference(domain, record, unknownScope)

	assert.EqualError(t, err, "unknown scope type: 100")
	assert.Nil(t, ref)
}

func TestReference_String(t *testing.T) {
	domain := "134"
	record := "156"
	ref, _ := NewReference(domain, record, GlobalScope)

	stringRef := ref.String()

	assert.Equal(t, fmt.Sprintf("#%s.#%s", domain, record), stringRef)
}

func TestReference_GetRecord(t *testing.T) {
	domain := "134"
	record := "156"
	ref, _ := NewReference(domain, record, GlobalScope)

	refRecord := ref.GetRecord()

	assert.Equal(t, record, refRecord)
}

func TestReference_GetDomain(t *testing.T) {
	domain := "134"
	record := "156"
	ref, _ := NewReference(domain, record, GlobalScope)

	refDomain := ref.GetDomain()

	assert.Equal(t, domain, refDomain)
}

func TestReference_GetScope(t *testing.T) {
	domain := "134"
	record := "156"
	ref, _ := NewReference(domain, record, GlobalScope)

	refScope := ref.GetScope()

	assert.Equal(t, GlobalScope, refScope)
}
