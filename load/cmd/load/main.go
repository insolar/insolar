package main

import (
	"github.com/insolar/insolar/load"
	"github.com/insolar/loadgen"
)

func main() {
	loadgen.Run(load.AttackerFromName, load.CheckFromName)
}
