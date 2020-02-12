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

package api

const (
	ParseError                     = -31700
	ParseErrorShort                = "ParseError"
	ParseErrorMessage              = "Parsing error on the server side: received an invalid JSON."
	InvalidRequestError            = -31600
	InvalidRequestErrorShort       = "InvalidRequest"
	InvalidRequestErrorMessage     = "The JSON received is not a valid request payload."
	MethodNotFoundError            = -31601
	MethodNotFoundErrorShort       = "MethodNotFound"
	MethodNotFoundErrorMessage     = "Method does not exist / is not available."
	InvalidParamsError             = -31602
	InvalidParamsErrorShort        = "InvalidParams"
	InvalidParamsErrorMessage      = "Invalid method parameter(s)."
	InternalError                  = -31603
	InternalErrorShort             = "Internal"
	InternalErrorMessage           = "Internal Platform error."
	TimeoutError                   = -31106
	TimeoutErrorShort              = "Timeout"
	TimeoutErrorMessage            = "Request's timeout has expired."
	UnauthorizedError              = -31401
	UnauthorizedErrorShort         = "Unauthorized"
	UnauthorizedErrorMessage       = "Action is not authorized."
	ExecutionError                 = -31103
	ExecutionErrorShort            = "Execution"
	ExecutionErrorMessage          = "Execution error."
	ServiceUnavailableError        = -31429
	ServiceUnavailableErrorShort   = "ServiceUnavailable"
	ServiceUnavailableErrorMessage = "Service unavailable, try again later."
)
