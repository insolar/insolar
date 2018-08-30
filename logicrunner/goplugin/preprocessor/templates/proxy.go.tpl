package {{ .PackageName }}

import (
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

{{ range $typeStruct := .Types }}
    {{- $typeStruct }}
{{ end }}

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
    ref, err := proxyctx.Current.RouteConstructorCall("typeRef", "{{ $func.Name }}", argsSerialized)
    if err != nil {
		panic(err)
    }

    return &{{ $.ContractType }}{Reference: ref}
}
{{ end }}

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
