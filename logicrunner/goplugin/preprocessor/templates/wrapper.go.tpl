package main

import (
    {{- range $import, $i := .Imports }}
        {{$import}}
    {{- end }}
)

{{ range $method := .Methods }}
func INSMETHOD_{{ $method.Name }}(object []byte, data []byte) ([]byte, []byte, error) {
    ph := proxyctx.Current

    self := new({{ $.ContractType }})

    err := ph.Deserialize(object, self)
    if err != nil {
        return nil, nil, err
    }

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
func INSCONSTRUCTOR_{{ $f.Name }}(data []byte) ([]byte, error) {
    ph := proxyctx.Current
    {{ $f.ArgumentsZeroList }}
    err := ph.Deserialize(data, &args)
    if err != nil {
        return nil, err
    }

    {{ $f.Results }} := {{ $f.Name }}( {{ $f.Arguments }} )
    if ret1 != nil {
        return nil, err
    }

    ret := []byte{}
    err = ph.Serialize(ret0, &ret)
    if err != nil {
        return nil, err
    }

    return ret, err
}
{{ end }}
