package gateway

import (
	"context"

	"github.com/insolar/insolar/insolar"
)

func NewVoid(c *commons) *Void {
	return &Void{c: c}
}

// VoidNetwork void network state
type Void struct {
	c *commons
}

func (g *Void) Run() {
	panic("implement me")
}

func (g *Void) GetState() insolar.NetworkState {
	return insolar.VoidNetworkState
}

func (g *Void) OnPulse(context.Context, insolar.Pulse) error {
	panic("implement me")
}
