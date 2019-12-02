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

package dbsv1

import (
	"fmt"
	"hash"
)

func calcCrc32(hasher hash.Hash32, b []byte) hash.Hash32 {
	switch n, err := hasher.Write(b); {
	case err != nil:
		panic(err)
	case n != len(b):
		panic(fmt.Errorf("internal error, crc calc failed: written=%d expected=%d", n, len(b)))
	}
	return hasher
}

func addCrc32(hash hash.Hash32, x uint32) {
	// byte order is according to crc32.appendUint32
	if n, err := hash.Write([]byte{byte(x >> 24), byte(x >> 16), byte(x >> 8), byte(x)}); err != nil || n != 4 {
		panic(fmt.Errorf("crc calc failure: %d, %v", n, err))
	}
}
