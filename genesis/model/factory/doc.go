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

/*
Package factory provides interfaces for factories of all objects and proxy objects.

Usage:

	package main

	type someFactory struct {}

	func (f *someFactory) GetClassID() string {
		return "someFactory"
	}

	func (f *someFactory) GetReference() *object.Reference {
		return f.Reference
	}

	func (f *someFactory) Create(parent object.Callable) (object.Callable, error)
		// do creation logic
	}

	func main() {
		factory := someFactory{}
		obj := factory.Create(parent)
	}
*/
package factory
