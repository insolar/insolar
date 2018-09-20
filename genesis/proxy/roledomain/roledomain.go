package roledomain

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/genesis/proxy/rolerecord"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

// Reference to class of this contract
var ClassReference = core.NewRefFromBase58("testRef")

// Contract proxy type
type RoleDomain struct {
	Reference core.RecordRef
}

type ContractHolder struct {
	data []byte
}

func (r *ContractHolder) AsChild(objRef core.RecordRef) *RoleDomain {
	ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.data)
	if err != nil {
		panic(err)
	}
	return &RoleDomain{Reference: ref}
}

func (r *ContractHolder) AsDelegate(objRef core.RecordRef) *RoleDomain {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.data)
	if err != nil {
		panic(err)
	}
	return &RoleDomain{Reference: ref}
}

// GetObject
func GetObject(ref core.RecordRef) (r *RoleDomain) {
	return &RoleDomain{Reference: ref}
}

func GetClass() core.RecordRef {
	return ClassReference
}

func GetImplementationFrom(object core.RecordRef) *RoleDomain {
	ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
	if err != nil {
		panic(err)
	}
	return GetObject(ref)
}

func NewRoleDomain() *ContractHolder {
	var args [0]interface{}

	var argsSerialized []byte
	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	data, err := proxyctx.Current.RouteConstructorCall(ClassReference, "NewRoleDomain", argsSerialized)
	if err != nil {
		panic(err)
	}

	return &ContractHolder{data: data}
}

// GetReference
func (r *RoleDomain) GetReference() core.RecordRef {
	return r.Reference
}

// GetClass
func (r *RoleDomain) GetClass() core.RecordRef {
	return ClassReference
}

func (r *RoleDomain) RegisterNode(pk string, role core.JetRole) core.RecordRef {
	var args [2]interface{}
	args[0] = pk
	args[1] = role

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, "RegisterNode", argsSerialized)
	if err != nil {
		panic(err)
	}

	resList := [1]interface{}{}
	var a0 core.RecordRef
	resList[0] = a0

	err = proxyctx.Current.Deserialize(res, &resList)
	if err != nil {
		panic(err)
	}

	return resList[0].(core.RecordRef)
}

func (r *RoleDomain) GetNodeRecord(ref core.RecordRef) *rolerecord.RoleRecord {
	var args [1]interface{}
	args[0] = ref

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, "GetNodeRecord", argsSerialized)
	if err != nil {
		panic(err)
	}

	resList := [1]interface{}{}
	var a0 *rolerecord.RoleRecord
	resList[0] = a0

	err = proxyctx.Current.Deserialize(res, &resList)
	if err != nil {
		panic(err)
	}

	return resList[0].(*rolerecord.RoleRecord)
}

func (r *RoleDomain) RemoveNode(nodeRef core.RecordRef) {
	var args [1]interface{}
	args[0] = nodeRef

	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	res, err := proxyctx.Current.RouteCall(r.Reference, "RemoveNode", argsSerialized)
	if err != nil {
		panic(err)
	}

	resList := [0]interface{}{}

	err = proxyctx.Current.Deserialize(res, &resList)
	if err != nil {
		panic(err)
	}

	return
}
