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

package refmap

import (
	"fmt"
	"math"
	"testing"

	"github.com/insolar/insolar/reference"
)

func TestBucketing(t *testing.T) {
	m := NewRefLocatorMap()
	//m.keys.SetHashSeed(0)

	const keyCount = int(1e6)
	for i := keyCount; i > 0; i-- {
		refLocal := makeLocal(i)
		for j := 1; j > 0; j-- {
			refBase := makeLocal(i + j*1e6)
			ref := reference.NewNoCopy(&refLocal, &refBase)
			m.Put(ref, ValueLocator(i))
		}
	}

	wb := m.FillLocatorBuckets(WriteBucketerConfig{
		ExpectedPerBucket: 128,
		UsePrimes:         false,
	})

	fmt.Println("Params", wb.AdjustedBucketSize(), wb.BucketCount())

	if len(wb.overflowEntries) > 0 {
		minOverflow := math.MaxInt32
		maxOverflow := 0
		for _, v := range wb.overflowEntries {
			n := len(v)
			if n < minOverflow {
				minOverflow = n
			}
			if n > maxOverflow {
				maxOverflow = n
			}
		}

		fmt.Println("Overflown", len(wb.overflowEntries), minOverflow, maxOverflow)
	}

	_ = writeLocatorBuckets(&wb)
}
