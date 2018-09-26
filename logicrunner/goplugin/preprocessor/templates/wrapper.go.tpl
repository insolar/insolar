package main

import (
    {{- range $import, $i := .Imports }}
        {{$import}}
    {{- end }}
)

{{ range $method := .Methods }}
func INSMETHOD_{{ $method.Name }}(ph proxyctx.ProxyHelper, object []byte,
                data []byte, context *core.LogicCallContext) ([]byte, []byte, error) {

    self := new({{ $.ContractType }})

    err := ph.Deserialize(object, self)
    if err != nil {
        return nil, nil, err
    }

    self.SetContext(context)

    {{ $method.ArgumentsZeroList }}
    err = ph.Deserialize(data, &args)
    if err != nil {
        return nil, nil, err
    }

{{ if $method.Results }}
    {{ $method.Results }} := self.{{ $method.Name }}( {{ $method.Arguments }} )
{{ else }}
    self.{{ $method.Name }}( {{ $method.Arguments }} )
{{ end }}

    state := []byte{}
    err = ph.Serialize(self, &state)
    if err != nil {
        return nil, nil, err
    }

{{ range $i := $method.ErrorInterfaceInRes }}
    ret{{ $i }} = ph.MakeErrorSerializable(ret{{ $i }})
{{ end }}

    ret := []byte{}
    err = ph.Serialize([]interface{} { {{ $method.Results }} }, &ret)

    return state, ret, err
}
{{ end }}


{{ range $f := .Functions }}
func INSCONSTRUCTOR_{{ $f.Name }}(ph proxyctx.ProxyHelper, data []byte) ([]byte, error) {

    {{ $f.ArgumentsZeroList }}
    err := ph.Deserialize(data, &args)
    if err != nil {
        return nil, err
    }

    {{ $f.Results }} := {{ $f.Name }}( {{ $f.Arguments }} )

    ret := []byte{}
    err = ph.Serialize({{ $f.Results }}, &ret)
    return ret, err
}
{{ end }}
