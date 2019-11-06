package injector

import (
	"sync"
)

func NewRegistry(m map[string]interface{}) ReadOnlyContainer {
	return NewRegistryWithParent(nil, m)
}

func NewRegistryWithParent(parent DependencyRegistryFunc, m map[string]interface{}) ReadOnlyContainer {
	n := len(m)
	if n == 0 {
		return ReadOnlyContainer{parent: parent}
	}

	c := make(map[string]interface{}, n)
	for k, v := range m {
		if k == "" {
			panic("illegal value")
		}
		c[k] = v
	}
	return ReadOnlyContainer{parent, c}
}

func NewContainer(parent DependencyRegistryFunc) DependencyContainer {
	return ReadWriteContainer{ReadOnlyContainer{parent, make(map[string]interface{})}}
}

func NewSyncContainer(parent DependencyRegistryFunc) DependencyContainer {
	return SyncReadWriteContainer{parent, &sync.Map{}}
}

/* ------------------------------------ */

type ReadOnlyContainer struct {
	parent DependencyRegistryFunc
	m      map[string]interface{}
}

func (r ReadOnlyContainer) IsEmpty() bool {
	return len(r.m) == 0
}

func (r ReadOnlyContainer) Count() int {
	return len(r.m)
}

func (r ReadOnlyContainer) parentDependency(id string) (interface{}, bool) {
	if r.parent == nil {
		return nil, false
	}
	return r.parent(id)
}

func (r ReadOnlyContainer) FindDependency(id string) (interface{}, bool) {
	if v, ok := r.FindLocalDependency(id); ok {
		return v, ok
	}
	return r.parentDependency(id)
}

func (r ReadOnlyContainer) FindLocalDependency(id string) (interface{}, bool) {
	v, ok := r.m[id]
	return v, ok
}

func (r ReadOnlyContainer) FilterLocalDependencies(fn func(string, interface{}) bool) bool {
	for id, v := range r.m {
		if fn(id, v) {
			return true
		}
	}
	return false
}

/* ------------------------------------ */

type roc = ReadOnlyContainer
type ReadWriteContainer struct {
	roc
}

func (r ReadWriteContainer) FindLocalDependency(id string) (interface{}, bool) {
	v, ok := r.m[id]
	return v, ok
}

func (r ReadWriteContainer) PutDependency(id string, v interface{}) {
	if id == "" {
		panic("illegal value")
	}
	r.m[id] = v
}

func (r ReadWriteContainer) TryPutDependency(id string, v interface{}) bool {
	if id == "" {
		panic("illegal value")
	}
	if _, ok := r.m[id]; ok {
		return false
	}
	r.m[id] = v
	return true
}

/* ------------------------------------ */

type SyncReadWriteContainer struct {
	parent DependencyRegistryFunc
	m      *sync.Map
}

func (r SyncReadWriteContainer) parentDependency(id string) (interface{}, bool) {
	if r.parent == nil {
		return nil, false
	}
	return r.parent(id)
}

func (r SyncReadWriteContainer) FindDependency(id string) (interface{}, bool) {
	if v, ok := r.FindLocalDependency(id); ok {
		return v, ok
	}
	return r.parentDependency(id)
}

func (r SyncReadWriteContainer) FindLocalDependency(id string) (interface{}, bool) {
	return r.m.Load(id)
}

func (r SyncReadWriteContainer) PutDependency(id string, v interface{}) {
	if id == "" {
		panic("illegal key")
	}
	r.m.Store(id, v)
}

func (r SyncReadWriteContainer) TryPutDependency(id string, v interface{}) bool {
	if id == "" {
		panic("illegal key")
	}
	_, loaded := r.m.LoadOrStore(id, v)
	return !loaded
}
