// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
