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
		"member.create":                   true,
		"first.New":                       true,
		"first.NewPanic":                  true,
		"first.Panic":                     true,
		"panicAsLogicalError.New":         true,
		"panicAsLogicalError.NewPanic":    true,
		"panicAsLogicalError.Panic":       true,
		"first.Recursive":                 true,
		"first.Test":                      true,
		"third.New":                       true,
		"first.NewZero":                   true,
		"second.NewWithOne":               true,
		"first.Get":                       true,
		"first.Inc":                       true,
		"first.Dec":                       true,
		"first.Hello":                     true,
		"first.Again":                     true,
		"first.GetFriend":                 true,
		"second.Hello":                    true,
		"first.TestPayload":               true,
		"first.ManyTimes":                 true,
		"first.NewSaga":                   true,
		"first.Transfer":                  true,
		"first.GetBalance":                true,
		"first.TransferWithRollback":      true,
		"first.TransferTwice":             true,
		"first.TransferToAnotherContract": true,
		"second.GetBalance":               true,
		"third.Transfer":                  true,
		"third.GetSagaCallsNum":           true,
		"first.SelfRef":                   true,
		"first.AnError":                   true,
		"first.NoError":                   true,
		"first.ReturnNil":                 true,
		"first.ConstructorReturnNil":      true,
		"first.ConstructorReturnError":    true,
		"first.GetChildPrototype":         true,
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
