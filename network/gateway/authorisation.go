package gateway

import (
	"context"

	"github.com/insolar/insolar/insolar"
)

func NewAuthorisation(c *commons) *Authorisation {
	return &Authorisation{c: c}
}

type Authorisation struct {
	c *commons
}

func (g *Authorisation) Run() {
	panic("implement me")
}

func (g *Authorisation) GetState() insolar.NetworkState {
	panic("implement me")
}

func (g *Authorisation) OnPulse(context.Context, insolar.Pulse) error {
	panic("implement me")
}
