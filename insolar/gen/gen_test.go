//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGen_Signature(t *testing.T) {
	zero := Signature(0)
	assert.NotNil(t, zero)
	assert.Equal(t, 0, len(zero))

	one := Signature(1)
	assert.NotNil(t, one)
	assert.Equal(t, 1, len(one))

	case256 := Signature(256)
	assert.NotNil(t, case256)
	assert.Equal(t, 256, len(case256))

	negative := Signature(-1)
	assert.Nil(t, negative)
}
