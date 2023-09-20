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
