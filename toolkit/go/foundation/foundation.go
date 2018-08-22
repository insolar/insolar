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
	GetImplementationFor(r *Reference) BaseContractInterface
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

var FakeLedger = make(map[*Reference]BaseContractInterface)
var FakeDelegates = make(map[*Reference]map[*Reference]BaseContractInterface)
var FakeChildren = make(map[*Reference]map[*Reference][]BaseContractInterface)

func (bc *BaseContract) GetImplementationFor(r *Reference) BaseContractInterface {
	return FakeDelegates[bc.context.Me][r]
}

func (bc *BaseContract) GetChildrenTyped(r *Reference) []BaseContractInterface {
	return FakeChildren[bc.context.Me][r]
}

func SaveToLedger(rec BaseContractInterface) *Reference {
	u2, _ := uuid.NewV4()
	key := Reference(u2.String())
	FakeLedger[&key] = rec
	return &key
}

func GetObject(ref *Reference) BaseContractInterface {
	return FakeLedger[ref].(BaseContractInterface)
}

func (bc *BaseContract) TakeDelegate(delegate BaseContractInterface, class *Reference) *Reference {
	me := bc.context.Me
	uid, _ := uuid.NewV4()
	key := Reference(uid.String())

	FakeLedger[&key] = delegate

	if FakeDelegates[me] == nil {
		FakeDelegates[me] = make(map[*Reference]BaseContractInterface)
	}
	FakeDelegates[me][class] = delegate

	if FakeChildren[me] == nil {
		FakeChildren[me] = make(map[*Reference][]BaseContractInterface)
	}
	if FakeChildren[me][class] == nil {
		FakeChildren[me][class] = make([]BaseContractInterface, 1)
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
		arr := []BaseContractInterface{}
		for _, v := range c[bc.context.Type] {
			if v.(BaseContractInterface).GetContext().Me != me {
				arr = append(arr, v)
			}
		}
		c[bc.context.Type] = arr
	}
}
