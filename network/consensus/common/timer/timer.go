///
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
///

package timer

import "time"

type Occasion interface {
	Deadline() time.Time
	NewTimer() Holder
	NewFunc(fn func()) Holder
	IsExpired() bool
}

type Holder interface {
	Channel() <-chan time.Time
	Stop()
}

func Never() Holder {
	return (*timerWithChan)(nil)
}

func New(d time.Duration) Holder {
	return &timerWithChan{time.NewTimer(d)}
}

func NewWithFunc(d time.Duration, fn func()) Holder {
	return &timerWithFn{time.AfterFunc(d, fn)}
}

type timerWithChan struct {
	t *time.Timer
}

func (p *timerWithChan) Channel() <-chan time.Time {
	if p == nil || p.t == nil {
		return nil
	}
	return p.t.C
}

func (p *timerWithChan) Stop() {
	if p == nil || p.t == nil {
		return
	}
	p.t.Stop()
}

type timerWithFn struct {
	t *time.Timer
}

func (p *timerWithFn) Channel() <-chan time.Time {
	panic("illegal state")
}

func (p *timerWithFn) Stop() {
	p.t.Stop()
}

func NewOccasion(deadline time.Time) Occasion {
	return &factory{deadline}
}

func NewOccasionAfter(d time.Duration) Occasion {
	return &factory{time.Now().Add(d)}
}

type factory struct {
	d time.Time
}

func (p *factory) IsExpired() bool {
	return p.d.Before(time.Now())
}

func (p *factory) Deadline() time.Time {
	return p.d
}

func (p *factory) NewTimer() Holder {
	return New(time.Until(p.d))
}

func (p *factory) NewFunc(fn func()) Holder {
	return NewWithFunc(time.Until(p.d), fn)
}

func NeverOccasion() Occasion {
	return &factoryNever{}
}

type factoryNever struct {
}

func (*factoryNever) IsExpired() bool {
	return false
}

func (*factoryNever) Deadline() time.Time {
	return time.Time{}
}

func (*factoryNever) NewTimer() Holder {
	return Never()
}

func (*factoryNever) NewFunc(fn func()) Holder {
	return Never()
}

func EverOccasion() Occasion {
	return &factoryEver{}
}

type factoryEver struct {
}

func (*factoryEver) IsExpired() bool {
	return true
}

func (*factoryEver) Deadline() time.Time {
	return time.Time{}
}

func (*factoryEver) NewTimer() Holder {
	return New(0)
}

func (*factoryEver) NewFunc(fn func()) Holder {
	go fn()
	return Never()
}
