package main

import (
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"time"
)

type One struct {
	foundation.BaseContract
	Number int
}

func New() (*One, error) {
	return &One{Number: 0}, nil
}

func NewWithNumber(num int) (*One, error) {
	return &One{Number: num}, nil
}

var INSATTR_GetAndIncrement_API = true

func (c *One) GetAndIncrement() (int, error) {
	time.Sleep(200 * time.Millisecond)
	c.Number++
	return c.Number, nil
}

var INSATTR_GetAndDecrement_API = true

func (c *One) GetAndDecrement() (int, error) {
	time.Sleep(200 * time.Millisecond)
	c.Number--
	return c.Number, nil
}

var INSATTR_Get_API = true

func (c *One) Get() (int, error) {
	return c.Number, nil
}

var INSATTR_DoNothing_API = true

func (r *One) DoNothing() error {
	return nil
}
