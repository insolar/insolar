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

// GlobalResolver is resolver for GlobalScope references.
type GlobalResolver struct {
	globalInstanceMap *map[*object.Reference]Proxy
}

// NewGlobalResolver creates new GlobalResolver instance.
// TODO: pass map?
func NewGlobalResolver() *GlobalResolver {
	instanceMap := new(map[*object.Reference]Proxy)
	return &GlobalResolver{
		globalInstanceMap: instanceMap,
	}
}

// GetObject reserve object by its reference and return its proxy.
func (r *GlobalResolver) GetObject(ref *object.Reference, classID string) (Proxy, error) {
	// TODO: check ref.Scope
	proxy, isExist := (*r.globalInstanceMap)[ref]
	if !isExist {
		return nil, fmt.Errorf("reference with address `%s` not found", ref)
	}

	if proxy.GetClassID() != classID {
		return nil, fmt.Errorf("instance type is not `%s`", classID)
	}
	return proxy, nil
}
