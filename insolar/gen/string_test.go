// Copyright 2020 Insolar Network Ltd.
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

package gen

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
)

func TestGen_StringFromBytes(t *testing.T) {
	symbolFuzzer := fuzz.New().NilChance(0).NumElements(1, 100)
	var symbols []byte
	symbolFuzzer.Fuzz(&symbols)
	for i := 0; i < 100; i++ {
		s := StringFromBytes(symbols, i)
		require.GreaterOrEqualf(t, i, len(s), "string length should not be greater than `maxcount`")
		for _, sym := range []byte(s) {
			require.Contains(t, symbols, sym, "byte should be in range")
		}
	}
}
