package {{ .PackageName }}

import (
    {{- range $import, $i := .Imports }}
        {{$import}}
    {{- end }}
)

{{ range $typeStruct := .Types }}
    {{- $typeStruct }}
{{ end }}

// ClassReference to class of this contract
var ClassReference = core.NewRefFromBase58("{{ .ClassReference }}")

// Contract proxy type
type {{ .ContractType }} struct {
    Reference core.RecordRef
}

type ContractConstructorHolder struct {
	constructorName string
	argsSerialized []byte
}

func (r *ContractConstructorHolder) AsChild(objRef core.RecordRef) *{{ .ContractType }} {
    ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.constructorName, r.argsSerialized)
    if err != nil {
        panic(err)
    }
    return &{{ .ContractType }}{Reference: ref}
}

func (r *ContractConstructorHolder) AsDelegate(objRef core.RecordRef) *{{ .ContractType }} {
    ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.constructorName, r.argsSerialized)
    if err != nil {
        panic(err)
    }
    return &{{ .ContractType }}{Reference: ref}
}

// GetObject
func GetObject(ref core.RecordRef) (r *{{ .ContractType }}) {
    return &{{ .ContractType }}{Reference: ref}
}

func GetClass() core.RecordRef {
    return ClassReference
}

func GetImplementationFrom(object core.RecordRef) *{{ .ContractType }} {
    ref, err := proxyctx.Current.GetDelegate(object, ClassReference)
    if err != nil {
        panic(err)
    }
    return GetObject(ref)
}

{{ range $func := .ConstructorsProxies }}
func {{ $func.Name }}( {{ $func.Arguments }} ) *ContractConstructorHolder {
    {{ $func.InitArgs }}

    var argsSerialized []byte
    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    return &ContractConstructorHolder{constructorName: "{{ $func.Name }}", argsSerialized: argsSerialized}
}
{{ end }}

// GetReference
func (r *{{ $.ContractType }}) GetReference() core.RecordRef {
    return r.Reference
}

// GetClass
func (r *{{ $.ContractType }}) GetClass() core.RecordRef {
    return ClassReference
}

{{ range $method := .MethodsProxies }}
func (r *{{ $.ContractType }}) {{ $method.Name }}( {{ $method.Arguments }} ) ( {{ $method.ResultsTypes }} ) {
    {{ $method.InitArgs }}
    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    res, err := proxyctx.Current.RouteCall(r.Reference, true, "{{ $method.Name }}", argsSerialized)
    if err != nil {
   		panic(err)
    }

    {{ $method.ResultZeroList }}
    err = proxyctx.Current.Deserialize(res, &resList)
    if err != nil {
        panic(err)
    }

    return {{ $method.Results }}
}

func (r *{{ $.ContractType }}) {{ $method.Name }}NoWait( {{ $method.Arguments }} ) {
    {{ $method.InitArgs }}
    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    _, err = proxyctx.Current.RouteCall(r.Reference, false, "{{ $method.Name }}", argsSerialized)
    if err != nil {
        panic(err)
    }
}
{{ end }}
