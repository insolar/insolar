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

package testutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecorderWriter(t *testing.T) {
	recorder := NewRecoder()
	fmt.Fprintln(recorder, "line1")
	fmt.Fprintf(recorder, "line2\nline3")
	fmt.Fprintf(recorder, "line4\nabc line5")
	expect := []string{
		"line1",
		"line2",
		"line3",
		"line4",
		"abc line5",
	}
	assert.Equal(t, expect, recorder.Items())
}
