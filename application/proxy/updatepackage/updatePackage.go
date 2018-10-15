package updatepackage

import (
	"github.com/insolar/insolar/application/contract/updateapproves"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
	"github.com/insolar/insolar/updater/request"
)

// ClassReference to class of this contract
var ClassReference = core.NewRefFromBase58("")

// Contract proxy type
type UpdatePackage struct {
	Reference core.RecordRef
}

type ContractConstructorHolder struct {
	constructorName string
	argsSerialized  []byte
}

func (r *ContractConstructorHolder) AsChild(objRef core.RecordRef) *UpdatePackage {
	ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		panic(err)
	}
	return &UpdatePackage{Reference: ref}
}

func (r *ContractConstructorHolder) AsDelegate(objRef core.RecordRef) *UpdatePackage {
	ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.constructorName, r.argsSerialized)
	if err != nil {
		panic(err)
	}
	return &UpdatePackage{Reference: ref}
}

// GetObject
func GetObject(ref core.RecordRef) (r *UpdatePackage) {
	return &UpdatePackage{Reference: ref}
}

func GetClass() core.RecordRef {
	return ClassReference
}

func GetImplementationFrom(object core.RecordRef) *UpdatePackage {
	ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
	if err != nil {
		panic(err)
	}
	return GetObject(ref)
}

func New(uv *request.Version, nodes map[core.RecordRef]*updateapproves.UpdateApproves, consensusCN int, currentCN int) *ContractConstructorHolder {
	var args [4]interface{}
	args[0] = uv
	args[1] = nodes
	args[2] = consensusCN
	args[2] = currentCN

	var argsSerialized []byte
	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}

	return &ContractConstructorHolder{constructorName: "New", argsSerialized: argsSerialized}
}

// GetReference
func (r *UpdatePackage) GetReference() core.RecordRef {
	return r.Reference
}

// GetClass
func (r *UpdatePackage) GetClass() core.RecordRef {
	return ClassReference
}

func (r *UpdatePackage) GetApproveResult() updateapproves.ApproveResult {
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
