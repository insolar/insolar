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

// Handler can resolver references from any allowed scopes.
type Handler struct {
	globalResolver  *globalResolver
	childResolver   *childResolver
	contextResolver *contextResolver
}

// NewHandler creates new resolverHandler instance.
func NewHandler(p interface{}) *Handler {
	parent, ok := p.(object.Parent)
	if !ok {
		parent = nil
	}
	return &Handler{
		globalResolver:  GlobalResolver,
		childResolver:   newChildResolver(parent),
		contextResolver: newContextResolver(parent),
	}
}

// GetObject resolve object by its reference and return its proxy.
func (r *Handler) GetObject(reference interface{}, classID interface{}) (interface{}, error) {
	ref, ok := reference.(*object.Reference)
	if !ok {
		return nil, fmt.Errorf("reference is not Reference class object")
	}
	switch ref.Scope {
	case object.GlobalScope:
		return r.globalResolver.GetObject(ref, classID)
	case object.ContextScope:
		return r.contextResolver.GetObject(ref, classID)
	case object.ChildScope:
		return r.childResolver.GetObject(ref, classID)
	default:
		return nil, fmt.Errorf("unknown scope type: %d", ref.Scope)
	}
}

// InitGlobalMap set globalInstanceMap into globalResolver.
func (r *Handler) InitGlobalMap(globalInstanceMap *map[string]object.Proxy) {
	r.globalResolver.InitGlobalMap(globalInstanceMap)
}
