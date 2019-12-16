///
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
///

package pgbuf

import (
	"math"
	"sync/atomic"
	"unsafe"
)

type pagePrepareFunc func(newPage *BufferPage)

func newPageLifo(prepareFn pagePrepareFunc) pageLifo {
	return pageLifo{prepareFn: prepareFn}
}

type pageLifo struct {
	prepareFn pagePrepareFunc
	next      unsafe.Pointer // atomic, *BufferPage
}

func (p *pageLifo) push(page *BufferPage) *BufferPage {
	for {
		cur := p.peek()
		if p.pushExpected(cur, page) {
			return cur
		}
	}
}

func (p *pageLifo) pushExpected(expected, page *BufferPage) bool {
	if p.prepareFn != nil {
		p.prepareFn(page)
	}
	page.bufferCleanupData.next = expected
	return atomic.CompareAndSwapPointer(&p.next, unsafe.Pointer(expected), unsafe.Pointer(page))
}

func (p *pageLifo) pop() *BufferPage {
	for {
		cur := p.peek()
		if cur == nil {
			return nil
		}
		if p.pushExpected(cur, cur.next) {
			return cur
		}
	}
}

func (p *pageLifo) peek() *BufferPage {
	return (*BufferPage)(atomic.LoadPointer(&p.next))
}

func (p *pageLifo) flush() *BufferPage {
	return p.replace(nil)
}

func (p *pageLifo) replace(page *BufferPage) *BufferPage {
	return (*BufferPage)(atomic.SwapPointer(&p.next, unsafe.Pointer(page)))
}

/* ============================ */

type BufferPage struct {
	active uint32 // atomic
	data   []byte
	bufferCleanupData
}

type bufferCleanupData struct {
	offset uint32 // atomic
	count  uint32

	head *BufferPage

	totalCapacity uint32
	next          *BufferPage // direction to the head, =nul for head

	prev unsafe.Pointer // atomic, *BufferPage, for trimHead==true only
	trim *BufferTrim    // for trimHead==true only and for head only
}

func (b *BufferPage) tryInitAccess() bool {
	return atomic.CompareAndSwapUint32(&b.active, 0, 1)
}

func (b *BufferPage) initAccess() {
	if !b.tryInitAccess() {
		panic("illegal state")
	}
}

func (b *BufferPage) startAccess() bool {
	for {
		c := atomic.LoadUint32(&b.active)
		switch c {
		case 0, math.MaxUint32:
			// this is no more an active page or is exclusively locked
			return false
		case math.MaxUint32 - 1:
			panic("overflow")
		}
		if atomic.CompareAndSwapUint32(&b.active, c, c+1) {
			return true
		}
	}
}

func (b *BufferPage) stopAccess() {
	for {
		c := atomic.LoadUint32(&b.active)
		if c == 0 || c == math.MaxUint32 {
			panic("illegal state")
		}
		if atomic.CompareAndSwapUint32(&b.active, c, c-1) {
			return
		}
	}
}

func (b *BufferPage) startExclusiveAccess() bool {
	for {
		if atomic.CompareAndSwapUint32(&b.active, 0, math.MaxUint32) {
			return true
		}
		if atomic.LoadUint32(&b.active) == math.MaxUint32 {
			panic("illegal state")
		}
	}
}

func (b *BufferPage) stopExclusiveAccess() {
	if !atomic.CompareAndSwapUint32(&b.active, math.MaxUint32, 0) {
		panic("illegal state")
	}
}

func (b *BufferPage) getBufLen() uint32 {
	pos := atomic.LoadUint32(&b.bufferCleanupData.offset)
	if pos == math.MaxUint32 {
		return 0
	}
	return pos
}

func (b *BufferPage) trimData(byteRecycle ByteBufferRecycler, recycleSize uint32) uint32 {
	data := b.data
	atomic.StoreUint32(&b.offset, math.MaxUint32)
	b.data = nil
	if byteRecycle != nil && uint32(len(data)) == recycleSize {
		byteRecycle.Put(data)
	}
	return uint32(len(data))
}

func (b *BufferPage) allocateSlice(reqLen uint32) []byte {
	for {
		pos := atomic.LoadUint32(&b.offset)

		if pos == math.MaxUint32 || uint32(len(b.data)) < pos+reqLen {
			return nil
		}

		if atomic.CompareAndSwapUint32(&b.offset, pos, pos+reqLen) {
			atomic.AddUint32(&b.bufferCleanupData.count, 1)
			return b.data[pos : pos+reqLen : pos+reqLen]
		}
	}
}
