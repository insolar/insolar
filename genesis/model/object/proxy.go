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

// Proxy marks instance as proxy object.
type Proxy interface {
	Child
}

// BaseProxy is a base implementation of Proxy.
type BaseProxy struct {
	Instance Child
}

// GetReference proxy call for instance method.
func (bp *BaseProxy) GetReference() *Reference {
	return bp.Instance.GetReference()
}

// GetParent proxy call for instance method.
func (bp *BaseProxy) GetParent() Parent {
	return bp.Instance.GetParent()
}

// GetClassID proxy call for instance method.
func (bp *BaseProxy) GetClassID() string {
	return bp.Instance.GetClassID()
}
