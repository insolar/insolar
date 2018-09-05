package {{ .PackageName }}

import (
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

{{ range $typeStruct := .Types }}
    {{- $typeStruct }}
{{ end }}

// Reference to class of this contract
var ClassReference = "{{ .ClassReference }}"

// Contract proxy type
type {{ .ContractType }} struct {
    Reference string
}

// GetObject
func GetObject(ref string) (r *{{ .ContractType }}) {
    return &{{ .ContractType }}{Reference: ref}
}

{{ range $func := .ConstructorsProxies }}
func {{ $func.Name }}( {{ $func.Arguments }} ) *{{ $.ContractType }} {
    {{ $func.InitArgs }}

    var argsSerialized []byte
    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

	// TODO: type
    ref, err := proxyctx.Current.RouteConstructorCall(ClassReference, "{{ $func.Name }}", argsSerialized)
    if err != nil {
		panic(err)
    }

    return &{{ $.ContractType }}{Reference: ref}
}
{{ end }}

// GetReference
// TODO replace return to Reference
func (r *{{ $.ContractType }}) GetReference() string {
    return r.Reference
}

// GetClass
// TODO replace return to Reference
func (r *{{ $.ContractType }}) GetClass() string {
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
