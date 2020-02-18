// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/application/genesisrefs"
	"github.com/pkg/errors"
)

// initAPIInfoResponse creates application-specific data,
// that will be included in response from /admin-api/rpc#network.getInfo
func initAPIInfoResponse() (map[string]interface{}, error) {
	rootDomain := genesisrefs.ContractRootDomain
	if rootDomain.IsEmpty() {
		return nil, errors.New("rootDomain ref is nil")
	}

	rootMember := genesisrefs.ContractRootMember
	if rootMember.IsEmpty() {
		return nil, errors.New("rootMember ref is nil")
	}

	migrationDaemonMembers := genesisrefs.ContractMigrationDaemonMembers
	migrationDaemonMembersStrs := make([]string, 0)
	for _, r := range migrationDaemonMembers {
		if r.IsEmpty() {
			return nil, errors.New("migration daemon members refs are nil")
		}
		migrationDaemonMembersStrs = append(migrationDaemonMembersStrs, r.String())
	}

	migrationAdminMember := genesisrefs.ContractMigrationAdminMember
	if migrationAdminMember.IsEmpty() {
		return nil, errors.New("migration admin member ref is nil")
	}
	feeMember := genesisrefs.ContractFeeMember
	if feeMember.IsEmpty() {
		return nil, errors.New("feeMember ref is nil")
	}
	return map[string]interface{}{
		"rootDomain":             rootDomain.String(),
		"rootMember":             rootMember.String(),
		"migrationAdminMember":   migrationAdminMember.String(),
		"feeMember":              feeMember.String(),
		"migrationDaemonMembers": migrationDaemonMembersStrs,
	}, nil
}

func createAPIOptions() (api.Options, error) {
	apiInfoResponse, err := initAPIInfoResponse()
	if err != nil {
		return api.Options{}, err
	}
	adminContractMethods := map[string]bool{
		"migration.deactivateDaemon": true,
		"migration.activateDaemon":   true,
		"migration.checkDaemon":      true,
		"migration.addAddresses":     true,
		"migration.getAddressCount":  true,
		"deposit.migration":          true,
		"member.getBalance":          true,
	}
	contractMethods := map[string]bool{
		"member.create":          true,
		"member.get":             true,
		"member.transfer":        true,
		"member.migrationCreate": true,
		"deposit.transfer":       true,
	}
	proxyToRootMethods := []string{"member.create", "member.migrationCreate", "member.get"}

	return api.Options{
		AdminContractMethods: adminContractMethods,
		ContractMethods:      contractMethods,
		InfoResponse:         apiInfoResponse,
		RootReference:        genesisrefs.ContractRootMember,
		ProxyToRootMethods:   proxyToRootMethods,
	}, nil
}
