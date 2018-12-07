/*
 *    Copyright 2018 Insolar
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

package proxyctx

import (
	"github.com/insolar/insolar/core"
)

// ProxyHelper interface with methods that are needed by contract proxies
type ProxyHelper interface {
	RouteCall(ref core.RecordRef, wait bool, method string, args []byte) ([]byte, error)
	SaveAsChild(parentRef, classRef core.RecordRef, constructorName string, argsSerialized []byte) (core.RecordRef, error)
	GetObjChildren(head core.RecordRef, prototype core.RecordRef) ([]core.RecordRef, error)
	GetObjChildrenIterator(head core.RecordRef, prototype core.RecordRef, iteratorId string) (*ChildrenTypedIterator, error)
	SaveAsDelegate(parentRef, classRef core.RecordRef, constructorName string, argsSerialized []byte) (core.RecordRef, error)
	GetDelegate(object, ofType core.RecordRef) (core.RecordRef, error)
	DeactivateObject(object core.RecordRef) error
	Serialize(what interface{}, to *[]byte) error
	Deserialize(from []byte, into interface{}) error
	MakeErrorSerializable(error) error
}

// Current - hackish way to give proxies access to the current environment
var Current ProxyHelper

type ChildrenTypedIterator struct {
	Parent         core.RecordRef
	ChildPrototype core.RecordRef // only child of specified prototype, if childPrototype.IsEmpty - ignored

	IteratorId string           // map key to iterators slice in logicrunner service
	Buff       []core.RecordRef // bucket of objects from previous RPC call to service
	buffIndex  int              // current element
	CanFetch   bool             // if true, we can call RPC again and get new objects
}

func (oi *ChildrenTypedIterator) HasNext() bool {
	return oi.hasInBuffer() || oi.CanFetch
}

func (oi *ChildrenTypedIterator) Next() (core.RecordRef, error) {
	if !oi.hasInBuffer() && oi.CanFetch {
		err := oi.fetch()
		if err != nil {
			return core.RecordRef{}, err
		}
	}

	return oi.nextFromBuffer(), nil
}

func (oi *ChildrenTypedIterator) hasInBuffer() bool {
	return oi.buffIndex < len(oi.Buff)
}

func (oi *ChildrenTypedIterator) nextFromBuffer() core.RecordRef {
	if !oi.hasInBuffer() {
		return core.RecordRef{}
	}

	result := oi.Buff[oi.buffIndex]
	oi.buffIndex++
	return result
}

func (oi *ChildrenTypedIterator) fetch() error {
	oi.buffIndex = 0
	oi.CanFetch = false

	temp, err := Current.GetObjChildrenIterator(oi.Parent, oi.ChildPrototype, oi.IteratorId)
	if err != nil {
		return err
	}
	oi.Buff = temp.Buff
	oi.IteratorId = temp.IteratorId
	oi.CanFetch = temp.CanFetch

	return nil
}
