package main

import (
    {{- range $import, $i := .Imports }}
        {{$import}}
    {{- end }}
)

{{ range $method := .Methods }}
func INSWRAPPER_{{ $method.Name }}(ph proxyctx.ProxyHelper, object []byte,
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
    {{ $method.Results }}:= self.{{ $method.Name }}( {{ $method.Arguments }} )
{{ else }}
    self.{{ $method.Name }}( {{ $method.Arguments }} )
{{ end }}

    state := []byte{}
    err = ph.Serialize(self, &state)
    if err != nil {
        return nil, nil, err
    }

    ret := []byte{}
    err = ph.Serialize([]interface{} { {{ $method.Results }}}, &ret)

    return state, ret, err
}
{{ end }}