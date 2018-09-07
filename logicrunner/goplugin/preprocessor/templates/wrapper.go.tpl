package {{ .PackageName}}

import (
    {{- range $import, $i := .Imports }}
        {{$import}}
    {{- end }}
)

{{ range $method := .Methods }}
func (self *{{ $.ContractType }}) INSWRAPER_{{ $method.Name }}(cbor foundation.CBORMarshaler, data []byte) ([]byte) {
    {{ $method.ArgumentsZeroList }}
    cbor.Unmarshal(&args, data)

{{ if $method.Results }}
    {{ $method.Results }}:= self.{{ $method.Name }}( {{ $method.Arguments }} )
{{ else }}
    self.{{ $method.Name }}( {{ $method.Arguments }} )
{{ end }}

    return cbor.Marshal([]interface{} { {{ $method.Results }}} )
}
{{ end }}

var INSEXPORT {{ .ContractType }}
