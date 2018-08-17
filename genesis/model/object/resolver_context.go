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

import (
	"fmt"
)

// contextResolver is resolver for ContextScope references.
type contextResolver struct {
	parent Parent
}

// newContextResolver creates new contextResolver instance.
func newContextResolver(parent Parent) *contextResolver {
	return &contextResolver{
		parent: parent,
	}
}

// GetObject resolves object by its reference and return its proxy.
func (r *contextResolver) GetObject(reference interface{}, cls interface{}) (interface{}, error) {
	ref, ok := reference.(Reference)
	if !ok {
		return nil, fmt.Errorf("reference is not Reference class object")
	}
	contextHolder := r.parent
	obj, err := contextHolder.GetContextStorage().Get(ref.GetRecord())

	if err != nil {
		return nil, err
	}

	proxy, ok := obj.(Proxy)
	if !ok {
		return nil, fmt.Errorf("object is not Proxy")
	}

	class := proxy.GetClass()
	for {

		if _, ok := class.(*ReferenceContainer); !ok {
			break
		}
		contextHolderWithChildInterface, isChild := contextHolder.(Child)
		if !isChild {
			return nil, fmt.Errorf("object with name %s does not exist", ref)
		}
		contextHolder = contextHolderWithChildInterface.GetParent()
		contextResolver := newContextResolver(contextHolder)
		newProxy, err := contextResolver.GetObject(proxy, cls)
		if err != nil {
			return nil, err
		}
		proxy = newProxy.(Proxy)
	}

	if err = checkClass(class, cls); err != nil {
		return nil, err
	}

	proxy.SetReference(ref)
	return proxy, nil
}
