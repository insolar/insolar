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
	"reflect"
	"time"
)

// Reference is an address of something on ledger.
type Reference string

// String - stringer interface
func (r *Reference) String() string {
	return string(*r)
}

// CallContext is a context of contract execution
type CallContext struct {
	Me     *Reference // My Reference.
	Caller *Reference // Reference of calling contract.
	Parent *Reference // Reference to parent or container contract.
	Type   *Reference // Reference to type record on ledger, we have just one type reference, yet.
	Time   time.Time  // Time of Calling side made call.
	Pulse  uint64     // Number of current pulse.
}

// BaseContract is a base class for all contracts.
type BaseContract struct {
	context *CallContext // context is hidden from everyone and not presented in real implementation.
}

// BaseContractInterface is an interface to deal with any contract same way
type BaseContractInterface interface {
	MyReference() *Reference
	GetImplementationFor(r *Reference) BaseContractInterface
	SetContext(c *CallContext)
}

// MyReference - Returns public reference of contract
func (bc *BaseContract) MyReference() *Reference {
	r := Reference(fmt.Sprintf("%x", reflect.ValueOf(bc).Pointer()))
	return &r
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

// SetContext sets context on contract in testing environment.
// It is not supposed to use it from contracts.
func (bc *BaseContract) SetContext(c *CallContext) {
	bc.context = c
}

var FakeLedger = make(map[string]BaseContractInterface)
var FakeDelegates = make(map[string]map[string]BaseContractInterface)
var FakeChildren = make(map[string]map[string][]BaseContractInterface)

var FakeContexts = make(map[uint]*CallContext)
var contextStep uint = 0

// InjectFakeContext - add mocked context to queue for substitution
func InjectFakeContext(step uint, ctx *CallContext, reset ...bool) {
	if len(reset) > 0 && reset[0] {
		FakeContexts = make(map[uint]*CallContext)
	}
	contextStep = 0
	FakeContexts[step] = ctx
}

func (bc *BaseContract) GetImplementationFor(r *Reference) BaseContractInterface {
	return FakeDelegates[bc.MyReference().String()][r.String()]
}

func GetImplementationFor(o *Reference, r *Reference) BaseContractInterface {
	return FakeDelegates[o.String()][r.String()]
}

func (bc *BaseContract) GetChildrenTyped(r *Reference) []BaseContractInterface {
	return FakeChildren[bc.MyReference().String()][r.String()]
}

func SaveToLedger(rec BaseContractInterface) *Reference {
	key := rec.MyReference()
	FakeLedger[key.String()] = rec
	return key
}

func GetObject(ref *Reference) BaseContractInterface {
	return FakeLedger[ref.String()].(BaseContractInterface)
}

func (bc *BaseContract) AddChild(child BaseContractInterface, class *Reference) *Reference {
	me := bc.MyReference()
	key := child.MyReference()
	child.SetContext(&CallContext{
		Me:     key,
		Parent: me,
		Type:   class,
	})
	FakeLedger[key.String()] = child

	if FakeChildren[me.String()] == nil {
		FakeChildren[me.String()] = make(map[string][]BaseContractInterface)
	}
	/*if FakeChildren[me][class] == nil {
		FakeChildren[me][class] = make([]BaseContractInterface, 1)
	}*/

	FakeChildren[me.String()][class.String()] = append(FakeChildren[me.String()][class.String()], child)

	return key
}

func (bc *BaseContract) TakeDelegate(delegate BaseContractInterface, class *Reference) *Reference {
	me := bc.MyReference()
	key := delegate.MyReference()

	delegate.SetContext(&CallContext{
		Me:     key,
		Parent: me,
		Type:   class,
	})
	FakeLedger[key.String()] = delegate

	if FakeDelegates[me.String()] == nil {
		FakeDelegates[me.String()] = make(map[string]BaseContractInterface)
	}
	FakeDelegates[me.String()][class.String()] = delegate

	if FakeChildren[me.String()] == nil {
		FakeChildren[me.String()] = make(map[string][]BaseContractInterface)
	}
	/*if FakeChildren[me][class] == nil {
		FakeChildren[me][class] = make([]BaseContractInterface, 1)
	}*/

	FakeChildren[me.String()][class.String()] = append(FakeChildren[me.String()][class.String()], delegate)

	return key
}

func (bc *BaseContract) SelfDestructRequest() {
	me := bc.MyReference()
	delete(FakeLedger, me.String())
	for _, v := range FakeDelegates {
		delete(v, me.String())
	}
	for _, c := range FakeChildren {
		arr := []BaseContractInterface{}
		for _, v := range c[bc.context.Type.String()] {
			if v.MyReference().String() != me.String() {
				arr = append(arr, v)
			}
		}
		c[bc.context.Type.String()] = arr
	}
}
