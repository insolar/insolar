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
	"math/big"
	"sync"

	"github.com/insolar/insolar/api/sdk"
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
