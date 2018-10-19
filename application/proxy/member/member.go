package member

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

// ClassReference to class of this contract
var ClassReference = core.NewRefFromBase58("")

// Member holds proxy type
type Member struct {
	Reference core.RecordRef
}

// ContractConstructorHolder holds logic with object construction
type ContractConstructorHolder struct {
	constructorName string
	argsSerialized  []byte
}

// AsChild saves object as child
func (r *ContractConstructorHolder) AsChild(objRef core.RecordRef) *Member {
	ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		panic(err)
	}
	return &Member{Reference: ref}
}

// AsDelegate saves object as delegate
func (r *ContractConstructorHolder) AsDelegate(objRef core.RecordRef) *Member {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		panic(err)
	}
	return &Member{Reference: ref}
}

// GetObject returns proxy object
func GetObject(ref core.RecordRef) (r *Member) {
	return &Member{Reference: ref}
}

// GetClass returns reference to the class
func GetClass() core.RecordRef {
	return ClassReference
}

// GetImplementationFrom returns proxy to delegate of given type
func GetImplementationFrom(object core.RecordRef) *Member {
	ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
	if err != nil {
		panic(err)
	}
	return GetObject(ref)
}

// New is constructor
func New(name string, key string) *ContractConstructorHolder {
	var args [2]interface{}
	args[0] = name
	args[1] = key

	var argsSerialized []byte
	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	return &ContractConstructorHolder{constructorName: "New", argsSerialized: argsSerialized}
}

// GetReference returns reference of the object
func (r *Member) GetReference() core.RecordRef {
	return r.Reference
}

// GetClass returns reference to the class
func (r *Member) GetClass() core.RecordRef {
	return ClassReference
}

// GetName is proxy generated method
func (r *Member) GetName() (string, error) {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetName", argsSerialized)
	if err != nil {
		panic(err)
	}

	ret := [2]interface{}{}
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetNameNoWait is proxy generated method
func (r *Member) GetNameNoWait() {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "GetName", argsSerialized)
	if err != nil {
		panic(err)
	}
}

// GetPublicKey is proxy generated method
func (r *Member) GetPublicKey() (string, error) {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetPublicKey", argsSerialized)
	if err != nil {
		panic(err)
	}

	ret := [2]interface{}{}
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetPublicKeyNoWait is proxy generated method
func (r *Member) GetPublicKeyNoWait() {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "GetPublicKey", argsSerialized)
	if err != nil {
		panic(err)
	}
}

// AuthorizedCall is proxy generated method
func (r *Member) AuthorizedCall(ref core.RecordRef, delegate core.RecordRef, method string, params []byte, seed []byte, sign []byte) ([]byte, error) {
	var args [6]interface{}
	args[0] = ref
	args[1] = delegate
	args[2] = method
	args[3] = params
	args[4] = seed
	args[5] = sign

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "AuthorizedCall", argsSerialized)
	if err != nil {
		panic(err)
	}

	ret := [2]interface{}{}
	var ret0 []byte
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// AuthorizedCallNoWait is proxy generated method
func (r *Member) AuthorizedCallNoWait(ref core.RecordRef, delegate core.RecordRef, method string, params []byte, seed []byte, sign []byte) {
	var args [6]interface{}
	args[0] = ref
	args[1] = delegate
	args[2] = method
	args[3] = params
	args[4] = seed
	args[5] = sign

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "AuthorizedCall", argsSerialized)
	if err != nil {
		panic(err)
	}
}
