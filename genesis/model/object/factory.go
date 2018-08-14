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

package object

// Factory allows to create new objects with reference.
type Factory interface {
	Proxy
	// Create returns new instance of specified type.
	Create(parent Parent) (Proxy, error)
}

type BaseFactory struct {
	BaseProxy
}

// Composite marks that instance have ability to be compose in another object.
type Composite interface {
	GetInterfaceKey() string // string ID of interface/type of Composite object; basically, GetClassID()
}

// CompositeFactory allows to create new composites.
type CompositeFactory interface {
	Proxy
	Create(parent Parent) (Composite, error)
	GetInterfaceKey() string // string ID of interface/type of Composite object; basically, GetClassID()
}

// ComposingContainer allows to store composites.
type ComposingContainer interface {
	CreateComposite(compositeFactory CompositeFactory) (Composite, error)
	GetComposite(interfaceKey string) (Composite, error)
	GetOrCreateComposite(compositeFactory CompositeFactory) (Composite, error)
}
