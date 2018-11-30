package recentstorage

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
)

// RecentStorage is a base interface for the storage of recent objects and indexes
//go:generate minimock -i github.com/insolar/insolar/ledger/recentstorage.RecentStorage -o ./ -s _mock.go
type RecentStorage interface {
	AddObject(id core.RecordID)
	AddObjectWithMeta(id core.RecordID, meta *message.RecentObjectsIndexMeta)

	AddPendingRequest(id core.RecordID)

	RemovePendingRequest(id core.RecordID)
	GetObjects() map[core.RecordID]*message.RecentObjectsIndexMeta
	GetRequests() []core.RecordID
	ClearZeroTTLObjects()
	ClearObjects()
}
