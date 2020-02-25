// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
