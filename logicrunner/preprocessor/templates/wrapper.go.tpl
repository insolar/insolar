//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package {{ .Package }}

import (
{{- range $import, $i := .Imports }}
	{{ $import }}
{{- end }}
    XXX_preprocessor "github.com/insolar/insolar/logicrunner/preprocessor"
)

type ExtendableError struct{
	S string
}

func ( e *ExtendableError ) Error() string{
	return e.S
}

func INSMETHOD_GetCode(object []byte, data []byte) ([]byte, []byte, error) {
	ph := proxyctx.Current
	self := new({{ $.ContractType }})

	if len(object) == 0 {
		return nil, nil, &ExtendableError{ S: "[ Fake GetCode ] ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &ExtendableError{ S: "[ Fake GetCode ] ( Generated Method ) Can't deserialize args.Data: " + err.Error() }
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
	ph := proxyctx.Current
	self := new({{ $.ContractType }})

	if len(object) == 0 {
		return nil, nil, &ExtendableError{ S: "[ Fake GetPrototype ] ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &ExtendableError{ S: "[ Fake GetPrototype ] ( Generated Method ) Can't deserialize args.Data: " + err.Error() }
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
func INSMETHOD_{{ $method.Name }}(object []byte, data []byte) ([]byte, []byte, error) {
	ph := proxyctx.Current

	self := new({{ $.ContractType }})

	if len(object) == 0 {
		return nil, nil, &ExtendableError{ S: "[ Fake{{ $method.Name }} ] ( INSMETHOD_* ) ( Generated Method ) Object is nil"}
	}

	err := ph.Deserialize(object, self)
	if err != nil {
		e := &ExtendableError{ S: "[ Fake{{ $method.Name }} ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Data: " + err.Error() }
		return nil, nil, e
	}

	{{ $method.ArgumentsZeroList }}
	err = ph.Deserialize(data, &args)
	if err != nil {
		e := &ExtendableError{ S: "[ Fake{{ $method.Name }} ] ( INSMETHOD_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error() }
		return nil, nil, e
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
		e := &ExtendableError{ S: "[ Fake{{ $f.Name }} ] ( INSCONSTRUCTOR_* ) ( Generated Method ) Can't deserialize args.Arguments: " + err.Error() }
		return nil, e
	}

	{{ $f.Results }} := {{ $f.Name }}( {{ $f.Arguments }} )
	if ret1 != nil {
		return nil, ret1
	}

	ret := []byte{}
	err = ph.Serialize(ret0, &ret)
	if err != nil {
		return nil, err
	}

	if ret0 == nil {
		e := &ExtendableError{ S: "[ Fake{{ $f.Name }} ] ( INSCONSTRUCTOR_* ) ( Generated Method ) Constructor returns nil" }
		return nil, e
	}

	return ret, err
}
{{ end }}

{{ if $.GenerateInitialize -}}
func Initialize() XXX_preprocessor.ContractWrapper {
    return XXX_preprocessor.ContractWrapper{
        GetCode: INSMETHOD_GetCode,
        GetPrototype: INSMETHOD_GetPrototype,
        Methods: XXX_preprocessor.ContractMethods{
            {{ range $method := .Methods -}}
                    "{{ $method.Name }}": INSMETHOD_{{ $method.Name }},
            {{- end }}
        },
        Constructors: XXX_preprocessor.ContractConstructors{
            {{ range $f := .Functions -}}
                    "{{ $f.Name }}": INSCONSTRUCTOR_{{ $f.Name }},
            {{ end }}
        },
    }
}
{{- end }}
