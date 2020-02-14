// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package server

import (
	"github.com/insolar/insolar/server/internal/heavy"
	"github.com/insolar/insolar/server/internal/light"
	"github.com/insolar/insolar/server/internal/virtual"
)

type Server interface {
	Serve()
}

func NewLightServer(cfgPath string) Server {
	return light.New(cfgPath)
}

func NewHeavyServer(cfgPath string, gensisCfgPath string, genesisOnly bool) Server {
	return heavy.New(cfgPath, gensisCfgPath, genesisOnly)
}

func NewVirtualServer(cfgPath string) Server {
	return virtual.New(cfgPath)
}
