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

package jet

import (
	"testing"

	"github.com/insolar/insolar/insolar/bits"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
)

func TestJet_Parent(t *testing.T) {
	var (
		parent = NewIDFromString("01010")
		child  = NewIDFromString("010101")
	)

	gotParent := Parent(child)
	require.Equal(t, parent, gotParent, "got proper parent")

	emptyChild := *insolar.NewJetID(0, nil)
	emptyParent := Parent(emptyChild)
	require.Equal(t, emptyChild, emptyParent, "for empty jet ID, got the same parent")
}

func TestJet_ParsePrefixAndResetBits(t *testing.T) {
	orig := []byte{0xFF}
	got := bits.ResetBits(orig, 5)
	require.Equal(t, parsePrefix("11111000"), got,
		"bit reset sucessfully %b == %b", parsePrefix("11111000"), got)
	require.NotEqual(t, &orig, &got, "without overflow returns a new slice")
}

func TestJet_SiblingParent(t *testing.T) {
	jetID := gen.JetID()
	left, right := Siblings(jetID)
	require.True(t, jetID.Equal(Parent(left)))
	require.True(t, jetID.Equal(Parent(right)))
}

func TestJet_NewJetIDSiblingParent(t *testing.T) {
	jetID := *insolar.NewJetID(5, gen.ID().Bytes())

	left, right := Siblings(jetID)
	require.True(t, jetID.Equal(Parent(left)))
	require.True(t, jetID.Equal(Parent(right)))
}
