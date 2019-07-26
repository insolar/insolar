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

	buf := insolar.MustSerialize(s)
	sum := md5.Sum(buf)

	for i := 0; i < 10000; i++ {
		buf = insolar.MustSerialize(s)
		require.Equal(t, sum, md5.Sum(buf))
	}
}

type TestMap struct {
	A string
	B map[string]string
	C int
}

func TestStableMap_common_map_is_not_deterministic(t *testing.T) {
	hashmap := make(map[[16]byte]uint)
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

	var buf []byte
	for i := 0; i < 10000; i++ {
		buf = insolar.MustSerialize(s)
		sum := md5.Sum(buf)
		hashmap[sum]++
	}

	assert.Greater(t, len(hashmap), 1)
}
