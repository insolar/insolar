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

package hash

import (
	"hash"

	"github.com/insolar/insolar/platformpolicy"
)

var policy = platformpolicy.NewPlatformPolicy()

// NewIDHash returns hash used for records ID generation.
func NewIDHash() hash.Hash {
	return hashAdapter.Hash224bits()
}

// IDHashBytes generates hash for record ID from byte slice.
func IDHashBytes(b []byte) []byte {
	return hashAdapter.Hash224bits().Hash(b)
}

// SHA3Bytes256 generates SHA3-256 hash for byte slice.
func SHA3Bytes256(b []byte) []byte {
	return hashAdapter.Hash256bits().Hash(b)
}
