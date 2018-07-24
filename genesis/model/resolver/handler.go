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

type resolverHandler struct {
	globalReserver  *globalResolver
	childReserver   *childResolver
	contextReserver *contextResolver
}

// NewResolverHandler creates new resolverHandler instance.
func NewResolverHandler(p object.Parent) Resolver {
	childReserver := newChildResolver(p)
	contextReserver := newContextResolver(p)
	return &resolverHandler{
		globalReserver:  GlobalResolver,
		childReserver:   childReserver,
		contextReserver: contextReserver,
	}
}

// GetObject resolve object by its reference and return its proxy.
func (r *resolverHandler) GetObject(ref *object.Reference, classID string) (object.Proxy, error) {
	switch ref.Scope {
	case object.GlobalScope:
		return r.globalReserver.GetObject(ref, classID)
	case object.ContextScope:
		return r.contextReserver.GetObject(ref, classID)
	case object.ChildScope:
		return r.childReserver.GetObject(ref, classID)
	default:
		return nil, fmt.Errorf("unknown scope type: %d", ref.Scope)
	}
}
