package gateway

import (
	"context"

	"github.com/insolar/insolar/insolar"
)

func NewJetless(c *commons) *Jetless {
	return &Jetless{c: c}
}

type Jetless struct {
	c *commons
}

func (g *Jetless) Run() {
	panic("implement me")
}

func (g *Jetless) GetState() insolar.NetworkState {
	panic("implement me")
}

func (g *Jetless) OnPulse(context.Context, insolar.Pulse) error {
	panic("implement me")
}
