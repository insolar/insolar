// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// Package foundation server implementation of smartcontract functions
package foundation

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/common"
)

// BaseContract is a base class for all contracts.
type BaseContract struct {
}

// ProxyInterface interface any proxy of a contract implements
type ProxyInterface interface {
	GetReference() insolar.Reference
	GetPrototype() (insolar.Reference, error)
	GetCode() (insolar.Reference, error)
}

// BaseContractInterface is an interface to deal with any contract same way
type BaseContractInterface interface {
	GetReference() insolar.Reference
	GetPrototype() insolar.Reference
	GetCode() insolar.Reference
}

// GetReference - Returns public reference of contract
func (bc *BaseContract) GetReference() insolar.Reference {
	ctx := bc.GetContext()
	if ctx.Callee == nil {
		panic("context has no callee set")
	}
	return *ctx.Callee
}

// GetPrototype - Returns prototype of contract
func (bc *BaseContract) GetPrototype() insolar.Reference {
	return *bc.GetContext().Prototype
}

// GetCode - Returns prototype of contract
func (bc *BaseContract) GetCode() insolar.Reference {
	return *bc.GetContext().Code
}

// GetContext returns current calling context OBSOLETED.
func (bc *BaseContract) GetContext() *insolar.LogicCallContext {
	return GetLogicalContext()
}

// SelfDestruct contract will be marked as deleted
func (bc *BaseContract) SelfDestruct() error {
	return common.CurrentProxyCtx.DeactivateObject(bc.GetReference())
}
