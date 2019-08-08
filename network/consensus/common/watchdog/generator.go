//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package watchdog

import (
	"math"
	"sync/atomic"
	"time"
)

type HeartbeatGeneratorFactory interface {
	CreateGenerator(name string) *HeartbeatGenerator
}

func NewHeartbeatGenerator(id HeartbeatID, heartbeatPeriod time.Duration, out chan<- Heartbeat) HeartbeatGenerator {
	return NewHeartbeatGeneratorWithRetries(id, heartbeatPeriod, 0, out)
}

func NewHeartbeatGeneratorWithRetries(id HeartbeatID, heartbeatPeriod time.Duration, retryCount uint8, out chan<- Heartbeat) HeartbeatGenerator {
	attempts := retryCount
	if out == nil {
		attempts = 0
	} else if attempts < math.MaxUint8 {
		attempts++
	}

	period := uint32(0)
	switch {
	case heartbeatPeriod < 0:
		panic("illegal value")
	case heartbeatPeriod == 0:
		break
	case heartbeatPeriod <= time.Millisecond:
		period = 1
	default:
		heartbeatPeriod /= time.Millisecond
		if heartbeatPeriod > math.MaxUint32 {
			period = math.MaxUint32
		} else {
			period = uint32(heartbeatPeriod)
		}
	}

	return HeartbeatGenerator{id: id, heartbeatPeriod: period, sendAttempts: attempts, out: out}
}

type HeartbeatGenerator struct {
	id              HeartbeatID
	heartbeatPeriod uint32
	atomicNano      int64
	sendAttempts    uint8
	//name string
	out chan<- Heartbeat
}

func (g *HeartbeatGenerator) Heartbeat() {
	g.ForcedHeartbeat(false)
}

func (g *HeartbeatGenerator) ForcedHeartbeat(forced bool) {

	lastNano := atomic.LoadInt64(&g.atomicNano)
	if lastNano == DisabledHeartbeat {
		//closed channel or generator
		return
	}

	currentNano := time.Now().UnixNano()
	if lastNano != 0 && !forced && currentNano-lastNano < int64(g.heartbeatPeriod)*int64(time.Millisecond) {
		return
	}

	if !atomic.CompareAndSwapInt64(&g.atomicNano, lastNano, currentNano) {
		return // there is no need to retry in case of contention
	}

	if g.send(Heartbeat{From: g.id, PreviousUnixTime: lastNano, UpdateUnixTime: currentNano}) {
		return
	}

	atomic.CompareAndSwapInt64(&g.atomicNano, currentNano, lastNano) // try to roll-back on failed send
}

func (g *HeartbeatGenerator) Cancel() {
	for {
		lastNano := atomic.LoadInt64(&g.atomicNano)
		if lastNano == DisabledHeartbeat {
			//closed channel or generator
			return
		}
		if atomic.CompareAndSwapInt64(&g.atomicNano, lastNano, DisabledHeartbeat) {
			g.send(Heartbeat{From: g.id, PreviousUnixTime: lastNano, UpdateUnixTime: DisabledHeartbeat})
			return // there is no need to retry in case of contention
		}
	}
}

func (g *HeartbeatGenerator) send(beat Heartbeat) bool {
	defer func() {
		err := recover() // just in case of the closed channel
		if err != nil {
			g.Disable()
		}
	}()

	if g.sendAttempts == 0 {
		return true
	}

	for i := g.sendAttempts; i > 0; i-- {
		select {
		case g.out <- beat:
			return true
		default:
			// avoid lock up
		}
	}
	return false
}

func (g *HeartbeatGenerator) Disable() {
	atomic.StoreInt64(&g.atomicNano, DisabledHeartbeat)
}

func (g *HeartbeatGenerator) IsEnabled() bool {
	return atomic.LoadInt64(&g.atomicNano) != DisabledHeartbeat
}
