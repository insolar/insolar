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

package bootstrap

import (
	"context"

	"github.com/insolar/insolar/instrumentation/inslogger"
)

type Initializer struct {
	cfgPath           string
	genesisConfigPath string
	genesisKeyOut     string
}

func NewInitializer(cfgPath string, genesisConfigPath, genesisKeyOut string) *Initializer {
	return &Initializer{
		cfgPath:           cfgPath,
		genesisConfigPath: genesisConfigPath,
		genesisKeyOut:     genesisKeyOut,
	}
}

func (s *Initializer) Run(ctx context.Context) {
	genesisConfig, err := ParseGenesisConfig(s.genesisConfigPath)
	checkError(ctx, err, "failed to create genesis Generator")

	genesisGenerator := NewGenerator(
		genesisConfig,
		s.genesisKeyOut,
	)
	err = genesisGenerator.Run(ctx)
	checkError(ctx, err, "failed to generate genesis")
}

func checkError(ctx context.Context, err error, message string) {
	if err == nil {
		return
	}
	inslogger.FromContext(ctx).Fatalf("%v: %v", message, err.Error())
}
