/*
 *    Copyright 2018 INS Ecosystem
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"errors"
	"fmt"
)

// @inscontract
// nolint
type HelloWorlder struct {
	Greeted int
}

// nolint
type FullName struct {
	First string
	Last  string
}

// nolint
type PersonalGreeting struct {
	Name    FullName
	Message string
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

func (hw *HelloWorlder) HelloHuman(Name FullName) PersonalGreeting {
	hw.Greeted++
	return PersonalGreeting{
		Name:    Name,
		Message: fmt.Sprintf("Dear %s %s, we specially say hello to you", Name.First, Name.Last),
	}
}

// nolint
func (hw *HelloWorlder) HelloHumanPointer(Name FullName) *PersonalGreeting {
	hw.Greeted++
	return &PersonalGreeting{
		Name:    Name,
		Message: fmt.Sprintf("Dear %s %s, we specially say hello to you", Name.First, Name.Last),
	}
}

// nolint
func (hw *HelloWorlder) MultiArgs(Name FullName, s string, i int) *PersonalGreeting {
	hw.Greeted++
	return &PersonalGreeting{
		Name:    Name,
		Message: fmt.Sprintf("Dear %s %s, we specially say hello to you", Name.First, Name.Last),
	}
}

// nolint
func (hw HelloWorlder) ConstEcho(s string) (string, error) {
	return s, nil
}

// nolint
func JustExportedStaticFunction(int, int) {}
