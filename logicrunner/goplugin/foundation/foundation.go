/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

// Package foundation server implementation of smartcontract functions
package foundation

import (
	"github.com/insolar/insolar/core"
)

// BaseContract is a base class for all contracts.
type BaseContract struct {
	context *core.LogicCallContext // Context hidden from anyone
}

// ProxyInterface interface any proxy of a contract implements
type ProxyInterface interface {
	GetReference() core.RecordRef
	GetClass() core.RecordRef
}

// BaseContractInterface is an interface to deal with any contract same way
type BaseContractInterface interface {
	GetReference() core.RecordRef
	GetClass() core.RecordRef
}

// GetReference - Returns public reference of contract
func (bc *BaseContract) GetReference() core.RecordRef {
	if bc.context == nil {
		panic("object has no context set before first use")
	}
	if bc.context.Callee == nil {
		panic("context has no callee set")
	}
	return *bc.context.Callee
}

// GetClass - Returns class of contract
func (bc *BaseContract) GetClass() core.RecordRef {
	if bc.context == nil {
		panic("object has no context set before first use")
	}
	return *bc.context.Class
}

// GetContext returns current calling context of this object.
// It exists only for currently called contract.
func (bc *BaseContract) GetContext() core.LogicCallContext {
	return *bc.context
}

// SetContext - do not use it in smartcontracts
func (bc *BaseContract) SetContext(cc *core.LogicCallContext) {
	if bc.context == nil {
		bc.context = cc
	} else {
		panic("context can not be set twice")
	}
}

// GetImplementationFor finds delegate typed r in object and returns it
func GetImplementationFor(o core.RecordRef, r core.RecordRef) ProxyInterface {
	panic("not implemented")
}

// GetChildrenTyped returns set of children objects with corresponding type
func (bc *BaseContract) GetChildrenTyped(r core.RecordRef) []ProxyInterface {
	panic("not implemented")
}

// GetObject create proxy by address
// unimplemented
func GetObject(ref core.RecordRef) ProxyInterface {
	panic("not implemented")
}

// SelfDestructRequest contract will be marked as deleted after call finishes
func (bc *BaseContract) SelfDestructRequest() {
}

/////// next code is system helper for wrappers generator //////

// CBORMarshaler is a special interface for serializer object
type CBORMarshaler interface {
	Marshal(interface{}) []byte
	Unmarshal(interface{}, []byte)
}

// Call other contract via network dispatcher
func Call(Reference core.RecordRef, MethodName string, Arguments []interface{}) ([]interface{}, error) {
	return nil, nil
}
