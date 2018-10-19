package rootdomain

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
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
func (r *ContractConstructorHolder) AsChild(objRef core.RecordRef) (*RootDomain, error) {
	ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		return nil, err
	}
	return &RootDomain{Reference: ref}, nil
}

// AsDelegate saves object as delegate
func (r *ContractConstructorHolder) AsDelegate(objRef core.RecordRef) (*RootDomain, error) {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		return nil, err
	}
	return &RootDomain{Reference: ref}, nil
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
func GetImplementationFrom(object core.RecordRef) (*RootDomain, error) {
	ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
	if err != nil {
		return nil, err
	}
	return GetObject(ref), nil
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
func (r *RootDomain) RegisterNode(publicKey string, role string) (string, error) {
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

// RegisterNodeNoWait is proxy generated method
func (r *RootDomain) RegisterNodeNoWait(publicKey string, role string) error {
	var args [2]interface{}
	args[0] = publicKey
	args[1] = role

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "RegisterNode", argsSerialized)
	if err != nil {
		return err
	}

	return nil
}

// Authorize is proxy generated method
func (r *RootDomain) Authorize() (string, core.NodeRole, error) {
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
	var ret2 *foundation.Error
	ret[2] = &ret2

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}

	if ret2 != nil {
		return ret0, ret1, ret2
	}
	return ret0, ret1, nil
}

// AuthorizeNoWait is proxy generated method
func (r *RootDomain) AuthorizeNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "Authorize", argsSerialized)
	if err != nil {
		return err
	}

	return nil
}

// CreateMember is proxy generated method
func (r *RootDomain) CreateMember(name string, key string) (string, error) {
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

// CreateMemberNoWait is proxy generated method
func (r *RootDomain) CreateMemberNoWait(name string, key string) error {
	var args [2]interface{}
	args[0] = name
	args[1] = key

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "CreateMember", argsSerialized)
	if err != nil {
		return err
	}

	return nil
}

// GetBalance is proxy generated method
func (r *RootDomain) GetBalance(reference string) (uint, error) {
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

	ret := [2]interface{}{}
	var ret0 uint
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

// GetBalanceNoWait is proxy generated method
func (r *RootDomain) GetBalanceNoWait(reference string) error {
	var args [1]interface{}
	args[0] = reference

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "GetBalance", argsSerialized)
	if err != nil {
		return err
	}

	return nil
}

// SendMoney is proxy generated method
func (r *RootDomain) SendMoney(from string, to string, amount uint) (bool, error) {
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

	ret := [2]interface{}{}
	var ret0 bool
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

// SendMoneyNoWait is proxy generated method
func (r *RootDomain) SendMoneyNoWait(from string, to string, amount uint) error {
	var args [3]interface{}
	args[0] = from
	args[1] = to
	args[2] = amount

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "SendMoney", argsSerialized)
	if err != nil {
		return err
	}

	return nil
}

// DumpUserInfo is proxy generated method
func (r *RootDomain) DumpUserInfo(reference string) ([]byte, error) {
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

// DumpUserInfoNoWait is proxy generated method
func (r *RootDomain) DumpUserInfoNoWait(reference string) error {
	var args [1]interface{}
	args[0] = reference

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "DumpUserInfo", argsSerialized)
	if err != nil {
		return err
	}

	return nil
}

// DumpAllUsers is proxy generated method
func (r *RootDomain) DumpAllUsers() ([]byte, error) {
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

// DumpAllUsersNoWait is proxy generated method
func (r *RootDomain) DumpAllUsersNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "DumpAllUsers", argsSerialized)
	if err != nil {
		return err
	}

	return nil
}

// GetNodeDomainRef is proxy generated method
func (r *RootDomain) GetNodeDomainRef() (core.RecordRef, error) {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetNodeDomainRef", argsSerialized)
	if err != nil {
		panic(err)
	}

	ret := [2]interface{}{}
	var ret0 core.RecordRef
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

// GetNodeDomainRefNoWait is proxy generated method
func (r *RootDomain) GetNodeDomainRefNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "GetNodeDomainRef", argsSerialized)
	if err != nil {
		return err
	}

	return nil
}
