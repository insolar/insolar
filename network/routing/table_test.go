// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package routing

import (
	"context"
	"strconv"
	"testing"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/network"
	mock "github.com/insolar/insolar/testutils/network"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newNode(ref insolar.Reference, id int) insolar.NetworkNode {
	address := "127.0.0.1:" + strconv.Itoa(id)
	result := node.NewNode(ref, insolar.StaticRoleUnknown, nil, address, "")
	result.(node.MutableNode).SetShortID(insolar.ShortNodeID(id))
	return result
}

func TestTable_Resolve(t *testing.T) {
	table := Table{}

	refs := gen.UniqueReferences(2)
	pulse := insolar.GenesisPulse
	nodeKeeperMock := mock.NewNodeKeeperMock(t)
	nodeKeeperMock.GetAccessorMock.Set(func(p1 insolar.PulseNumber) network.Accessor {
		n := newNode(refs[0], 123)
		return node.NewAccessor(node.NewSnapshot(pulse.PulseNumber, []insolar.NetworkNode{n}))
	})

	pulseAccessorMock := mock.NewPulseAccessorMock(t)
	pulseAccessorMock.GetLatestPulseMock.Set(func(ctx context.Context) (p1 insolar.Pulse, err error) {
		return *pulse, nil
	})

	table.PulseAccessor = pulseAccessorMock
	table.NodeKeeper = nodeKeeperMock

	h, err := table.Resolve(refs[0])
	require.NoError(t, err)
	assert.EqualValues(t, 123, h.ShortID)
	assert.Equal(t, "127.0.0.1:123", h.Address.String())

	_, err = table.Resolve(refs[1])
	assert.Error(t, err)
}
