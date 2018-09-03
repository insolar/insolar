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

// Package foundation is a base package for writing smartcontracts in go language.
// This is client side to use in standalone tests. It have the same signatures
// as a real realization, but all methods is intended to simulate real ledger behavior in tests.
package foundation

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
)

// Reference is an address of something on ledger.
type Reference string

// String - stringer interface
func (r Reference) String() string {
	return string(r)
}

// Equal is equaler
func (r Reference) Equal(o Reference) bool {
	return r == o
}

// CallContext is a context of contract execution
type CallContext struct {
	Me     Reference // My Reference.
	Caller Reference // Reference of calling contract.
	Parent Reference // Reference to parent or container contract.
	Class  Reference // Reference to type record on ledger, we have just one type reference, yet.
	Time   time.Time // Time of Calling side made call.
	Pulse  uint64    // Number of current pulse.
}

// BaseContract is a base class for all contracts.
type BaseContract struct {
	context *CallContext // context is hidden from everyone and not presented in real implementation.
}

// BaseContractInterface is an interface to deal with any contract same way
type ProxyInterface interface {
	GetReference() Reference
	GetClass() Reference
}

// BaseContractInterface is an interface to deal with any contract same way
type BaseContractInterface interface {
	SetContext(ctx *CallContext)
	GetReference() Reference
	GetClass() Reference
}

// GetReference - Returns public reference of contract
func (bc *BaseContract) GetReference() Reference {
	return bc.context.Me
}

// GetClass
func (bc *BaseContract) GetClass() Reference {
	return bc.context.Class
}

// GetContext returns current calling context of this object.
// It exists only for currently called contract.
func (bc *BaseContract) GetContext(debug ...string) *CallContext {
	contextStep++
	if len(debug) > 0 && debug[0] != "" {
		fmt.Printf("%s: %d\n", debug[0], contextStep)
	}
	if FakeContexts[contextStep] != nil {
		return FakeContexts[contextStep]
	}
	if bc.context != nil {
		return bc.context
	}
	return &CallContext{}
}

func (bc *BaseContract) SetContext(ctx *CallContext) {
	bc.context = ctx
}

var FakeLedger = make(map[string]ProxyInterface)
var FakeDelegates = make(map[string]map[string]ProxyInterface)
var FakeChildren = make(map[string]map[string][]ProxyInterface)

var FakeContexts = make(map[uint]*CallContext)
var contextStep uint

// InjectFakeContext - add mocked context to queue for substitution
func InjectFakeContext(step uint, ctx *CallContext, reset ...bool) {
	if len(reset) > 0 && reset[0] {
		FakeContexts = make(map[uint]*CallContext)
	}
	contextStep = 0
	FakeContexts[step] = ctx
}

func GetImplementationFor(o Reference, r Reference) ProxyInterface {
	return FakeDelegates[o.String()][r.String()]
}

func (bc *BaseContract) GetChildrenTyped(r Reference) []ProxyInterface {
	return FakeChildren[bc.GetReference().String()][r.String()]
}

func SaveToLedger(rec BaseContractInterface, class Reference) Reference {
	key, _ := uuid.NewV4()
	rec.SetContext(&CallContext{
		Me:    Reference(key.String()),
		Class: class,
	})
	FakeLedger[key.String()] = rec.(ProxyInterface)
	return Reference(key.String())
}

func GetObject(ref Reference) BaseContractInterface {
	return FakeLedger[ref.String()].(BaseContractInterface)
}

func (bc *BaseContract) AddChild(child BaseContractInterface, class Reference) Reference {
	me := bc.GetReference()
	key, _ := uuid.NewV4()

	child.SetContext(&CallContext{
		Parent: me,
		Me:     Reference(key.String()),
		Class:  class,
	})
	FakeLedger[key.String()] = child.(ProxyInterface)

	if FakeChildren[me.String()] == nil {
		FakeChildren[me.String()] = make(map[string][]ProxyInterface)
	}

	FakeChildren[me.String()][class.String()] = append(FakeChildren[me.String()][class.String()], child)
	return Reference(key.String())
}

func (bc *BaseContract) InjectDelegate(delegate BaseContractInterface, class Reference) Reference {
	me := bc.GetReference()
	key, _ := uuid.NewV4()

	delegate.SetContext(&CallContext{
		Parent: me,
		Me:     Reference(key.String()),
		Class:  class,
	})

	FakeLedger[key.String()] = delegate.(ProxyInterface)

	if FakeDelegates[me.String()] == nil {
		FakeDelegates[me.String()] = make(map[string]ProxyInterface)
	}
	FakeDelegates[me.String()][class.String()] = delegate.(ProxyInterface)

	if FakeChildren[me.String()] == nil {
		FakeChildren[me.String()] = make(map[string][]ProxyInterface)
	}

	FakeChildren[me.String()][class.String()] = append(FakeChildren[me.String()][class.String()], delegate.(ProxyInterface))
	return Reference(key.String())
}

func (bc *BaseContract) SelfDestructRequest() {
	me := bc.GetReference()
	delete(FakeLedger, me.String())
	for _, v := range FakeDelegates {
		delete(v, me.String())
	}
	for _, c := range FakeChildren {
		arr := []ProxyInterface{}
		for _, v := range c[bc.context.Class.String()] {
			if v.GetReference().String() != me.String() {
				arr = append(arr, v)
			}
		}
		c[bc.context.Class.String()] = arr
	}
}
