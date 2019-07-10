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

package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func lessTestFn(v1 interface{}, v2 interface{}) bool {
	return v1.(int) < v2.(int)
}

func TestNewHeadedLazySortedList(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, -1)
	require.Equal(t, len(hl.data.data), 0)

	require.Equal(t, cap(hl.data.data), 2)

	hl = NewHeadedLazySortedList(1, lessTestFn, 0)
	require.Equal(t, cap(hl.data.data), 2)

	hl = NewHeadedLazySortedList(3, lessTestFn, 1)
	require.Equal(t, cap(hl.data.data), 3)

	hl = NewHeadedLazySortedList(1, lessTestFn, 3)
	require.Equal(t, cap(hl.data.data), 3)
}

func TestInnerLen(t *testing.T) {
	inl := innerHeadedLazySortedList{data: make([]interface{}, 0), less: lessTestFn}
	inl.Add(1)

	require.Equal(t, inl.Len(), 1)
}

func TestInnerLess(t *testing.T) {
	inl := innerHeadedLazySortedList{data: make([]interface{}, 0), less: lessTestFn}
	inl.Add(1)
	inl.Add(2)
	require.Equal(t, inl.Less(0, 1), lessTestFn(0, 1))
}

func TestInnerSwap(t *testing.T) {
	inl := innerHeadedLazySortedList{data: make([]interface{}, 0), less: lessTestFn}
	inl.Add(1)
	inl.Add(2)
	inl.Swap(0, 1)
	require.Equal(t, inl.Get(0), 2)

	require.Equal(t, inl.Get(1), 1)
}

func TestInnerAdd(t *testing.T) {
	inl := innerHeadedLazySortedList{data: make([]interface{}, 0), less: lessTestFn}
	inl.Add(1)
	require.Equal(t, inl.Get(0), 1)

	require.Equal(t, inl.len, 1)

	inl = innerHeadedLazySortedList{data: make([]interface{}, 2), less: lessTestFn}
	inl.Add(1)
	require.Equal(t, inl.Get(0), 1)

	require.Equal(t, inl.len, 1)

	inl.Add(3)
	require.Equal(t, inl.Get(0), 1)

	require.Equal(t, inl.Get(1), 3)

	require.Equal(t, inl.len, 2)
}

func TestInnerGet(t *testing.T) {
	inl := innerHeadedLazySortedList{data: make([]interface{}, 0), less: lessTestFn}
	inl.Add(2)
	inl.Add(3)
	require.Equal(t, inl.Get(0), 2)

	require.Equal(t, inl.Get(1), 3)

	require.Panics(t, func() { inl.Get(-1) })

	require.Panics(t, func() { inl.Get(2) })
}

func TestLen(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, 1)
	require.Equal(t, hl.Len(), 0)

	hl.Add(2)
	require.Equal(t, hl.Len(), 1)

	hl.Add(3)
	require.Equal(t, hl.Len(), 2)
}

func TestGet(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, 1)
	hl.Add(2)
	hl.Add(3)
	require.Equal(t, hl.Get(0), 2)

	require.Equal(t, hl.Get(1), 3)

	require.Panics(t, func() { hl.Get(-1) })

	require.Panics(t, func() { hl.Get(2) })
}

func TestAdd(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, 1)
	hl.Add(1)
	require.Equal(t, hl.sorted, Sorted)

	hl.Add(2)
	require.Equal(t, hl.sorted, UnsortedTail)

	hl.Add(0)
	require.Equal(t, hl.sorted, UnsortedTail)

	hl = NewHeadedLazySortedList(2, lessTestFn, 2)
	hl.Add(1)
	hl.sorted = UnsortedAll
	hl.Add(2)
	require.Equal(t, hl.sorted, Sorted)

	hl.Add(1)
	require.Equal(t, hl.sorted, Sorted)

	hl = NewHeadedLazySortedList(4, lessTestFn, 3)
	hl.Add(1)
	hl.Add(2)
	require.Equal(t, hl.sorted, UnsortedAll)

	hl.Add(1)
	require.Equal(t, hl.sorted, UnsortedAll)
}

func TestSortAll(t *testing.T) {
	hl := NewHeadedLazySortedList(1, lessTestFn, 1)
	hl.Add(2)
	hl.SortAll()
	require.Equal(t, hl.sorted, Sorted)

	hl.Add(3)
	require.Equal(t, hl.sorted, UnsortedTail)

	hl.SortAll()
	require.Equal(t, hl.sorted, Sorted)
}

// TODO Tests
func TestLazyTailSort(t *testing.T) {
	s := NewHeadedLazySortedList(3, lessTestFn, 0)

	t.Logf("%+v", s)
}
