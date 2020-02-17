// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package configuration

import (
	"fmt"
)

// APIRunner holds configuration for api
type APIRunner struct {
	Address string
	RPC     string
	// IsAdmin indicates status of api (internal or external)
	IsAdmin     bool
	SwaggerPath string
}

// NewAPIRunner creates new api config
func NewAPIRunner(admin bool) APIRunner {
	if admin {
		return APIRunner{
			Address:     "localhost:19001",
			RPC:         "/admin-api/rpc",
			SwaggerPath: "application/api/spec/api-exported.yaml",
			IsAdmin:     true,
		}
	}
	return APIRunner{
		Address:     "localhost:19101",
		RPC:         "/api/rpc",
		SwaggerPath: "application/api/spec/api-exported.yaml",
		IsAdmin:     false,
	}
}

func (ar *APIRunner) String() string {
	res := fmt.Sprintf("Addr -> %s, RPC -> %s, IsAdmin -> %t, SwaggerPath -> %s\n", ar.Address, ar.RPC, ar.IsAdmin, ar.SwaggerPath)
	return res
}
