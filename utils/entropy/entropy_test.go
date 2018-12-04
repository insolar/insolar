/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package entropy

import (
	"crypto/rand"
	"fmt"
	mrand "math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/platformpolicy"
)

func randslice(size int) []byte {
	b := make([]byte, size)
	rand.Read(b)
	return b
}

func TestSelectByEntropy(t *testing.T) {
	mrand.Seed(time.Now().UTC().UnixNano())

	scheme := platformpolicy.NewPlatformCryptographyScheme()
	entropy := randslice(64)

	valuescount := 10
	values := make([][]byte, 0, valuescount)
	seen := map[string]bool{}
	for i := 0; i < valuescount; {
		value := randslice(64)
		if seen[string(value)] {
			continue
		}
		values = append(values, value)
		seen[string(value)] = true
		i++
		// fmt.Printf(">> gen value: %b\n", value)
	}

	// fmt.Printf("entropy => %b\n", entropy)
	// fmt.Println(strings.Repeat("-", 77))

	count := 10
	// fmt.Printf("values => %v\n", values)
	result1, err := SelectByEntropy(scheme, entropy, values, count)
	require.NoError(t, err)
	assert.Equal(t, count, len(result1))
	// fmt.Printf("result1 => %v\n", result1)

	// fmt.Println(strings.Repeat("-", 77))

	// fmt.Printf("values => %v\n", values)
	result2, err := SelectByEntropy(scheme, entropy, values, count)
	assert.Equal(t, result1, result2)
	// fmt.Printf("result2 => %v\n", result2)

	seencount := map[string]int{}
	for _, val := range result2 {
		n, _ := seencount[string(val)]
		n++
		seencount[string(val)] = n
	}
	for k, v := range seencount {
		if v < 2 {
			delete(seencount, k)
		}
	}
	if !assert.Equal(t, 0, len(seencount), "values should not repeat") {
		fmt.Printf("repeats: %#v\n", seencount)
	}
}

// go test -v ./utils/entropy/ -bench=. -cpu=1 -benchmem -run=NONE > ./utils/entropy/benchresults/$(git rev-parse --short HEAD).txt
// benchcmp ./utils/entropy/benchresults/{d06b172a,cd73cf1e}.txt
func BenchmarkSelectByEntropy(b *testing.B) {
	benches := []struct {
		values int
		count  int
	}{
		{10, 1},
		{10, 5},
		{10, 10},
		{100, 1},
		{100, 50},
		{100, 100},
		{1000, 1},
		{1000, 500},
		{1000, 1000},
	}
	for _, bench := range benches {
		b.Run(
			fmt.Sprintf("%v_from_%v", bench.count, bench.values),
			func(b *testing.B) {
				benchSelectByEntropy(b, bench.values, bench.count, SelectByEntropy)
			})
	}
}

type entropyfunc func(
	scheme core.PlatformCryptographyScheme,
	entropy []byte,
	values [][]byte,
	count int,
) ([][]byte, error)

// compiler should avoid to optimize call of benched function
var results [][]byte

func benchSelectByEntropy(b *testing.B, valuescount int, count int, fn entropyfunc) {
	scheme := platformpolicy.NewPlatformCryptographyScheme()
	entropy := randslice(64)

	values := make([][]byte, 0, valuescount)
	for i := 0; i < valuescount; i++ {
		values = append(values, randslice(64))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results, _ = fn(scheme, entropy, values, count)
	}
}
