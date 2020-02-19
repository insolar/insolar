// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package goplugintestutils

import (
	"go/build"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

const insolarImportPath = "github.com/insolar/insolar"

var testdataDir = testdataPath()

func buildCLI(name string) (string, error) {
	binPath := filepath.Join(testdataDir, name)

	out, err := exec.Command(
		"go", "build",
		"-o", binPath,
		filepath.Join(insolarImportPath, name),
	).CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "can't build preprocessor. buildPrototypes output: %s", string(out))
	}
	return binPath, nil
}

func buildInsiderCLI() (string, error) {
	return buildCLI("cmd/insgorund")
}

func BuildPreprocessor() (string, error) {
	return buildCLI("cmd/insgocc")
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
	icc, err := buildInsiderCLI()
	if err != nil {
		return "", "", err
	}

	insgocc, err := BuildPreprocessor()
	if err != nil {
		return "", "", err
	}
	return icc, insgocc, nil
}
