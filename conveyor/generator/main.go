/*
 *    Copyright 2019 Insolar Technologies
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

package main

import (
	"os"
	"github.com/insolar/insolar/conveyor/generator/generator"
	"github.com/insolar/insolar/log"
	"bufio"
)

func main() {
	g := generator.Generator{}
	g.ParseFile("conveyor/generator/sample/sample_state_machine.go")
	file, err := os.Create("conveyor/generator/sample/sample_state_machine_generated.go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	g.GenerateStateMachine(w, 0)
	g.GenerateRawHandlers(w, 0)
	w.Flush()
}
