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
	"github.com/pkg/errors"

{{- range $contract := .Contracts }}
    {{ $contract.ImportName }} "{{ $contract.ImportPath }}"
{{- end }}

    XXX_insolar "github.com/insolar/insolar/insolar"
    XXX_preprocessor "github.com/insolar/insolar/logicrunner/preprocessor"
)

func InitializeContractMethods() map[string]XXX_preprocessor.ContractWrapper {
    return map[string]XXX_preprocessor.ContractWrapper{
{{- range $contract := .Contracts }}
        "{{ $contract.Name }}": {{ $contract.ImportName }}.Initialize(),
{{- end }}
    }
}

func shouldLoadRef(strRef string) XXX_insolar.Reference {
    ref, err := XXX_insolar.NewReferenceFromBase58(strRef)
    if err != nil {
        panic(errors.Wrap(err, "Unexpected error, bailing out"))
    }
    return *ref
}

func InitializeContractRefs() map[XXX_insolar.Reference]string {
    rv := make(map[XXX_insolar.Reference]string, 0)

    {{ range $contract := .Contracts -}}
    rv[shouldLoadRef("{{ $contract.Reference }}")] = "{{ $contract.Name }}"
    {{- end }}

    return rv
}
