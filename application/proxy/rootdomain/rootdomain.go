package rootdomain

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

// ClassReference to class of this contract
var ClassReference = core.NewRefFromBase58("")

// RootDomain holds proxy type
type RootDomain struct {
	Reference core.RecordRef
}

// ContractConstructorHolder holds logic with object construction
type ContractConstructorHolder struct {
	constructorName string
	argsSerialized  []byte
}

// AsChild saves object as child
func (r *ContractConstructorHolder) AsChild(objRef core.RecordRef) *RootDomain {
	ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		panic(err)
	}
	return &RootDomain{Reference: ref}
}

// AsDelegate saves object as delegate
func (r *ContractConstructorHolder) AsDelegate(objRef core.RecordRef) *RootDomain {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		panic(err)
	}
	return &RootDomain{Reference: ref}
}

// GetObject returns proxy object
func GetObject(ref core.RecordRef) (r *RootDomain) {
	return &RootDomain{Reference: ref}
}

// GetClass returns reference to the class
func GetClass() core.RecordRef {
	return ClassReference
}

// GetImplementationFrom returns proxy to delegate of given type
func GetImplementationFrom(object core.RecordRef) *RootDomain {
	ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
	if err != nil {
		panic(err)
	}
	return GetObject(ref)
}

// NewRootDomain is constructor
func NewRootDomain() *ContractConstructorHolder {
	var args [0]interface{}

	var argsSerialized []byte
	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	return &ContractConstructorHolder{constructorName: "NewRootDomain", argsSerialized: argsSerialized}
}

// GetReference returns reference of the object
func (r *RootDomain) GetReference() core.RecordRef {
	return r.Reference
}

// GetClass returns reference to the class
func (r *RootDomain) GetClass() core.RecordRef {
	return ClassReference
}

// RegisterNode is proxy generated method
func (r *RootDomain) RegisterNode(publicKey string, role string) string {
	var args [2]interface{}
	args[0] = publicKey
	args[1] = role

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "RegisterNode", argsSerialized)
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

// RegisterNodeNoWait is proxy generated method
func (r *RootDomain) RegisterNodeNoWait(publicKey string, role string) {
	var args [2]interface{}
	args[0] = publicKey
	args[1] = role

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "RegisterNode", argsSerialized)
	if err != nil {
		panic(err)
	}
}

// Authorize is proxy generated method
func (r *RootDomain) Authorize() (string, core.NodeRole, string) {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "Authorize", argsSerialized)
	if err != nil {
		panic(err)
	}

	ret := [3]interface{}{}
	var ret0 string
	ret[0] = &ret0
	var ret1 core.NodeRole
	ret[1] = &ret1
	var ret2 string
	ret[2] = &ret2

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return ret0, ret1, ret2
}

// AuthorizeNoWait is proxy generated method
func (r *RootDomain) AuthorizeNoWait() {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "Authorize", argsSerialized)
	if err != nil {
		panic(err)
	}
}

// CreateMember is proxy generated method
func (r *RootDomain) CreateMember(name string, key string) string {
	var args [2]interface{}
	args[0] = name
	args[1] = key

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "CreateMember", argsSerialized)
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

// CreateMemberNoWait is proxy generated method
func (r *RootDomain) CreateMemberNoWait(name string, key string) {
	var args [2]interface{}
	args[0] = name
	args[1] = key

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "CreateMember", argsSerialized)
	if err != nil {
		panic(err)
	}
}

// GetBalance is proxy generated method
func (r *RootDomain) GetBalance(reference string) uint {
	var args [1]interface{}
	args[0] = reference

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetBalance", argsSerialized)
	if err != nil {
		panic(err)
	}

	ret := [1]interface{}{}
	var ret0 uint
	ret[0] = &ret0

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return ret0
}

// GetBalanceNoWait is proxy generated method
func (r *RootDomain) GetBalanceNoWait(reference string) {
	var args [1]interface{}
	args[0] = reference

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "GetBalance", argsSerialized)
	if err != nil {
		panic(err)
	}
}

// SendMoney is proxy generated method
func (r *RootDomain) SendMoney(from string, to string, amount uint) bool {
	var args [3]interface{}
	args[0] = from
	args[1] = to
	args[2] = amount

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "SendMoney", argsSerialized)
	if err != nil {
		panic(err)
	}

	ret := [1]interface{}{}
	var ret0 bool
	ret[0] = &ret0

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return ret0
}

// SendMoneyNoWait is proxy generated method
func (r *RootDomain) SendMoneyNoWait(from string, to string, amount uint) {
	var args [3]interface{}
	args[0] = from
	args[1] = to
	args[2] = amount

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "SendMoney", argsSerialized)
	if err != nil {
		panic(err)
	}
}

// DumpUserInfo is proxy generated method
func (r *RootDomain) DumpUserInfo(reference string) []byte {
	var args [1]interface{}
	args[0] = reference

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "DumpUserInfo", argsSerialized)
	if err != nil {
		panic(err)
	}

	ret := [1]interface{}{}
	var ret0 []byte
	ret[0] = &ret0

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return ret0
}

// DumpUserInfoNoWait is proxy generated method
func (r *RootDomain) DumpUserInfoNoWait(reference string) {
	var args [1]interface{}
	args[0] = reference

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "DumpUserInfo", argsSerialized)
	if err != nil {
		panic(err)
	}
}

// DumpAllUsers is proxy generated method
func (r *RootDomain) DumpAllUsers() []byte {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "DumpAllUsers", argsSerialized)
	if err != nil {
		panic(err)
	}

	ret := [1]interface{}{}
	var ret0 []byte
	ret[0] = &ret0

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	return ret0
}

// DumpAllUsersNoWait is proxy generated method
func (r *RootDomain) DumpAllUsersNoWait() {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "DumpAllUsers", argsSerialized)
	if err != nil {
		panic(err)
	}
}
