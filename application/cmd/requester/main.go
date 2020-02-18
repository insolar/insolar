//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

package main

import (
	"github.com/insolar/insolar/application/cmd/requester/cmd"
	"github.com/insolar/insolar/log"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal("requester execution failed:", err)
	}
}
