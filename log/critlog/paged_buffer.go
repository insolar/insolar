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

package critlog

import (
	"encoding/binary"
	"math"
	"sync/atomic"
	"unsafe"
)

func NewPagedBuffer(pageSize, sizeLimit int, trimHead bool) PagedBuffer {
	if pageSize <= 0 {
		panic("illegal value")
	}
	if sizeLimit <= 0 || sizeLimit >= math.MaxInt32 {
		trimHead = false
		sizeLimit = math.MaxInt32
	}

	return PagedBuffer{defaultPageSize: uint32(pageSize), capacityLimit: uint32(sizeLimit), trimHead: trimHead}
}

type PagedBuffer struct {
	defaultPageSize uint32
	next            unsafe.Pointer // atomic, *BufferPage

	capacityLimit uint32
	trimHead      bool
	//
	//capacityReleased uint32
	//lastKnownHead *BufferPage
}

type pageFifo struct {
	next unsafe.Pointer // atomic, *BufferPage
}

func (p *pageFifo) add(page *BufferPage) {

}

type BufferPage struct {
	active uint32 // atomic
	offset uint32 // atomic

	head *BufferPage

	data []byte

	totalCapacity uint32
	next          *BufferPage // direction to the head, =nul for head

	prev unsafe.Pointer // atomic, *BufferPage, for trimHead==true only
	trim *BufferTrim    // for trimHead==true only and for head only
}

type BufferTrim struct {
	active          uint32 // atomic
	trimmedCapacity uint32
	target          *BufferPage
}

func (p *PagedBuffer) currentPage() *BufferPage {
	return (*BufferPage)(atomic.LoadPointer(&p.next))
}

func (p *PagedBuffer) addPage(checkExpected bool, expectedPage, newPage *BufferPage) bool {
	newPage.initAccess() // make sure that the tail is always "active"
	for {
		current := p.currentPage()
		if checkExpected && current != expectedPage {
			return false
		}

		newPage.totalCapacity = uint32(len(newPage.data))
		if current == nil {
			newPage.head = newPage
			newPage.trim = &BufferTrim{}
		} else {
			newPage.totalCapacity += current.totalCapacity
			newPage.head = current.head
			newPage.trim = nil
		}

		newPage.next = current
		if atomic.CompareAndSwapPointer(&p.next, unsafe.Pointer(current), unsafe.Pointer(newPage)) {
			if p.trimHead {
				atomic.StorePointer(&current.prev, unsafe.Pointer(newPage))
			}
			current.stopAccess()
			return true
		}
	}
}

func (p *PagedBuffer) createPage(capacity uint32) *BufferPage {
	return &BufferPage{data: make([]byte, capacity)}
}

func (p *PagedBuffer) flushPages() *BufferPage {
	return (*BufferPage)(atomic.SwapPointer(&p.next, nil))
}

func (p *PagedBuffer) FlushPages() []*BufferPage {
	tail := p.flushPages()
	if tail == nil {
		return nil
	}

	pageCount := 0
	for n := tail; n != nil; n = n.next {
		pageCount++
	}

	pages := make([]*BufferPage, pageCount)

	for n := tail; n != nil; n = n.next {
		pageCount--
		pages[pageCount] = n
	}

	return pages
}

const serviceHeaderLen = 4

func (p *PagedBuffer) allocateBuffer(reqLen uint32) (*BufferPage, []byte) {
	if reqLen == 0 {
		panic("illegal value")
	}

	reqLen += serviceHeaderLen

	if reqLen >= p.defaultPageSize {
		newPage := p.createPage(reqLen)
		newPage.offset = reqLen
		writeServiceHeader(reqLen, newPage.data[0:serviceHeaderLen:serviceHeaderLen])
		result := newPage.data[serviceHeaderLen:reqLen:reqLen]

		p.addPage(false, nil, newPage)
		return newPage, result
	}

	current := p.currentPage()
	for {
		if current != nil {

			if !current.startAccess() {
				current = p.currentPage()
				continue
			}

			// TODO need to back-sync for possible current page swap

			pos := atomic.LoadUint32(&current.offset)

			if uint32(len(current.data)) >= pos+reqLen {
				if atomic.CompareAndSwapUint32(&current.offset, pos, pos+reqLen) {
					writeServiceHeader(reqLen, current.data[pos:pos+serviceHeaderLen:pos+serviceHeaderLen])
					return current, current.data[pos+serviceHeaderLen : pos+reqLen : pos+reqLen]
				}
			}
			atomic.AddUint32(&current.active, math.MaxUint32)

			if p.trimHead {
				if current.next != nil && current.totalCapacity >= p.capacityLimit {
					p.trimFromHead(current)
				}
				current.stopAccess()
			} else {
				current.stopAccess()
				if current.totalCapacity >= p.capacityLimit {
					return nil, nil
				}
			}
		}

		newPage := p.createPage(p.defaultPageSize)
		if p.addPage(true, current, newPage) {
			current = newPage
		} else {
			current = p.currentPage()
		}
	}
}

var byteOrder = binary.LittleEndian

func writeServiceHeader(allocatedLen uint32, buf []byte) {
	byteOrder.PutUint32(buf, allocatedLen)
}

func readServiceHeader(buf []byte) (allocatedLen uint32) {
	return byteOrder.Uint32(buf)
}

func (p *PagedBuffer) Write(b []byte) (int, error) {
	ln := uint32(len(b))
	if ln == 0 {
		return 0, nil
	}
	pg, buf := p.allocateBuffer(ln)

	copy(buf, b)
	pg.stopAccess()
	return int(ln), nil
}

func (p *PagedBuffer) WriteWithPrefix(prefix, data []byte) (int, error) {
	ln := uint32(len(data) + len(prefix))
	if ln == 0 {
		return 0, nil
	}
	pg, buf := p.allocateBuffer(ln)

	pLen := copy(buf, prefix)
	copy(buf[pLen:], data)
	pg.stopAccess()
	return int(ln), nil
}

func (p *PagedBuffer) trimFromHead(current *BufferPage) {
	if current.head == current {
		// unable to trim current
		panic("illegal state")
	}

	trim := current.head.trim
	if !trim.start() {
		return
	}
	defer trim.stop()

	if trim.target == nil {
		trim.target = current.head
	}

	for p.capacityLimit+trim.trimmedCapacity < current.totalCapacity {
		if !trim.target.startExclusiveAccess() {
			return
		}
		trim.trimmedCapacity += uint32(len(trim.target.data))
		atomic.StoreUint32(&trim.target.offset, math.MaxUint32)
		trim.target.data = nil
		trim.target.stopExclusiveAccess()
		trim.target = (*BufferPage)(atomic.LoadPointer(&trim.target.prev))
	}
}

/* ============================ */

func (b *BufferPage) initAccess() {
	if !atomic.CompareAndSwapUint32(&b.active, 0, 1) {
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
	atomic.AddUint32(&b.active, math.MaxUint32)
}

func (b *BufferPage) startExclusiveAccess() bool {
	return atomic.CompareAndSwapUint32(&b.active, 0, math.MaxUint32)
}

func (b *BufferPage) stopExclusiveAccess() {
	if !atomic.CompareAndSwapUint32(&b.active, math.MaxUint32, 0) {
		panic("illegal state")
	}
}

func (b *BufferPage) StartRead(idleFn func()) {
	for !b.startExclusiveAccess() {
		idleFn()
	}
}

func (b *BufferPage) StopRead() {
	b.stopExclusiveAccess()
}

func (b *BufferPage) ReadNext(offset uint32) (nextOffset uint32, data []byte) {
	limit := atomic.LoadUint32(&b.offset)

	if limit == math.MaxUint32 || offset >= limit {
		return 0, nil
	}

	chunkLen := readServiceHeader(b.data[offset:])
	nextOffset = chunkLen + offset
	if nextOffset > limit {
		panic("illegal state")
	}
	chunk := b.data[offset+serviceHeaderLen : offset+chunkLen : offset+chunkLen]
	if nextOffset == limit {
		return 0, chunk
	}
	return nextOffset, chunk
}

/* ============================ */

func (p *BufferTrim) start() bool {
	return atomic.CompareAndSwapUint32(&p.active, 0, 1)
}

func (p *BufferTrim) stop() {
	if !atomic.CompareAndSwapUint32(&p.active, 1, 0) {
		panic("illegal state")
	}
}
