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
		holder := configuration.NewHolder()
		fmt.Printf("Default configuration:\n %+v\n", holder.Configuration)
		holder.SaveAs("insolar.yml")
	}

*/
package configuration
