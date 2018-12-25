package storage

import (
	"crypto"

	"github.com/insolar/insolar/core"
)

type Node struct {
	FID   core.RecordRef
	FRole core.StaticRole
}

func (Node) GetGlobuleID() core.GlobuleID {
	panic("implement me")
}

func (n Node) ID() core.RecordRef {
	return n.FID
}

func (Node) PhysicalAddress() string {
	panic("implement me")
}

func (Node) PublicKey() crypto.PublicKey {
	panic("implement me")
}

func (n Node) Role() core.StaticRole {
	return n.FRole
}

func (Node) ShortID() core.ShortNodeID {
	panic("implement me")
}

func (Node) Version() string {
	panic("implement me")
}
