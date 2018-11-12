package index

import (
	"sync"
)

// RecentObjectsIndexProvider is a base interface for recent objects provider
type RecentObjectsIndexProvider interface {
	AddToFetchedObjects(indexId []byte)
	AddToUpdatedObjects(indexId []byte)

	GetCacheCopy() *RecentObjects

	IncreaseTtl()
	RemoveWithTtlMoreThen(minTtl int)
}

// RecentObjectsIndexInMemoryProvider is an in-memory impl
type RecentObjectsIndexInMemoryProvider struct {
	cache *RecentObjects

	fetchedLock sync.RWMutex
	updatedLock sync.RWMutex
}

func NewRecentObjectsIndexInMemoryProvider() *RecentObjectsIndexInMemoryProvider {
	return &RecentObjectsIndexInMemoryProvider{
		cache: &RecentObjects{
			Updated: map[string]*RecentObjectsMeta{},
			Fetched: map[string]*RecentObjectsMeta{},
		},
	}
}

// AddToFetchedObjects adds index to fetched collection in a thread-safe way
func (p *RecentObjectsIndexInMemoryProvider) AddToFetchedObjects(indexId []byte) {
	p.fetchedLock.Lock()
	defer p.fetchedLock.Unlock()

	key := string(indexId)
	_, ok := p.cache.Fetched[key]
	if !ok {
		p.cache.Fetched[key] = &RecentObjectsMeta{TTL: 0}
	}
}

// AddToUpdatedObjects adds index to updated collection in a thread-safe way
func (p *RecentObjectsIndexInMemoryProvider) AddToUpdatedObjects(indexId []byte) {
	p.updatedLock.Lock()
	defer p.updatedLock.Unlock()

	key := string(indexId)
	_, ok := p.cache.Updated[key]
	if !ok {
		p.cache.Updated[key] = &RecentObjectsMeta{TTL: 0}
	}
}

// GetCacheCopy returns a copy of the cache
func (p *RecentObjectsIndexInMemoryProvider) GetCacheCopy() *RecentObjects {
	result := &RecentObjects{
		Fetched: map[string]*RecentObjectsMeta{},
		Updated: map[string]*RecentObjectsMeta{},
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		p.fetchedLock.RLock()
		defer func() {
			p.fetchedLock.RUnlock()
			wg.Done()
		}()

		for fetchedKey, fetchedMeta := range p.cache.Fetched {
			result.Fetched[fetchedKey] = fetchedMeta
		}
	}()
	go func() {
		p.updatedLock.RLock()
		defer func() {
			p.updatedLock.RUnlock()
			wg.Done()
		}()

		for updatedKey, updatedMeta := range p.cache.Updated {
			result.Fetched[updatedKey] = updatedMeta
		}
	}()
	wg.Wait()

	return result
}

// IncreaseTtl increases all ttl by 1
func (p *RecentObjectsIndexInMemoryProvider) IncreaseTtl() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		p.fetchedLock.Lock()
		defer func() {
			p.fetchedLock.Unlock()
			wg.Done()
		}()

		for _, fetchedMeta := range p.cache.Fetched {
			fetchedMeta.TTL++
		}
	}()
	go func() {
		p.updatedLock.Lock()
		defer func() {
			p.updatedLock.Unlock()
			wg.Done()
		}()

		for _, updatedMeta := range p.cache.Updated {
			updatedMeta.TTL++
		}
	}()
	wg.Wait()
}

// RemoveWithTtlMoreThen remove all objects with Ttl more then provided number
func (p *RecentObjectsIndexInMemoryProvider) RemoveWithTtlMoreThen(minTtl int) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		p.fetchedLock.Lock()
		defer func() {
			p.fetchedLock.Unlock()
			wg.Done()
		}()

		for key, fetchedMeta := range p.cache.Fetched {
			if fetchedMeta.TTL > minTtl{
				delete(p.cache.Fetched, key)
			}
		}
	}()
	go func() {
		p.updatedLock.Lock()
		defer func() {
			p.updatedLock.Unlock()
			wg.Done()
		}()

		for key, updatedMeta := range p.cache.Updated {
			if updatedMeta.TTL > minTtl{
				delete(p.cache.Fetched, key)
			}
		}
	}()
	wg.Wait()
}
