// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
