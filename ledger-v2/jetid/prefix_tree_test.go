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

package jetid

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrefixTree_Print(t *testing.T) {
	pt := PrefixTree{}
	pt.PrintTable()

	pt.Split(0, 0)
	pt.PrintTable()

	pt.Merge(0, 1)
	pt.PrintTable()

	pt.Split(0, 0)
	pt.PrintTable()

	pt.Split(0, 1)
	pt.PrintTable()

	pt.Split(0, 2)
	pt.PrintTable()

	pt.Split(1, 1)
	pt.PrintTable()

	pt.Split(3, 2)
	pt.PrintTable()

	pt.Merge(0, 3)
	pt.PrintTable()

	pt.Merge(0, 2)
	pt.PrintTable()

	pt.Merge(3, 3)
	pt.PrintTable()

	pt.Merge(1, 2)
	pt.PrintTable()

	pt.Merge(0, 1)
	pt.PrintTable()
}

func TestPrefixTree_Propagate(t *testing.T) {
	pt := NewPrefixTree(true)
	//pt.Split(0, 0)
	//pt.Split(0, 1)
	//pt.Split(0, 2)
	//pt.Split(0, 3)
	//pt.Split(0, 4)
	//pt.Split(0, 5)
	//pt.Split(0, 6)
	//pt.Split(0, 7)
	//pt.Split(0, 8)
	//pt.Split(0, 9)
	//pt.Split(0, 10)
	//pt.Split(0, 11)
	//pt.Split(0, 12)
	//pt.Split(0, 13)
	//pt.Split(0, 14)
	//pt.Split(0, 15)
	//pt.PrintTable()
	//pt.Merge(0, 16)
	pt.PrintTable()
}

func TestPrefixTree_SplitMax0(t *testing.T) {
	pt := PrefixTree{}
	pt.Split(0, 0)
	pt.Split(0, 1)
	pt.Split(0, 2)
	pt.Split(0, 3)
	pt.Split(0, 4)
	pt.Split(0, 5)
	pt.Split(0, 6)
	pt.Split(0, 7)
	pt.Split(0, 8)
	pt.Split(0, 9)
	pt.Split(0, 10)
	pt.Split(0, 11)
	pt.Split(0, 12)
	pt.Split(0, 13)
	pt.Split(0, 14)
	pt.Split(0, 15)
	pt.PrintTable()
	pt.Merge(0, 16)
	pt.PrintTable()
}

func TestPrefixTree_SplitMax1(t *testing.T) {
	pt := PrefixTree{}
	pt.Split(0, 0)
	pt.Split(1, 1)
	pt.Split(3, 2)
	pt.Split(7, 3)
	pt.Split(15, 4)
	pt.Split(31, 5)
	pt.Split(63, 6)
	pt.Split(127, 7)
	pt.Split(255, 8)
	pt.Split(511, 9)
	pt.Split(1023, 10)
	pt.Split(2047, 11)
	pt.Split(4095, 12)
	pt.Split(8191, 13)
	pt.Split(16383, 14)
	pt.Split(32767, 15)
	pt.PrintTable()
	pt.Merge(32767, 16)
	pt.PrintTable()
}

func TestPrefixTree_Serialize(t *testing.T) {

	pt := PrefixTree{}
	pt.Init() // to make it properly comparable

	pt.Split(0, 0)
	//
	pt.Split(0, 1)
	pt.Split(0, 2)
	pt.Split(0, 3)
	pt.Split(0, 4)
	pt.Split(0, 5)
	pt.Split(0, 6)
	pt.Split(0, 7)
	pt.Split(0, 8)
	pt.Split(0, 9)
	pt.Split(0, 10)
	pt.Split(0, 11)
	pt.Split(0, 12)
	pt.Split(0, 13)
	pt.Split(0, 14)
	pt.Split(0, 15)
	//
	pt.Split(1, 1)
	pt.Split(3, 2)
	pt.Split(7, 3)
	pt.Split(15, 4)
	pt.Split(31, 5)
	pt.Split(63, 6)
	pt.Split(127, 7)
	pt.Split(255, 8)
	pt.Split(511, 9)
	pt.Split(1023, 10)
	pt.Split(2047, 11)
	pt.Split(4095, 12)
	pt.Split(8191, 13)
	pt.Split(16383, 14)
	pt.Split(32767, 15)
	pt.PrintTable()

	buf := bytes.Buffer{}
	require.NoError(t, pt.CompactSerialize(&buf))
	bufCopy := buf.Bytes() // will be ok as we don't write into it further

	fmt.Printf("Compact: %5d bytes\n", len(bufCopy))
	fmt.Println(hex.Dump(bufCopy))

	pt2 := PrefixTree{}
	require.NoError(t, pt2.CompactDeserialize(&buf))

	buf2 := bytes.Buffer{}
	require.NoError(t, pt.CompactSerialize(&buf2))
	if !bytes.Equal(bufCopy, buf2.Bytes()) {
		pt2.PrintTable()
	}
	require.Equal(t, bufCopy, buf2.Bytes())
	require.Equal(t, pt, pt2)
}
