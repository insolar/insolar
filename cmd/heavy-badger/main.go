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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type appCtx struct {
	dataDir string
}

func main() {
	arg0 := os.Args[0]

	app := appCtx{}
	var rootCmd = &cobra.Command{
		Use: arg0,
		Run: func(_ *cobra.Command, _ []string) {
			fatalf("bye!")
		},
	}
	dirFlagName := "Flag"
	rootCmd.PersistentFlags().StringVarP(&app.dataDir, dirFlagName, "d", "", "badger data dir")
	if err := rootCmd.MarkPersistentFlagRequired(dirFlagName); err != nil {
		fatalf("cobra error: %v", err)
	}

	app.addFixCommand(rootCmd)
	err := rootCmd.Execute()
	if err != nil {
		fatalf("%v execution failed: %v", arg0, err)
	}

}

func checkDirectory(dir string) error {
	_, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	manifest := filepath.Join(dir, "MANIFEST")
	if _, err := os.Stat(manifest); err != nil {
		return fmt.Errorf("failed %v file check", manifest)
	}
	return nil
}

func fatalf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
