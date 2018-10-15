package member

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

// ClassReference to class of this contract
var ClassReference = core.NewRefFromBase58("")

// Contract proxy type
type Member struct {
	Reference core.RecordRef
}

type ContractConstructorHolder struct {
	constructorName string
	argsSerialized  []byte
}

func (r *ContractConstructorHolder) AsChild(objRef core.RecordRef) *Member {
	ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		panic(err)
	}
	return &Member{Reference: ref}
}

func (r *ContractConstructorHolder) AsDelegate(objRef core.RecordRef) *Member {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		panic(err)
	}
	return &Member{Reference: ref}
}

// GetObject
func GetObject(ref core.RecordRef) (r *Member) {
	return &Member{Reference: ref}
}

func GetClass() core.RecordRef {
	return ClassReference
}

func GetImplementationFrom(object core.RecordRef) *Member {
	ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
	if err != nil {
		panic(err)
	}
	return GetObject(ref)
}

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

// GetReference
func (r *Member) GetReference() core.RecordRef {
	return r.Reference
}

// GetClass
func (r *Member) GetClass() core.RecordRef {
	return ClassReference
}

func (r *Member) GetName() string {
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

	ret := [1]interface{}{}
	var ret0 string
	ret[0] = &ret0

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return ret0
}

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

func (r *Member) GetPublicKey() string {
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

	ret := [1]interface{}{}
	var ret0 string
	ret[0] = &ret0

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return ret0
}

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

func (r *Member) AuthorizedCall(ref core.RecordRef, delegate core.RecordRef, method string, params []byte, seed []byte, sign []byte) ([]byte, *foundation.Error) {
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

	return ret0, ret1
}

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
