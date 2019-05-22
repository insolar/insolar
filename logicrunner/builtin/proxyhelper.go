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

package builtin

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

type ProxyHelper struct {
	methods logicrunner.RPCMethods
}

func NewProxyHelper(runner insolar.LogicRunner) *ProxyHelper {
	return nil
}

func (h *ProxyHelper) RouteCall(ref insolar.Reference, wait bool, immutable bool, method string, args []byte, proxyPrototype insolar.Reference) ([]byte, error) {
	panic("implement me")
}

func (h *ProxyHelper) SaveAsChild(parentRef, classRef insolar.Reference, constructorName string, argsSerialized []byte) (insolar.Reference, error) {
	panic("implement me")
}

func (h *ProxyHelper) GetObjChildrenIterator(head insolar.Reference, prototype insolar.Reference, iteratorID string) (*proxyctx.ChildrenTypedIterator, error) {
	panic("implement me")
}

func (h *ProxyHelper) SaveAsDelegate(parentRef, classRef insolar.Reference, constructorName string, argsSerialized []byte) (insolar.Reference, error) {
	panic("implement me")
}

func (h *ProxyHelper) GetDelegate(object, ofType insolar.Reference) (insolar.Reference, error) {
	panic("implement me")
}

func (h *ProxyHelper) DeactivateObject(object insolar.Reference) error {
	panic("implement me")
}

func (h *ProxyHelper) Serialize(what interface{}, to *[]byte) error {
	panic("implement me")
}

func (h *ProxyHelper) Deserialize(from []byte, into interface{}) error {
	panic("implement me")
}

func (h *ProxyHelper) MakeErrorSerializable(err error) error {
	panic("implement me")
}
