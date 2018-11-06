/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package version

import (
	"fmt"
	"runtime"
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
