package recentstorage

import (
	"github.com/insolar/insolar/core"
)

// Provider provides a recent storage for jet
//go:generate minimock -i github.com/insolar/insolar/ledger/recentstorage.Provider -o ./ -s _mock.go
type Provider interface {
	GetStorage(jetID core.RecordID) RecentStorage
}

// RecentStorage is a base interface for the storage of recent objects and indexes
//go:generate minimock -i github.com/insolar/insolar/ledger/recentstorage.RecentStorage -o ./ -s _mock.go
type RecentStorage interface {
	AddObject(id core.RecordID, isMine bool)
	AddObjectWithTLL(id core.RecordID, ttl int, isMine bool)

	AddPendingRequest(id core.RecordID)
	RemovePendingRequest(id core.RecordID)

	IsMine(id core.RecordID) bool

	GetObjects() map[core.RecordID]int
	GetRequests() []core.RecordID

	ClearZeroTTLObjects()
	ClearObjects()
}
