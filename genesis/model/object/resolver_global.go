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

// globalResolver is resolver for GlobalScope references.
type globalResolver struct {
	globalInstanceMap *map[string]Proxy
}

// newGlobalResolver creates new globalResolver instance with empty map.
func newGlobalResolver() *globalResolver {
	return &globalResolver{}
}

// GetObject resolves object by its reference and return its proxy.
func (r *globalResolver) GetObject(reference interface{}, cls interface{}) (interface{}, error) {
	ref, ok := reference.(Reference)
	if !ok {
		return nil, fmt.Errorf("reference is not Reference class object")
	}
	parentProxy, exist := (*r.globalInstanceMap)[ref.GetDomain()]
	if !exist {
		return nil, fmt.Errorf("reference with address `%s` not found", ref)
	}
	parent, ok := parentProxy.(Parent)
	if !ok {
		return nil, fmt.Errorf("object with domain `%s` can not have children", ref.GetDomain())
	}
	proxy, err := parent.GetChild(ref.GetRecord())
	if err != nil {
		return nil, err
	}

	class := proxy.GetClass()

	if err = checkClass(class, cls); err != nil {
		return nil, err
	}

	proxy.(Proxy).SetReference(ref)
	return proxy, nil
}

// InitGlobalMap sets globalInstanceMap into globalResolver.
func (r *globalResolver) InitGlobalMap(globalInstanceMap *map[string]Proxy) {
	if r.globalInstanceMap == nil {
		r.globalInstanceMap = globalInstanceMap
	}
}

// GlobalResolver is a public globalResolver instance for resolving all global references.
var GlobalResolver = newGlobalResolver()
