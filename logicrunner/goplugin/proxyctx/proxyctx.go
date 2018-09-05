/*
 *    Copyright 2018 INS Ecosystem
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

package proxyctx

// ProxyHelper interface with methods that are needed by contract proxies
type ProxyHelper interface {
	RouteCall(ref string, method string, args []byte) ([]byte, error)
	RouteConstructorCall(classRef string, name string, args []byte) ([]byte, error)
	SaveAsChild(parentRef, classRef string, data []byte) (string, error)
	Serialize(what interface{}, to *[]byte) error
	Deserialize(from []byte, into interface{}) error
}

// Current - hackish way to give proxies access to the current environment
var Current ProxyHelper
