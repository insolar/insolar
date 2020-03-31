// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/application/genesisrefs/contracts"
	"github.com/insolar/insolar/applicationbase/genesis"
	"github.com/insolar/insolar/log"
)

func initStates(genesisConfigPath string) ([]genesis.ContractState, error) {
	b, err := ioutil.ReadFile(genesisConfigPath)
	if err != nil {
		log.Fatalf("failed to load genesis configuration from file: %v", genesisConfigPath)
	}
	var config struct {
		ContractsConfig application.GenesisContractsConfig
	}
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatalf("failed to parse genesis configuration from file: %v", genesisConfigPath)
	}

	contractsConfig := config.ContractsConfig
	// Hint: order matters, because of dependency contracts on each other.
	states := []genesis.ContractState{
		contracts.RootDomain(),
		contracts.GetMemberGenesisContractState(contractsConfig.RootPublicKey, application.GenesisNameRootMember, application.GenesisNameRootDomain),
	}

	return states, nil
}
