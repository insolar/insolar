// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build debug

package insolard

import (
	"github.com/google/gops/agent"
)

func init() {
	// starts gops agent https://github.com/google/gops on default addr (127.0.0.1:0)
	psAgentLauncher = func() error {
		if err := agent.Listen(agent.Options{}); err != nil {
			return err
		}
		return nil
	}
}
