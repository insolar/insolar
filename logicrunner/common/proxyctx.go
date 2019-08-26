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
	SystemError
	Serializer
	RouteCall(
		ref insolar.Reference,
		wait bool, immutable bool, saga bool,
		method string, args []byte, proxyPrototype insolar.Reference,
	) (result []byte, err error)
	SaveAsChild(
		parentRef, classRef insolar.Reference, constructorName string, argsSerialized []byte,
	) (objRef *insolar.Reference, result []byte, err error)
	DeactivateObject(object insolar.Reference) error
	MakeErrorSerializable(error) error
}

// CurrentProxyCtx - hackish way to give proxies access to the current environment. Also,
// to avoid compiling in whole Insolar platform into every contract based on GoPlugin.
var CurrentProxyCtx ProxyHelper
