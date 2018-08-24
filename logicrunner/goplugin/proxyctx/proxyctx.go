package proxyctx

type ProxyHelper interface {
	RouteCall(ref string, method string, args []byte) ([]byte, error)
	Serialize(what interface{}, to *[]byte) error
	Deserialize(from []byte, into interface{}) error
}

// Current - hackish way to give proxies access to the current environment
var Current ProxyHelper
