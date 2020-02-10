// Copyright 2020 Insolar Network Ltd.
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

package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/insolar/insolar/application/bootstrap"
	basebootstrap "github.com/insolar/insolar/applicationbase/bootstrap"
)

func bootstrapCommand() *cobra.Command {
	var (
		configPath         string
		certificatesOutDir string
	)
	c := &cobra.Command{
		Use:   "bootstrap",
		Short: "creates files required for new network (keys, genesis config)",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			contractsGen, err := bootstrap.NewGenesisContractsGenerator(configPath)
			check("failed to start genesis contracts config generator", err)

			contractsConfig, err := contractsGen.CreateGenesisContractsConfig(ctx)
			check("failed to create genesis contracts config", err)

			gen, err := basebootstrap.NewGenerator(configPath, certificatesOutDir, contractsConfig)
			check("base bootstrap failed to start", err)

			err = gen.Run(ctx)
			check("base bootstrap failed", err)
		},
	}
	c.Flags().StringVarP(
		&configPath, "config", "c", "bootstrap.yaml", "path to bootstrap config")
	c.Flags().StringVarP(
		&certificatesOutDir, "certificates-out-dir", "o", "", "dir with certificate files")
	return c
}
