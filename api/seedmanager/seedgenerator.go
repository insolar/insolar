package seedmanager

import (
	"github.com/insolar/x-crypto/rand"

	"github.com/pkg/errors"
)

// SeedSize is size of seed
const SeedSize uint = 32

// Seed is a type of seed
type Seed = [SeedSize]byte

// SeedGenerator holds logic with seed generation
type SeedGenerator struct {
}

// Next returns next random seed
func (sg *SeedGenerator) Next() (*Seed, error) {
	seed := Seed{}
	_, err := rand.Read(seed[:])
	if err != nil {
		return nil, errors.Wrap(err, "failed to get next seed")
	}

	return &seed, nil
}
