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
	"fmt"

	"github.com/insolar/insolar/genesis/model/object"
)

// globalResolver is resolver for GlobalScope references.
type globalResolver struct {
	globalInstanceMap *map[string]object.Proxy
}

// newGlobalResolver creates new globalResolver instance with empty map.
func newGlobalResolver() *globalResolver {
	return &globalResolver{}
}

// GetObject resolve object by its reference and return its proxy.
func (r *globalResolver) GetObject(reference interface{}, cls interface{}) (interface{}, error) {
	ref, ok := reference.(*object.Reference)
	if !ok {
		return nil, fmt.Errorf("reference is not Reference class object")
	}
	parentProxy, isExist := (*r.globalInstanceMap)[ref.Domain]
	if !isExist {
		return nil, fmt.Errorf("reference with address `%s` not found", ref)
	}
	parent, ok := parentProxy.(object.Parent)
	if !ok {
		return nil, fmt.Errorf("object with domain `%s` can not have children", ref.Domain)
	}
	proxy, err := parent.GetChild(ref.Record)
	if err != nil {
		return nil, err
	}
	classID, ok := cls.(string)
	if !ok {
		return nil, fmt.Errorf("classID is not string")
	}
	if proxy.GetClassID() != classID {
		return nil, fmt.Errorf("instance class is not `%s`", classID)
	}
	return proxy, nil
}

// InitGlobalMap set globalInstanceMap into globalResolver.
func (r *globalResolver) InitGlobalMap(globalInstanceMap *map[string]object.Proxy) {
	if r.globalInstanceMap == nil {
		r.globalInstanceMap = globalInstanceMap
	}
}

// GlobalResolver is a public globalResolver instance for resolving all global references.
var GlobalResolver = newGlobalResolver()
