package network2

type NoNetwork struct {
	BaseGateway
}

func (g *NoNetwork) run() {
	vn := new(VoidNetwork)
	vn.BaseGateway = g.BaseGateway
	g.network.switchGateway(vn)
}

func (g *NoNetwork) onPulse() {
	panic("implement me")
}

func (g *NoNetwork) onMessage() {
	panic("implement me")
}
