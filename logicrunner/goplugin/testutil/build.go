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
	"os/exec"

	"github.com/pkg/errors"
)

// ICC shared path to compiler binary.
var ICC = "../cmd/insgocc/insgocc"

func buildCLI(name string) error {
	out, err := exec.Command("go", "build", "-o", "./goplugin/"+name+"/"+name, "./goplugin/"+name+"/").CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "can't build %s: %s", name, string(out))
	}
	return nil
}

func buildInciderCLI() error {
	return buildCLI("ginsider-cli")
}

func buildPreprocessor() error {
	out, err := exec.Command("go", "build", "-o", ICC, "../cmd/insgocc/").CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "can't build %s: %s", ICC, string(out))
	}
	return nil
}

// Build compiles ginsider-cli
func Build() error {
	err := buildInciderCLI()
	if err != nil {
		return err
	}

	err = buildPreprocessor()
	if err != nil {
		return err
	}
	return nil
}
