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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/artifacts"
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

func NewHeavyServer(cfgPath string, gensisCfgPath string) Server {
	return heavy.New(cfgPath, gensisCfgPath)
}

func NewVirtualServer(cfgPath string, codeRegistry map[string]insolar.ContractWrapper,
	codeRefRegistry map[insolar.Reference]string, codeDescriptors []artifacts.CodeDescriptor,
	prototypeDescriptors []artifacts.PrototypeDescriptor) Server {
	return virtual.New(cfgPath, codeRegistry,
		codeRefRegistry,
		codeDescriptors,
		prototypeDescriptors,
	)
}
