package load

import (
	"github.com/skudasov/loadgen"
)

func CheckFromName(name string) loadgen.RuntimeCheckFunc {
	switch name {
	default:
		return nil
	}
}
