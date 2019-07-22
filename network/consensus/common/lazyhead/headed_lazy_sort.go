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

package lazyhead

import (
	"sort"
)

type LessFunc func(v1, v2 interface{}) bool

type HeadedLazySortedList struct {
	headLen int
	sorted  sortState
	data    innerHeadedLazySortedList
}

func NewHeadedLazySortedList(headLen int, less LessFunc, capacity int) HeadedLazySortedList {
	r := HeadedLazySortedList{data: innerHeadedLazySortedList{less: less}, headLen: headLen}

	switch {
	case capacity <= 0:
		capacity = headLen << 1
	case capacity < headLen:
		capacity = headLen
	}
	r.data.data = make([]interface{}, 0, capacity)
	return r
}

type innerHeadedLazySortedList struct {
	len  int
	data []interface{}
	less LessFunc
}

type sortState int

const (
	Sorted sortState = iota
	UnsortedTail
	UnsortedAll
)

func (r *innerHeadedLazySortedList) Len() int {
	return r.len
}

func (r *innerHeadedLazySortedList) Less(i, j int) bool {
	return r.less(r.data[i], r.data[j])
}

func (r *innerHeadedLazySortedList) Swap(i, j int) {
	r.data[i], r.data[j] = r.data[j], r.data[i]
}

func (r *innerHeadedLazySortedList) Add(item interface{}) {
	if r.len == len(r.data) {
		r.data = append(r.data, item)
	} else {
		r.data[r.len] = item
	}
	r.len++
}

func (r *innerHeadedLazySortedList) Get(i int) interface{} {
	if i >= r.len {
		panic("index out of range")
	}
	return r.data[i]
}

func (r *HeadedLazySortedList) Len() int {
	return r.data.len
}

func (r *HeadedLazySortedList) Get(index int) interface{} {
	return r.data.Get(index)
}

func (r *HeadedLazySortedList) Add(item interface{}) {
	r.data.Add(item)
	switch {
	case r.data.len == 1:
		// one-size array is always sorted
		r.sorted = Sorted
	case r.data.len == r.headLen: // TODO check actual performance impact of this tweak
		// Perf tweak - head part of longer arrays should be sorted to ease insertion
		sort.Sort(&r.data)
		r.sorted = Sorted
	case r.sorted == UnsortedAll:
		// unsorted remains unsorted
		return
	case r.data.len < r.headLen:
		// shorter arrays are sorted on demand
		r.sorted = UnsortedAll
	case r.data.less(r.data.data[r.headLen-1], item):
		// Sorted or UnsortedTail with more items than headLen
		// and the new item doesn't fit into head by weight
		// so we leave it in the tail
		r.sorted = UnsortedTail
	default:
		// Sorted or UnsortedTail with more items than headLen
		// and the new item should be fit into head by weight
		// so we find an insertion point in the head and move it there
		var insertPos int
		if r.headLen == 1 {
			insertPos = 0
		} else {
			insertPos = sort.Search(r.headLen, func(i int) bool { return r.data.less(item, r.data.data[i]) })
		}
		copy(r.data.data[insertPos+1:r.data.len], r.data.data[insertPos:r.data.len])
		r.data.data[insertPos] = item
	}
}

func (r *HeadedLazySortedList) SortAll() {
	if r.sorted == Sorted {
		return
	}
	sort.Sort(&r.data)
	r.sorted = Sorted
}

func (r *innerHeadedLazySortedList) cutOffHeadByLen(headCutLen int, to []interface{}) []interface{} {

	if to == nil {
		to = make([]interface{}, headCutLen)
		copy(to, r.data)
	} else {
		to = append(to, r.data[:headCutLen]...)
	}
	copy(r.data, r.data[headCutLen:r.len])
	for ; headCutLen > 0; headCutLen-- {
		r.len--
		r.data[r.len] = nil
	}
	return to
}

func (r *HeadedLazySortedList) GetReversedHead(relIndex int) interface{} {
	return r.Get(r.checkAndGetAdjustedHeadLen(relIndex) - 1)
}

func (r *HeadedLazySortedList) HasFullHead(headLenReduction int) bool {
	if headLenReduction < 0 || headLenReduction > r.headLen {
		panic("index out of range")
	}
	return r.headLen <= headLenReduction+r.data.len
}

func (r *HeadedLazySortedList) checkAndGetAdjustedHeadLen(headLenReduction int) int {
	if headLenReduction < 0 || headLenReduction > r.headLen {
		panic("index out of range")
	}
	headCutLen := r.headLen - headLenReduction
	if headCutLen > r.data.len {
		return r.data.len
	}
	return headCutLen
}

func (r *HeadedLazySortedList) CutOffHeadInto(headLenReduction int, to []interface{}) []interface{} {

	headCutLen := r.checkAndGetAdjustedHeadLen(headLenReduction)
	return r.CutOffHeadByLenInto(headCutLen, to)
}

func (r *HeadedLazySortedList) CutOffHeadByLenInto(headCutLen int, to []interface{}) []interface{} {

	if headCutLen == 0 {
		return to
	}

	if r.sorted == UnsortedAll {
		r.SortAll()
		return r.data.cutOffHeadByLen(headCutLen, to)
	}

	res := r.data.cutOffHeadByLen(headCutLen, to)
	if r.data.len <= 1 {
		r.sorted = Sorted
	} else if r.sorted == UnsortedTail {
		r.sorted = UnsortedAll
	}
	return res
}

func (r *HeadedLazySortedList) CutOffHead(headLenReduction int) []interface{} {

	res := r.CutOffHeadInto(headLenReduction, nil)
	if res == nil {
		return []interface{}{}
	}
	return res
}

func (r *HeadedLazySortedList) CutOffHeadByLen(headLen int) []interface{} {

	res := r.CutOffHeadByLenInto(headLen, nil)
	if res == nil {
		return []interface{}{}
	}
	return res
}

func (r *HeadedLazySortedList) Flush() []interface{} {
	// return the old array to avoid cleanup
	res := r.data.data[:r.data.len]
	r.data.data = make([]interface{}, 0, r.headLen<<1)
	r.sorted = Sorted
	return res
}

func (r *HeadedLazySortedList) GetHeadLen() int {

	return r.headLen
}

func (r *HeadedLazySortedList) GetAvailableHeadLen(headLenReduction int) int {

	return r.checkAndGetAdjustedHeadLen(headLenReduction)
}
