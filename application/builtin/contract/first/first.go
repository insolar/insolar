// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package first

import (
	"errors"

	recursive "github.com/insolar/insolar/application/builtin/proxy/first"
	"github.com/insolar/insolar/application/builtin/proxy/second"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

type One struct {
	foundation.BaseContract
	Number int
	Friend insolar.Reference
}

func New() (*One, error) {
	return &One{}, nil
}

func (r *One) Panic() error {
	panic("AAAAAAAA!")
	return nil
}
func NewPanic() (*One, error) {
	panic("BBBBBBBB!")
}

func (r *One) Recursive() error {
	remoteSelf := recursive.GetObject(r.GetReference())
	err := remoteSelf.Recursive()
	return err
}

func (r *One) Test(firstRef *insolar.Reference) (string, error) {
	return second.GetObject(*firstRef).GetName()
}

func NewZero() (*One, error) {
	return &One{Number: 0}, nil
}
func NewWithNumber(num int) (*One, error) {
	return &One{Number: num}, nil
}

func (r *One) DoNothing() error {
	return nil
}

func (r *One) Inc() (int, error) {
	r.Number++
	return r.Number, nil
}

func (r *One) Get() (int, error) {
	return r.Number, nil
}

func (r *One) Dec() (int, error) {
	r.Number--
	return r.Number, nil
}

func (r *One) Hello(s string) (string, error) {
	holder := second.NewWithX()
	friend, err := holder.AsChild(r.GetReference())
	if err != nil {
		return "1", err
	}
	res, err := friend.Hello(s)
	if err != nil {
		return "2", err
	}
	r.Friend = friend.GetReference()
	return "Hi, " + s + "! Two said: " + res, nil
}

func (r *One) Again(s string) (string, error) {
	res, err := second.GetObject(r.Friend).Hello(s)
	if err != nil {
		return "", err
	}
	return "Hi, " + s + "! Two said: " + res, nil
}

func (r *One) GetFriend() (string, error) {
	return r.Friend.String(), nil
}

func (r *One) TestPayload() (second.Payload, error) {
	f := second.GetObject(r.Friend)
	err := f.SetPayload(second.Payload{Int: 10, Str: "HiHere"})
	if err != nil {
		return second.Payload{}, err
	}
	p, err := f.GetPayload()
	if err != nil {
		return second.Payload{}, err
	}
	str, err := f.GetPayloadString()
	if err != nil {
		return second.Payload{}, err
	}
	if p.Str != str {
		return second.Payload{}, errors.New("Oops")
	}
	return p, nil
}

func (r *One) ManyTimes() error {
	holder := second.New()
	friend, err := holder.AsChild(r.GetReference())
	if err != nil {
		return err
	}
	for i := 0; i < 100; i++ {
		_, err := friend.Hello("some")
		if err != nil {
			return err
		}
	}
	return nil
}

func NewSaga() (*One, error) {
	return &One{Number: 100}, nil
}

func (r *One) Transfer(n int) (string, error) {
	rec := recursive.NewSaga()
	w2, err := rec.AsChild(r.GetReference())
	if err != nil {
		return "1", err
	}
	r.Number -= n
	err = w2.Accept(n)
	if err != nil {
		return "2", err
	}
	return w2.GetReference().String(), nil
}

func (r *One) GetBalance() (int, error) {
	return r.Number, nil
}

//ins:saga(Rollback)
func (r *One) Accept(amount int) error {
	r.Number += amount
	return nil
}
func (r *One) Rollback(amount int) error {
	r.Number -= amount
	return nil
}

type StepOneArgs struct {
	CallerRef insolar.Reference
	Amount    int
}

func (r *One) TransferWithRollback(n int) (string, error) {
	second := recursive.NewSaga()
	w2, err := second.AsChild(r.GetReference())
	if err != nil {
		return "1", err
	}
	// second saga call
	args := &recursive.StepOneArgs{
		CallerRef: r.GetReference(),
		Amount:    n,
	}
	err = w2.AcceptStepOne(args)
	if err != nil {
		return "2", err
	}
	return w2.GetReference().String(), nil
}

//ins:saga(RollbackStepOne)
func (r *One) AcceptStepOne(a *StepOneArgs) error {
	r.Number += a.Amount
	// second saga call from the accept method
	first := recursive.GetObject(a.CallerRef)
	return first.AcceptStepTwo(a.Amount)
}

func (r *One) RollbackStepOne(a *StepOneArgs) error {
	r.Number -= a.Amount
	return nil
}

//ins:saga(RollbackStepTwo)
func (r *One) AcceptStepTwo(amount int) error {
	r.Number -= amount
	return nil
}
func (r *One) RollbackStepTwo(amount int) error {
	r.Number += amount
	return nil
}

func (r *One) TransferTwice(n int) (string, error) {
	second := recursive.NewSaga()
	w2, err := second.AsChild(r.GetReference())
	if err != nil {
		return "1", err
	}
	r.Number -= n
	// second saga call
	fst := n / 2
	err = w2.Accept(fst)
	if err != nil {
		return "2", err
	}
	// second saga call
	err = w2.Accept(n - fst)
	if err != nil {
		return "3", err
	}
	return w2.GetReference().String(), nil
}

func (r *One) TransferToAnotherContract(n int) (string, error) {
	second := second.NewSaga()
	w2, err := second.AsChild(r.GetReference())
	if err != nil {
		return "1", err
	}
	r.Number -= n
	err = w2.Accept(n)
	if err != nil {
		return "2", err
	}
	return w2.GetReference().String(), nil
}

func (r *One) SelfRef() (string, error) {
	return r.GetReference().String(), nil
}

func (r *One) AnError() error {
	holder := second.New()
	friend, err := holder.AsChild(r.GetReference())
	if err != nil {
		return err
	}
	return friend.AnError()
}

func (r *One) NoError() error {
	holder := second.New()
	friend, err := holder.AsChild(r.GetReference())
	if err != nil {
		return err
	}
	return friend.NoError()
}

func (r *One) ReturnNil() (*string, error) {
	holder := second.New()
	friend, err := holder.AsChild(r.GetReference())
	if err != nil {
		return nil, err
	}
	return friend.ReturnNil()
}

func (r *One) ConstructorReturnNil() (*string, error) {
	holder := second.NewNil()
	_, err := holder.AsChild(r.GetReference())
	if err != nil {
		return nil, err
	}
	ok := "all was well"
	return &ok, nil
}

func (r *One) ConstructorReturnError() (*string, error) {
	holder := second.NewWithErr()
	_, err := holder.AsChild(r.GetReference())
	if err != nil {
		return nil, err
	}
	ok := "all was well"
	return &ok, nil
}

func (r *One) GetChildPrototype() (string, error) {
	holder := second.New()
	child, err := holder.AsChild(r.GetReference())
	if err != nil {
		return "", err
	}
	ref, err := child.GetPrototype()
	return ref.String(), err
}
