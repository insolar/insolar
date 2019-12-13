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
	"hash/fnv"

	"github.com/insolar/insolar/ledger-v2/unsafekit"
	"github.com/insolar/insolar/longbits"
)

func hash32(v longbits.ByteString, seed uint32) uint32 {
	// FNV-1a has a better avalanche property vs FNV-1
	// use of 64 bits improves distribution
	h := fnv.New64a()
	if seed != 0 {
		_, _ = h.Write([]byte{byte(seed), byte(seed >> 8), byte(seed >> 16), byte(seed >> 24)})
	}
	unsafekit.Hash(v, h)
	sum := h.Sum64()
	return uint32(sum) ^ uint32(sum>>32)
}
