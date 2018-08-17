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

// ResolverHandler should resolve references from any allowed scopes.
type ResolverHandler struct {
	globalResolver  *globalResolver
	childResolver   *childResolver
	contextResolver *contextResolver
}

// NewResolverHandler creates new ResolverHandler instance.
func NewResolverHandler(p interface{}) *ResolverHandler {
	parent, ok := p.(Parent)
	if !ok {
		parent = nil
	}
	return &ResolverHandler{
		globalResolver:  GlobalResolver,
		childResolver:   newChildResolver(parent),
		contextResolver: newContextResolver(parent),
	}
}

// GetObject resolves object by its reference and return its proxy.
func (r *ResolverHandler) GetObject(reference interface{}, class interface{}) (interface{}, error) {
	ref, ok := reference.(Reference)
	if !ok {
		return nil, fmt.Errorf("reference is not Reference class object")
	}
	switch ref.GetScope() {
	case GlobalScope:
		return r.globalResolver.GetObject(ref, class)
	case ContextScope:
		return r.contextResolver.GetObject(ref, class)
	case ChildScope:
		return r.childResolver.GetObject(ref, class)
	default:
		return nil, fmt.Errorf("unknown scope type: %d", ref.GetScope())
	}
}

// InitGlobalMap sets globalInstanceMap into globalResolver.
func (r *ResolverHandler) InitGlobalMap(globalInstanceMap *map[string]Proxy) {
	r.globalResolver.InitGlobalMap(globalInstanceMap)
}
