/*
 *    Copyright 2019 Insolar Technologies
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

package sequence

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGenerator(t *testing.T) {
	gen := NewGenerator()

	require.Equal(t, gen, &generator{sequence: new(uint64)})
}

func TestGenerator_Generate(t *testing.T) {
	gen := NewGenerator()

	seq1 := gen.Generate()
	assert.Equal(t, seq1, Sequence(1))

	seq2 := gen.Generate()
	assert.Equal(t, seq2, Sequence(2))

	seq3 := gen.Generate()
	assert.Equal(t, seq3, Sequence(3))
}
