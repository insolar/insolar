// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

package storage

import (
	"context"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStorage()
	startPulse := *insolar.GenesisPulse

	ks := platformpolicy.NewKeyProcessor()
	p1, err := ks.GeneratePrivateKey()
	assert.NoError(t, err)
	n := node.NewNode(gen.Reference(), insolar.StaticRoleVirtual, ks.ExtractPublicKey(p1), "127.0.0.1:22", "ver2")
	nodes := []insolar.NetworkNode{n}

	for i := 0; i < entriesCount+2; i++ {
		p := startPulse
		p.PulseNumber += insolar.PulseNumber(i)

		snap := node.NewSnapshot(p.PulseNumber, nodes)
		err = s.Append(p.PulseNumber, snap)
		assert.NoError(t, err)

		err = s.AppendPulse(ctx, p)
		assert.NoError(t, err)

		p1, err := s.GetLatestPulse(ctx)
		assert.NoError(t, err)
		assert.Equal(t, p, p1)

		snap1, err := s.ForPulseNumber(p1.PulseNumber)
		assert.NoError(t, err)
		assert.True(t, snap1.Equal(snap), "snapshots should be equal")
	}

	// first pulse and snapshot should be truncated
	assert.Len(t, s.entries, entriesCount)
	assert.Len(t, s.snapshotEntries, entriesCount)

	p2, err := s.GetPulse(ctx, startPulse.PulseNumber)
	assert.EqualError(t, err, ErrNotFound.Error())
	assert.Equal(t, p2, *insolar.GenesisPulse)

	snap2, err := s.ForPulseNumber(startPulse.PulseNumber)
	assert.EqualError(t, err, ErrNotFound.Error())
	assert.Nil(t, snap2)

}
