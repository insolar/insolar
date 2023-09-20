package sequence

import (
	"sync/atomic"
)

type Sequence uint64

type Generator interface {
	Generate() Sequence
}

type generator struct {
	sequence *uint64
}

func NewGenerator() Generator {
	return &generator{
		sequence: new(uint64),
	}
}

func (g *generator) Generate() Sequence {
	return Sequence(atomic.AddUint64(g.sequence, 1))
}
