// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package contracts

import (
	"fmt"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/applicationbase/genesis"

	"github.com/insolar/insolar/application/builtin/contract/member"
)

func RootDomain() genesis.ContractState {

	return genesis.ContractState{
		Name:       application.GenesisNameRootDomain,
		Prototype:  application.GenesisNameRootDomain,
		ParentName: "",
	}
}

func GetMemberGenesisContractState(publicKey string, name string, parent string) genesis.ContractState {
	m, err := member.New(publicKey)
	if err != nil {
		panic(fmt.Sprintf("'%s' member constructor failed", name))
	}

	return genesis.ContractState{
		Name:       name,
		Prototype:  application.GenesisNameMember,
		ParentName: parent,
		Memory:     genesis.MustGenMemory(m),
	}
}
