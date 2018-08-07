package main

import "errors"

// nolint
type HelloWorlder struct {
	Greeted int
}

// nolint
func (hw *HelloWorlder) Hello() (string, error) {
	hw.Greeted++
	return "Hello world 2", nil
}

// nolint
func (hw *HelloWorlder) Fail() (string, error) {
	hw.Greeted++
	return "", errors.New("We failed 2")
}

// nolint
func (hw *HelloWorlder) Echo(s string) (string, error) {
	hw.Greeted++
	return s, nil
}

// nolint
func (hw HelloWorlder) ConstEcho(s string) (string, error) {
	return s, nil
}

// nolint
func JustExportedStaticFunction(int, int) {}

var INSEXPORT HelloWorlder //nolint

/// generated
