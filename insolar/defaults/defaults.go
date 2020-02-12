// Copyright 2020 Insolar Network Ltd.
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

package defaults

import (
	"os"
	"path/filepath"
)

// ArtifactsDir returns path of artifacts dir.
func ArtifactsDir() string {
	return envVarWithDefault("INSOLAR_ARTIFACTS_DIR", ".artifacts")
}

// LaunchnetDir returns path of launchnet's artifacts dir.
func LaunchnetDir() string {
	return envVarWithDefault("LAUNCHNET_BASE_DIR", filepath.Join(ArtifactsDir(), "launchnet"))
}

// LaunchnetConfigDir returns path of launchnet's configs dir.
func LaunchnetConfigDir() string {
	return filepath.Join(LaunchnetDir(), "configs")
}

// LaunchnetDiscoveryNodesLogsDir returns path to dir with launchnet's discovery nodes logs.
func LaunchnetDiscoveryNodesLogsDir() string {
	return filepath.Join(LaunchnetDir(), "logs", "discoverynodes")
}

func envVarWithDefault(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value != "" {
		return value
	}
	return defaultValue
}

// PathWithBaseDir adds base path to path if path is not absolute.
func PathWithBaseDir(path string, base string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(base, path)
}
