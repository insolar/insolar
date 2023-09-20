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
