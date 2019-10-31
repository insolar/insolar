//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package conveyor

import (
	"fmt"
	"sync/atomic"

	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/pulse"
)

type PulseDataManager struct {
	// set at construction, immutable
	execPulseDataService smachine.ExecutionAdapter // for PulseDataService

	cache PulseDataCache

	// set at init, immutable
	maxPastPulseAge uint32
	futureCycles    uint8

	// mutable
	presentAndFuturePulse uint64 //atomic
	earliestCacheBound    uint32 // atomic
	preparingPulseFlag    uint32 // atomic
}

type PulseDataService interface {
	LoadPulseData(pulse.Number) (pulse.Data, bool)
}

func (p *PulseDataManager) Init(minCachePulseAge, maxPastPulseAge uint32, maxFutureCycles uint8) {
	if minCachePulseAge == 0 || minCachePulseAge > pulse.MaxTimePulse {
		panic("illegal value")
	}
	if maxPastPulseAge < minCachePulseAge || maxPastPulseAge > pulse.MaxTimePulse {
		panic("illegal value")
	}
	p.maxPastPulseAge = maxPastPulseAge
	p.futureCycles = maxFutureCycles
	p.cache.Init(p, minCachePulseAge, 2) // any pulse data stays in cache for at least 2 pulse cycles
}

const uninitializedFuture = pulse.LocalRelative

func (p *PulseDataManager) GetPresentPulse() (present pulse.Number, nearestFuture pulse.Number) {
	v := atomic.LoadUint64(&p.presentAndFuturePulse)
	if v == 0 {
		return pulse.Unknown, uninitializedFuture
	}
	return p._split(v)
}

func (*PulseDataManager) _split(v uint64) (present pulse.Number, nearestFuture pulse.Number) {
	return pulse.Number(v), pulse.Number(v >> 32)
}

func (p *PulseDataManager) setPresentPulse(pd pulse.Data) {
	presentPN := pd.PulseNumber
	futurePN := pd.NextPulseNumber()

	if epd, ok := p.cache.Check(presentPN); ok {
		if epd != pd {
			panic("illegal state")
		}
	}

	for {
		prev := atomic.LoadUint64(&p.presentAndFuturePulse)
		if prev != 0 {
			expectedPN := pulse.Number(prev >> 32)
			if pd.PulseNumber < expectedPN {
				panic(fmt.Errorf("illegal pulse data: pn=%v, expected=%v", presentPN, expectedPN))
			}
		}
		if atomic.CompareAndSwapUint64(&p.presentAndFuturePulse, prev, uint64(presentPN)|uint64(futurePN)<<32) {
			if prev == 0 {
				atomic.CompareAndSwapUint32(&p.earliestCacheBound, 0, uint32(presentPN))
			}
			break
		}
	}

	p.cache.EvictAndRotate(presentPN)
}

func (p *PulseDataManager) getEarliestCacheBound() pulse.Number {
	return pulse.Number(atomic.LoadUint32(&p.earliestCacheBound))
}

func (p *PulseDataManager) isPreparingPulse() bool {
	return atomic.LoadUint32(&p.preparingPulseFlag) != 0
}

func (p *PulseDataManager) setPreparingPulse(out PreparePulseChangeChannel) {
	atomic.StoreUint32(&p.preparingPulseFlag, 1)
}

func (p *PulseDataManager) unsetPreparingPulse() {
	atomic.StoreUint32(&p.preparingPulseFlag, 0)
}

func (p *PulseDataManager) GetPulseData(pn pulse.Number) (pulse.Data, bool) {
	return p.cache.Get(pn)
}

func (p *PulseDataManager) getCachedPulseSlot(pn pulse.Number) (*PulseSlot, bool) {
	return p.cache.getPulseSlot(pn)
}

// for non-recent past HasPulseData() can be incorrect / incomplete
func (p *PulseDataManager) HasPulseData(pn pulse.Number) bool {
	return p.cache.Contains(pn)
}

func (p *PulseDataManager) TouchPulseData(pn pulse.Number) bool {
	return p.cache.Touch(pn)
}

// Returns true when the given PN can be accepted into Future pulse slot, otherwise must be rejected
func (p *PulseDataManager) IsAllowedFutureSpan(futurePN pulse.Number) bool {
	presentPN, expectedPN := p.GetPresentPulse()
	return p.isAllowedFutureSpan(presentPN, expectedPN, futurePN)
}

func (p *PulseDataManager) isAllowedFutureSpan(presentPN, expectedPN pulse.Number, futurePN pulse.Number) bool {
	if futurePN < expectedPN {
		return false
	}
	return p.futureCycles == 0 || futurePN <= (expectedPN+(expectedPN-presentPN)*pulse.Number(p.futureCycles))
}

func (p *PulseDataManager) IsAllowedPastSpan(pastPN pulse.Number) bool {
	presentPN, _ := p.GetPresentPulse()
	return p.isAllowedPastSpan(presentPN, pastPN)
}

func (p *PulseDataManager) isAllowedPastSpan(presentPN pulse.Number, pastPN pulse.Number) bool {
	return pastPN < presentPN && pastPN+pulse.Number(p.maxPastPulseAge) >= presentPN
}

func (p *PulseDataManager) IsRecentPastRange(pastPN pulse.Number) bool {
	presentPN, _ := p.GetPresentPulse()
	return p.isRecentPastRange(presentPN, pastPN)
}

//Returns true when the given PN is within a mandatory retention interval for the cache. So we don't need to populate it
func (p *PulseDataManager) isRecentPastRange(presentPN pulse.Number, pastPN pulse.Number) bool {
	return pastPN < presentPN &&
		(pastPN+pulse.Number(p.cache.GetMinRange())) >= presentPN &&
		pastPN >= p.getEarliestCacheBound() // this interval can be much narrower for a recently started node
}

func (p *PulseDataManager) prepareAsync(ctx smachine.ExecutionContext, fn func(svc PulseDataService) smachine.AsyncResultFunc) smachine.AsyncCallRequester {
	return p.execPulseDataService.PrepareAsync(ctx, func(svc interface{}) smachine.AsyncResultFunc {
		fn(svc.(PulseDataService))
		return nil
	})
}

func (p *PulseDataManager) PreparePulseDataRequest(ctx smachine.ExecutionContext,
	pn pulse.Number,
	resultFn func(isAvailable bool, pd pulse.Data),
) smachine.AsyncCallRequester {
	if resultFn == nil {
		panic("illegal value")
	}
	if pd, ok := p.GetPulseData(pn); ok {
		resultFn(ok, pd)
	}

	return p.prepareAsync(ctx, func(svc PulseDataService) smachine.AsyncResultFunc {
		pd, ok := svc.LoadPulseData(pn)

		return func(ctx smachine.AsyncResultContext) {
			if ok && pd.IsValidPulsarData() {
				p.putPulseData(pd)
				resultFn(ok, pd)
			} else {
				resultFn(false, pulse.Data{})
			}
		}
	}).WithFlags(smachine.AutoWakeUp)
}

func (p *PulseDataManager) putPulseData(data pulse.Data) {
	p.cache.Put(data)
}
