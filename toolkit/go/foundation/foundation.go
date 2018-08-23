package foundation

import (
	"fmt"
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
	GetContext(debug ...string) *CallContext
	GetImplementationFor(r *Reference) BaseContractInterface
	SetContext(c *CallContext)
}

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

func (bc *BaseContract) SetContext(c *CallContext) {
	if bc.context != nil {
		return
	}
	bc.context = c
}

var FakeLedger = make(map[*Reference]BaseContractInterface)
var FakeDelegates = make(map[*Reference]map[*Reference]BaseContractInterface)
var FakeChildren = make(map[*Reference]map[*Reference][]BaseContractInterface)

var FakeContexts = make(map[uint]*CallContext)
var contextStep uint = 0

func InjectFakeContext(step uint, ctx *CallContext, reset ...bool) {
	if len(reset) > 0 && reset[0] {
		FakeContexts = make(map[uint]*CallContext)
	}
	contextStep = 0
	FakeContexts[step] = ctx
}

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

func (bc *BaseContract) AddChild(child BaseContractInterface, class *Reference) *Reference {
	me := bc.context.Me
	uid, _ := uuid.NewV4()
	key := Reference(uid.String())

	child.SetContext(&CallContext{
		Me: &key,
	})
	FakeLedger[&key] = child

	if FakeChildren[me] == nil {
		FakeChildren[me] = make(map[*Reference][]BaseContractInterface)
	}
	/*if FakeChildren[me][class] == nil {
		FakeChildren[me][class] = make([]BaseContractInterface, 1)
	}*/

	FakeChildren[me][class] = append(FakeChildren[me][class], child)

	return &key
}

func (bc *BaseContract) TakeDelegate(delegate BaseContractInterface, class *Reference) *Reference {
	me := bc.context.Me
	uid, _ := uuid.NewV4()
	key := Reference(uid.String())

	delegate.SetContext(&CallContext{
		Me: &key,
	})
	FakeLedger[&key] = delegate

	if FakeDelegates[me] == nil {
		FakeDelegates[me] = make(map[*Reference]BaseContractInterface)
	}
	FakeDelegates[me][class] = delegate

	if FakeChildren[me] == nil {
		FakeChildren[me] = make(map[*Reference][]BaseContractInterface)
	}
	/*if FakeChildren[me][class] == nil {
		FakeChildren[me][class] = make([]BaseContractInterface, 1)
	}*/

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
