package updateapproves

import (
	"github.com/insolar/insolar/application/contract/updateapproves"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

// ClassReference to class of this contract
var ClassReference = core.NewRefFromBase58("")

// Contract proxy type
type UpdateApproves struct {
	Reference core.RecordRef
}

type ContractConstructorHolder struct {
	constructorName string
	argsSerialized  []byte
}

func (r *ContractConstructorHolder) AsChild(objRef core.RecordRef) *UpdateApproves {
	ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		panic(err)
	}
	return &UpdateApproves{Reference: ref}
}

func (r *ContractConstructorHolder) AsDelegate(objRef core.RecordRef) *UpdateApproves {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		panic(err)
	}
	return &UpdateApproves{Reference: ref}
}

// GetObject
func GetObject(ref core.RecordRef) (r *UpdateApproves) {
	return &UpdateApproves{Reference: ref}
}

func GetClass() core.RecordRef {
	return ClassReference
}

func GetImplementationFrom(object core.RecordRef) *UpdateApproves {
	ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
	if err != nil {
		panic(err)
	}
	return GetObject(ref)
}

func New(nodeRec *core.RecordRef, result updateapproves.ApproveResult, signature []byte) *ContractConstructorHolder {
	var args [3]interface{}
	args[0] = nodeRec
	args[1] = result
	args[2] = signature

	var argsSerialized []byte
	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	return &ContractConstructorHolder{constructorName: "New", argsSerialized: argsSerialized}
}

// GetReference
func (r *UpdateApproves) GetReference() core.RecordRef {
	return r.Reference
}

// GetClass
func (r *UpdateApproves) GetClass() core.RecordRef {
	return ClassReference
}

func (r *UpdateApproves) GetApproveResult() updateapproves.ApproveResult {
	var args [0]interface{}
	var argsSerialized []byte

	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}
	res, err := proxyctx.Current.RouteCall(r.Reference, true, "GetApproveResult", argsSerialized)
	if err != nil {
		panic(err)
	}

	ret := [1]interface{}{}
	var ret0 updateapproves.ApproveResult
	ret[0] = &ret0

	err = proxyctx.Current.Deserialize(res, &ret)
	if err != nil {
		panic(err)
	}
	return ret0
}
