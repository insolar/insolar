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

package goplugintestutils

import (
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const insolarImportPath = "github.com/insolar/insolar"

var testdataDir = testdataPath()

func buildCLI(name string) (string, error) {
	binPath := filepath.Join(testdataDir, name)
	cmd := exec.Command(
		"go", "build",
		"-o", binPath,
		filepath.Join(insolarImportPath, "cmd", name),
	)
	for _, ev := range os.Environ() {
		if strings.HasPrefix(ev, "CGO_ENABLED=") {
			continue
		}
		cmd.Env = append(cmd.Env, ev)
	}
	cmd.Env = append(cmd.Env, "CGO_ENABLED=1")

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "can't build preprocessor. Build output: %s", string(out))
	}
	return binPath, nil
}

func buildInciderCLI() (string, error) {
	return buildCLI("insgorund")
}

func buildPreprocessor() (string, error) {
	return buildCLI("insgocc")
}

func testdataPath() string {
	p, err := build.Default.Import("github.com/insolar/insolar", "", build.FindOnly)
	if err != nil {
		panic(err)
	}
	return filepath.Join(p.Dir, "testdata", "logicrunner")
}

// Build compiles and return path to insgorund and insgocc binaries.
func Build() (string, string, error) {
	icc, err := buildInciderCLI()
	if err != nil {
		return "", "", err
	}

	insgocc, err := buildPreprocessor()
	if err != nil {
		return "", "", err
	}
	return icc, insgocc, nil
}
