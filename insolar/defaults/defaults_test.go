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

package defaults

import (
	"os"
	"testing"

	"github.com/magiconair/properties/assert"
)

type tCase struct {
	env    map[string]string
	defFn  func() string
	expect string
}

var cases = []tCase{

	// ArtifactsDir checks
	{
		defFn:  ArtifactsDir,
		expect: ".artifacts",
	},
	{
		env: map[string]string{
			"INSOLAR_ARTIFACTS_DIR": "blah/bla",
		},
		defFn:  ArtifactsDir,
		expect: "blah/bla",
	},

	// LaunchnetDir checks
	{
		defFn:  LaunchnetDir,
		expect: ".artifacts/launchnet",
	},
	{
		env: map[string]string{
			"INSOLAR_ARTIFACTS_DIR": "blah/bla",
		},
		defFn:  LaunchnetDir,
		expect: "blah/bla/launchnet",
	},
	{
		env: map[string]string{
			"LAUNCHNET_BASE_DIR": "blah/bla",
		},
		defFn:  LaunchnetDir,
		expect: "blah/bla",
	},
}

func TestDefaults(t *testing.T) {
	for _, tc := range cases {
		for name, value := range tc.env {
			os.Setenv(name, value)
		}

		assert.Equal(t, tc.defFn(), tc.expect)

		for name := range tc.env {
			os.Setenv(name, "")
		}
	}
}
