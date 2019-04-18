/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package object

import (
	"sync"

	"github.com/insolar/insolar/insolar"
)

// PulseIndexModifier provides methods for adding data to an index
type PulseIndexModifier interface {
	// Add adds an id to a pulse slot of index
	// If id is added twice, only last pulse will be saved
	Add(id insolar.ID, pn insolar.PulseNumber)
}

// PulseIndexCleaner is a set of methods for cleaning an index
type PulseIndexCleaner interface {
	// DeleteForPulse removes all ids from a pulse slot for a provided pulse
	DeleteForPulse(pn insolar.PulseNumber)
}

// PulseIndexAccessor returns methods for fetching data from pulse slots
type PulseIndexAccessor interface {
	// ForPN returns map of ids from a specified pulse slot
	ForPN(pn insolar.PulseNumber) map[insolar.ID]struct{}
	// LastUsage returns a pulse slot for specific ID
	// Second argument is a status of an operation
	LastUsage(id insolar.ID) (insolar.PulseNumber, bool)
}

// PulseIndex is a union of PulseIndexModifier, PulseIndexCleaner and PulseIndexAccessor
type PulseIndex interface {
	PulseIndexModifier
	PulseIndexCleaner
	PulseIndexAccessor
}

type pulseIndex struct {
	lock        sync.RWMutex
	idsByPulse  map[insolar.PulseNumber]map[insolar.ID]struct{}
	lastUsagePn map[insolar.ID]insolar.PulseNumber
}

// NewPulseIndex creates a new instance of PulseIndex
func NewPulseIndex() PulseIndex {
	return &pulseIndex{
		idsByPulse:  map[insolar.PulseNumber]map[insolar.ID]struct{}{},
		lastUsagePn: map[insolar.ID]insolar.PulseNumber{},
	}
}

// Add adds an id to a pulse slot of index
// If id is added twice, only last pulse will be saved
func (p *pulseIndex) Add(id insolar.ID, pn insolar.PulseNumber) {
	p.lock.Lock()
	defer p.lock.Unlock()

	ids, ok := p.idsByPulse[pn]
	if !ok {
		ids = map[insolar.ID]struct{}{}
		p.idsByPulse[pn] = ids
	}

	lstPN, ok := p.lastUsagePn[id]
	if ok {
		delete(p.idsByPulse[lstPN], id)
	}

	ids[id] = struct{}{}
	p.lastUsagePn[id] = pn
}

// DeleteForPulse removes all ids from a pulse slot for a provided pulse
func (p *pulseIndex) DeleteForPulse(pn insolar.PulseNumber) {
	p.lock.Lock()
	defer p.lock.Unlock()

	_, ok := p.idsByPulse[pn]
	if !ok {
		return
	}

	delete(p.idsByPulse, pn)
	for id, lpn := range p.lastUsagePn {
		if lpn == pn {
			delete(p.lastUsagePn, id)
		}
	}
}

// ForPN returns map of ids from a specified pulse slot
func (p *pulseIndex) ForPN(pn insolar.PulseNumber) map[insolar.ID]struct{} {
	p.lock.RLock()
	defer p.lock.RUnlock()

	ids, ok := p.idsByPulse[pn]
	if !ok {
		return nil
	}

	res := map[insolar.ID]struct{}{}
	for id := range ids {
		res[id] = struct{}{}
	}

	return res
}

// LastUsage returns a pulse slot for specific ID
// Second argument is a status of an operation
func (p *pulseIndex) LastUsage(id insolar.ID) (insolar.PulseNumber, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	res, ok := p.lastUsagePn[id]
	return res, ok
}
