// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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

			contractsConfig, err := bootstrap.CreateGenesisContractsConfig(ctx, configPath)
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
