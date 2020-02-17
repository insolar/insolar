// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/insolar/insolar/application/api/sdk"
)

func oneSimpleRequest(insSDK *sdk.SDK) {
	fmt.Println("Try to create new member:")
	m, traceID, err := insSDK.CreateMember()
	check("Can not create member, error: ", err)
	fmt.Println("Success! New member ref: ", m.GetReference(), ". TraceId: ", traceID)
	fmt.Print("oneSimpleRequest done just fine\n\n")
}

func severalSimpleRequestToRootMember(insSDK *sdk.SDK) {
	fmt.Println("Try to create several new members:")
	for i := 0; i < 10; i++ {
		m, traceID, err := insSDK.CreateMember()
		check("Can not create member, error: ", err)
		fmt.Println("Success! New member ref: ", m.GetReference(), ". TraceId: ", traceID)
	}
	fmt.Print("severalSimpleRequestToRootMember done just fine\n\n")
}

func severalSimpleRequestToDifferentMembers(insSDK *sdk.SDK) {
	fmt.Println("Try to transfer:")
	fmt.Println("Creating some members for transfer ...")
	var members []sdk.Member
	for i := 0; i < 20; i++ {
		m, traceID, err := insSDK.CreateMember()
		check("Can not create member, error: ", err)
		members = append(members, m)
		fmt.Println("Success! New member ref: ", m.GetReference(), ". TraceId: ", traceID)
	}

	for i := 0; i < 10; i++ {
		traceID, err := insSDK.Transfer(big.NewInt(1).String(), members[i], members[i+10])
		check("Can not transfer money, error: ", err)
		fmt.Println("Transfer success. TraceId: ", traceID)
	}
	fmt.Print("severalSimpleRequestToDifferentMembers done just fine\n\n")
}

func severalParallelRequestToRootMember(insSDK *sdk.SDK) {
	fmt.Println("Try to create several new members in parallel:")
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			m, traceID, err := insSDK.CreateMember()
			check("Can not create member, error: ", err)
			fmt.Println("Success! New member ref: ", m.GetReference(), ". TraceId: ", traceID)
		}()
	}
	wg.Wait()
	fmt.Print("severalParallelRequestToRootMember done just fine\n\n")
}

func severalParallelRequestToDifferentMembers(insSDK *sdk.SDK) {
	fmt.Println("Try to transfer in parallel:")
	fmt.Println("Creating some members for transfer ...")
	var members []sdk.Member
	for i := 0; i < 20; i++ {
		m, traceID, err := insSDK.CreateMember()
		check("Can not create member, error: ", err)
		fmt.Println("Success! New member ref: ", m.GetReference(), ". TraceId: ", traceID)
		members = append(members, m)
	}
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			traceID, err := insSDK.Transfer(big.NewInt(1).String(), members[i], members[i+10])
			check("Can not transfer money, error: ", err)
			fmt.Println("Transfer success. TraceId: ", traceID)
		}(i)
	}
	wg.Wait()
	fmt.Print("severalParallelRequestToDifferentMembers done just fine\n\n")
}
