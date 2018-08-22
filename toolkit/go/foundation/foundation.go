package foundation

import (
	"time"

	"github.com/satori/go.uuid"
)

type Reference string

type CallContext struct {
	Me     *Reference
	Caller *Reference
	Parent *Reference
	Type   *Reference
	Time   time.Time
	Pulse  uint64
}

type BaseContract struct {
	context *CallContext
}

type BaseContractInterface interface {
	GetContext() *CallContext
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

func (bc *BaseContract) GetImplementationFor(r *Reference) interface{} {
	return FakeDelegates[bc.context.Me][r]
}

func (bc *BaseContract) GetChildrenTyped(r *Reference) []interface{} {
	return FakeChildren[bc.context.Me][r]
}

func SaveToLedger(rec interface{}) *Reference {
	u2, _ := uuid.NewV4()
	key := Reference(u2.String())
	FakeLedger[&key] = rec
	return &key
}

func SetDelegate(to *Reference, class *Reference, delegate interface{}) {
	if FakeDelegates[to] == nil {
		FakeDelegates[to] = make(map[*Reference]interface{})
	}
	FakeDelegates[to][class] = delegate
}

func (bc *BaseContract) SelfDestructRequest() {
	me := bc.context.Me
	delete(FakeLedger, me)
	for _, v := range FakeDelegates {
		delete(v, me)
	}
	for _, c := range FakeChildren {
		arr := []interface{}{}
		for _, v := range c[bc.context.Type] {
			if v.(BaseContractInterface).GetContext().Me != me {
				arr = append(arr, v)
			}
		}
		c[bc.context.Type] = arr
	}
}
