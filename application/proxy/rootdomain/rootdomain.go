package rootdomain

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

// PrototypeReference to prototype of this contract
// error checking hides in generator
var PrototypeReference, _ = core.NewRefFromBase58("11112VdzPPzfsmgTMBDFpo6KRWbFXX5hUzBihadPPKW.11111111111111111111111111111111")

// RootDomain holds proxy type
type RootDomain struct {
	Reference core.RecordRef
	Prototype core.RecordRef
	Code      core.RecordRef
}

// ContractConstructorHolder holds logic with object construction
type ContractConstructorHolder struct {
	constructorName string
	argsSerialized  []byte
}

// AsChild saves object as child
func (r *ContractConstructorHolder) AsChild(objRef core.RecordRef) (*RootDomain, error) {
	ref, err := proxyctx.Current.SaveAsChild(objRef, *PrototypeReference, r.constructorName, r.argsSerialized)
	if err != nil {
		return nil, err
	}
	return &RootDomain{Reference: ref}, nil
}

// AsDelegate saves object as delegate
func (r *ContractConstructorHolder) AsDelegate(objRef core.RecordRef) (*RootDomain, error) {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, *PrototypeReference, r.constructorName, r.argsSerialized)
	if err != nil {
		return nil, err
	}
	return &RootDomain{Reference: ref}, nil
}

// GetObject returns proxy object
func GetObject(ref core.RecordRef) (r *RootDomain) {
	return &RootDomain{Reference: ref}
}

// GetPrototype returns reference to the prototype
func GetPrototype() core.RecordRef {
	return *PrototypeReference
}

// GetImplementationFrom returns proxy to delegate of given type
func GetImplementationFrom(object core.RecordRef) (*RootDomain, error) {
	ref, err := proxyctx.Current.GetDelegate(object, *PrototypeReference)
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

// GetPrototype returns reference to the code
func (r *RootDomain) GetPrototype() (core.RecordRef, error) {
	if r.Prototype.IsEmpty() {
		ret := [2]interface{}{}
		var ret0 core.RecordRef
		ret[0] = &ret0
		var ret1 *foundation.Error
		ret[1] = &ret1

		res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetPrototype", make([]byte, 0), *PrototypeReference)
		if err != nil {
			return ret0, err
		}

		err = proxyctx.Current.Deserialize(res, &ret)
		if err != nil {
			return ret0, err
		}

		if ret1 != nil {
			return ret0, ret1
		}

		r.Prototype = ret0
	}

	return r.Prototype, nil

}

// GetCode returns reference to the code
func (r *RootDomain) GetCode() (core.RecordRef, error) {
	if r.Code.IsEmpty() {
		ret := [2]interface{}{}
		var ret0 core.RecordRef
		ret[0] = &ret0
		var ret1 *foundation.Error
		ret[1] = &ret1

		res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetCode", make([]byte, 0), *PrototypeReference)
		if err != nil {
			return ret0, err
		}

		err = proxyctx.Current.Deserialize(res, &ret)
		if err != nil {
			return ret0, err
		}

		if ret1 != nil {
			return ret0, ret1
		}

		r.Code = ret0
	}

	return r.Code, nil
}

// CreateMember is proxy generated method
func (r *RootDomain) CreateMember(name string, key string) (string, error) {
	var args [2]interface{}
	args[0] = name
	args[1] = key

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 string
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "CreateMember", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
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

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "CreateMember", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetRootMemberRef is proxy generated method
func (r *RootDomain) GetRootMemberRef() (*core.RecordRef, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 *core.RecordRef
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetRootMemberRef", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// GetRootMemberRefNoWait is proxy generated method
func (r *RootDomain) GetRootMemberRefNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "GetRootMemberRef", argsSerialized, *PrototypeReference)
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

	ret := [2]interface{}{}
	var ret0 []byte
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "DumpUserInfo", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
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

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "DumpUserInfo", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// DumpAllUsers is proxy generated method
func (r *RootDomain) DumpAllUsers() ([]byte, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 []byte
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "DumpAllUsers", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
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

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "DumpAllUsers", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// Info is proxy generated method
func (r *RootDomain) Info() (interface{}, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 interface{}
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "Info", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
	}

	if ret1 != nil {
		return ret0, ret1
	}
	return ret0, nil
}

// InfoNoWait is proxy generated method
func (r *RootDomain) InfoNoWait() error {
	var args [0]interface{}

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return err
	}

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "Info", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}

// GetNodeDomainRef is proxy generated method
func (r *RootDomain) GetNodeDomainRef() (core.RecordRef, error) {
	var args [0]interface{}

	var argsSerialized []byte

	ret := [2]interface{}{}
	var ret0 core.RecordRef
	ret[0] = &ret0
	var ret1 *foundation.Error
	ret[1] = &ret1

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		return ret0, err
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetNodeDomainRef", argsSerialized, *PrototypeReference)
	if err != nil {
		return ret0, err
	}

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		return ret0, err
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

	_, err = proxyctx.Current.RouteCall(r.Reference, false, "GetNodeDomainRef", argsSerialized, *PrototypeReference)
	if err != nil {
		return err
	}

	return nil
}
