//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package aeshash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	require.Equal(t, uint(0x9898989898989898), Str(""))
	require.Equal(t, uint(0xa0be041b1a8572be), Str("insolar"))
	require.Equal(t, uint(0xa0be041b1a8572be), StrWithSeed("insolar", 0))
	require.Equal(t, uint(0xafa17a6f4a4ed3ef), StrWithSeed("insolar", 1))

	require.Equal(t, uint(0x9898989898989898), Slice(nil))
	require.Equal(t, uint(0x9898989898989898), Slice([]byte{}))
	require.Equal(t, uint(0xa0be041b1a8572be), Slice([]byte("insolar")))
	require.Equal(t, uint(0xa0be041b1a8572be), SliceWithSeed([]byte("insolar"), 0))
	require.Equal(t, uint(0xafa17a6f4a4ed3ef), SliceWithSeed([]byte("insolar"), 1))
}
