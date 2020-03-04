// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package server

import (
	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/applicationbase/genesis"
	"github.com/insolar/insolar/logicrunner/builtin"
	"github.com/insolar/insolar/server/internal/heavy"
	"github.com/insolar/insolar/server/internal/light"
	"github.com/insolar/insolar/server/internal/virtual"
)

type Server interface {
	Serve()
}

// NewLightServer creates instance of Server for node with Light role
func NewLightServer(cfgPath string, apiOptions api.Options) Server {
	return light.New(cfgPath, apiOptions)
}

// NewHeavyServer creates instance of Server for node with Heavy role
func NewHeavyServer(
	cfgPath string,
	genesisCfgPath string,
	genesisOptions genesis.Options,
	genesisOnly bool,
	apiOptions api.Options,
) Server {
	return heavy.New(cfgPath, genesisCfgPath, genesisOptions, genesisOnly, apiOptions)
}

// NewVirtualServer creates instance of Server for node with Virtual role
func NewVirtualServer(cfgPath string, builtinContracts builtin.BuiltinContracts, apiOptions api.Options,
) Server {
	return virtual.New(cfgPath, builtinContracts, apiOptions)
}
