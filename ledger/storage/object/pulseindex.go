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

type PulseIndexModifier interface {
	Add(id insolar.ID, pn insolar.PulseNumber)
}

type PulseIndexCleaner interface {
	DeleteForPulse(pn insolar.PulseNumber)
}

type PulseIndexAccessor interface {
	ForPN(pn insolar.PulseNumber) map[insolar.ID]struct{}
}

type PulseIndex interface {
	PulseIndexModifier
	PulseIndexCleaner
	PulseIndexAccessor
}

type pulseIndex struct {
	lock             sync.RWMutex
	indexByPulseStor map[insolar.PulseNumber]map[insolar.ID]struct{}
}

func NewPulseIndex() PulseIndex {
	return &pulseIndex{indexByPulseStor: map[insolar.PulseNumber]map[insolar.ID]struct{}{}}
}

func (p *pulseIndex) Add(id insolar.ID, pn insolar.PulseNumber) {
	p.lock.Lock()
	defer p.lock.Unlock()

	ids, ok := p.indexByPulseStor[pn]
	if !ok {
		ids = map[insolar.ID]struct{}{}
		p.indexByPulseStor[pn] = ids
	}
	ids[id] = struct{}{}
}

func (p *pulseIndex) DeleteForPulse(pn insolar.PulseNumber) {
	p.lock.Lock()
	defer p.lock.Unlock()

	_, ok := p.indexByPulseStor[pn]
	if !ok {
		return
	}

	delete(p.indexByPulseStor, pn)
}

func (p *pulseIndex) ForPN(pn insolar.PulseNumber) map[insolar.ID]struct{} {
	p.lock.RLock()
	defer p.lock.RUnlock()

	ids, ok := p.indexByPulseStor[pn]
	if !ok {
		return nil
	}

	res := map[insolar.ID]struct{}{}
	for id := range ids {
		res[id] = struct{}{}
	}

	return res
}
