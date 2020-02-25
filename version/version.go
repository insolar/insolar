// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package version

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Version is release semantic version.
	Version = "unset"
	// BuildNumber is CI build number.
	BuildNumber = "unset"
	// BuildDate is build date.
	BuildDate = "unset"
	// BuildTime is build date.
	BuildTime = "unset"
	// CITool is a continuous integration tool(Travis, DockerCloud, etc.).
	CITool = "unset"
	// GitHash is short git commit hash.
	GitHash = "unset"
)

// GetFullVersion returns multi line full version information
func GetFullVersion() string {

	result := fmt.Sprintf(`
 Version      : %s
 Build number : %s
 Build date   : %s %s
 Git hash     : %s
 Go version   : %s
 Go compiler  : %s
 Platform     : %s/%s`, Version, BuildNumber, BuildDate, BuildTime, GitHash, runtime.Version(),
		runtime.Compiler, runtime.GOOS, runtime.GOARCH)

	return result
}

func GetCommand(cmdName string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Print the version info of %s", cmdName),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(GetFullVersion())
			os.Exit(0)
		},
	}
}
