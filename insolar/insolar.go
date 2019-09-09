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

package insolar

import (
	"go/build"
	"log"
)

// RootModule holds root module name.
var RootModule = "github.com/insolar/insolar"

// RootModuleDir returns abs path to root module for any package it runned.
func RootModuleDir() string {
	p, err := build.Default.Import(RootModule, ".", build.FindOnly)
	if err != nil {
		log.Fatal("failed to resolve", RootModule)
	}
	return p.Dir
}
