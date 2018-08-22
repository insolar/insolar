package foundation

import (
	"time"
)

type Reference string

type CallContext struct {
	Me     *Reference
	Caller *Reference
	Parent *Reference
	Time   time.Time
	Pulse  uint64
}

type BaseContract struct {
	context *CallContext
}

func (bc BaseContract) GetContext() *CallContext {
	return bc.context
}

func (bc BaseContract) SetContext(c *CallContext) {
	if bc.context != nil {
		return
	}
	bc.context = c
}

var FakeLedger = make(map[*Reference]interface{})
var FakeDelegates = make(map[*Reference]map[*Reference]interface{})
var FakeChildren = make(map[*Reference]map[*Reference][]interface{})

func (bc BaseContract) GetImplementationFor(r *Reference) interface{} {
	return FakeDelegates[bc.context.Me][r]
}

func (bc BaseContract) GetChildrenTyped(r *Reference) []interface{} {
	return FakeChildren[bc.context.Me][r]
}

func (bc BaseContract) SelfDestructRequest() {}
