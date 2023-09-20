package lazyhead

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func lessTestFn(v1 interface{}, v2 interface{}) bool {
	return v1.(int) < v2.(int)
}

func TestNewHeadedLazySortedList(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, -1)
	require.Len(t, hl.data.data, 0)

	require.Equal(t, 2, cap(hl.data.data))

	hl = NewHeadedLazySortedList(1, lessTestFn, 0)
	require.Equal(t, 2, cap(hl.data.data))

	hl = NewHeadedLazySortedList(3, lessTestFn, 1)
	require.Equal(t, 3, cap(hl.data.data))

	hl = NewHeadedLazySortedList(1, lessTestFn, 3)
	require.Equal(t, 3, cap(hl.data.data))
}

func TestInnerLen(t *testing.T) {
	inl := innerHeadedLazySortedList{data: make([]interface{}, 0), less: lessTestFn}
	inl.Add(1)

	require.Equal(t, 1, inl.Len())
}

func TestInnerLess(t *testing.T) {
	inl := innerHeadedLazySortedList{data: make([]interface{}, 0), less: lessTestFn}
	inl.Add(1)
	inl.Add(2)
	require.Equal(t, lessTestFn(0, 1), inl.Less(0, 1))
}

func TestInnerSwap(t *testing.T) {
	inl := innerHeadedLazySortedList{data: make([]interface{}, 0), less: lessTestFn}
	inl.Add(1)
	inl.Add(2)
	inl.Swap(0, 1)
	require.Equal(t, 2, inl.Get(0))

	require.Equal(t, 1, inl.Get(1))
}

func TestInnerAdd(t *testing.T) {
	inl := innerHeadedLazySortedList{data: make([]interface{}, 0), less: lessTestFn}
	inl.Add(1)
	require.Equal(t, 1, inl.Get(0))

	require.Equal(t, 1, inl.len)

	inl = innerHeadedLazySortedList{data: make([]interface{}, 2), less: lessTestFn}
	inl.Add(1)
	require.Equal(t, 1, inl.Get(0))

	require.Equal(t, 1, inl.len)

	inl.Add(3)
	require.Equal(t, 1, inl.Get(0))

	require.Equal(t, 3, inl.Get(1))

	require.Equal(t, 2, inl.len)
}

func TestInnerGet(t *testing.T) {
	inl := innerHeadedLazySortedList{data: make([]interface{}, 0), less: lessTestFn}
	inl.Add(2)
	inl.Add(3)
	require.Equal(t, 2, inl.Get(0))

	require.Equal(t, 3, inl.Get(1))

	require.Panics(t, func() { inl.Get(-1) })

	require.Panics(t, func() { inl.Get(2) })
}

func TestLen(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, 1)
	require.Zero(t, hl.Len())

	hl.Add(2)
	require.Equal(t, 1, hl.Len())

	hl.Add(3)
	require.Equal(t, 2, hl.Len())
}

func TestGet(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, 1)
	hl.Add(2)
	hl.Add(3)
	require.Equal(t, 2, hl.Get(0))

	require.Equal(t, 3, hl.Get(1))

	require.Panics(t, func() { hl.Get(-1) })

	require.Panics(t, func() { hl.Get(2) })
}

func TestAdd(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, 1)
	hl.Add(1)
	require.Equal(t, Sorted, hl.sorted)

	hl.Add(2)
	require.Equal(t, UnsortedTail, hl.sorted)

	hl.Add(0)
	require.Equal(t, UnsortedTail, hl.sorted)

	hl = NewHeadedLazySortedList(2, lessTestFn, 2)
	hl.Add(1)
	hl.sorted = UnsortedAll
	hl.Add(2)
	require.Equal(t, Sorted, hl.sorted)

	hl.Add(1)
	require.Equal(t, Sorted, hl.sorted)

	hl = NewHeadedLazySortedList(4, lessTestFn, 3)
	hl.Add(1)
	hl.Add(2)
	require.Equal(t, UnsortedAll, hl.sorted)

	hl.Add(1)
	require.Equal(t, UnsortedAll, hl.sorted)
}

func TestSortAll(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, 1)
	hl.Add(2)
	hl.SortAll()
	require.Equal(t, Sorted, hl.sorted)

	hl.Add(3)
	require.Equal(t, UnsortedTail, hl.sorted)

	hl.SortAll()
	require.Equal(t, Sorted, hl.sorted)
}

func TestInnerCutOffHeadByLen(t *testing.T) {
	inl := innerHeadedLazySortedList{data: make([]interface{}, 0), less: lessTestFn}
	inl.Add(2)
	inl.Add(3)
	to := inl.cutOffHeadByLen(0, nil)
	require.Len(t, to, 0)

	require.Len(t, inl.data, 2)

	to = inl.cutOffHeadByLen(1, nil)
	require.Len(t, to, 1)

	require.Equal(t, 1, inl.len)

	require.Nil(t, inl.data[1])

	inl.Add(4)
	to2 := make([]interface{}, 1)
	to = inl.cutOffHeadByLen(2, to2)
	require.Len(t, to, 3)

	require.Zero(t, inl.len)

	require.Nil(t, inl.data[1])
}

func TestGetReversedHead(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, 1)
	hl.Add(2)
	require.Panics(t, func() { hl.GetReversedHead(-1) })

	require.Panics(t, func() { hl.GetReversedHead(2) })

	k := hl.GetReversedHead(0)
	require.Equal(t, 2, k.(int))
}

func TestHasFullHead(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, 1)
	require.Panics(t, func() { hl.HasFullHead(-1) })

	require.Panics(t, func() { hl.HasFullHead(2) })

	require.True(t, hl.HasFullHead(1))

	require.False(t, hl.HasFullHead(0))
}

func TestCheckHeadLen(t *testing.T) {
	hl := NewHeadedLazySortedList(3, lessTestFn, 1)
	hl.Add(2)
	require.Panics(t, func() { hl.checkAndGetAdjustedHeadLen(-1) })

	require.Panics(t, func() { hl.checkAndGetAdjustedHeadLen(4) })

	hl.Add(3)
	hlr := 1
	n := hl.checkAndGetAdjustedHeadLen(hlr)
	require.Equal(t, hl.headLen-hlr, n)

	n = hl.checkAndGetAdjustedHeadLen(0)
	require.Equal(t, hl.data.len, n)
}

func TestCutOffHeadInto(t *testing.T) {
	hl := NewHeadedLazySortedList(3, lessTestFn, 1)
	item := 2
	hl.Add(item)
	to := hl.CutOffHeadInto(1, nil)
	require.Len(t, to, 1)

	require.Equal(t, []interface{}{item}, to)
}

func TestCutOffHeadByLenInto(t *testing.T) {
	hl := NewHeadedLazySortedList(3, lessTestFn, 1)
	item1 := 2
	hl.Add(item1)
	to := hl.CutOffHeadByLenInto(0, nil)
	require.Nil(t, to)

	to2 := make([]interface{}, 1)
	to = hl.CutOffHeadByLenInto(0, to2)
	require.Len(t, to, len(to2))

	item2 := 3
	hl.Add(item2)
	hl.sorted = UnsortedAll
	to = hl.CutOffHeadByLenInto(1, to2)
	require.Equal(t, Sorted, hl.sorted)

	require.Equal(t, 1, hl.data.len)

	require.Equal(t, item2, hl.Get(0))

	require.Equal(t, item1, to[1].(int))

	to = hl.CutOffHeadByLenInto(1, to2)
	require.Equal(t, Sorted, hl.sorted)

	require.Equal(t, item2, to[1].(int))

	hl.Add(3)
	hl.Add(2)
	hl.Add(3)
	hl.Add(5)
	hl.Add(7)
	require.Equal(t, UnsortedTail, hl.sorted)

	to = hl.CutOffHeadByLenInto(2, to2)
	require.Equal(t, UnsortedAll, hl.sorted)
}

func TestCutOffHead(t *testing.T) {
	hl := NewHeadedLazySortedList(2, lessTestFn, 1)
	require.Equal(t, []interface{}{}, hl.CutOffHead(0))

	item1 := 4
	hl.Add(item1)
	item2 := 2
	hl.Add(item2)
	require.Equal(t, []interface{}{item2}, hl.CutOffHead(1))

	hl.Add(item1)
	hl.Add(item2)
	item3 := 3
	hl.Add(item3)
	require.Equal(t, []interface{}{item2, item3}, hl.CutOffHead(0))

	require.Panics(t, func() { hl.CutOffHead(3) })
}

func TestCutOffHeadByLen(t *testing.T) {
	hl := NewHeadedLazySortedList(2, lessTestFn, 1)
	require.Equal(t, []interface{}{}, hl.CutOffHeadByLen(0))

	item1 := 4
	hl.Add(item1)
	item2 := 2
	hl.Add(item2)
	require.Equal(t, []interface{}{item2, item1}, hl.CutOffHeadByLen(2))

	hl.Add(item1)
	hl.Add(item2)
	item3 := 3
	hl.Add(item3)
	require.Equal(t, []interface{}{item2, item3}, hl.CutOffHeadByLen(2))

	require.Panics(t, func() { hl.CutOffHeadByLen(3) })
}

func TestFlush(t *testing.T) {
	hl := NewHeadedLazySortedList(2, lessTestFn, 1)
	item1 := 4
	hl.Add(item1)
	item2 := 2
	hl.Add(item2)
	res := hl.Flush()
	require.Equal(t, []interface{}{item2, item1}, res)

	require.Equal(t, Sorted, hl.sorted)
}

func TestGetHeadLen(t *testing.T) {
	headLen := 2
	hl := NewHeadedLazySortedList(headLen, lessTestFn, 1)
	require.Equal(t, headLen, hl.GetHeadLen())
}

func TestGetAvailableHeadLen(t *testing.T) {
	headLen := 2
	hl := NewHeadedLazySortedList(headLen, lessTestFn, 1)
	item1 := 4
	hl.Add(item1)
	item2 := 2
	hl.Add(item2)
	hlr := 1
	require.Equal(t, headLen-hlr, hl.GetAvailableHeadLen(hlr))

	require.Panics(t, func() { hl.GetAvailableHeadLen(-1) })

	require.Panics(t, func() { hl.GetAvailableHeadLen(3) })

	headLen = 3
	hl = NewHeadedLazySortedList(headLen, lessTestFn, 1)
	hl.Add(item1)
	require.Equal(t, 1, hl.GetAvailableHeadLen(1))
}
