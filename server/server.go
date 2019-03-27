package server

import (
	"github.com/insolar/insolar/server/internal/genesis"
	"github.com/insolar/insolar/server/internal/light"
	"github.com/insolar/insolar/server/internal/virtual"
)

type Server interface {
	Serve()
}

func NewLightServer(cfgPath string, trace bool) Server {
	return light.New(cfgPath, trace)
}

func NewVirtualServer(cfgPath string, trace bool) Server {
	return virtual.New(cfgPath, trace)
}

func NewGenesisServer(cfgPath string, trace bool, genesisConfigPath, genesisKeyOut string) Server {
	return genesis.New(cfgPath, trace, genesisConfigPath, genesisKeyOut)
}
