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
	GetImplementationFor(r *Reference) interface{}
}

func (bc *BaseContract) GetContext() *CallContext {
	return bc.context
}

func (bc *BaseContract) SetContext(c *CallContext) {
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

func GetObject(ref *Reference) BaseContractInterface {
	return FakeLedger[ref].(BaseContractInterface)
}

func (bc *BaseContract) SetYourDelegate(delegate interface{}, class *Reference) *Reference {
	me := bc.context.Me
	uid, _ := uuid.NewV4()
	key := Reference(uid.String())

	FakeLedger[&key] = delegate

	if FakeDelegates[me] == nil {
		FakeDelegates[me] = make(map[*Reference]interface{})
	}
	FakeDelegates[me][class] = delegate

	if FakeChildren[me] == nil {
		FakeChildren[me] = make(map[*Reference][]interface{})
	}
	if FakeChildren[me][class] == nil {
		FakeChildren[me][class] = make([]interface{}, 1)
	}

	FakeChildren[me][class] = append(FakeChildren[me][class], delegate)

	return &key
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
