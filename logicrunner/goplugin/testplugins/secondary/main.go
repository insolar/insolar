package main

// command for build this file go build -buildmode=plugin main.go

import "errors"

type HelloWorlder struct {
	Greeted int
}

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

var EXP HelloWorlder
