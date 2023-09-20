package rwlock

import "sync"

type RWLocker interface {
	sync.Locker
	RLock()
	RUnlock()
}

func DummyLocker() RWLocker {
	return &dummyLock
}

var dummyLock = dummyLocker{}

type dummyLocker struct {
}

func (*dummyLocker) Lock() {
}

func (*dummyLocker) Unlock() {
}

func (*dummyLocker) RUnlock() {
}

func (*dummyLocker) RLock() {
}

func (*dummyLocker) String() string {
	return "dummyLocker"
}
