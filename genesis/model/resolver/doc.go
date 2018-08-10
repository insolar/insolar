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
Package resolver provides interface and default implementation of resolvers for getting objects from references.
Interface Resolver uses interface{} type for reference, class and proxy (which GetObject returns),
because in future implementation its going to be plugin. Virtual machine will be use it and provide resolving logic.

Usage:
	package main

	import (
		"github.com/insolar/insolar/genesis/model/class"
		"github.com/insolar/insolar/genesis/model/object"
	}

	func main() {
		resolver := NewChildResolver(parent)
		obj, err := resolver.GetObject(ref, class.ObjectID)
		res := obj.(object.Object)
	}


Proxy provides interface and default implementation of proxy. It inherited by SmartContractProxy and Factory

Usage:

	proxy.SetReference(Reference) sets reference to proxy.
	proxy.GetReference() gets reference from proxy.
	proxy.GetParent() always returns nil.
	proxy.GetClassID() is a proxy call for instance method.


ReferenceContainer provides methods for store Reference as BaseProxy

Usage:

	ref, _ := object.NewReference(domain, record, object.GlobalScope)
	container = NewReferenceContainer(ref)

	container.GetClassID()         // return string representation of object's class.
	container.GetStoredReference() // returns stored reference.

*/
package resolver
