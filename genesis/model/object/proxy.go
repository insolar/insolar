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

package object

import (
	"github.com/insolar/insolar/genesis/model/class"
)

// Proxy marks instance as proxy object.
type Proxy interface {
	Child
	GetReference() Reference
	SetReference(reference Reference)
}

// BaseProxy is a base implementation of Proxy.
type BaseProxy struct {
	reference Reference
}

// GetClassID is a proxy call for instance method.
func (bp *BaseProxy) GetClassID() string {
	return class.ProxyID
}

func (bp *BaseProxy) GetClass() Proxy {
	return nil
}

// GetParent always returns nil.
func (bp *BaseProxy) GetParent() Parent {
	return nil
}

// GetReference is a proxy call for instance method.
func (bp *BaseProxy) GetReference() Reference {
	return bp.reference
}

// SetReference is a proxy call for instance method.
func (bp *BaseProxy) SetReference(reference Reference) {
	bp.reference = reference
}
