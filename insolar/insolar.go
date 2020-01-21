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

package insolar

import (
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/insolar/insolar/insolar/defaults"
)

// RootModule holds root module name.
var RootModule = "github.com/insolar/insolar"

// RootModuleDir returns abs path to root module for any package where it's called.
func RootModuleDir() string {
	p, err := build.Default.Import(RootModule, ".", build.FindOnly)
	if err != nil {
		log.Fatal("failed to resolve", RootModule)
	}
	return p.Dir
}

func ContractBuildTmpDir(prefix string) string {
	dir := filepath.Join(RootModuleDir(), defaults.ArtifactsDir(), "tmp")
	// create if not exist
	if err := os.MkdirAll(dir, 0777); err != nil {
		panic(err)
	}

	tmpDir, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		panic(err)
	}
	return tmpDir
}
