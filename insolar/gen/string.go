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
	fuzz "github.com/google/gofuzz"
)

// StringFromBytes generates random id with length from 0 to maxcount randomly filled by provided symbols.
func StringFromBytes(symbols []byte, maxcount int) string {
	if maxcount == 0 {
		return ""
	}
	f := fuzz.New().Funcs(func(b *[]byte, c fuzz.Continue) {
		count := c.Intn(maxcount)
		for i := 0; i < count; i++ {
			*b = append(*b, symbols[c.Intn(len(symbols))])
		}
	})
	var bstr []byte
	f.NilChance(0).Fuzz(&bstr)
	return string(bstr)
}
