// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build ignore

package logbuf

import (
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

const maxCapacity = math.MaxInt32 // not math.MaxUint32 to avoid overflow

func NewPagedBufferTrimFromOldest(pageSize, sizeLimit int, recycler ByteBufferRecycler) PagedBuffer {
	if sizeLimit <= 0 || sizeLimit >= maxCapacity {
		return NewPagedBufferTrimFromLatest(pageSize, sizeLimit)
	}

	return newPagedBuffer(pageSize, sizeLimit, true, recycler, func(newPage *BufferPage) {
		newPage.bufferCleanupData.totalCapacity = uint32(len(newPage.data))
		n := newPage.next
		if n == nil {
			newPage.head = newPage
			newPage.trim = &BufferTrim{}
		} else {
			newPage.totalCapacity += n.totalCapacity
			newPage.head = n.head
			newPage.trim = nil
		}
	})
}

func NewPagedBufferTrimFromLatest(pageSize, sizeLimit int) PagedBuffer {
	return newPagedBuffer(pageSize, sizeLimit, false, nil, func(newPage *BufferPage) {
		newPage.totalCapacity = uint32(len(newPage.data))
		n := newPage.next
		newPage.trim = nil
		if n == nil {
			newPage.head = newPage
		} else {
			newPage.totalCapacity += n.totalCapacity
			newPage.head = n.head
		}
	})
}

func newPagedBuffer(pageSize, sizeLimit int, trimHead bool, recycler ByteBufferRecycler, bufferPrepareFn func(newPage *BufferPage)) PagedBuffer {
	if pageSize <= 256 {
		panic("illegal value")
	}

	if sizeLimit <= 0 || sizeLimit >= maxCapacity {
		sizeLimit = maxCapacity
	}

	return PagedBuffer{
		defaultPageSize:  uint32(pageSize),
		oneShotThreshold: uint32(pageSize*4) / 5,
		capacityLimit:    uint32(sizeLimit),
		trimHead:         trimHead,
		buffer:           newPageLifo(bufferPrepareFn),
		recycle:          newPageLifo(nil),
		byteRecycle:      recycler,
	}
}

type PagedBuffer struct {
	defaultPageSize  uint32
	oneShotThreshold uint32
	capacityLimit    uint32
	trimHead         bool

	buffer  pageLifo
	recycle pageLifo

	byteRecycle ByteBufferRecycler
}

type ByteBufferRecycler interface {
	Get() []byte
	Put([]byte)
}

type BufferTrim struct {
	active          uint32 // atomic
	mutex           sync.Mutex
	trimmedCapacity uint32
	target          *BufferPage
}

func (p *PagedBuffer) currentPage() *BufferPage {
	return p.buffer.peek()
}

func (p *PagedBuffer) addPageWithExpected(replaced, newPage *BufferPage) bool {
	if !p.buffer.pushExpected(replaced, newPage) {
		return false
	}
	if replaced != nil {
		if p.trimHead {
			atomic.StorePointer(&replaced.prev, unsafe.Pointer(newPage))
		}
		replaced.stopAccess()
	}
	return true
}

func (p *PagedBuffer) addPage(newPage *BufferPage) {
	replaced := p.buffer.push(newPage)
	if replaced != nil {
		if p.trimHead {
			atomic.StorePointer(&replaced.prev, unsafe.Pointer(newPage))
		}
		replaced.stopAccess()
	}
}

func (p *PagedBuffer) createPageWithSize(capacity uint32) *BufferPage {
	if capacity == p.defaultPageSize {
		return p.createPage()
	}
	return &BufferPage{data: make([]byte, capacity)}
}

func (p *PagedBuffer) createPage() *BufferPage {
	for {
		recycledPage := p.recycle.pop()
		if recycledPage == nil {
			break
		}
		if recycledPage.tryInitAccess() && uint32(len(recycledPage.data)) == p.defaultPageSize {
			//recycledPages.bufferCleanupData = bufferCleanupData{}
			recycledPage.next = nil
			recycledPage.stopAccess()
			return recycledPage
		}
	}

	if p.byteRecycle != nil {
		recycled := p.byteRecycle.Get()
		if uint32(len(recycled)) == p.defaultPageSize {
			return &BufferPage{data: recycled}
		}
	}
	return &BufferPage{data: make([]byte, p.defaultPageSize)}
}

func (p *PagedBuffer) flushPages() *BufferPage {
	page := p.buffer.flush()
	if page != nil {
		page.stopAccess()
	}
	return page
}

func (p *PagedBuffer) FlushPages() FlushedPages {
	return FlushedPages{buffer: p, flushedPages: p.flushPages()}
}

// MUST: only pages with len(data)==defaultPageSize can be put here
func (p *PagedBuffer) _recycleFlushedPages(pages *BufferPage) {
	if pages == nil {
		return
	}
	p.recycle.replace(pages)
}

func (p *PagedBuffer) allocateBuffer(reqLen uint32) (*BufferPage, []byte) {
	if reqLen == 0 {
		panic("illegal value")
	}

	if reqLen >= p.defaultPageSize {
		return p.allocateOneShotBuffer(reqLen)
	}

	current := p.currentPage()
	var newPage *BufferPage
	for {
		if current != nil {
			if !current.startAccess() {
				// page was flushed before we've got it
				current = p.currentPage()
				continue
			}
			buf := current.allocateSlice(reqLen)
			if buf != nil {
				return current, buf
			}

			// we need more space than available in the page, but have to check limits and try to trim

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

		if reqLen >= p.oneShotThreshold {
			return p.allocateOneShotBuffer(reqLen)
		}

		if newPage == nil {
			newPage = p.createPage()
			newPage.initAccess()
		}
		if p.addPageWithExpected(current, newPage) {
			current = newPage
			newPage = nil
		} else {
			current = p.currentPage()
		}
	}
}

func (p *PagedBuffer) allocateOneShotBuffer(reqLen uint32) (*BufferPage, []byte) {
	newPage := p.createPageWithSize(reqLen)
	newPage.initAccess()

	newPage.offset = reqLen
	p.addPage(newPage)
	return newPage, newPage.data
}

func (p *PagedBuffer) trimFromHead(current *BufferPage) {
	if current.head == current {
		// unable to trim current
		panic("illegal state")
	}

	trim := current.head.bufferCleanupData.trim
	if !trim.start() {
		return
	}
	defer trim.stop()

	if trim.target == nil {
		trim.target = current.bufferCleanupData.head
	}

	for p.capacityLimit+trim.trimmedCapacity < current.totalCapacity {
		if trim.target == nil || !trim.target.startExclusiveAccess() {
			return
		}
		trim.trimmedCapacity += trim.target.trimData(p.byteRecycle, p.defaultPageSize)
		trim.target.stopExclusiveAccess()
		trim.target = (*BufferPage)(atomic.LoadPointer(&trim.target.bufferCleanupData.prev))
	}
}

/* ============================ */

func (p *BufferTrim) start() bool {
	if atomic.CompareAndSwapUint32(&p.active, 0, 1) {
		p.mutex.Lock()
		return true
	}
	return false
}

func (p *BufferTrim) stop() {
	if !atomic.CompareAndSwapUint32(&p.active, 1, 0) {
		panic("illegal state")
	}
	p.mutex.Unlock()
}

/* ============================ */

type FlushedPages struct {
	buffer       *PagedBuffer
	flushedPages *BufferPage
	pages        []*BufferPage
}

func (p *FlushedPages) StartAccess(idleFn func()) {
	if idleFn == nil {
		idleFn = runtime.Gosched
	}

	if p.flushedPages == nil {
		return
	}

	pageCount := 0
	for n := p.flushedPages; n != nil; n = n.next {
		pageCount++
	}

	p.pages = make([]*BufferPage, pageCount)
	for n := p.flushedPages; n != nil; n = n.next {
		for !n.startExclusiveAccess() {
			idleFn()
		}
		pageCount--
		p.pages[pageCount] = n
	}

	p.flushedPages = nil
}

func (p *FlushedPages) StopAccess() {
	var recyclePages *BufferPage

	for _, pg := range p.pages {
		if uint32(len(pg.data)) != p.buffer.defaultPageSize {
			continue
		}
		pg.stopExclusiveAccess()
		pg.bufferCleanupData = bufferCleanupData{}
		pg.next = recyclePages
		recyclePages = pg
	}

	if recyclePages != nil {
		p.buffer._recycleFlushedPages(recyclePages)
	}
}

func (p *FlushedPages) Count() int {
	return len(p.pages)
}

func (p *FlushedPages) Page(i int) (uint32, []byte) {
	pg := p.pages[i]
	data := pg.data
	if data == nil {
		return pg.count, nil
	}
	return pg.count, data[:pg.getBufLen()]
}
