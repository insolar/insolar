// Copyright 2020 Insolar Network Ltd.
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

package entropygenerator

import (
	"crypto/rand"

	"github.com/insolar/insolar/insolar"
)

// EntropyGenerator is the base interface for generation of entropy for pulses
type EntropyGenerator interface {
	GenerateEntropy() insolar.Entropy
}

// StandardEntropyGenerator is the base impl of EntropyGenerator with using of crypto/rand
type StandardEntropyGenerator struct {
}

// GenerateEntropy generate entropy with using of EntropyGenerator
func (generator *StandardEntropyGenerator) GenerateEntropy() insolar.Entropy {
	entropy := make([]byte, insolar.EntropySize)
	_, err := rand.Read(entropy)
	if err != nil {
		panic(err)
	}
	var result insolar.Entropy
	copy(result[:], entropy[:insolar.EntropySize])
	return result
}
