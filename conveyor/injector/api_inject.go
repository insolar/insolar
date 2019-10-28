//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package injector

type DependencyRegistry interface {
	FindDependency(id string) (interface{}, bool)
}

type LocalDependencyRegistry interface {
	DependencyRegistry
	FindLocalDependency(id string) (interface{}, bool)
}

type DependencyRegistryFunc func(id string) (interface{}, bool)

type DependencyContainer interface {
	DependencyRegistry
	PutDependency(id string, v interface{})
	TryPutDependency(id string, v interface{}) bool
}

type DependencyProviderFunc func(target interface{}, id string, resolveFn DependencyRegistryFunc) interface{}
