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

type PrefixCalc struct {
}

func (p PrefixCalc) FromSlice(prefixLen int, overlapOfs int, data []byte) Prefix {
	switch {
	case prefixLen < 0 || prefixLen > 32:
		panic("illegal value")
	case overlapOfs < 0:
		panic("illegal value")
	case prefixLen == 0:
		return 0
	}

	return p.fromSlice(prefixLen, overlapOfs, data)
}

func (p PrefixCalc) FromReader(prefixLen int, overlapOfs int, data io.Reader) (Prefix, error) {
	switch {
	case prefixLen < 0 || prefixLen > 32:
		panic("illegal value")
	case overlapOfs < 0:
		panic("illegal value")
	case data == nil:
		panic("illegal value")
	case prefixLen == 0:
		return 0, nil
	}

	dataBuf := make([]byte, overlapOfs+(prefixLen+1)>>1)
	switch n, err := data.Read(dataBuf); {
	case err != nil:
		return 0, err
	case n != len(dataBuf):
		return 0, fmt.Errorf("insufficient data length")
	}

	return p.fromSlice(prefixLen, overlapOfs, dataBuf), nil
}

func (p PrefixCalc) fromSlice(prefixLen int, overlapOfs int, data []byte) Prefix {
	result := Prefix(0)
	bit := Prefix(1)

	for i, d := range data {
		if overlapOfs > 0 {
			d ^= data[i+overlapOfs]
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
