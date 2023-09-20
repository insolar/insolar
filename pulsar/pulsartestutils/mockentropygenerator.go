package pulsartestutils

import (
	"github.com/insolar/insolar/insolar"
)

// MockEntropy for pulsar's tests
var MockEntropy = [64]byte{1, 2, 3, 4, 5, 6, 7, 8}

// MockEntropyGenerator implements EntropyGenerator and is being used for tests
type MockEntropyGenerator struct {
}

// GenerateEntropy returns mocked entropy
func (MockEntropyGenerator) GenerateEntropy() insolar.Entropy {
	return MockEntropy
}
