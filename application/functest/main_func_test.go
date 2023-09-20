// +build functest bloattest functest_endless_abandon

package functest

import (
	"fmt"
	"os"
	"testing"

	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
)

func TestMain(m *testing.M) {
	os.Exit(launchnet.Run(
		func() int {
			err := LoadAllMembersKeys()
			if err != nil {
				fmt.Println("[ setup ] error while loading keys: ", err.Error())
				return 1
			}
			fmt.Println("[ setup ] all keys successfully loaded")
			return m.Run()
		},
		AppPath,
		SetInfo,
		AfterSetup,
		"-gw",
	))
}
