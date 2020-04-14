// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package preprocessor

var proxyTmpl = `
// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// Code generated by insgocc. DO NOT EDIT.
// source template in logicrunner/preprocessor/templates

package {{ .PackageName }}

import (
{{- range $import, $i := .Imports }}
{{ $import }}
{{- end }}
)

{{ range $typeStruct := .Types }}
{{- $typeStruct }}
{{ end }}

// PrototypeReference to prototype of this contract
// error checking hides in generator
var PrototypeReference, _ = insolar.NewObjectReferenceFromString("{{ .ClassReference }}")


// {{ .ContractType }} holds proxy type
type {{ .ContractType }} struct {
Reference insolar.Reference
Prototype insolar.Reference
Code insolar.Reference
}

// ContractConstructorHolder holds logic with object construction
type ContractConstructorHolder struct {
	constructorName string
	argsSerialized []byte
}

// AsChild saves object as child
func (r *ContractConstructorHolder) AsChild(objRef insolar.Reference) (*{{ .ContractType }}, error) {
ret, err := common.CurrentProxyCtx.SaveAsChild(objRef, *PrototypeReference, r.constructorName, r.argsSerialized)
if err != nil {
return nil, err
}

var ref insolar.Reference
var constructorError *foundation.Error
resultContainer := foundation.Result{
Returns: []interface{}{ &ref, &constructorError },
}
err = common.CurrentProxyCtx.Deserialize(ret, &resultContainer)
if err != nil {
return nil, err
}

if resultContainer.Error != nil {
return nil, resultContainer.Error
}

if constructorError != nil {
return nil, constructorError
}

return &{{ .ContractType }}{Reference: ref}, nil
}

// GetObject returns proxy object
func GetObject(ref insolar.Reference) *{{ .ContractType }} {
if !ref.IsObjectReference() {
return nil
}
return &{{ .ContractType }}{Reference: ref}
}

// GetPrototype returns reference to the prototype
func GetPrototype() insolar.Reference {
	return *PrototypeReference
}

{{ range $func := .ConstructorsProxies }}
// {{ $func.Name }} is constructor
func {{ $func.Name }}( {{ $func.Arguments }} ) *ContractConstructorHolder {
{{ $func.InitArgs }}

var argsSerialized []byte
err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
if err != nil {
panic(err)
}

return &ContractConstructorHolder{constructorName: "{{ $func.Name }}", argsSerialized: argsSerialized}
}
{{ end }}

// GetReference returns reference of the object
func (r *{{ $.ContractType }}) GetReference() insolar.Reference {
return r.Reference
}

// GetPrototype returns reference to the code
func (r *{{ $.ContractType }}) GetPrototype() (insolar.Reference, error) {
if r.Prototype.IsEmpty() {
ret := [2]interface{}{}
var ret0 insolar.Reference
ret[0] = &ret0
var ret1 *foundation.Error
ret[1] = &ret1

res, err := common.CurrentProxyCtx.RouteCall(r.Reference, false, false, "GetPrototype", make([]byte, 0), *PrototypeReference)
if err != nil {
return ret0, err
}

err = common.CurrentProxyCtx.Deserialize(res, &ret)
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
func (r *{{ $.ContractType }}) GetCode() (insolar.Reference, error) {
if r.Code.IsEmpty() {
ret := [2]interface{}{}
var ret0 insolar.Reference
ret[0] = &ret0
var ret1 *foundation.Error
ret[1] = &ret1

res, err := common.CurrentProxyCtx.RouteCall(r.Reference, false, false, "GetCode", make([]byte, 0), *PrototypeReference)
if err != nil {
return ret0, err
}

err = common.CurrentProxyCtx.Deserialize(res, &ret)
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

{{ range $method := .MethodsProxies }}
// {{ $method.Name }} is proxy generated method
func (r *{{ $.ContractType }}) {{ $method.Name }}{{if $method.Immutable}}AsMutable{{end}}( {{ $method.Arguments }} ) ( {{ $method.ResultsTypes }} ) {
{{ $method.InitArgs }}
var argsSerialized []byte

{{ $method.ResultZeroList }}

err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
if err != nil {
return {{ $method.ResultsWithErr }}
}

{{/* Saga call doesn't has a reply (it's nil), thus we shouldn't try to deserialize it. */}}
{{if $method.SagaInfo.IsSaga }}
_, err = common.CurrentProxyCtx.RouteCall(r.Reference, false, {{ $method.SagaInfo.IsSaga }}, "{{ $method.Name }}", argsSerialized, *PrototypeReference)
if err != nil {
return {{ $method.ResultsWithErr }}
}
{{else}}
res, err := common.CurrentProxyCtx.RouteCall(r.Reference, false, {{ $method.SagaInfo.IsSaga }}, "{{ $method.Name }}", argsSerialized, *PrototypeReference)
if err != nil {
return {{ $method.ResultsWithErr }}
}

resultContainer := foundation.Result{
Returns: ret,
}
err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
if err != nil {
return {{ $method.ResultsWithErr }}
}
if resultContainer.Error != nil {
err = resultContainer.Error
return {{ $method.ResultsWithErr }}
}
if {{ $method.ErrorVar }} != nil {
return {{ $method.Results }}
}
{{end -}}

return {{ $method.ResultsNilError }}
}

{{if not $method.SagaInfo.IsSaga}}

// {{ $method.Name }}AsImmutable is proxy generated method
func (r *{{ $.ContractType }}) {{ $method.Name }}{{if not $method.Immutable}}AsImmutable{{end}}( {{ $method.Arguments }} ) ( {{ $method.ResultsTypes }} ) {
{{ $method.InitArgs }}
var argsSerialized []byte

{{ $method.ResultZeroList }}

err := common.CurrentProxyCtx.Serialize(args, &argsSerialized)
if err != nil {
return {{ $method.ResultsWithErr }}
}

res, err := common.CurrentProxyCtx.RouteCall(r.Reference, true, false, "{{ $method.Name }}", argsSerialized, *PrototypeReference)
if err != nil {
return {{ $method.ResultsWithErr }}
}

resultContainer := foundation.Result{
Returns: ret,
}
err = common.CurrentProxyCtx.Deserialize(res, &resultContainer)
if err != nil {
return {{ $method.ResultsWithErr }}
}
if resultContainer.Error != nil {
err = resultContainer.Error
return {{ $method.ResultsWithErr }}
}
if {{ $method.ErrorVar }} != nil {
return {{ $method.Results }}
}
return {{ $method.ResultsNilError }}
}
{{ end }}
{{ end }}
`
var wrapperTmpl = `
// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// Code generated by insgocc. DO NOT EDIT.
// source template in logicrunner/preprocessor/templates

package {{ .Package }}

import (
{{- range $import, $i := .Imports }}
	{{ $import }}
{{- end }}
)

const PanicIsLogicalError = {{ .PanicIsLogicalError }}

func INS_META_INFO() ([] map[string]string) {
	result := make([]map[string] string, 0)
	{{ range $method := .Methods }}
		{{ if $method.SagaInfo.IsSaga }}
		{
		info := make(map[string] string, 3)
		info["Type"] = "SagaInfo"
		info["MethodName"] = "{{ $method.Name }}"
		info["RollbackMethodName"] = "{{ $method.SagaInfo.RollbackMethodName }}"
		result = append(result, info)
		}
		{{end}}
	{{end}}
	return result
}

func INSMETHOD_GetCode(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	self := new({{ $.ContractType }})

	if len(object) == 0 {
		return nil, nil, &foundation.Error{S: "[ Fake GetCode ] ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{ S: "[ Fake GetCode ] ( Generated Method ) Can't deserialize args.Data: " + err.Error() }
		return nil, nil, e
	}

	state := []byte{}
	err = ph.Serialize(self, &state)
	if err != nil {
		return nil, nil, err
	}

	ret := []byte{}
	err = ph.Serialize([]interface{} { self.GetCode().Bytes() }, &ret)

	return state, ret, err
}

func INSMETHOD_GetPrototype(object []byte, data []byte) ([]byte, []byte, error) {
	ph := common.CurrentProxyCtx
	self := new({{ $.ContractType }})

	if len(object) == 0 {
		return nil, nil, &foundation.Error{ S: "[ Fake GetPrototype ] ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &foundation.Error{ S: "[ Fake GetPrototype ] ( Generated Method ) Can't deserialize args.Data: " + err.Error() }
		return nil, nil, e
	}

	state := []byte{}
	err = ph.Serialize(self, &state)
	if err != nil {
		return nil, nil, err
	}

	ret := []byte{}
	err = ph.Serialize([]interface{} { self.GetPrototype().Bytes() }, &ret)

	return state, ret, err
}

{{ range $method := .Methods }}
func INSMETHOD_{{ $method.Name }}(object []byte, data []byte) (newState []byte, result []byte, err error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)

	self := new({{ $.ContractType }})

	if len(object) == 0 {
		err = &foundation.Error{ S: "[ Fake{{ $method.Name }} ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
		return
	}

	err = ph.Deserialize(object, self)
	if err != nil {
		err = &foundation.Error{ S: "[ Fake{{ $method.Name }} ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error() }
		return
	}

	{{ $method.ArgumentsZeroList }}
	err = ph.Deserialize(data, &args)
	if err != nil {
		err = &foundation.Error{ S: "[ Fake{{ $method.Name }} ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error() }
		return
	}

	{{ $method.ResultDefinitions }}

	serializeResults := func() error {
		return ph.Serialize(
			foundation.Result{Returns:[]interface{}{ {{ $method.Results }} }},
			&result,
		)
	}

	needRecover := true
	defer func() {
		if !needRecover {
			return
		}
		if r := recover(); r != nil {
			recoveredError := errors.Wrap(errors.Errorf("%v", r), "Failed to execute method (panic)")
			recoveredError = ph.MakeErrorSerializable(recoveredError)

			if PanicIsLogicalError {
				ret{{ $method.LastErrorInRes }} = recoveredError

				newState = object
				err = serializeResults()
				if err == nil {
				    newState = object
				}
			} else {
				err = recoveredError
			}
		}
	}()

	{{ $method.Results }} = self.{{ $method.Name }}( {{ $method.Arguments }} )

	needRecover = false

	if ph.GetSystemError() != nil {
		return nil, nil, ph.GetSystemError()
	}

	err = ph.Serialize(self, &newState)
	if err != nil {
		return nil, nil, err
	}

{{ range $i := $method.ErrorInterfaceInRes }}
	ret{{ $i }} = ph.MakeErrorSerializable(ret{{ $i }})
{{ end }}

	err = serializeResults()
	if err != nil {
		return
	}

	return
}
{{ end }}


{{ range $f := .Functions }}
func INSCONSTRUCTOR_{{ $f.Name }}(ref insolar.Reference, data []byte) (state []byte, result []byte, err error) {
	ph := common.CurrentProxyCtx
	ph.SetSystemError(nil)

	{{ $f.ArgumentsZeroList }}
	err = ph.Deserialize(data, &args)
	if err != nil {
		err = &foundation.Error{ S: "[ Fake{{ $f.Name }} ] ( INSCONSTRUCTOR_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error() }
		return
	}

	{{ $f.ResultDefinitions }}

	serializeResults := func() error {
		return ph.Serialize(
			foundation.Result{Returns:[]interface{}{ ref, ret1 }},
			&result,
		)
	}

	needRecover := true
	defer func() {
		if !needRecover {
			return
		}
		if r := recover(); r != nil {
			recoveredError := errors.Wrap(errors.Errorf("%v", r), "Failed to execute constructor (panic)")
			recoveredError = ph.MakeErrorSerializable(recoveredError)

			if PanicIsLogicalError {
				ret1 = recoveredError

				err = serializeResults()
				if err== nil {
				    state = data
				}
			} else {
				err = recoveredError
			}
		}
	}()

	{{ $f.Results }} = {{ $f.Name }}( {{ $f.Arguments }} )

	needRecover = false

	ret1 = ph.MakeErrorSerializable(ret1)
	if ret0 == nil && ret1 == nil {
		ret1 = &foundation.Error{ S: "constructor returned nil" }
	}

	if ph.GetSystemError() != nil {
		err = ph.GetSystemError()
		return
	}

	err = serializeResults()
	if err != nil {
		return
	}

	if ret1 != nil {
		// logical error, the result should be registered with type RequestSideEffectNone
		state = nil
		return
	}

	err = ph.Serialize(ret0, &state)
	if err != nil {
		return
	}

	return
}
{{ end }}

{{ if $.GenerateInitialize -}}
func Initialize() insolar.ContractWrapper {
	return insolar.ContractWrapper{
		Methods: insolar.ContractMethods{
			{{ range $method := .Methods -}}
					"{{ $method.Name }}": INSMETHOD_{{ $method.Name }},
			{{ end }}
			"GetCode": INSMETHOD_GetCode,
			"GetPrototype": INSMETHOD_GetPrototype,
		},
		Constructors: insolar.ContractConstructors{
			{{ range $f := .Functions -}}
					"{{ $f.Name }}": INSCONSTRUCTOR_{{ $f.Name }},
			{{ end }}
		},
	}
}
{{- end }}
`
var initializationTmpl = `
// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// Code generated by insgocc. DO NOT EDIT.
// source template in logicrunner/preprocessor/templates

package {{ .Package }}

import (
	"github.com/pkg/errors"
{{ range $contract := .Contracts }}
    {{ $contract.ImportName }} "{{ $contract.ImportPath }}"
{{- end }}

    XXX_insolar "github.com/insolar/insolar/insolar"
    XXX_artifacts "github.com/insolar/insolar/logicrunner/artifacts"
)

func InitializeContractMethods() map[string]XXX_insolar.ContractWrapper {
    return map[string]XXX_insolar.ContractWrapper{
{{- range $contract := .Contracts }}
        "{{ $contract.Name }}": {{ $contract.ImportName }}.Initialize(),
{{- end }}
    }
}

func shouldLoadRef(strRef string) XXX_insolar.Reference {
    ref, err := XXX_insolar.NewReferenceFromString(strRef)
    if err != nil {
        panic(errors.Wrap(err, "Unexpected error, bailing out"))
    }
    return *ref
}

func InitializeCodeRefs() map[XXX_insolar.Reference]string {
    rv := make(map[XXX_insolar.Reference]string, {{ len .Contracts }})

    {{ range $contract := .Contracts -}}
    rv[shouldLoadRef("{{ $contract.CodeReference }}")] = "{{ $contract.Name }}"
    {{ end }}

    return rv
}

func InitializePrototypeRefs() map[XXX_insolar.Reference]string {
    rv := make(map[XXX_insolar.Reference]string, {{ len .Contracts }})

    {{ range $contract := .Contracts -}}
    rv[shouldLoadRef("{{ $contract.PrototypeReference }}")] = "{{ $contract.Name }}"
    {{ end }}

    return rv
}

func InitializeCodeDescriptors() []XXX_artifacts.CodeDescriptor {
    rv := make([]XXX_artifacts.CodeDescriptor, 0, {{ len .Contracts }})

    {{ range $contract := .Contracts -}}
    // {{ $contract.Name }}
    rv = append(rv, XXX_artifacts.NewCodeDescriptor(
        /* code:        */ nil,
        /* machineType: */ XXX_insolar.MachineTypeBuiltin,
        /* ref:         */ shouldLoadRef("{{ $contract.CodeReference }}"),
    ))
    {{ end }}
    return rv
}

func InitializePrototypeDescriptors() []XXX_artifacts.PrototypeDescriptor {
    rv := make([]XXX_artifacts.PrototypeDescriptor, 0, {{ len .Contracts }})

    {{ range $contract := .Contracts }}
    { // {{ $contract.Name }}
        pRef := shouldLoadRef("{{ $contract.PrototypeReference }}")
        cRef := shouldLoadRef("{{ $contract.CodeReference }}")
        rv = append(rv, XXX_artifacts.NewPrototypeDescriptor(
            /* head:         */ pRef,
            /* state:        */ *pRef.GetLocal(),
            /* code:         */ cRef,
        ))
    }
    {{ end }}
    return rv
}
`
