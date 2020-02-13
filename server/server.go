// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
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
func NewLightServer(cfgPath string, apiInfoResponse map[string]interface{}) Server {
	return light.New(cfgPath, apiInfoResponse)
}

// NewHeavyServer creates instance of Server for node with Heavy role
func NewHeavyServer(
	cfgPath string,
	genesisCfgPath string,
	genesisOptions genesis.Options,
	genesisOnly bool,
	apiInfoResponse map[string]interface{},
) Server {
	return heavy.New(cfgPath, genesisCfgPath, genesisOptions, genesisOnly, apiInfoResponse)
}

// NewVirtualServer creates instance of Server for node with Virtual role
func NewVirtualServer(cfgPath string, builtinContracts builtin.BuiltinContracts, apiInfoResponse map[string]interface{}) Server {
	return virtual.New(cfgPath, builtinContracts, apiInfoResponse)
}
