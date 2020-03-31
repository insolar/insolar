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

	return map[string]interface{}{
		"rootDomain": rootDomain.String(),
		"rootMember": rootMember.String(),
	}, nil
}

// initAPIOptions creates options object, that contains application-specific settings for api component.
func initAPIOptions() (api.Options, error) {
	apiInfoResponse, err := initAPIInfoResponse()
	if err != nil {
		return api.Options{}, err
	}
	adminContractMethods := map[string]bool{}
	contractMethods := map[string]bool{
		"member.create": true,
	}
	proxyToRootMethods := []string{"member.create"}

	return api.Options{
		AdminContractMethods: adminContractMethods,
		ContractMethods:      contractMethods,
		InfoResponse:         apiInfoResponse,
		RootReference:        genesisrefs.ContractRootMember,
		ProxyToRootMethods:   proxyToRootMethods,
	}, nil
}
