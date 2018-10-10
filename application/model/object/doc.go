/*
 *    Copyright 2018 Insolar
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
Package object provides basic interfaces and default implementations of them.
Provides reference interface and types of reference scope.

Usage:

	domain := "1"
	record := "1"
	ref, err := NewReference(domain, record, GlobalScope)

Provides factory interface for factories of all objects and proxy objects.

Usage:

	type someFactory struct {}

	func (f *someFactory) GetClassID() string {
		return "someFactory"
	}

	func (f *someFactory) GetReference() object.Reference {
		return f.Reference
	}

	func (f *someFactory) Create(parent object.Callable) (object.Callable, error)
		// do creation logic
	}

	func main() {
		factory := someFactory{}
		obj := factory.Create(parent)
	}

Provides resolver interface and default implementation of resolvers for getting objects from references.
Interface Resolver uses interface{} type for reference, class and proxy (which GetObject returns),
because in future implementation its going to be plugin. Virtual machine will be use it and provide resolving logic.

Usage:

	resolver := NewChildResolver(parent)
	obj, err := resolver.GetObject(ref, class.ObjectID)
	res := obj.(object.Object)


Proxy provides interface and default implementation of proxy. It inherited by SmartContractProxy and Factory

Usage:

	proxy := &BaseProxy{}

	proxy.SetReference(Reference) sets reference to proxy.
	proxy.GetReference() gets reference from proxy.
	proxy.GetParent() always returns nil.
	proxy.GetClassID() is a proxy call for instance method.


ReferenceContainer provides methods for store Reference as Proxy

Usage:

	ref, _ := object.NewReference(domain, record, object.GlobalScope)
	container = NewReferenceContainer(ref)

	container.GetClassID()         // return string representation of object's class.
	container.GetStoredReference() // returns stored reference.


*/
package object
