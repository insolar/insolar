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

package jet

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/core"
)

func TestJet_Parent(t *testing.T) {
	var (
		parent = NewIDFromString("01010")
		child  = NewIDFromString("010101")
	)

	gotParent := Parent(child)
	require.Equal(t, parent, gotParent, "got proper parent")

	emptyChild := *core.NewJetID(0, nil)
	emptyParent := Parent(emptyChild)
	require.Equal(t, emptyChild, emptyParent, "for empty jet ID, got the same parent")
}

func TestJet_ResetBits(t *testing.T) {
	orig := []byte{0xFF}
	got := ResetBits(orig, 5)
	require.Equal(t, parsePrefix("11111000"), got,
		"bit reset sucessfully %b == %b", parsePrefix("11111000"), got)
	require.NotEqual(t, &orig, &got, "without overflow returns a new slice")

	gotWithOverflow := ResetBits(orig, 9)
	require.Equal(t, []byte{0xFF}, gotWithOverflow, "returns equals slice on overflow")
	require.Equal(t, &orig, &gotWithOverflow, "on overflow returns the same slice")
	require.Equal(t, []byte{0xFF}, orig, "original unchanged after ResetBits")
}
