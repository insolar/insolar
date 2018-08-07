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
	"testing"

	"github.com/stretchr/testify/assert"
)

var ref = &reference{
	domain: "domain1",
	record: "record1",
	scope:  GlobalScope,
}

func TestBaseCallable_SetReference(t *testing.T) {
	callable := &BaseCallable{}
	callable.SetReference(ref)

	assert.Equal(t, ref, callable.reference)
}

func TestBaseCallable_GetReference(t *testing.T) {
	callable := &BaseCallable{
		reference: ref,
	}
	assert.Equal(t, ref, callable.GetReference())
}
