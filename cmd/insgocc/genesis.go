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
	cmd.Flags().StringVarP(&srcDir, "sources-dir", "s", "", "sources dir (default is application/contract)")
	cmd.Flags().StringVarP(&tmpDir, "temp-dir", "t", "", "temporary dir")
	// default value for bool flags is not displayed automatically, thus it's done manually here
	cmd.Flags().BoolVarP(&keepTemp, "keep-temp", "k", false, "keep temp directory (default \"false\")")
	cmd.Flags().BoolVarP(&noProxy, "no-proxy", "", false, "skip proxy compilation (default \"false\")")

	return cmd
}
