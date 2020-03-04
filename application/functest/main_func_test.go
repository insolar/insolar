// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"fmt"
	"os"
	"testing"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/application/testutils/launchnet"
	baselaunchnet "github.com/insolar/insolar/applicationbase/testutils/launchnet"
)

func TestMain(m *testing.M) {
	os.Exit(baselaunchnet.Run(
		func() int {
			err := setMigrationDaemonsRef()
			if err != nil {
				fmt.Println(errors.Wrap(err, "[ setup ] get reference daemons by public key failed ").Error())
			}

			return m.Run()
		},
		launchnet.AppPath,
		launchnet.LoadAllMembersKeys,
		launchnet.SetInfo,
		launchnet.AfterSetup))
}
