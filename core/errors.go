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

package core

import "github.com/pkg/errors"

var (
	// ErrUnknown returned when error type cannot be defined.
	ErrUnknown = errors.New("unknown error")
	// ErrDeactivated returned when requested object is deactivated.
	ErrDeactivated = errors.New("object is deactivated")
	// ErrStateNotAvailable returned when requested object is deactivated.
	ErrStateNotAvailable = errors.New("object state is not available")
)
