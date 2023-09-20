package server

import (
	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/applicationbase/genesis"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/logicrunner/builtin"
	"github.com/insolar/insolar/server/internal/heavy"
	"github.com/insolar/insolar/server/internal/light"
	"github.com/insolar/insolar/server/internal/virtual"
)

type Server interface {
	Serve()
}

// NewLightServer creates instance of Server for node with Light role
func NewLightServer(cfgHolder *configuration.LightHolder, apiOptions api.Options) Server {
	return light.New(cfgHolder, apiOptions)
}

// NewHeavyServer creates instance of Server for node with Heavy role
func NewHeavyServer(
	cfgHolder configuration.ConfigHolder,
	genesisCfgPath string,
	genesisOptions genesis.Options,
	genesisOnly bool,
	apiOptions api.Options,
	allowedContractVersion int64,
) Server {
	return heavy.New(cfgHolder, genesisCfgPath, genesisOptions, genesisOnly, apiOptions, allowedContractVersion)
}

// NewVirtualServer creates instance of Server for node with Virtual role
func NewVirtualServer(cfgHolder *configuration.VirtualHolder, builtinContracts builtin.BuiltinContracts, apiOptions api.Options) Server {
	return virtual.New(cfgHolder, builtinContracts, apiOptions)
}
