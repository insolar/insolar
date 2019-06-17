//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package common

import (
	"github.com/insolar/insolar/insolar"
)

// ProxyHelper interface with methods that are needed by contract proxies
type ProxyHelper interface {
	RouteCall(ref insolar.Reference, wait bool, immutable bool, method string, args []byte, proxyPrototype insolar.Reference) ([]byte, error)
	SaveAsChild(parentRef, classRef insolar.Reference, constructorName string, argsSerialized []byte) (insolar.Reference, error)
	GetObjChildrenIterator(head insolar.Reference, prototype insolar.Reference, iteratorID string) (*ChildrenTypedIterator, error)
	SaveAsDelegate(parentRef, classRef insolar.Reference, constructorName string, argsSerialized []byte) (insolar.Reference, error)
	GetDelegate(object, ofType insolar.Reference) (insolar.Reference, error)
	DeactivateObject(object insolar.Reference) error
	Serialize(what interface{}, to *[]byte) error
	Deserialize(from []byte, into interface{}) error
	MakeErrorSerializable(error) error
}

// Current - hackish way to give proxies access to the current environment
var CurrentProxyCtx ProxyHelper

// ChildrenTypedIterator iterator over children of object with specified type
// it uses cache on insolard service side, provided by IteratorID
type ChildrenTypedIterator struct {
	Parent         insolar.Reference
	ChildPrototype insolar.Reference // only child of specified prototype, if childPrototype.IsEmpty - ignored

	IteratorID string              // map key to iterators slice in logicrunner service
	Buff       []insolar.Reference // bucket of objects from previous RPC call to service
	buffIndex  int                 // current element
	CanFetch   bool                // if true, we can call RPC again and get new objects
}

// HasNext return true if iterator has element in cache or can fetch data again
func (oi *ChildrenTypedIterator) HasNext() bool {
	return oi.hasInBuffer() || oi.CanFetch
}

// Next return next element from iterator cache or fetching new from service
// return error only if fetch() fails
func (oi *ChildrenTypedIterator) Next() (insolar.Reference, error) {
	if !oi.hasInBuffer() && oi.CanFetch {
		err := oi.fetch()
		if err != nil {
			oi.CanFetch = false
			return insolar.Reference{}, err
		}
	}

	return oi.nextFromBuffer(), nil
}

func (oi *ChildrenTypedIterator) hasInBuffer() bool {
	return oi.buffIndex < len(oi.Buff)
}

func (oi *ChildrenTypedIterator) nextFromBuffer() insolar.Reference {
	if !oi.hasInBuffer() {
		return insolar.Reference{}
	}

	result := oi.Buff[oi.buffIndex]
	oi.buffIndex++
	return result
}

func (oi *ChildrenTypedIterator) fetch() error {
	oi.buffIndex = 0
	oi.CanFetch = false
	oi.Buff = nil

	temp, err := CurrentProxyCtx.GetObjChildrenIterator(oi.Parent, oi.ChildPrototype, oi.IteratorID)
	if err != nil {
		oi.IteratorID = ""
		return err
	}
	oi.Buff = temp.Buff
	oi.IteratorID = temp.IteratorID
	oi.CanFetch = temp.CanFetch

	return nil
}
