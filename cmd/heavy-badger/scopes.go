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
