package second

import (
	"errors"
	"fmt"

	one "github.com/insolar/insolar/application/builtin/proxy/first"
	"github.com/insolar/insolar/application/builtin/proxy/third"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

type Second struct {
	foundation.BaseContract
	Number int
	OneRef insolar.Reference
	X      int
	P      Payload
}

func (r *Second) GetName() (string, error) {
	return "first", nil
}

func New() (*Second, error) {
	return &Second{Number: 10, OneRef: *insolar.NewEmptyReference()}, nil
}
func NewWithOne(oneNumber int) (*Second, error) {
	holder := one.NewWithNumber(oneNumber)
	objOne, err := holder.AsChild(foundation.GetNodeDomain())
	if err != nil {
		return nil, err
	}
	return &Second{Number: oneNumber, OneRef: objOne.GetReference()}, nil
}

func (r *Second) DoNothing() error {
	return nil
}

func (r *Second) Get() (int, error) {
	return r.Number, nil
}

type Payload struct {
	Int int
	Str string
}

func NewWithX() (*Second, error) {
	return &Second{X: 0}, nil
}

func (r *Second) Hello(s string) (string, error) {
	r.X++
	return fmt.Sprintf("Hello you too, %s. %d times!", s, r.X), nil
}

func (r *Second) GetPayload() (Payload, error) {
	return r.P, nil
}

func (r *Second) SetPayload(P Payload) error {
	r.P = P
	return nil
}

func (r *Second) GetPayloadString() (string, error) {
	return r.P.Str, nil
}

func NewSaga() (*Second, error) {
	return &Second{Number: 100}, nil
}
func (r *Second) GetBalance() (int, error) {
	return r.Number, nil
}

//ins:saga(INS_FLAG_NO_ROLLBACK_METHOD)
func (r *Second) Accept(amount int) error {
	r.Number += amount
	return nil
}

func (r *Second) AnError() error {
	return errors.New("an error")
}

func (r *Second) NoError() error {
	return nil
}

func (r *Second) ReturnNil() (*string, error) {
	return nil, nil
}

func (r *Second) ExternalCallDoNothing() error {
	holder := third.New()
	objThree, err := holder.AsChild(r.GetReference())
	if err != nil {
		return err
	}
	return objThree.DoNothing()
}

func (r *Second) GetParent() (string, error) {
	return r.GetContext().Parent.String(), nil
}

func NewNil() (*Second, error) {
	// nil, nil is considered a logical error in the constructor
	return nil, nil
}

func NewWithErr() (*Second, error) {
	return nil, errors.New("Epic fail in NewWithErr()")
}
