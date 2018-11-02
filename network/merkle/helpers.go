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

package merkle

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptohelpers/hash"
)
func hashConcat(args ...[]byte) []byte {
	var result []byte
	for _, arg := range args {
		argHash := hash.SHA3Bytes256(arg)
		result = append(result, argHash...)
	}

	return hash.SHA3Bytes256(result)
}
func pulseHash(pulse *core.Pulse) []byte {
	return hashConcat(pulse.PulseNumber.Bytes(), pulse.Entropy[:])
}

func nodeInfoHash(pulseHash, stateHash []byte) []byte {
	return hashConcat(pulseHash, stateHash)
}
