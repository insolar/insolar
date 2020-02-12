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
	"bytes"
	"testing"
)

func TestStandardEntropyGenerator_GenerateEntropy_EntropySize(t *testing.T) {
	generator := &StandardEntropyGenerator{}

	first := generator.GenerateEntropy()

	if len(first) != 64 {
		t.Errorf("Length of entropy should be equal to 64, got %v", len(first))
	}
}

func TestStandardEntropyGenerator_GenerateEntropy_EntropyShouldBeUnique(t *testing.T) {
	generator := &StandardEntropyGenerator{}
	first := generator.GenerateEntropy()
	second := generator.GenerateEntropy()

	result := bytes.Equal(first[:], second[:])

	if result {
		t.Errorf("Entropies shouldn't be the same, got - %v, wanted - %v", first, second)
	}
}
