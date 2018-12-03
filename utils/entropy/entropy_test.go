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

package entropy

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/platformpolicy"
)

func TestSelectByEntropy(t *testing.T) {
	scheme := platformpolicy.NewPlatformCryptographyScheme()
	entropy := [64]byte{'X', 'Y', 'Z'}
	values := [][]byte{
		{'A'},
		{'B'},
		{'C'},
	}
	count := 1
	result1, err := SelectByEntropy(scheme, entropy[:], values, count)
	require.NoError(t, err)
	fmt.Printf("%#v\n", result1)
	assert.Equal(t, count, len(result1))

	result2, err := SelectByEntropy(scheme, entropy[:], values, count)
	assert.Equal(t, result1, result2)
}
