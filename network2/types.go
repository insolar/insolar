package network2

type Gateway interface {
	run()
	onPulse()
	onMessage()
}

type BaseGateway struct {
	network Network
}
