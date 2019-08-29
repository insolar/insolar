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

package jetcoordinator

import (
	"crypto/rand"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/utils/entropy"
)

// In reality compares no sort vs with sort + in/out conversions of array of empty interfaces
// This benchamark results would be suitable for analyzing how much we lost on input/output
// conversion only after sorting removal
//
// prepare benchmarks results:
// go test -v ./ledger/jetcoordinator/ -bench=SelectByEntropy -cpu=1 -benchmem -run=NONE > wrapped.txt
// SelectByEntropyBench=orig go test -v ./ledger/jetcoordinator/ -bench=SelectByEntropy -cpu=1 -benchmem -run=NONE > orig.txt
//
// measure overhead:
// benchcmp orig.txt wrapped.txt
//
func BenchmarkSelectByEntropy(b *testing.B) {
	benchtype := strings.ToLower(os.Getenv("SelectByEntropyBench"))
	switch benchtype {
	case "orig", "wrapped":
		// all ok
	case "":
		benchtype = "wrapped"
	default:
		panic(fmt.Sprintf("Unknown benchtype %v", benchtype))
	}

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
	fmt.Printf("# Bench: %v\n", benchtype)
	for _, bench := range benches {
		b.Run(
			fmt.Sprintf("%v_from_%v", bench.count, bench.values),
			func(b *testing.B) {
				if benchtype == "orig" {
					benchSelectByEntropy(b, bench.values, bench.count)
					return
				}
				benchSelectByEntropyWrapped(b, bench.values, bench.count)
			})
	}
}

// compiler should avoid to optimize call of benched function
var resultsI []interface{}
var resultsB [][]byte

func benchSelectByEntropy(b *testing.B, valuescount int, count int) {
	scheme := platformpolicy.NewPlatformCryptographyScheme()
	entropybytes := randslice(64)

	values := make([]interface{}, 0, valuescount)
	for i := 0; i < valuescount; i++ {
		values = append(values, interface{}(randslice(64)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// sort.SliceStable(valuesB, )
		// bytes.Compare(a, b) < 0
		resultsI, _ = entropy.SelectByEntropy(scheme, entropybytes, values, count)
	}
}

// compiler should avoid to optimize call of benched function
var refresults []insolar.Reference

func benchSelectByEntropyWrapped(b *testing.B, valuescount int, count int) {
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	var e insolar.Entropy
	copy(e[:], randslice(64))

	values := make([]insolar.Node, 0, valuescount)
	for i := 0; i < valuescount; i++ {
		values = append(values, insolar.Node{ID: gen.Reference()})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		refresults, _ = getRefs(scheme, e[:], values, count)
	}
}

func randslice(size int) []byte {
	b := make([]byte, size)
	rand.Read(b)
	return b
}
