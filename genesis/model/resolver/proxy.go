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

package resolver

import (
	"github.com/insolar/insolar/genesis/model/object"
)

// Proxy marks instance as proxy object.
type Proxy interface {
	object.Child
}

// BaseProxy is a base implementation of Proxy.
type BaseProxy struct {
	Instance object.Child
}

// GetReference is a proxy call for instance method.
func (bp *BaseProxy) GetReference() *object.Reference {
	return bp.Instance.GetReference()
}

// GetParent is a proxy call for instance method.
func (bp *BaseProxy) GetParent() object.Parent {
	return bp.Instance.GetParent()
}

// GetResolver always returns nil.
func (bp *BaseProxy) GetResolver() Resolver {
	return nil
}

// GetClassID is a proxy call for instance method.
func (bp *BaseProxy) GetClassID() string {
	return bp.Instance.GetClassID()
}
