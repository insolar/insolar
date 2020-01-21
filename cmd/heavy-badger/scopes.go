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
	"fmt"

	"github.com/spf13/cobra"

	"github.com/insolar/insolar/insolar/store"
)

func allScopes() []store.Scope {
	var result []store.Scope
	start := store.ScopePulse
	end := store.ScopeNodeHistory
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	return result
}

func scopeFromName(name string) (store.Scope, error) {
	for _, s := range allScopes() {
		if s.String() == name {
			return s, nil
		}
	}
	return 0, fmt.Errorf("scope with name %v not found", name)
}

func scopesListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "scopes",
		Short: "print all scope names",
		Run: func(_ *cobra.Command, _ []string) {
			for _, v := range allScopes() {
				fmt.Printf("%s: %d (b%08b)\n", v.String(), v, v)
			}
		},
	}
}
