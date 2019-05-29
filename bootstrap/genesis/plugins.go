//
// Copyright 2019 Insolar Technologies GmbH
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
//

package genesis

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (g *Generator) generatePlugins() error {
	insgoccBin := g.config.Contracts.Insgocc
	args := []string{
		"compile-genesis-plugins",
		"-o", g.config.Contracts.OutDir,
	}

	fmt.Println(insgoccBin, strings.Join(args, " "))
	gocc := exec.Command(insgoccBin, args...)
	gocc.Stderr = os.Stderr
	gocc.Stdout = os.Stdout
	return gocc.Run()
}
