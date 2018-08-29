package {{ .PackageName }}

import (
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

{{ .Types }}

// Contract proxy type
type {{ .ContractType }} struct {
    Reference string
}

// GetObject
func GetObject(ref string) (r *{{ .ContractType }}) {
    return &{{ .ContractType }}{Reference: ref}
}

{{ .MethodsProxies }}
