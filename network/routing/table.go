// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package routing

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/storage"

	"github.com/pkg/errors"
)

type Table struct {
	NodeKeeper    network.NodeKeeper    `inject:""`
	PulseAccessor storage.PulseAccessor `inject:""`
}

func (t *Table) isLocalNode(insolar.Reference) bool {
	return true
}

func (t *Table) resolveRemoteNode(insolar.Reference) (*host.Host, error) {
	return nil, errors.New("not implemented")
}

// Resolve NodeID -> ShortID, Address. Can initiate network requests.
func (t *Table) Resolve(ref insolar.Reference) (*host.Host, error) {
	if t.isLocalNode(ref) {
		p, err := t.PulseAccessor.GetLatestPulse(context.Background())
		if err != nil {
			return nil, errors.Wrap(err, "failed to get latest pulse --==-- ")
		}

		node := t.NodeKeeper.GetAccessor(p.PulseNumber).GetActiveNode(ref)
		if node == nil {
			return nil, errors.New("no such local node with NodeID: " + ref.String())
		}
		return host.NewHostNS(node.Address(), node.ID(), node.ShortID())
	}
	return t.resolveRemoteNode(ref)
}
