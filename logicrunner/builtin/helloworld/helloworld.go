/*
 *    Copyright 2019 Insolar
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

package helloworld

import "github.com/insolar/insolar/core"

// HelloWorld contract
type HelloWorld struct {
	// Greeted - how many callers we "greated"
	Greeted int
}

// CodeRef returns something strange
func CodeRef() core.RecordRef {
	var ref core.RecordRef
	ref[core.RecordRefSize-1] = 1
	return ref
}

// NewHelloWorld returns a new empty contract
func NewHelloWorld() *HelloWorld {
	return &HelloWorld{}
}

// Greet greats the caller
func (hw *HelloWorld) Greet(name string) string {
	hw.Greeted++
	return "Hello " + name + "'s world"
}
