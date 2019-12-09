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

package unsafekit

import (
	"fmt"
	"reflect"
	"runtime"
	"unsafe"

	"github.com/insolar/insolar/longbits"
)

func KeepAliveWhile(p unsafe.Pointer, fn func(unsafe.Pointer) uintptr) uintptr {
	res := fn(p)
	runtime.KeepAlive(p)
	return res
}

const keepAlivePageLen = 4096/PtrSize - 4

type keepAlivePage [keepAlivePageLen]unsafe.Pointer

type KeepAliveList struct {
	list          []unsafe.Pointer
	lastPageCount uint16
	listOfPages   bool
}

func (p *KeepAliveList) Reset() {
	*p = KeepAliveList{}
}

func (p *KeepAliveList) Keep(v interface{}) {
	vv := reflect.ValueOf(v)
	switch vv.Kind() {
	case reflect.Chan, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Slice:
		p.keep(vv.Pointer())
	case reflect.String:
		p.keepStrData(vv.String())
	default:
		panic(fmt.Sprint("unsupported value", vv))
	}
	runtime.KeepAlive(v)
}

func (p *KeepAliveList) KeepDataOf(v longbits.ByteString) {
	p.keepStrData(string(v))
	runtime.KeepAlive(v)
}

func (p *KeepAliveList) KeepDataOfBytes(v []byte) {
	pSlice := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	p.keep(pSlice.Data)
	runtime.KeepAlive(v)
}

func (p *KeepAliveList) keepStrData(s string) {
	pString := (*reflect.StringHeader)(unsafe.Pointer(&s))
	p.keep(pString.Data)
}

func (p *KeepAliveList) keep(v uintptr) {
	switch {
	case p.listOfPages:
		if p.lastPageCount == keepAlivePageLen {
			break
		}

		switch {
		case len(p.list) < 2:
			panic("illegal state")
		case p.lastPageCount > keepAlivePageLen:
			panic("illegal state")
		}
		lastPage := (*keepAlivePage)(p.list[len(p.list)-1])
		lastPage[p.lastPageCount] = unsafe.Pointer(v)
		p.lastPageCount++
		return

	case p.lastPageCount != 0:
		panic("illegal state")

	case len(p.list) == keepAlivePageLen:
		lastPage := new(keepAlivePage)
		copy(lastPage[:], p.list)
		p.list = make([]unsafe.Pointer, 1, 2)
		p.list[0] = unsafe.Pointer(lastPage)

	case len(p.list) > keepAlivePageLen:
		panic("illegal state")

	default:
		p.list = append(p.list, unsafe.Pointer(v))
		return
	}

	switch {
	case p.lastPageCount != keepAlivePageLen:
		panic("illegal state")
	case len(p.list) < 1:
		panic("illegal state")
	}

	lastPage := new(keepAlivePage)
	p.list = append(p.list, (unsafe.Pointer)(lastPage))
	lastPage[0] = unsafe.Pointer(p)
	p.lastPageCount = 1
}
