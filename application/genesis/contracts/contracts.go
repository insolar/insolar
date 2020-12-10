// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package contracts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/insolar/insolar/application/builtin/contract/member"
	appgenesis "github.com/insolar/insolar/application/genesis"
	"github.com/insolar/insolar/applicationbase/genesis"
	"github.com/insolar/insolar/log"
)

func InitStates(genesisConfigPath string) ([]genesis.ContractState, error) {
	b, err := ioutil.ReadFile(genesisConfigPath)
	if err != nil {
		log.Fatalf("failed to load genesis configuration from file: %v", genesisConfigPath)
	}
	var config struct {
		ContractsConfig appgenesis.GenesisContractsConfig
	}
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatalf("failed to parse genesis configuration from file: %v", genesisConfigPath)
	}

	contractsConfig := config.ContractsConfig
	// Hint: order matters, because of dependency contracts on each other.
	states := []genesis.ContractState{
		rootDomain(),
		getMemberGenesisContractState(contractsConfig.RootPublicKey, appgenesis.GenesisNameRootMember, appgenesis.GenesisNameRootDomain),
	}

	return states, nil
}

func rootDomain() genesis.ContractState {
	return genesis.ContractState{
		Name:       appgenesis.GenesisNameRootDomain,
		Prototype:  appgenesis.GenesisNameRootDomain,
		ParentName: "",
	}
}

func getMemberGenesisContractState(publicKey string, name string, parent string) genesis.ContractState {
	m, err := member.New(publicKey)
	if err != nil {
		panic(fmt.Sprintf("'%s' member constructor failed", name))
	}

	return genesis.ContractState{
		Name:       name,
		Prototype:  appgenesis.GenesisNameMember,
		ParentName: parent,
		Memory:     genesis.MustGenMemory(m),
	}
}
