// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"github.com/pkg/errors"

	"github.com/insolar/insolar/application/appfoundation"
	appbuiltin "github.com/insolar/insolar/application/builtin"
	"github.com/insolar/insolar/application/genesis"
	"github.com/insolar/insolar/application/genesis/contracts"
	"github.com/insolar/insolar/insolard"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/builtin"
)

func main() {
	parentDomain := genesis.GenesisNameRootDomain
	scVersion := int64(appfoundation.AllowedVersionSmartContract)
	builtinContracts := builtin.BuiltinContracts{
		CodeRegistry:         appbuiltin.InitializeContractMethods(),
		CodeRefRegistry:      appbuiltin.InitializeCodeRefs(),
		CodeDescriptors:      appbuiltin.InitializeCodeDescriptors(),
		PrototypeDescriptors: appbuiltin.InitializePrototypeDescriptors(),
	}
	apiOptions, err := genesis.InitAPIOptions()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to get API info response"))
	}
	initStates := contracts.InitStates

	insolard.RunInsolarNode(parentDomain, scVersion, apiOptions, initStates, builtinContracts)
}
