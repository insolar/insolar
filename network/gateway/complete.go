package gateway

import (
	"context"

	"github.com/insolar/insolar/insolar"
)

func NewComple(c *commons) *Complete {
	return &Complete{c: c}
}

type Complete struct {
	c *commons
}

func (g *Complete) Run() {
	panic("implement me")
}

func (g *Complete) GetState() insolar.NetworkState {
	panic("implement me")
}

func (g *Complete) OnPulse(context.Context, insolar.Pulse) error {
	panic("implement me")
}
