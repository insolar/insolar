package {{ .PackageName }}

import (
    {{- range $import, $i := .Imports }}
        {{$import}}
    {{- end }}
)

{{ range $typeStruct := .Types }}
    {{- $typeStruct }}
{{ end }}

// Reference to class of this contract
var ClassReference = core.String2Ref("{{ .ClassReference }}")

// Contract proxy type
type {{ .ContractType }} struct {
    Reference core.RecordRef
}

type ContractHolder struct {
	data []byte
}

func (r *ContractHolder) AsChild(objRef core.RecordRef) *{{ .ContractType }} {
    ref, err := proxyctx.Current.SaveAsChild(objRef, ClassReference, r.data)
    if err != nil {
        panic(err)
    }
    return &{{ .ContractType }}{Reference: ref}
}

func (r *ContractHolder) AsDelegate(objRef core.RecordRef) *{{ .ContractType }} {
    ref, err := proxyctx.Current.SaveAsDelegate(objRef, ClassReference, r.data)
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

{{ range $func := .ConstructorsProxies }}
func {{ $func.Name }}( {{ $func.Arguments }} ) *ContractHolder {
    {{ $func.InitArgs }}

    var argsSerialized []byte
    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    data, err := proxyctx.Current.RouteConstructorCall(ClassReference, "{{ $func.Name }}", argsSerialized)
    if err != nil {
		panic(err)
    }

    return &ContractHolder{data: data}
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

    res, err := proxyctx.Current.RouteCall(r.Reference, "{{ $method.Name }}", argsSerialized)
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
{{ end }}
