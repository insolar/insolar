package main

import (
	"github.com/insolar/insolar/load"
	"github.com/skudasov/loadgen"
)

func main() {
	loadgen.Run(load.AttackerFromName, load.CheckFromName)
}
