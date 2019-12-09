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

package filaments

import (
	"sync/atomic"
	"unsafe"
)

func NewAtomicEntry(entry *WriteEntry) AtomicEntry {
	return AtomicEntry{entry}
}

type AtomicEntry struct {
	entry *WriteEntry
}

func (p *AtomicEntry) _ptr() *unsafe.Pointer {
	return (*unsafe.Pointer)((unsafe.Pointer)(&p.entry))
}

func (p *AtomicEntry) Get() *WriteEntry {
	return (*WriteEntry)(atomic.LoadPointer(p._ptr()))
}

func (p *AtomicEntry) Set(v *WriteEntry) {
	atomic.StorePointer(p._ptr(), (unsafe.Pointer)(v))
}

func (p *AtomicEntry) Swap(v *WriteEntry) *WriteEntry {
	return (*WriteEntry)(atomic.SwapPointer(p._ptr(), (unsafe.Pointer)(v)))
}

func (p *AtomicEntry) CmpAndSwap(expected, new *WriteEntry) bool {
	return atomic.CompareAndSwapPointer(p._ptr(), (unsafe.Pointer)(expected), (unsafe.Pointer)(new))
}
