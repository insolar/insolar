package recentstorage

import (
	"github.com/insolar/insolar/core"
)

// RecentStorage is a base interface for the storage of recent objects and indexes
//go:generate minimock -i github.com/insolar/insolar/ledger/recentstorage.RecentStorage -o ./ -s _mock.go
type RecentStorage interface {
	AddObject(id core.RecordID)
	AddObjectWithTLL(id core.RecordID, ttl int)

	AddPendingRequest(id core.RecordID)
	RemovePendingRequest(id core.RecordID)

	IsMine(id core.RecordID) bool
	MarkAsMine(id core.RecordID) error

	GetObjects() map[core.RecordID]int
	GetRequests() []core.RecordID

	ClearZeroTTLObjects()
	ClearObjects()
}
