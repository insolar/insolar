package pool

import (
	"fmt"
	"sync"
)

type iterateFunc func(entry *entry)

type entryHolder struct {
	sync.RWMutex
	entries map[string]*entry
}

func newEntryHolder() *entryHolder {
	return &entryHolder{
		entries: make(map[string]*entry),
	}
}

func (eh *entryHolder) key(host fmt.Stringer) string {
	return host.String()
}

func (eh *entryHolder) get(host fmt.Stringer) (*entry, bool) {
	eh.RLock()
	defer eh.RUnlock()
	e, ok := eh.entries[eh.key(host)]
	return e, ok
}

func (eh *entryHolder) delete(host fmt.Stringer) bool {
	eh.Lock()
	defer eh.Unlock()

	e, ok := eh.entries[eh.key(host)]
	if ok {
		e.close()
		delete(eh.entries, eh.key(host))
		return true
	}
	return false
}

func (eh *entryHolder) add(host fmt.Stringer, entry *entry) {
	eh.Lock()
	defer eh.Unlock()
	eh.entries[eh.key(host)] = entry
}

func (eh *entryHolder) clear() {
	eh.Lock()
	defer eh.Unlock()
	for key := range eh.entries {
		delete(eh.entries, key)
	}
}

func (eh *entryHolder) iterate(iterateFunc iterateFunc) {
	eh.Lock()
	defer eh.Unlock()
	for _, h := range eh.entries {
		iterateFunc(h)
	}
}

func (eh *entryHolder) size() int {
	eh.RLock()
	defer eh.RUnlock()
	return len(eh.entries)
}
