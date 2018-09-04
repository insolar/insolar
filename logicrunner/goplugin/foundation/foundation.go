/*
 *    Copyright 2018 INS Ecosystem
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
	"bytes"
	"time"
)

// Reference is an address of something on ledger.
type Reference []byte

// String - stringer interface
func (r Reference) String() string {
	return string(r)
}

// Equal is equaler
func (r Reference) Equal(o Reference) bool {
	return bytes.Equal(r, o)
}

// CallContext is a context of contract execution
type CallContext struct {
	Me     Reference // My Reference.
	Caller Reference // Reference of calling contract.
	Parent Reference // Reference to parent or container contract.
	Type   Reference // Reference to type record on ledger, we have just one type reference, yet.
	Time   time.Time // Time of Calling side made call.
	Pulse  uint64    // Number of current pulse.
}

// BaseContract is a base class for all contracts.
type BaseContract struct {
	context *CallContext // Context hidden from anyone
}

type ProxyInterface interface {
	GetReference() Reference
	GetClass() Reference
}

// BaseContractInterface is an interface to deal with any contract same way
type BaseContractInterface interface {
	GetReference() Reference
	GetClass() Reference
}

// MyReference - Returns public reference of contract
func (bc *BaseContract) MyReference() Reference {
	if bc.context == nil {
		return nil
	}
	return bc.context.Me
}

// GetContext returns current calling context of this object.
// It exists only for currently called contract.
func (bc *BaseContract) GetContext() *CallContext {
	return bc.context
}

// SetContext - do not use it in smartcontracts
func (bc *BaseContract) SetContext(cc *CallContext) {
	if bc.context == nil {
		bc.context = cc
	}
}

// GetImplementationFor finds delegate typed r in object and returns it
// unimplemented
func GetImplementationFor(o Reference, r Reference) ProxyInterface {
	return nil
}

// GetChildrenTyped returns set of children objects with corresponding type
func (bc *BaseContract) GetChildrenTyped(r Reference) []ProxyInterface {
	return nil
}

// GetObject create proxy by address
// unimplemented
func GetObject(ref Reference) ProxyInterface {
	return nil
}

// TakeDelegate injects delegate to object
func (bc *BaseContract) InjectDelegate(delegate BaseContractInterface, class Reference) Reference {
	return nil
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
func Call(Reference Reference, MethodName string, Arguments []interface{}) ([]interface{}, error) {
	return nil, nil
}
