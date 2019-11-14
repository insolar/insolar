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
	"fmt"
	"io"
)

const SplitMedian = 7 // makes 0 vs 1 ratio like [0..6] vs [7..15]
// this enables left branches of jets to be ~23% less loaded

type PrefixCalc struct {
	OverlapOfs uint8
}

func (p PrefixCalc) FromSlice(prefixLen int, data []byte) Prefix {
	switch {
	case prefixLen < 0 || prefixLen > 32:
		panic("illegal value")
	case prefixLen == 0:
		return 0
	}

	return p.fromSlice(prefixLen, data)
}

func (p PrefixCalc) FromReader(prefixLen int, data io.Reader) (Prefix, error) {
	switch {
	case prefixLen < 0 || prefixLen > 32:
		panic("illegal value")
	case data == nil:
		panic("illegal value")
	case prefixLen == 0:
		return 0, nil
	}

	dataBuf := make([]byte, int(p.OverlapOfs)+(prefixLen+1)>>1)
	switch n, err := data.Read(dataBuf); {
	case err != nil:
		return 0, err
	case n != len(dataBuf):
		return 0, fmt.Errorf("insufficient data length")
	}

	return p.fromSlice(prefixLen, dataBuf), nil
}

func (p PrefixCalc) fromSlice(prefixLen int, data []byte) Prefix {
	result := Prefix(0)
	bit := Prefix(1)

	for i, d := range data {
		if p.OverlapOfs > 0 {
			d ^= data[i+int(p.OverlapOfs)]
		}

		if d&0xF >= SplitMedian {
			result |= bit
		}
		if prefixLen == 1 {
			return result
		}
		bit <<= 1

		if (d >> 4) >= SplitMedian {
			result |= bit
		}
		prefixLen -= 2
		if prefixLen == 0 {
			return result
		}
		bit <<= 1
	}

	panic(fmt.Errorf("insufficient data length"))
}
