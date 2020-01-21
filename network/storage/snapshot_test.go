// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

package storage

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/assert"
)

func TestNewMemorySnapshotStorage(t *testing.T) {
	ss := NewMemoryStorage()

	ks := platformpolicy.NewKeyProcessor()
	p1, err := ks.GeneratePrivateKey()
	n := node.NewNode(gen.Reference(), insolar.StaticRoleVirtual, ks.ExtractPublicKey(p1), "127.0.0.1:22", "ver2")

	pulse := insolar.Pulse{PulseNumber: 15}
	snap := node.NewSnapshot(pulse.PulseNumber, []insolar.NetworkNode{n})

	err = ss.Append(pulse.PulseNumber, snap)
	assert.NoError(t, err)

	snapshot2, err := ss.ForPulseNumber(pulse.PulseNumber)
	assert.NoError(t, err)

	assert.True(t, snap.Equal(snapshot2))
}
