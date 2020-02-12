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

package bits

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResetBits(t *testing.T) {
	orig := []byte{0xFF}
	got := ResetBits(orig, 5)
	require.NotEqual(t, &orig, &got, "without overflow returns a new slice")

	gotWithOverflow := ResetBits(orig, 9)
	require.Equal(t, []byte{0xFF}, gotWithOverflow, "returns equals slice on overflow")
	require.Equal(t, &orig, &gotWithOverflow, "on overflow returns the same slice")
	require.Equal(t, []byte{0xFF}, orig, "original unchanged after resetBits")
}
