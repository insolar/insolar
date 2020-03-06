// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func genesisCompile() *cobra.Command {
	var (
		srcDir   = ""
		outDir   = ""
		tmpDir   = ""
		keepTemp = false
		noProxy  = false
	)

	var cmd = &cobra.Command{
		Use:   "compile-genesis-plugins",
		Short: "Compile genesis plugins",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			cb := newContractBuilder(tmpDir, noProxy)
			cb.setSourcesDir(srcDir)
			cb.setOutputDir(outDir)
			if !keepTemp {
				defer cb.clean()
			}
			results, err := cb.build(ctx, args...)
			checkError(err)
			for _, res := range results {
				fmt.Printf("make %v for contract %v\n", res.SoFilePath, res.ContractName)
			}
		},
	}
	cmd.Flags().StringVarP(&outDir, "output-dir", "o", ".", "output dir")
	cmd.Flags().StringVarP(&srcDir, "sources-dir", "s", "", "sources dir (default is applicationbase/contract)")
	cmd.Flags().StringVarP(&tmpDir, "temp-dir", "t", "", "temporary dir")
	// default value for bool flags is not displayed automatically, thus it's done manually here
	cmd.Flags().BoolVarP(&keepTemp, "keep-temp", "k", false, "keep temp directory (default \"false\")")
	cmd.Flags().BoolVarP(&noProxy, "no-proxy", "", false, "skip proxy compilation (default \"false\")")

	return cmd
}
