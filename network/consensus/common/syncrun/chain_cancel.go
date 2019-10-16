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

package syncrun

import (
	"context"
	"sync/atomic"
)

func NewChainedCancel() *ChainedCancel {
	return &ChainedCancel{}
}

type ChainedCancel struct {
	state uint32 // atomic
	chain atomic.Value
}

const (
	stateCancelled = 1 << iota
	stateChainHandlerBeingSet
	stateChainHandlerSet
)

func (p *ChainedCancel) Cancel() {
	if p == nil {
		return
	}
	for {
		lastState := atomic.LoadUint32(&p.state)
		switch {
		case lastState&stateCancelled != 0:
			return
		case !atomic.CompareAndSwapUint32(&p.state, lastState, lastState|stateCancelled):
			continue
		case lastState == stateChainHandlerSet:
			p.runChain()
		}
		return
	}
}

func (p *ChainedCancel) runChain() {
	// here is a potential problem, because Go spec doesn't provide ANY ordering on atomic operations
	// but Go compiler does provide some guarantees, so lets hope for the best

	fn := (p.chain.Load()).(context.CancelFunc)
	if fn == nil {
		// this can only happen when atomic ordering is broken
		panic("unexpected atomic ordering")
	}

	// prevent repeated calls as well as retention of references & possible memory leaks
	p.chain.Store(context.CancelFunc(func() {}))
	fn()
}

func (p *ChainedCancel) IsCancelled() bool {
	return p != nil && atomic.LoadUint32(&p.state)&stateCancelled != 0
}

/*
	SetChain sets a chained function once.
	The chained function can only be set once to a non-null value, further calls will panic.
    But if the chained function was not set, the SetChain(nil) can be called multiple times.

	The chained function is guaranteed to be called only once. And it will be called is IsCancelled is already true.
*/
func (p *ChainedCancel) SetChain(chain context.CancelFunc) {
	if chain == nil {
		if p.chain.Load() == nil {
			return
		}
		panic("illegal state")
	}
	for {
		lastState := atomic.LoadUint32(&p.state)
		switch {
		case lastState&^stateCancelled != 0: // chain is set or being set
			panic("illegal state")
			return
		case !atomic.CompareAndSwapUint32(&p.state, lastState, lastState|stateChainHandlerBeingSet): //
			continue
		}
		break
	}

	p.chain.Store(chain)

	for {
		lastState := atomic.LoadUint32(&p.state)
		switch {
		case lastState&^stateCancelled != stateChainHandlerBeingSet:
			// this can only happen when atomic ordering is broken
			panic("unexpected atomic ordering")
		case !atomic.CompareAndSwapUint32(&p.state, lastState, (lastState&stateCancelled)|stateChainHandlerSet):
			continue
		case lastState&stateCancelled != 0:
			// if cancel was set then call the chained cancel here
			// otherwise, the cancelling process will be responsible to call the chained cancel
			p.runChain()
		}
		return
	}
}
