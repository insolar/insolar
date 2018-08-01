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
Package object provides basic interfaces and default implementations of them. Provides reference interface and types of reference scope.

Usage:
	package main

	func main() {
		domain := "1"
		record := "1"
		ref, err := NewReference(domain, record, GlobalScope)
	}


Proxy is public interface to call object's methods. If you want to make proxy for your object inherit BaseProxy

Usage:

	// make your custom domain proxy

	type customDomainProxy struct {
		object.BaseProxy
	}

	// create proxy for your custom domain

	func newCustomDomainProxy(parent object.Parent) (*customDomainProxy, error) {
		instance, err := newCustomDomain(parent)
		if err != nil {
			return nil, err
		}
		return &customDomainProxy{
			BaseProxy: object.BaseProxy{
				Instance: instance,
			},
		}, nil
	}

    proxy, err := newCustomDomainProxy(...)

    proxy.GetReference() is a proxy call for instance method.
    proxy.GetParent() is a proxy call for instance method.
    proxy.GetClassID() is a proxy call for instance method.

*/
package object
