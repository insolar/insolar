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

// childResolver is resolver for ChildScope references.
type childResolver struct {
	parent Parent
}

// newChildResolver creates new childResolver instance.
func newChildResolver(parent Parent) *childResolver {
	return &childResolver{
		parent: parent,
	}
}

// GetObject resolves object by its reference and return its proxy.
func (r *childResolver) GetObject(reference interface{}, cls interface{}) (interface{}, error) {
	ref, ok := reference.(Reference)
	if !ok {
		return nil, fmt.Errorf("reference is not Reference class object")
	}
	obj, err := r.parent.GetChild(ref.GetRecord())
	if err != nil {
		return nil, err
	}

	proxy, ok := obj.(Proxy)
	if !ok {
		return nil, fmt.Errorf("object is not Proxy")
	}

	classID, ok := cls.(string)
	if !ok {
		return nil, fmt.Errorf("classID is not string")
	}
	if proxy.GetClassID() != classID {
		return nil, fmt.Errorf("instance class is not `%s`", classID)
	}
	proxy.SetReference(ref)
	return proxy, nil
}
