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

package thread

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewController(t *testing.T) {
	t.Parallel()
	c := NewController()
	require.NotNil(t, c)
	require.NotNil(t, c.cancel)
}

func TestController_Cancel(t *testing.T) {
	t.Parallel()
	ch := make(chan struct{})
	controller := Controller{
		cancel: ch,
	}
	var expected <-chan struct{} = ch
	require.Equal(t, expected, controller.Cancel())
}

func TestController_Pulse(t *testing.T) {
	t.Parallel()
	ch := make(chan struct{})
	controller := Controller{
		cancel: ch,
	}
	var unexpected <-chan struct{} = ch
	controller.Pulse()
	require.NotEqual(t, unexpected, controller.cancel)
	select {
	case <-ch:
	default:
		t.Fatal("cancel channel should be closed")
	}
}
