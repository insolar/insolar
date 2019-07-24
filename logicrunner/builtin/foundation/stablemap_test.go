//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package foundation

import (
	"crypto/md5"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStableMap_Len(t *testing.T) {
	sm := StableMap{}
	assert.Equal(t, 0, sm.Len())
	sm.Set("foo", 123)
	assert.Equal(t, 1, sm.Len())
	sm.Set("bar", 456)
	assert.Equal(t, 2, sm.Len())
	sm.Delete("foo")
	assert.Equal(t, 1, sm.Len())
	sm.Delete("bar")
	assert.Equal(t, 0, sm.Len())
}

func TestStableMap_Get(t *testing.T) {
	sm := StableMap{}

	val, ok := sm.Get("foo")
	assert.Nil(t, val)
	assert.False(t, ok)

	sm.Set("foo", 123)
	val, ok = sm.Get("foo")
	assert.Equal(t, 123, val)
	assert.True(t, ok)
	val, ok = sm.Get("bar")
	assert.Nil(t, val)
	assert.False(t, ok)

	sm.Set("bar", 456)
	val, ok = sm.Get("bar")
	assert.Equal(t, 456, val)
	assert.True(t, ok)

	sm.Delete("foo")
	val, ok = sm.Get("foo")
	assert.Nil(t, val)
	assert.False(t, ok)
}

func TestStableMap_Set(t *testing.T) {
	sm := StableMap{}

	sm.Set("foo", 123)
	val, ok := sm.Get("foo")
	assert.Equal(t, 123, val)
	assert.True(t, ok)

	sm.Set("bar", 456)
	val, ok = sm.Get("bar")
	assert.Equal(t, 456, val)
	assert.True(t, ok)

	sm.Set("bar", "baz")
	val, ok = sm.Get("bar")
	assert.Equal(t, "baz", val)
	assert.True(t, ok)
}

func TestStableMap_Delete(t *testing.T) {
	sm := StableMap{}

	sm.Set("foo", 123)
	val, ok := sm.Get("foo")
	assert.Equal(t, 123, val)
	assert.True(t, ok)

	sm.Delete("foo")
	val, ok = sm.Get("foo")
	assert.Nil(t, val)
	assert.False(t, ok)

	sm.Delete("bar")
	val, ok = sm.Get("bar")
	assert.Nil(t, val)
	assert.False(t, ok)
}

func TestStableMap_Keys(t *testing.T) {
	sm := StableMap{}

	assert.Empty(t, sm.GetKeys())

	sm.Set("foo", 123)
	sm.Set("bar", 456)
	sm.Set("baz", 789)
	assert.Equal(t, []interface{}{"foo", "bar", "baz"}, sm.GetKeys())

	sm.Delete("bar")
	sm.Set("bar", 456)
	sm.Set(123, "foobar")
	assert.Equal(t, []interface{}{"foo", "baz", "bar", 123}, sm.GetKeys())
}

func TestStableMap_Values(t *testing.T) {
	sm := StableMap{}

	assert.Empty(t, sm.GetValues())

	sm.Set("foo", 123)
	sm.Set("bar", 456)
	sm.Set("baz", 789)
	assert.Equal(t, []interface{}{123, 456, 789}, sm.GetValues())

	sm.Delete("bar")
	sm.Set("bar", 456)
	sm.Set(123, "foobar")
	assert.Equal(t, []interface{}{123, 789, 456, "foobar"}, sm.GetValues())
}

func TestStableMap_Pairs(t *testing.T) {
	sm := StableMap{}

	assert.Empty(t, sm.Pairs())

	sm.Set("foo", 123)
	sm.Set("bar", 456)
	sm.Set("baz", 789)
	assert.Equal(
		t,
		[][2]interface{}{
			{"foo", 123},
			{"bar", 456},
			{"baz", 789},
		},
		sm.Pairs(),
	)

	sm.Delete("bar")
	sm.Set("bar", 456)
	sm.Set(123, "foobar")
	assert.Equal(
		t,
		[][2]interface{}{
			{"foo", 123},
			{"baz", 789},
			{"bar", 456},
			{123, "foobar"},
		},
		sm.Pairs(),
	)
}

type TestStable struct {
	A string
	B StableMap
	C int
}

func TestStableMap_serialization(t *testing.T) {
	s := TestStable{}

	s.A = "foobar"
	s.B.Set("foo", "123")
	s.B.Set("bar", "456")
	s.B.Set("baz", "789")
	s.B.Set("123", "foo")
	s.B.Set("456", "bar")
	s.B.Set("789", "baz")
	s.C = 123

	buf := insolar.MustSerialize(&s)

	s2 := TestStable{}
	insolar.MustDeserialize(buf, &s2)

	assert.Equal(t, s, s2)
}

func TestStableMap_is_deterministic(t *testing.T) {
	s := TestStable{}

	s.A = "foobar"
	s.B.Set("foo", 123)
	s.B.Set("bar", 456)
	s.B.Set("baz", 789)
	s.B.Set(123, "foo")
	s.B.Set(456, "bar")
	s.B.Set(789, "baz")
	s.C = 123

	buf := insolar.MustSerialize(s)
	sum := md5.Sum(buf)

	for i := 0; i < 10000; i++ {
		buf = insolar.MustSerialize(s)
		require.Equal(t, sum, md5.Sum(buf))
	}
}

type TestMap struct {
	A string
	B map[interface{}]interface{}
	C int
}

func TestStableMap_common_map_is_not_deterministic(t *testing.T) {
	hashmap := make(map[[16]byte]uint)
	s := TestMap{}

	s.A = "foobar"
	s.B = make(map[interface{}]interface{})
	s.B["foo"] = 123
	s.B["bar"] = 456
	s.B["baz"] = 789
	s.B[123] = "foo"
	s.B[456] = "bar"
	s.B[789] = "baz"
	s.C = 123

	var buf []byte
	for i := 0; i < 10000; i++ {
		buf = insolar.MustSerialize(s)
		sum := md5.Sum(buf)
		hashmap[sum]++
	}

	assert.Greater(t, len(hashmap), 1)
}
