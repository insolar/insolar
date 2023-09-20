/*
Package configuration holds configuration for all components in Insolar host binary
It allows also helps to manage config resources using Holder

Usage:

	package main

	import (
		"github.com/insolar/insolar/configuration"
		"fmt"
	)

	func main() {
		holder := configuration.LightHolder()
		fmt.Printf("Default configuration:\n %+v\n", holder.VirtualConfig)
		holder.SaveAs("insolar.yml")
	}

*/
package configuration
