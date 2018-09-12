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

package pulsar

import "math/rand"

type EntropyGenerator interface {
	GenerateEntropy() [8]byte
}

type StandardEntropyGenerator struct {
}

func (generator *StandardEntropyGenerator) GenerateEntropy() [8]byte {
	entropy := make([]byte, 8)
	rand.Read(entropy)
	var result [8]byte
	copy(result[:], entropy[:8])
	return result
}
