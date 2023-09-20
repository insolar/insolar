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
		immutable bool, saga bool,
		method string, args []byte, proxyPrototype insolar.Reference,
	) (result []byte, err error)
	SaveAsChild(
		parentRef, classRef insolar.Reference, constructorName string, argsSerialized []byte,
	) (result []byte, err error)
	DeactivateObject(object insolar.Reference) error
	MakeErrorSerializable(error) error
}

// CurrentProxyCtx - hackish way to give proxies access to the current environment.
var CurrentProxyCtx ProxyHelper
