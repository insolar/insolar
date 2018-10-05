/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

/*
Package rpc allows higher level components to register methods that can be called by other network hosts.

Usage:

	r := rpc.NewRPC()

	r.RegisterMethod("hello_world", func(sender *host.Host, args [][]byte) ([]byte, error) {
		fmt.Println("Hello World")
		return nil, nil
	})

	r.Invoke(&host.Host{}, "hello_world", [][]byte{})

*/
package rpc
