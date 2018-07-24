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
	globalInstanceMap *map[*object.Reference]object.Proxy
}

// newGlobalResolver creates new globalResolver instance.
// TODO: pass map?
func newGlobalResolver() *globalResolver {
	instanceMap := make(map[*object.Reference]object.Proxy)
	return &globalResolver{
		globalInstanceMap: &instanceMap,
	}
}

// GetObject reserve object by its reference and return its proxy.
func (r *globalResolver) GetObject(ref *object.Reference, classID string) (object.Proxy, error) {
	// TODO: check ref.Scope
	proxy, isExist := (*r.globalInstanceMap)[ref]
	if !isExist {
		return nil, fmt.Errorf("reference with address `%s` not found", ref)
	}

	if proxy.GetClassID() != classID {
		return nil, fmt.Errorf("instance class is not `%s`", classID)
	}
	return proxy, nil
}

// GlobalResolver is a public globalResolver instance for reserving all global references
var GlobalResolver = newGlobalResolver()
