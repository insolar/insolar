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

// +build debug

package main

import (
	"github.com/google/gops/agent"
)

func init() {
	// starts gops agent https://github.com/google/gops on default addr (127.0.0.1:0)
	psAgentLauncher = func() error {
		if err := agent.Listen(agent.Options{}); err != nil {
			return err
		}
		return nil
	}
}
