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

package unsafekit

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnwrapAsSlice(t *testing.T) {
	if BigEndian {
		t.SkipNow()
	}

	bytes := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	s := WrapBytes(bytes)
	st := MustMMapSliceType(reflect.TypeOf([]uint32(nil)), false)
	slice := UnwrapAsSlice(s, st).([]uint32)
	require.Equal(t, []uint32{0x03020100, 0x07060504, 0x0B0A0908}, slice)
	bytes[0] = 0xFF
	bytes[11] = 0xEE
	require.Equal(t, []uint32{0x030201FF, 0x07060504, 0xEE0A0908}, slice)

	runtime.KeepAlive(bytes)
}
