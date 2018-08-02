package main

import "errors"

// @inscontract
type HelloWorlder struct { //Docccc
	Greeted int
}

// @method
func (hw *HelloWorlder) Hello() (string, error) {
	hw.Greeted++
	return "Hello world 2", nil
}

func (hw *HelloWorlder) Fail() (string, error) {
	hw.Greeted++
	return "", errors.New("We failed 2")
}

func (hw *HelloWorlder) Echo(s string) (string, error) {
	hw.Greeted++
	return s, nil
}

func (hw HelloWorlder) ConstEcho(s string) (string, error) {
	return s, nil
}

func JustExportedStaticFunction(int, int) {}

/// generated
