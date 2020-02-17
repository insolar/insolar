// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package foundation

import (
	"bytes"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/assert"
)

type TestStable struct {
	A string
	B StableMap
	C int
}

func TestStableMap_serialization(t *testing.T) {
	s := TestStable{}

	s.A = "foobar"
	s.B = make(StableMap)
	s.B["foo"] = "123"
	s.B["bar"] = "456"
	s.B["baz"] = "789"
	s.B["123"] = "foo"
	s.B["456"] = "bar"
	s.B["789"] = "baz"
	s.C = 123

	buf := insolar.MustSerialize(&s)

	s2 := TestStable{}
	insolar.MustDeserialize(buf, &s2)

	assert.Equal(t, s, s2)
}

func TestStableMap_is_deterministic(t *testing.T) {
	s1 := TestStable{}
	s1.A = "foobar"
	s1.B = make(StableMap)
	s1.B["foo"] = "123"
	s1.B["bar"] = "456"
	s1.B["baz"] = "789"
	s1.B["123"] = "foo"
	s1.B["456"] = "bar"
	s1.B["789"] = "baz"
	s1.C = 123

	s2 := TestStable{}
	s2.C = 123
	s2.A = "foobar"
	s2.B = make(StableMap)
	s2.B["789"] = "baz"
	s2.B["456"] = "bar"
	s2.B["123"] = "foo"
	s2.B["baz"] = "789"
	s2.B["bar"] = "456"
	s2.B["foo"] = "123"

	buf1 := insolar.MustSerialize(s1)

	buf2 := insolar.MustSerialize(s2)

	assert.Equal(t, buf1, buf2)
}

type TestMap struct {
	A string
	B map[string]string
	C int
}

func TestStableMap_common_map_is_not_deterministic(t *testing.T) {
	s := TestMap{}

	s.A = "foobar"
	s.B = make(map[string]string)
	s.B["foo"] = "123"
	s.B["bar"] = "456"
	s.B["baz"] = "789"
	s.B["123"] = "foo"
	s.B["456"] = "bar"
	s.B["789"] = "baz"
	s.C = 123

	firstBuf := insolar.MustSerialize(s)
	var buf []byte
	for i := 0; i < 10000; i++ {
		buf = insolar.MustSerialize(s)
		if !bytes.Equal(firstBuf, buf) {
			break
		}
	}
	assert.NotEqual(t, firstBuf, buf)
}
