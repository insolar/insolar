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
	svc  PulseDataService
	exec smachine.ExecutionAdapter

	// mutable
	presentAndFuturePulse uint64 //atomic
	preparingPulse        uint32 //atomic
}

type PulseDataService interface {
	LoadPulseData(pulse.Number) (pulse.Data, bool)
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
	futurePN := pd.GetNextPulseNumber()

	for {
		prev := atomic.LoadUint64(&p.presentAndFuturePulse)
		if prev != 0 {
			expectedPN := pulse.Number(prev >> 32)
			if pd.PulseNumber < expectedPN {
				panic(fmt.Errorf("illegal pulse data: pn=%v, expected=%v", presentPN, expectedPN))
			}
		}
		if atomic.CompareAndSwapUint64(&p.presentAndFuturePulse, prev, uint64(presentPN)|uint64(futurePN)<<32) {
			return
		}
	}
}

func (p *PulseDataManager) isPreparingPulse() bool {
	return atomic.LoadUint32(&p.preparingPulse) != 0
}

func (p *PulseDataManager) setPreparingPulse(out PreparePulseChangeChannel) {
	atomic.StoreUint32(&p.preparingPulse, 1)
}

func (p *PulseDataManager) unsetPreparingPulse() {
	atomic.StoreUint32(&p.preparingPulse, 0)
}

func (p *PulseDataManager) GetPulseData(pn pulse.Number) (pulse.Data, bool) {
	panic("unimplemented")
}

// for non-recent past HasPulseData() can be incorrect / incomplete
func (p *PulseDataManager) HasPulseData(pn pulse.Number) bool {
	return true // TODO HasPulseData
}

func (p *PulseDataManager) IsAllowedFutureSpan(expectedPN pulse.Number, futurePN pulse.Number) bool {
	// TODO limit how much we can handle as future
	return futurePN >= expectedPN
}

func (p *PulseDataManager) IsAllowedPastSpan(presentPN pulse.Number, pastPN pulse.Number) bool {
	// TODO limit how much we can handle as future
	return pastPN < presentPN
}

func (p *PulseDataManager) IsRecentPastRange(pastPN pulse.Number) bool {
	// TODO limit how much we can handle as future
	return pastPN.IsTimePulse()
}

func (p *PulseDataManager) prepareAsync(ctx smachine.ExecutionContext, fn func(svc PulseDataService) smachine.AsyncResultFunc) smachine.AsyncCallRequester {
	return p.exec.PrepareAsync(ctx, func() smachine.AsyncResultFunc {
		fn(p.svc)
		return nil
	})
}

func (p *PulseDataManager) RequestPulseData(ctx smachine.ExecutionContext,
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
	// TODO implement cache
}
