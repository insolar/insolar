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

package testutil

import (
	"go/build"
	"os/exec"
	"path/filepath"

	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

const insolarImportPath = "github.com/insolar/insolar"

var testdataDir = testdataPath()

func buildCLI(name string) (string, error) {
	clipath := filepath.Join(testdataDir, name)
	out, err := exec.Command(
		"go", "build",
		"-o", clipath,
		filepath.Join(insolarImportPath, "logicrunner", "goplugin", name),
	).CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "can't build %s: %s", name, string(out))
	}
	return clipath, nil
}

func buildInciderCLI() (string, error) {
	return buildCLI("ginsider-cli")
}

func buildPreprocessor() (string, error) {
	insgocc := filepath.Join(testdataDir, "insgocc")
	out, err := exec.Command(
		"go", "build",
		"-o", insgocc,
		filepath.Join(insolarImportPath, "cmd", "insgocc"),
	).CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "can't build preprocessor. Build output: %s", string(out))
	}
	return insgocc, nil
}

func testdataPath() string {
	p, err := build.Default.Import("github.com/insolar/insolar", "", build.FindOnly)
	log.Fatal("import found dir:", p.Dir)
	if err != nil {
		panic(err)
	}
	return filepath.Join(p.Dir, "testdata", "logicrunner")
}

// Build compiles and return path to ginsider-cli and insgocc binaries.
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
