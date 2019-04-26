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
