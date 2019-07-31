//
// Copyright 2019 Insolar Technologies GmbH
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
//

package errors

// nonRepeatableError is error that will not be retried
type nonRetryableError struct {
	reason string
}

// NewNonRetryable creates new nonRepeatableError that will not be retried
func NewNonRetryable(reason string) *nonRetryableError {
	return &nonRetryableError{
		reason: reason,
	}
}

func (err *nonRetryableError) Error() string {
	return err.reason
}

func (*nonRetryableError) DoNotRepeat() {}

// IsNonRetryable checks if error is non retryable
func IsNonRetryable(err error) bool {
	type nonRepeater interface {
		DoNotRepeat()
	}
	_, ok := err.(nonRepeater)
	return ok
}
