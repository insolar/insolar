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

func TestController_BeginPulse(t *testing.T) {
	t.Parallel()
	chCancel := make(chan struct{})
	chBegin := make(chan struct{})
	controller := Controller{
		cancel:     chCancel,
		canBegin:   chBegin,
		canProcess: make(chan struct{}),
	}

	controller.BeginPulse()
	require.NotEqual(t, chBegin, controller.canBegin)
	require.NotEqual(t, chCancel, controller.cancel)
	select {
	case <-chBegin:
	default:
		t.Fatal("canBegin channel should be closed")
	}
	select {
	case <-controller.canProcess:
	default:
		t.Fatal("canProcess channel should be closed")
	}
}

func TestController_ClosePulse(t *testing.T) {
	t.Parallel()
	chCancel := make(chan struct{})
	chBegin := make(chan struct{})
	controller := Controller{
		cancel:     chCancel,
		canBegin:   chBegin,
		canProcess: make(chan struct{}),
	}

	controller.ClosePulse()
	require.Equal(t, chBegin, controller.canBegin)
	require.Equal(t, chCancel, controller.cancel)
	select {
	case <-chCancel:
	default:
		t.Fatal("close channel should be closed")
	}
	select {
	case <-controller.canProcess:
		t.Fatal("close channel should be closed")
	default:
	}
}
